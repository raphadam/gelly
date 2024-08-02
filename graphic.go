package gelly

import (
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Key struct {
	id  int
	gen uint
	ok  bool
}

type Value[T any] struct {
	K Key
	V T
}

type cell struct {
	nfr int
	itv int
	gen uint
}

type Pool[T any] struct {
	nfr      int
	gen      uint
	len      int
	creates  []Value[T]
	destroys []Key
	cells    []cell
	Values   []Value[T]
}

func (p *Pool[T]) Create(v T) Key {
	k := Key{id: p.nfr, gen: p.gen, ok: true}
	p.creates = append(p.creates, Value[T]{
		K: k,
		V: v,
	})

	if p.nfr >= p.len {
		p.cells = append(p.cells, cell{itv: -1, gen: p.gen})
		p.len++
		p.nfr++
		return k
	}

	p.cells[p.nfr].itv = -1
	p.cells[p.nfr].gen = p.gen
	p.nfr = p.cells[p.nfr].nfr
	return k
}

func (p *Pool[T]) Get(k Key) (*T, bool) {
	if !p.valid(k) {
		return nil, false
	}

	return &p.Values[p.cells[k.id].itv].V, true
}

func (p *Pool[T]) Destroy(k Key) bool {
	if !p.valid(k) {
		return false
	}

	p.gen++
	p.destroys = append(p.destroys, k)
	return true
}

func (p *Pool[T]) For(h func(i int, k Key, v *T) bool) {
	for i := range p.Values {
		br := h(i, p.Values[i].K, &p.Values[i].V)
		if br {
			break
		}
	}
}

func (p *Pool[T]) Apply() {
	for _, k := range p.destroys {
		p.cells[k.id].nfr = p.nfr
		p.cells[k.id].gen = p.gen
		p.nfr = k.id

		last := p.Values[len(p.Values)-1]
		p.Values[p.cells[k.id].itv] = last
		p.cells[last.K.id].itv = p.cells[k.id].itv
		p.Values = p.Values[:len(p.Values)-1]
	}
	p.destroys = p.destroys[:0]

	for _, create := range p.creates {
		p.cells[create.K.id].itv = len(p.Values)
		p.Values = append(p.Values, create)
	}
	p.creates = p.creates[:0]
}

func (p *Pool[T]) valid(k Key) bool {
	if !k.ok {
		return false
	}

	if k.id >= p.len {
		return false
	}

	if p.cells[k.id].gen != k.gen {
		return false
	}

	if p.cells[k.id].itv < 0 {
		return false
	}

	return true
}

// // Pool
// type Key[T any] struct {
// 	id  int
// 	gen uint
// 	ok  bool
// }

// type Value[T any] struct {
// 	K Key
// 	V T
// 	T Transform
// }

// type cell struct {
// 	nfr int
// 	itv int
// 	gen uint
// }

// type Pool[T any] struct {
// 	nfr      int
// 	gen      uint
// 	len      int
// 	creates  []Value[T]
// 	destroys []Key
// 	cells    []cell
// 	Values   []Value[T]
// }

// // func (p *Pool[T]) SetTransform(k Key, t Transform) {
// // }

// func (p *Pool[T]) Create(v T) Key {
// 	k := Key{id: p.nfr, gen: p.gen, ok: true}
// 	p.creates = append(p.creates, Value[T]{
// 		K: k,
// 		V: v,
// 	})

// 	if p.nfr >= p.len {
// 		p.cells = append(p.cells, cell{itv: -1, gen: p.gen})
// 		p.len++
// 		p.nfr++
// 		return k
// 	}

// 	p.cells[p.nfr].itv = -1
// 	p.cells[p.nfr].gen = p.gen
// 	p.nfr = p.cells[p.nfr].nfr
// 	return k
// }

// func (p *Pool[T]) Get(k Key) (*T, bool) {
// 	if !p.valid(k) {
// 		return nil, false
// 	}

// 	return &p.Values[p.cells[k.id].itv].V, true
// }

// func (p *Pool[T]) Destroy(k Key) bool {
// 	if !p.valid(k) {
// 		return false
// 	}

// 	p.gen++
// 	p.destroys = append(p.destroys, k)
// 	return true
// }

// func (p *Pool[T]) For(h func(k Key, v *T) bool) {
// 	for i := range p.Values {
// 		br := h(p.Values[i].K, &p.Values[i].V)
// 		if br {
// 			break
// 		}
// 	}
// }

// func (p *Pool[T]) Apply() {
// 	for _, k := range p.destroys {
// 		p.cells[k.id].nfr = p.nfr
// 		p.cells[k.id].gen = p.gen
// 		p.nfr = k.id

// 		last := p.Values[len(p.Values)-1]
// 		p.Values[p.cells[k.id].itv] = last
// 		p.cells[last.K.id].itv = p.cells[k.id].itv
// 		p.Values = p.Values[:len(p.Values)-1]
// 	}
// 	p.destroys = p.destroys[:0]

// 	for _, create := range p.creates {
// 		p.cells[create.K.id].itv = len(p.Values)
// 		p.Values = append(p.Values, create)
// 	}
// 	p.creates = p.creates[:0]
// }

// func (p *Pool[T]) valid(k Key) bool {
// 	if !k.ok {
// 		return false
// 	}

// 	if k.id >= p.len {
// 		return false
// 	}

// 	if p.cells[k.id].gen != k.gen {
// 		return false
// 	}

// 	if p.cells[k.id].itv < 0 {
// 		return false
// 	}

// 	return true
// }

type RenderableType int

const (
	ASPRITE RenderableType = iota
)

type Renderable struct {
	Type     RenderableType
	Material Material
	Asprite  Asprite
}

type Material struct {
}

// type Renderer struct {
// 	pool Pool[Renderable]
// }

// func (r *Renderer) Init() {
// }

// func (r *Renderer) Create(v Renderable) Key[Renderable] {
// 	return r.pool.Create(v)
// }

// func (r *Renderer) Get(k Key[Renderable]) (*Renderable, bool) {
// 	return r.pool.Get(k)
// }

// func (r *Renderer) Update(dt time.Duration) {
// 	for i := range r.pool.Values {
// 		switch r.pool.Values[i].V.Type {
// 		case ASPRITE:
// 			r.pool.Values[i].V.Asprite.Update(dt)
// 		}
// 	}
// }

// func (r *Renderer) Draw(i *ebiten.Image) {
// 	// for i := range r.pool.Values {
// 	// 	switch r.pool.Values[i].V.Type {
// 	// 	case ASPRITE:
// 	// 		r.pool.Values[i].V.Asprite.Draw(r, r.pool.Values[i].V.Transform)
// 	// 	}
// 	// }
// }

// PHYSIC
type Physic struct {
	pool Pool[Rigidbody]
}

type BodyType int

const (
	STATIC BodyType = iota
	KINECTIC
)

type Rigidbody struct {
	BodyType BodyType
	Property Property
}

type Property struct {
	Acceleration float64
	Velocity     float64
	Friction     float64
}

// func (p *Physic) Init(spaceWidth int, spaceHeight int, cellWidth int, cellheight int) {
// 	p.space = resolv.NewSpace(spaceWidth, spaceHeight, cellWidth, cellheight)
// }

// func (p *Physic) Create(r Rect, pr Property) Key[Rigidbody] {
// 	object := resolv.NewObject(r.X, r.Y, r.W, r.H)
// 	p.space.Add(object)

// 	return p.pool.Create(Rigidbody{
// 		object:   object,
// 		Property: pr,
// 	})
// }

// func (p *Physic) Get(k Key[Rigidbody]) (*Rigidbody, bool) {
// 	return p.pool.Get(k)
// }

// func (p *Physic) Update(dt time.Duration) {
// 	p.pool.Apply()

// 	for i := range p.pool.Values {
// 		switch p.pool.Values[i].V.BodyType {
// 		case STATIC:
// 			log.Println("STATIC OBJECT")

// 		case KINECTIC:
// 			log.Println("KINETIC OBJECT")
// 		}
// 	}

// 	// if collision := entity.V.Physic.Object.Check(dy, 0); collision != nil {
// 	// 	dy = collision.ContactWithObject(collision.Objects[0]).Y
// 	// }
// 	// entity.V.Physic.Object.Position.Y += dy
// 	// entity.V.Physic.Object.Update()
// 	// entity.V.Transform.Position = gelly.Vector2(entity.V.Physic.Object.Position)
// }

// func (p *Physic) DrawDebug(r *ebiten.Image) {
// 	objects := p.space.Objects()
// 	for _, obj := range objects {
// 		vector.StrokeRect(r,
// 			float32(obj.Position.X), float32(obj.Position.Y),
// 			float32(obj.Size.X), float32(obj.Size.Y),
// 			1, ColorCyan, false,
// 		)
// 	}
// }

type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

func (r Rect) Intersect(o Rect) bool {
	return (r.X+r.W > o.X) && (r.X < o.X+o.W) && (r.Y+r.H > o.Y) && (r.Y < o.Y+o.H)
}

func DegToRad(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}

func RadToDeg(rad float64) float64 {
	return rad * (180.0 / math.Pi)
}

type Vector2 struct {
	X float64
	Y float64
}

func (p Vector2) Inside(r Rect) bool {
	return p.X >= r.X && p.X <= (r.X+r.W) && p.Y >= r.Y && p.Y <= (r.Y+r.H)
}

func (p Vector2) MulC(mul float64) Vector2 {
	return Vector2{
		X: p.X * mul,
		Y: p.Y * mul,
	}
}

func (p Vector2) Add(v Vector2) Vector2 {
	return Vector2{
		X: v.X + p.X,
		Y: v.Y + p.Y,
	}
}

// func (p Vector2) Minus(other Vector2) Vector2 {
// 	other.X -= p.X
// 	other.Y -= p.Y
// 	return other
// }

type Path struct {
	A Vector2
	B Vector2
	V float64
	// total float64
}

// func (p *Path) Approach(v float64) Vector2 {
// 	v *= v

// 	p.total += v
// 	progress := p.total / p.V
// 	dir := p.A.Minus(p.B)

// 	dir.X *= progress
// 	dir.Y *= progress

// 	return dir
// }

func MaxAbs(val float64, max float64) float64 {
	if val > max {
		val = max
	} else if val < -max {
		val = -max
	}

	return val
}

func Lerp(start float64, end float64, fraction float64) float64 {
	return start + (end-start)*Clamp(0, 1, fraction)
}

func Clamp(from float64, to float64, value float64) float64 {
	if value <= from {
		return from
	}

	if value >= to {
		return to
	}

	return value
}

// func InvLerp(start float64, end float64, fraction float64) float64 {
// 	return start + (end-start)*fraction
// }

// func Remap(start float64, end float64, fraction float64) float64 {
// 	log.Fatal("TODO: empl")
// 	return start + (end-start)*fraction
// }

// type Follow struct {
// 	Paths    []Path
// 	total    float64
// 	currPath int
// }

// func (f *Follow) Approach(v float64) gelly.Vector2 {
// 	return gelly.Vector2{}
// }

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

type Transform struct {
	Rotation float64
	Position Vector2
	Scale    Vector2
}

func (t Transform) Add(o Transform) Transform {
	t.Rotation += o.Rotation
	t.Position.X += o.Position.X
	t.Position.Y += o.Position.Y
	t.Scale.X += o.Scale.Y
	t.Scale.X += o.Scale.Y
	return t
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

var (
	ColorWhite = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	ColorBlack = color.RGBA{R: 0, G: 0, B: 0, A: 255}

	ColorLightGrey = color.RGBA{R: 160, G: 160, B: 160, A: 255}
	ColorGrey      = color.RGBA{R: 127, G: 127, B: 127, A: 255}
	ColorDarkGrey  = color.RGBA{R: 60, G: 60, B: 60, A: 255}

	ColorRed   = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	ColorGreen = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	ColorBlue  = color.RGBA{R: 0, G: 0, B: 255, A: 255}

	ColorYellow  = color.RGBA{R: 255, G: 255, B: 0, A: 255}
	ColorMagenta = color.RGBA{R: 255, G: 0, B: 255, A: 255}
	ColorCyan    = color.RGBA{R: 0, G: 255, B: 255, A: 255}

	BackgroundColor = color.RGBA{G: 100, B: 120, A: 255}
)
