package gelly

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

type Material struct {
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
