package gelly

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

var (
	ErrUserDisconnected error = fmt.Errorf("user disconnected")
	ErrUserKicked       error = fmt.Errorf("user kicked")
)

type internalMessageType int

const (
	userReceiveMessage internalMessageType = iota
	userSendMessage
	userConnected
	userDisconnected
	userKickFromRoom
	userJoinRoom
)

type internalMessage struct {
	t      internalMessageType
	conn   *Conn
	m      Message
	r      *Room
	reason error
}

type RoomBlueprint interface {
	Init(r *Room)
	Join(r *Room, conn *Conn)
	Message(r *Room, conn *Conn, msg Message)
	Tick(r *Room, dt time.Duration)
	Left(r *Room, conn *Conn)
	Dispose(r *Room)
}

type Room struct {
	ctx      context.Context
	conns    map[*Conn]context.CancelFunc
	joinc    chan *Conn
	rb       RoomBlueprint
	readc    chan internalMessage
	servc    chan<- internalMessage
	tickrate time.Duration
}

// TODO: impl
func (r *Room) ChangeTickrate(new time.Duration) {
}

func (r *Room) Join(conn *Conn) {
	// should called the internal room data
	r.joinc <- conn
}

// TODO: send back to the server the conn
func (r *Room) Kick(conn *Conn) {
	cancel, ok := r.conns[conn]
	if !ok {
		log.Fatal("conn not found")
	}

	delete(r.conns, conn)
	r.rb.Left(r, conn)
	cancel()
}

func (r *Room) Emit(conn *Conn, msg Message) {
	for c := range r.conns {
		if c == conn {
			continue
		}

		conn.Write(msg)
	}
}

// TODO: what happen if the server try to broadcast?
func (r *Room) Broadcast(msg Message) {
	for conn := range r.conns {
		conn.Write(msg)
	}
}

func (r *Room) Close() {
}

func (r *Room) handleJoin(conn *Conn) {
	ctx, cancel := context.WithCancel(r.ctx)
	r.conns[conn] = cancel
	r.rb.Join(r, conn)

	go func() {
		defer log.Println("HANDLE JOIN CLOSED")

		for {
			select {
			case <-ctx.Done():
				r.servc <- internalMessage{
					t:      userKickFromRoom,
					conn:   conn,
					r:      r,
					reason: ErrUserKicked,
				}
				return

			case msg, ok := <-conn.readc:
				if !ok {
					decoMsg := internalMessage{
						t:      userDisconnected,
						conn:   conn,
						reason: ErrUserDisconnected,
					}

					r.readc <- decoMsg
					r.servc <- decoMsg
					return
				}

				r.readc <- msg
			}
		}
	}()
}

// TODO: add tickrate
// TODO: add cancel
func (r *Room) serve() {
	r.joinc = make(chan *Conn)
	r.readc = make(chan internalMessage)
	r.conns = make(map[*Conn]context.CancelFunc)
	tick := time.NewTicker(r.tickrate)

	r.rb.Init(r)

	go func() {
		defer close(r.joinc)
		defer close(r.readc)

		for {
			select {
			case <-tick.C:
				// TODO: pass the good delta time
				r.rb.Tick(r, time.Millisecond*20)

			// TODO: maybe check if already connected
			case conn := <-r.joinc:
				r.handleJoin(conn)

			case msg := <-r.readc:
				switch msg.t {
				case userReceiveMessage:
					r.rb.Message(r, msg.conn, msg.m)

				case userDisconnected:
					delete(r.conns, msg.conn)
					r.rb.Left(r, msg.conn)

				default:
					log.Fatal("room default error")
				}
			}
		}
	}()
}

type ServerBlueprint interface {
	Init(s *Server)
	FindRoom(s *Server, conn *Conn)
	Tick(s *Server, dt time.Duration)
	LeftRoom(s *Server, conn *Conn, room *Room, reason error)
	Dispose(s *Server)
}

type Server struct {
	tick    time.Duration
	addr    string
	sb      ServerBlueprint
	ctx     context.Context
	internc chan internalMessage
}

func (s *Server) Broadcast(msg Message) {
}

func (s *Server) CreateRoom(tickrate time.Duration, rb RoomBlueprint) *Room {
	room := &Room{
		rb:       rb,
		ctx:      s.ctx,
		servc:    s.internc,
		tickrate: tickrate,
	}

	room.serve()
	return room
}

func (s *Server) Close() {
}

func (s *Server) serve() {
	s.internc = make(chan internalMessage)

	// TODO: maybe do it somehere else
	// TODO: must call the dispose
	s.sb.Init(s)

	// TODO: should close the ticker
	// TODO: should be able to change the ticker time
	ticker := time.NewTicker(s.tick)

	go func() {
		defer close(s.internc)

		// TODO: case cancel
		for {
			select {
			case <-ticker.C:
				// TODO: should give real dt
				s.sb.Tick(s, 20*time.Millisecond)

			case msg := <-s.internc:
				switch msg.t {
				case userConnected:
					s.sb.FindRoom(s, msg.conn)

				// TODO: check this one
				case userDisconnected:
					s.sb.LeftRoom(s, msg.conn, msg.r, msg.reason)

				case userKickFromRoom:
					s.sb.LeftRoom(s, msg.conn, msg.r, msg.reason)

				default:
					log.Fatal("not handled message type")
				}
			}
		}
	}()
}

func Serve(ctx context.Context, addr string, tick time.Duration, sb ServerBlueprint) error {
	server := &Server{
		addr: addr,
		sb:   sb,
		ctx:  ctx,
		tick: tick,
	}
	server.serve()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			log.Fatal("error on accept", err)
		}

		gconn := &Conn{
			wsc: conn,
		}
		gconn.handle(ctx)
		server.internc <- internalMessage{
			t:    userConnected,
			conn: gconn,
		}
	})

	return http.ListenAndServe(addr, nil)
}

// TODO: maybe should not display the close method
// TODO: maybe use the context
type Conn struct {
	wsc    *websocket.Conn
	writec chan internalMessage
	readc  chan internalMessage
	cancel context.CancelFunc
}

func (c *Conn) Write(m Message) {
	c.writec <- internalMessage{
		t: userSendMessage,
		m: m,
	}
}

// TODO: implement close conn
// TODO: check if calling close from room is working
func (c *Conn) Close() {
	c.cancel()
}

func (c *Conn) handle(ctx context.Context) {
	c.readc = make(chan internalMessage)
	c.writec = make(chan internalMessage)

	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	go func() {
		defer func() {
			close(c.readc)
			c.cancel()
			log.Println("READ CHANNEL CLOSED")
		}()

		for {
			_, data, err := c.wsc.Read(ctx)
			if err != nil {
				// TODO: should treat the different king of disconnection
				// here when the client disconnected
				return
			}

			decoder := gob.NewDecoder(bytes.NewReader(data))
			packet := &packet{}

			err = decoder.Decode(packet)
			if err != nil {
				log.Fatal("unable to decode")
			}

			// TODO: may never stop cause trying to write to close channel
			c.readc <- internalMessage{
				t:    userReceiveMessage,
				conn: c,
				m:    packet.M,
			}
		}

	}()

	go func() {
		defer func() {
			close(c.writec)
			c.wsc.CloseNow()
			log.Println("WRITE CHANNEL CLOSED")
		}()

		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-c.writec:
				switch msg.t {
				case userSendMessage:
					if msg.m == nil {
						log.Fatalf("cannot send nil message")
					}

					var buff bytes.Buffer

					encoder := gob.NewEncoder(&buff)
					packet := packet{M: msg.m}

					err := encoder.Encode(packet)
					if err != nil {
						log.Fatalf("unable to marshal the message %s", err)
					}

					err = c.wsc.Write(ctx, websocket.MessageBinary, buff.Bytes())
					if err != nil {
						// TODO: check the different ways of stopping
						log.Fatal("error on writing goroutine")
					}

				default:
					log.Fatal("don't now this one")
				}
			}
		}

	}()
}
