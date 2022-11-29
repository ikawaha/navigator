package event

// Button represents the mouse button event.
type Button int

const (
	// LeftButton is the event when the left button is clicked.
	LeftButton Button = iota
	// MiddleButton is the event when the middle button is clicked.
	MiddleButton
	// RightButton is the event when the right button iss clicked.
	RightButton
)

// String implements the Stringer interface.
func (b Button) String() string {
	switch b {
	case LeftButton:
		return "left mouse button"
	case MiddleButton:
		return "middle mouse button"
	case RightButton:
		return "right mouse button"
	}
	return "unknown"
}
