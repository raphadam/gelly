package main

import (
	"context"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/raphadam/gelly"
)

func main() {
	ctx := context.Background()

	log.Fatal(gelly.Serve(
		ctx,                 // your own context
		"127.0.0.1:9090",    // addr to serve
		time.Millisecond*20, // initial tickrate
		&GameServer{},       // your game server
	))

}

// the server manage the different game rooms and dispatch the incomming connections
type GameServer struct {
	Room1 *gelly.Room
	Room2 *gelly.Room
}

// called once on startup to setup state and rooms
func (g *GameServer) Init(s *gelly.Server) {

	// each room is running in its own goroutine
	g.Room1 = s.CreateRoom(time.Millisecond*30, &GameRoom{})
	g.Room2 = s.CreateRoom(time.Millisecond*30, &GameRoom{})

	// you can create any type of room
	// g.WaitingRoom = s.CreateRoom(time.Second*5, &WaitingRoom{})
	// g.AfkRoom = s.CreateRoom(time.Second*5, &AfkRoom{})
}

// when an user connect to the server this method is called
func (g *GameServer) FindRoom(s *gelly.Server, conn *gelly.Conn) {
	// here you can make decision to close the incomming conn
	// if reason conn.Close()

	// or join a room
	g.Room1.Join(conn)
}

func (g *GameServer) Tick(s *gelly.Server, dt time.Duration) {
	// you can broadcast a message to any room
	// g.Room1.Broadcast(&gelly.CustomUserMessage{})

	// you can change the tickrate of a room
	// g.Room1.ChangeTickrate(time.Hour * 3)
}

// when an user leave a room or is kicked this method is called
func (g *GameServer) LeftRoom(s *gelly.Server, conn *gelly.Conn, room *gelly.Room, reason error) {
}

// called once on close to dispose data
func (g *GameServer) Dispose(s *gelly.Server) {
}

// The room is where you would write the actual game logic
// but can be used for anything
type GameRoom struct {
	users map[*gelly.Conn]any
}

func (gr *GameRoom) Init(r *gelly.Room) {
	gr.users = make(map[*gelly.Conn]any)
}

func (gr *GameRoom) Join(r *gelly.Room, conn *gelly.Conn) {
	gr.users[conn] = struct{}{}

	// you can still decide here to kick the conn at any moment
	// r.Kick(conn)
}

// when any user that joined the room this method is called
func (gr *GameRoom) Message(r *gelly.Room, conn *gelly.Conn, msg gelly.Message) {
	// you can type switch to get the underlying message type
	// switch m := msg.(type) {
	// case *CustomUserMessage:
	// 	m.Something = 32
	// }

	// you can write a message to any conn conccurently
	// conn.Write(&CustomMessageResponse{text: "hello"})

	// you can update your internal game state here
}

func (gr *GameRoom) Tick(r *gelly.Room, dt time.Duration) {
}

func (gr *GameRoom) Left(r *gelly.Room, conn *gelly.Conn) {
}

func (gr *GameRoom) Dispose(r *gelly.Room) {
}

// var MainScene = gelly.Scene{
// 	Name:   "MainScene",
// 	Layers: []gelly.Layer{&GameUI{}, &GameClient{}},
// }

// log.Fatal(gelly.Run(MainScene))

type GameClient struct {
	camera gelly.Camera
	player gelly.Vector2
}

// called once on startup to load and setup game
func (g *GameClient) Init(c *gelly.Client) {
	g.camera = gelly.NewFollowingCamera(1280, 720, 1280, 720)
	g.player = gelly.Vector2{X: 100, Y: 100}

	// load file ...
	// setup game state...
}

func (g *GameClient) Message(c *gelly.Client, msg gelly.Message) bool {
	// handle user custom messages
	// switch msg.(type) {
	// }

	// write message to other layers
	// c.Write(gelly.Local, &CustomUserMessage{})

	// return true to capture the message
	return false
}

func (g *GameClient) Update(c *gelly.Client, dt time.Duration) {
	g.player.X += 1

	g.camera.Follow(g.player)
}

func (g *GameClient) Draw(r *ebiten.Image) {
	sprite := gelly.Sprite{
		Centered: true,
		FlipH:    true,
		Transform: gelly.Transform{
			Position: g.player,
		},
		// Image: CustomUserImg,
	}

	// draw on the camera surface
	sprite.Draw(g.camera.Surface)

	// draw the camera on the ebiten image
	g.camera.Draw(r)
}

// called once on closing for cleaning
func (g *GameClient) Dispose(c *gelly.Client) {
}
