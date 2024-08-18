package gelly

import "math"

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

type Path struct {
	A Vector2
	B Vector2
	V float64
	// total float64
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

// func (p Vector2) Minus(other Vector2) Vector2 {
// 	other.X -= p.X
// 	other.Y -= p.Y
// 	return other
// }

// func (p *Path) Approach(v float64) Vector2 {
// 	v *= v

// 	p.total += v
// 	progress := p.total / p.V
// 	dir := p.A.Minus(p.B)

// 	dir.X *= progress
// 	dir.Y *= progress

// 	return dir
// }
