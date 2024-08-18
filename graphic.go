package gelly

import (
	"image"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	Viewport  Vector2
	Transform Transform
	Surface   *ebiten.Image
}

func NewFollowingCamera(width int, height int, vWidth int, vHeight int) Camera {
	return Camera{
		Viewport: Vector2{X: float64(vWidth), Y: float64(vHeight)},
		Surface:  ebiten.NewImage(width, height),
	}
}

func (c *Camera) Follow(position Vector2) {
	c.Transform.Position.X = position.X - (c.Viewport.X / 2)
	c.Transform.Position.Y = position.Y - (c.Viewport.Y / 2)
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}

	m.Translate(-c.Transform.Position.X, -c.Transform.Position.Y)
	m.Translate(-c.Viewport.X*0.5, -c.Viewport.Y*0.5)
	m.Scale(
		math.Pow(1.01, float64(c.Transform.Scale.X)),
		math.Pow(1.01, float64(c.Transform.Scale.Y)),
	)
	m.Rotate(DegToRad(c.Transform.Rotation))
	m.Translate(c.Viewport.X*0.5, c.Viewport.Y*0.5)

	return m
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		// When scaling it can happened that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}

func (c *Camera) CursorPosition() (float64, float64) {
	return c.ScreenToWorld(ebiten.CursorPosition())
}

func (c *Camera) Draw(r *ebiten.Image) {
	r.DrawImage(c.Surface, &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	})
}

type Sprite struct {
	Centered bool
	FlipH    bool
	FlipV    bool
	Region   Rect
	// Offset    Vector2
	Image     *ebiten.Image
	Transform Transform
}

func (s Sprite) Draw(r *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	if s.FlipH {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(s.Region.W), 0)
	}

	if s.FlipV {
		op.GeoM.Scale(1, -1)
		op.GeoM.Translate(0, float64(s.Region.H))
	}

	if s.Centered {
		op.GeoM.Translate(-float64(s.Region.W)/2, -float64(s.Region.H)/2)
	}

	var sub *ebiten.Image

	if s.Region.W > 0 && s.Region.H > 0 {
		sub = s.Image.SubImage(image.Rect(
			int(s.Region.X),
			int(s.Region.Y),
			int(s.Region.X)+int(s.Region.W),
			int(s.Region.Y)+int(s.Region.H),
		)).(*ebiten.Image)
	} else {
		sub = s.Image
	}

	op.GeoM.Scale(math.Pow(1.01, s.Transform.Scale.X), math.Pow(1.01, s.Transform.Scale.Y))
	op.GeoM.Rotate(DegToRad(s.Transform.Rotation))
	op.GeoM.Translate(s.Transform.Position.X, s.Transform.Position.Y)

	r.DrawImage(sub, op)
}

type Tile struct {
	Src Rect
	Dst Rect
	// TODO: add param flip
}

type Tilemap struct {
	Tilesize  int
	Tiles     []Tile
	Image     *ebiten.Image
	Transform Transform
}

func (tm *Tilemap) Draw(r *ebiten.Image) {

	for _, tile := range tm.Tiles {
		op := &ebiten.DrawImageOptions{}

		scaleX := math.Pow(1.01, tm.Transform.Scale.X)
		scaleY := math.Pow(1.01, tm.Transform.Scale.Y)

		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Rotate(DegToRad(tm.Transform.Rotation))
		op.GeoM.Translate(float64(tile.Dst.X)*scaleX, float64(tile.Dst.Y)*scaleY)

		sub := tm.Image.SubImage(image.Rect(
			int(tile.Src.X),
			int(tile.Src.Y),
			int(tile.Src.X+tile.Src.W),
			int(tile.Src.Y+tile.Src.H),
		)).(*ebiten.Image)

		r.DrawImage(sub, op)
	}
}

type Animation struct {
	Frames    []Rect
	FrameRate time.Duration
}

type Asprite struct {
	animations       map[string]Animation
	img              *ebiten.Image
	currentName      string
	currentAnimation Animation
	currentFrame     int
	elapsedTime      time.Duration
	FlipH            bool
	FlipV            bool
	Centered         bool
	Transform        Transform
}

func NewAsprite(img *ebiten.Image, initial string, animations map[string]Animation) Asprite {
	return Asprite{
		animations:       animations,
		img:              img,
		currentName:      initial,
		currentAnimation: animations[initial],
	}
}

func (a *Asprite) Change(animation string) {
	if a.currentName == animation {
		return
	}

	anim, ok := a.animations[animation]
	if !ok {
		log.Fatalf("no animation define for %s", animation)
	}

	a.currentName = animation
	a.currentAnimation = anim
	a.currentFrame = 0
	a.elapsedTime = 0
}

func (a *Asprite) Update(dt time.Duration) {
	a.elapsedTime += dt
	if a.elapsedTime > a.currentAnimation.FrameRate {
		a.elapsedTime -= a.currentAnimation.FrameRate
		a.currentFrame++
		if a.currentFrame >= len(a.currentAnimation.Frames) {
			a.currentFrame = 0
		}
	}
}

func (a *Asprite) Draw(r *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	rect := a.currentAnimation.Frames[a.currentFrame]

	if a.FlipH {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(rect.W), 0)
	}

	if a.FlipV {
		op.GeoM.Scale(1, -1)
		op.GeoM.Translate(0, float64(rect.H))
	}

	if a.Centered {
		op.GeoM.Translate(-float64(rect.W)/2, -float64(rect.H)/2)
	}

	sub := a.img.SubImage(image.Rect(
		int(rect.X),
		int(rect.Y),
		int(rect.X)+int(rect.W),
		int(rect.Y)+int(rect.H),
	)).(*ebiten.Image)

	op.GeoM.Scale(math.Pow(1.01, a.Transform.Scale.X), math.Pow(1.01, a.Transform.Scale.Y))
	op.GeoM.Rotate(DegToRad(a.Transform.Rotation))
	op.GeoM.Translate(a.Transform.Position.X, a.Transform.Position.Y)

	r.DrawImage(sub, op)
}
