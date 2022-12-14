package session

import "math"

// Speed represents the vector/scalar speed.
type Speed interface {
	vector() (x, y int)
	scalar() uint
}

// VectorSpeed represents the vector speed.
type VectorSpeed struct {
	X int
	Y int
}

func (s VectorSpeed) vector() (x, y int) {
	return s.X, s.Y
}

func (s VectorSpeed) scalar() uint {
	return uint(math.Hypot(float64(s.X), float64(s.Y)))
}

// ScalarSpeed represents the scalar speed.
type ScalarSpeed uint

func (s ScalarSpeed) vector() (x, y int) {
	scalar := int(float64(s) / math.Sqrt2)
	return scalar, scalar
}

func (s ScalarSpeed) scalar() uint {
	return uint(s)
}
