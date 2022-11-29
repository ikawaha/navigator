package event

// Touch represents the touch event.
type Touch int

const (
	// HoldFinger is the hold finger event.
	HoldFinger Touch = iota
	// ReleaseFinger is the release finger event.
	ReleaseFinger
	// MoveFinger is the move finger event.
	MoveFinger
)

// String implements the Stringer interface.
func (t Touch) String() string {
	switch t {
	case HoldFinger:
		return "hold finger down"
	case ReleaseFinger:
		return "release finger"
	case MoveFinger:
		return "move finger"
	}
	return "perform touch"
}
