package event

type Button int

const (
	LeftButton Button = iota
	MiddleButton
	RightButton
)

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
