package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"io"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/raphadam/gelly"
	"github.com/solarlune/ldtkgo"
	"golang.org/x/image/font/gofont/goregular"
)

// //go:embed assets/Traps/Spikes/Idle.png
// var SpikeBytes []byte
// var SpikeSprite gelly.Sprite

//go:embed assets/maskman.png
var MaskManBytes []byte
var MaskManImg *ebiten.Image
var MaskManAsprite gelly.Asprite

//go:embed assets/strawberries.png
var StrawberriesBytes []byte
var StrawberriesImg *ebiten.Image
var StrawberriesAsprite gelly.Asprite

//go:embed assets/warinja.ldtk
var WarinjaBytes []byte
var WarinjaProject *ldtkgo.Project

//go:embed assets/background.png
var BackgroundBytes []byte
var BackgroundImg *ebiten.Image
var Background gelly.Tilemap

//go:embed assets/keyboard_arrow_left.png
var KeyboardArrowLeftBytes []byte
var KeyboardArrowLeftImg *ebiten.Image

//go:embed assets/keyboard_arrow_right.png
var KeyboardArrowRightBytes []byte
var KeyboardArrowRightImg *ebiten.Image

//go:embed assets/keyboard_arrow_up.png
var KeyboardArrowUpBytes []byte
var KeyboardArrowUpImg *ebiten.Image

//go:embed assets/checkpoints.png
var CheckpointBytes []byte
var CheckpointImg *ebiten.Image
var CheckpointAsprite gelly.Asprite

//go:embed assets/collected.png
var CollectedBytes []byte
var CollectedImg *ebiten.Image
var CollectedAsprite gelly.Asprite

//go:embed assets/jump.wav
var JumpBytes []byte
var JumpPlayer *audio.Player

//go:embed assets/pickup.wav
var PickUpBytes []byte
var PickUpPlayer *audio.Player

//go:embed assets/win.wav
var WinBytes []byte
var WinPlayer *audio.Player

// var CollectedImg *ebiten.Image
// var CollectedAsprite gelly.Asprite

var IntGrid []*ldtkgo.Integer

var Goface *text.GoTextFaceSource

func init() {
	audioContext := audio.NewContext(44000)

	stream, err := wav.DecodeWithSampleRate(44000, bytes.NewReader(JumpBytes))
	if err != nil {
		log.Fatal("unable to make fo text")
	}
	data, err := io.ReadAll(stream)
	if err != nil {
		log.Fatal("unable to read music")
	}
	JumpPlayer = audioContext.NewPlayerFromBytes(data)

	stream, err = wav.DecodeWithSampleRate(44000, bytes.NewReader(PickUpBytes))
	if err != nil {
		log.Fatal("unable to make fo text")
	}
	data, err = io.ReadAll(stream)
	if err != nil {
		log.Fatal("unable to read music")
	}
	PickUpPlayer = audioContext.NewPlayerFromBytes(data)

	stream, err = wav.DecodeWithSampleRate(44000, bytes.NewReader(WinBytes))
	if err != nil {
		log.Fatal("unable to make fo text")
	}
	data, err = io.ReadAll(stream)
	if err != nil {
		log.Fatal("unable to read music")
	}
	WinPlayer = audioContext.NewPlayerFromBytes(data)

	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal("unable to make fo text")
	}
	Goface = s

	// img, _, err := image.Decode(bytes.NewReader(SpikeBytes))
	// if err != nil {
	// 	log.Fatal("unable to decode image", img)
	// }
	// SpikeSprite = gelly.Sprite{
	// 	Image:  ebiten.NewImageFromImage(img),
	// 	Region: gelly.Rect{W: 16, H: 16},
	// }

	// keyboards
	img, _, err := image.Decode(bytes.NewReader(KeyboardArrowLeftBytes))
	if err != nil {
		log.Fatal("unable to decode image", img)
	}
	KeyboardArrowLeftImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(KeyboardArrowRightBytes))
	if err != nil {
		log.Fatal("unable to decode image", img)
	}
	KeyboardArrowRightImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(KeyboardArrowUpBytes))
	if err != nil {
		log.Fatal("unable to decode image", img)
	}
	KeyboardArrowUpImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(StrawberriesBytes))
	if err != nil {
		log.Fatal("unable to decode image", img)
	}
	StrawberriesImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(CollectedBytes))
	if err != nil {
		log.Fatal("unable to decode image", img)
	}
	CollectedImg = ebiten.NewImageFromImage(img)

	rects := [38]gelly.Rect{}
	for i := range 38 {
		rects[i] = gelly.Rect{
			X: float64(i * 32), Y: 0,
			W: 32, H: 32,
		}
	}

	CollectedAsprite = gelly.NewAsprite(CollectedImg, "idle", map[string]gelly.Animation{
		"idle": {FrameRate: 100 * time.Millisecond, Frames: rects[:5]},
	})

	StrawberriesAsprite = gelly.NewAsprite(StrawberriesImg, "idle", map[string]gelly.Animation{
		"idle": {FrameRate: 100 * time.Millisecond, Frames: rects[:17]}, // 11
	})

	img, _, err = image.Decode(bytes.NewReader(MaskManBytes))
	if err != nil {
		log.Fatal("unable to decode image", img)
	}
	MaskManImg = ebiten.NewImageFromImage(img)

	MaskManAsprite = gelly.NewAsprite(MaskManImg, "idle", map[string]gelly.Animation{
		"doubleJump": {FrameRate: 100 * time.Millisecond, Frames: rects[:6]},    // 6
		"walk":       {FrameRate: 100 * time.Millisecond, Frames: rects[26:38]}, // 12
		"fall":       {FrameRate: 100 * time.Millisecond, Frames: rects[6:7]},   // 1
		"hit":        {FrameRate: 100 * time.Millisecond, Frames: rects[7:14]},  // 7
		"idle":       {FrameRate: 100 * time.Millisecond, Frames: rects[14:25]}, // 11
		"jump":       {FrameRate: 400 * time.Millisecond, Frames: rects[25:26]}, // 1
	})

	img, _, err = image.Decode(bytes.NewReader(BackgroundBytes))
	if err != nil {
		log.Fatal("unable to decode image", img)
	}
	BackgroundImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(CheckpointBytes))
	if err != nil {
		log.Fatal("unable to decode image", img)
	}
	CheckpointImg = ebiten.NewImageFromImage(img)

	checkpointRects := [37]gelly.Rect{}
	for i := range 37 {
		checkpointRects[i] = gelly.Rect{X: float64(i * 64), W: 64, H: 64}
	}

	img, _, err = image.Decode(bytes.NewReader(StrawberriesBytes))
	if err != nil {
		log.Fatal("unable to decode image", img)
	}
	StrawberriesImg = ebiten.NewImageFromImage(img)
	// ddd

	CheckpointAsprite = gelly.NewAsprite(CheckpointImg, "idle", map[string]gelly.Animation{
		"opening": {FrameRate: 100 * time.Millisecond, Frames: checkpointRects[:26]},   // 10
		"idle":    {FrameRate: 100 * time.Millisecond, Frames: checkpointRects[26:27]}, // 1
		"finish":  {FrameRate: 100 * time.Millisecond, Frames: checkpointRects[27:]},   // 26
	})

	proj, err := ldtkgo.Read(WarinjaBytes)
	if err != nil {
		log.Fatal("unable to read ldtk", img)
	}
	WarinjaProject = proj

	level := WarinjaProject.LevelByIdentifier("Level_0")

	background := level.LayerByIdentifier("Background")
	Background = LoadTilemapFromLdtk(16, background, BackgroundImg)

	IntGrid = level.LayerByIdentifier("IntGrid").IntGrid
}

func LoadTilemapFromLdtk(tilesize int, layer *ldtkgo.Layer, tileset *ebiten.Image) gelly.Tilemap {
	if layer == nil {
		log.Fatal("terrain 0 not found in the ldtk file")
	}

	tiles := make([]gelly.Tile, len(layer.Tiles))

	for i, tile := range layer.Tiles {
		tiles[i] = gelly.Tile{
			Src: gelly.Rect{
				// TODO: make this customisable
				X: float64(tile.Src[0]),
				Y: float64(tile.Src[1]),
				W: float64(tilesize), H: float64(tilesize),
			},
			Dst: gelly.Rect{
				X: float64(tile.Position[0]),
				Y: float64(tile.Position[1]),
				W: float64(tilesize), H: float64(tilesize),
			},
		}
	}

	return gelly.Tilemap{
		Tilesize: tilesize,
		Image:    tileset,
		Tiles:    tiles,
	}
}
