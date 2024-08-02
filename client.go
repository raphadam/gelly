package gelly

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"nhooyr.io/websocket"
)

// type ComponentType int

// type Component interface {
// 	Type() ComponentType
// }

// const (
// 	SPRITE_COMPONENT ComponentType = iota
// 	MAX_COMPONENT
// )

// type Sprite struct {
// }

// func (s Sprite) Change(name string) {
// }

// func (s Sprite) FlipH(ok bool) {
// }

// func (s Sprite) Type() ComponentType {
// 	return SPRITE_COMPONENT
// }

// type EntityBase struct {
// 	name string
// }

// type Player struct {
// 	EntityBase

// 	HP int
// }

// func CreateEntity[E EntityBase]() {
// }

// func Exec() {
// 	p := Player{}
// 	CreateEntity[Player]()
// }

// type signature uint64

// type profile struct {
// 	e Entity
// 	s signature
// 	t Transform
// }

// type EventType int

// type Event interface {
// 	Type() EventType
// }

// type Handler func(cmd *Command, e Entity, evt Event)

// type Command struct {
// 	count      int
// 	profiles   []profile
// 	components [MAX_COMPONENT][]Component
// 	subs       map[EventType][]Handler
// }

// func NewCommand() *Command {
// 	return &Command{
// 		subs: make(map[EventType][]Handler),
// 	}
// }

// func Child(cmd *Command, e Entity, setup Setup) {
// }

// func Create(cmd *Command, setup Setup) {
// 	e := Entity{id: cmd.count}
// 	cmd.profiles = append(cmd.profiles, profile{e: e})
// 	cmd.count++
// 	setup(cmd, e)
// }

// func Destroy(cmd *Command, e Entity) {
// }

// func Add[C Component](cmd *Command, e Entity, c C) {
// 	cmd.components[c.Type()] = append(cmd.components[c.Type()], c)
// }

// func Subscribe[E Event](cmd *Command, e Entity, h Handler) {
// 	var d E
// 	cmd.subs[d.Type()] = append(cmd.subs[d.Type()], h)
// }

// func GetTransform(cmd *Command, e Entity) Transform {
// 	return Transform{}
// }

// func SetTransform(cmd *Command, e Entity, t Transform) {
// }

// func Publish(cmd *Command, e Entity, evt Event) {
// }

// var handlers []Handler[Event]

// type HandlerId int

// func Register[T Event](h Handler[T]) HandlerId {
// 	return HandlerId(1)
// }

// const (
// 	PLAYER_MOVED gelly.EventType = iota
// 	TICK_UPDATED
// )

// type PlayerMoved struct {
// 	Pos gelly.Vector2
// }

// func (e PlayerMoved) Type() gelly.EventType {
// 	return PLAYER_MOVED
// }

// type TickUpdated struct {
// 	Dt time.Duration
// }

// func (e TickUpdated) Type() gelly.EventType {
// 	return TICK_UPDATED
// }

// func Player(cmd *gelly.Command, e gelly.Entity) {
// 	gelly.Add(cmd, e, gelly.Sprite{})

// 	gelly.Subscribe[PlayerMoved](cmd, e, OnPlayerMoved)
// 	gelly.Subscribe[TickUpdated](cmd, e, OnTickPlayer)
// }

// func OnPlayerMoved(cmd *gelly.Command, e gelly.Entity, evt gelly.Event) {
// 	// transform := gelly.GetTransform(cmd, e)
// 	// transform.Position.X += 10

// 	// gelly.Publish[PlayerMoved](cmd, e, PlayerMoved{})
// }

// func OnTickPlayer(cmd *gelly.Command, e gelly.Entity, evt gelly.Event) {
// 	// transform := gelly.GetTransform(cmd, e)
// 	// transform.Position.X += 10
// 	// gelly.Publish[PlayerMoved](cmd, e, PlayerMoved{})
// }

// type Entity struct {
// 	id int
// 	// gen uint
// }

// type Setup func(cmd *Command, e Entity)

// NORMAL

const (
	clientDefaultTitle        = "Gelly - Go game library"
	clientDefaultWindowWidth  = 1280
	clientDefaultWindowHeight = 720
	clientDefaultLayoutWidth  = 1280 / 4
	clientDefaultLayoutHeight = 720 / 4
	clientDefaultVsync        = false
	clientDefaultFullscreen   = false
	clientDefaultResizable    = false
)

type Message interface {
	IsMessage()
}

func init() {
	gob.Register(packet{})
	gob.Register(SocketConnected{})
	gob.Register(SocketDisconnected{})
}

type packet struct {
	M Message
}

func RegisterMessage(m Message) {
	gob.Register(m)
}

type SocketConnected struct {
	Addr string
}

func (m *SocketConnected) IsMessage() {
}

type SocketDisconnected struct {
	Reason string
	Err    error
}

func (m *SocketDisconnected) IsMessage() {
}

// TODO: maybe tell if the message is from the network or local
type Layer interface {
	Init(c *Client)
	Message(c *Client, msg Message) bool
	Update(c *Client, dt time.Duration)
	Draw(r *ebiten.Image)
	Dispose(c *Client)
}

type Scene struct {
	Name   string
	Layers []Layer

	front []Message
	back  []Message
}

func (s *Scene) write(msg Message) {
	s.front = append(s.front, msg)
}

func (s *Scene) update(c *Client, dt time.Duration) {
	temp := s.front
	s.front = s.back
	s.back = temp

	for _, event := range s.back {

	innerLoop:
		for _, l := range s.Layers {
			shouldCapture := l.Message(c, event)

			if shouldCapture {
				break innerLoop
			}
		}

	}

	s.back = s.back[:0]

	for _, layer := range s.Layers {
		layer.Update(c, dt)
	}
}

func (s *Scene) draw(r *ebiten.Image) {
	for i := len(s.Layers) - 1; i >= 0; i-- {
		// TODO: give a plain image instead of the screen for performance
		s.Layers[i].Draw(r)
	}
}

func (s *Scene) init(c *Client) {
	for _, layer := range s.Layers {
		layer.Init(c)
	}
}

func (s *Scene) dispose(c *Client) {
	for _, layer := range s.Layers {
		layer.Dispose(c)
	}
}

// TODO: cancel ?
type socket struct {
	wsc   *websocket.Conn
	mux   sync.Mutex
	front []Message
	back  []Message
}

func (s *socket) connect(addr string) error {
	if s.wsc != nil {
		return fmt.Errorf("already conncted")
	}

	wsc, _, err := websocket.Dial(context.Background(), addr, nil)
	if err != nil {
		return fmt.Errorf("unable to dial the server %w", err)
	}
	s.wsc = wsc

	s.handleClientConnection(s.wsc)
	return nil
}

func (s *socket) handleClientConnection(conn *websocket.Conn) {
	go func() {
		defer conn.CloseNow()

		for {
			_, data, err := conn.Read(context.Background())
			if err != nil {
				s.wsc = nil

				s.mux.Lock()
				s.front = append(s.front, &SocketDisconnected{Reason: "TODO: make a reson", Err: err})
				s.mux.Unlock()
				return
			}

			decoder := gob.NewDecoder(bytes.NewReader(data))
			packet := &packet{}

			err = decoder.Decode(packet)
			if err != nil {
				log.Fatal("unable to gob decode the incomming data", err)
			}

			// TODO: maybe switch to a channel
			s.mux.Lock()
			s.front = append(s.front, packet.M)
			s.mux.Unlock()
		}
	}()
}

func (s *socket) write(m Message) error {
	if m == nil {
		return fmt.Errorf("cannot send nil message")
	}

	if s.wsc == nil {
		return fmt.Errorf("must be connected before sending any message")
	}

	// // TODO: maybe make a channel for that
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	packet := packet{M: m}

	err := encoder.Encode(packet)
	if err != nil {
		return fmt.Errorf("unable to marshal the message %w", err)
	}

	err = s.wsc.Write(context.Background(), websocket.MessageBinary, buff.Bytes())
	if err != nil {
		return fmt.Errorf("unable to write a message %w", err)
	}

	return nil
}

func (s *socket) disconnect() error {
	// if c.wsc == nil {
	// 	return fmt.Errorf("must be connected before disconnecting")
	// }

	// err := c.wsc.Close(websocket.StatusNormalClosure, "")
	// if err != nil {
	// 	return fmt.Errorf("error on close %w", err)
	// }

	// c.wsc = nil
	return nil
}

// TODO: maybe need to provide a connection closed or open with messages
func (s *socket) swapQueues() {
	s.mux.Lock()
	temp := s.front
	s.front = s.back
	s.back = temp
	s.mux.Unlock()
}

type delatatime struct {
	previousTime time.Time
	currentTime  time.Time
	dt           time.Duration
}

func (t *delatatime) update() {
	t.currentTime = time.Now()
	t.dt = t.currentTime.Sub(t.previousTime)
	t.previousTime = t.currentTime
}

type Client struct {
	scene  Scene
	input  input
	socket socket
	dt     delatatime

	changeSceneRequested bool
	nextSceneToChange    Scene

	shoudClose   bool
	layoutHeight int
	layoutWidth  int
}

type WriteMethod int

const (
	Local WriteMethod = iota
	Online
	Both
)

// TODO: maybe give back the error message
// TODO: maybe make a write sure method
func (c *Client) Write(method WriteMethod, m Message) {
	switch method {

	case Local:
		c.scene.write(m)

	case Online:
		err := c.socket.write(m)
		if err != nil {
			log.Fatal("unable to write a message")
		}

	case Both:
		c.scene.write(m)
		err := c.socket.write(m)
		if err != nil {
			log.Fatal("unable to write a message")
		}

	default:
		log.Fatal("unknown write message method")
	}

}

// TODO: Transition ?
// TODO: maybe done in a middle of a layer update ?
// TODO: maybe need to transfer the events ?
// TODO: maybe payload with an interface
// TODO: may some time change to a delayed command pattern
// Next Scene will be launch at the end of the update
func (c *Client) ChangeScene(next Scene) {
	c.changeSceneRequested = true
	c.nextSceneToChange = next
}

func (c *Client) Connect(addr string) error {
	return c.socket.connect(addr)
}

func (c *Client) Disconnect() error {
	return c.socket.disconnect()
}

func (c *Client) IsMouseMoved() bool {
	return c.input.mouseMoved
}

// Activate / Deactivate ?
func (c *Client) Show(l Layer) {
}

func (c *Client) SetLayoutSize(width int, height int) {
	c.layoutWidth = width
	c.layoutHeight = height
}

func (c *Client) Hide(l Layer) {
}

func (c *Client) Up(l Layer) {
}

func (c *Client) Down(l Layer) {
}

func (c *Client) Close() {
	c.shoudClose = true
}

func (c *Client) update() error {
	if c.shoudClose {
		// TODO: not working when pressing the exit cross with cursor
		c.scene.dispose(c)
		return errors.New("close requested normally")
	}

	// Deleayed scene changing should read the ChangeScene method for more explanation
	if c.changeSceneRequested {
		c.scene.dispose(c)
		c.nextSceneToChange.init(c)
		c.scene = c.nextSceneToChange
		c.changeSceneRequested = false
	}

	// CLOCK TIME
	c.dt.update()

	// // PROCESS INPUTS
	c.input.update()

	// PROCESS MESSAGES
	c.socket.swapQueues()

	// TODO: need to send the packet
	for _, msg := range c.socket.back {
		c.scene.write(msg)
	}
	c.socket.back = c.socket.front[:0]

	// UPDATE Scene
	c.scene.update(c, c.dt.dt)

	return nil
}

func (c *Client) draw(screen *ebiten.Image) {
	// DRAW MOVIE
	c.scene.draw(screen)
}

func Run(s Scene) error {
	// TODO: make more verifications
	// if s.Layer == nil {
	// 	return fmt.Errorf("an layer must exist")
	// }

	ebiten.SetWindowTitle(clientDefaultTitle)
	ebiten.SetWindowSize(clientDefaultWindowWidth, clientDefaultWindowHeight)
	ebiten.SetVsyncEnabled(clientDefaultVsync)
	ebiten.SetFullscreen(clientDefaultFullscreen)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	c := Client{
		scene: s,
		dt: delatatime{
			previousTime: time.Now(),
			currentTime:  time.Now(),
		},
		layoutWidth:  clientDefaultLayoutWidth,
		layoutHeight: clientDefaultLayoutHeight,
	}

	// TODO: check if the init should be done in other place
	c.scene.init(&c)

	return ebiten.RunGameWithOptions(&ebitenClient{client: &c}, &ebiten.RunGameOptions{
		InitUnfocused: true,
	})
}

type ebitenClient struct {
	client *Client
}

func (eb *ebitenClient) Update() error {
	return eb.client.update()
}

func (eb *ebitenClient) Draw(screen *ebiten.Image) {
	eb.client.draw(screen)
}

func (eb *ebitenClient) Layout(outsideWidth, outsideHeight int) (int, int) {
	return eb.client.layoutWidth, eb.client.layoutHeight
}

type input struct {
	xPreviousMousePos int
	yPreviousMousePos int
	mouseMoved        bool
}

func (i *input) update() {
	// Cursor position
	xpos, ypos := ebiten.CursorPosition()
	i.mouseMoved = i.xPreviousMousePos != xpos || i.yPreviousMousePos != ypos
	i.xPreviousMousePos = xpos
	i.yPreviousMousePos = ypos
}
