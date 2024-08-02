package main

// import (
// 	"log"
// 	"time"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/inpututil"
// 	"github.com/raphadam/gelly"
// )

// func main() {
// 	gelly.Run(gelly.Scene{
// 		Name:   "Base",
// 		Layers: []gelly.Layer{&Example{}},
// 	})
// }

// type Example struct {
// }

// func (e *Example) Init(c *gelly.Client) {
// 	log.Println("INIT")
// }

// func (e *Example) Message(c *gelly.Client, msg gelly.Message) bool {
// 	return false
// }

// func (e *Example) Update(c *gelly.Client, dt time.Duration) {
// 	log.Println("Update")

// 	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
// 		c.Close()
// 	}
// }

// func (e *Example) Draw(r *ebiten.Image) {
// 	log.Println("Draw")
// }

// func (e *Example) Dispose(c *gelly.Client) {
// 	log.Println("DISPOSE")
// }
