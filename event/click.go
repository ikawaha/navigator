package event

// Click represents the mouse click event.
type Click int

const (
	// SingleClick is the single mouse click event.
	SingleClick Click = iota
	// HoldClick is the event when a click is held.
	HoldClick
	// ReleaseClick is the event when a click is released.
	ReleaseClick
)

// String implements the Stringer interface.
func (c Click) String() string {
	switch c {
	case SingleClick:
		return "single click"
	case HoldClick:
		return "hold"
	case ReleaseClick:
		return "release"
	}
	return "unknown"
}
