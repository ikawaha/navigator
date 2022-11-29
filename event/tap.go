package event

// Tap represents the tap event.
type Tap int

const (
	// SingleTap is the single tap event.
	SingleTap Tap = iota
	// DoubleTap is the double tap event.
	DoubleTap
	// LongTap is the long tap event.
	LongTap
)

// String implements the Stringer interface.
func (t Tap) String() string {
	switch t {
	case SingleTap:
		return "tap"
	case DoubleTap:
		return "double tap"
	case LongTap:
		return "long tap"
	}
	return "perform tap"
}
