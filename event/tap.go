package event

type Tap int

const (
	SingleTap Tap = iota
	DoubleTap
	LongTap
)

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
