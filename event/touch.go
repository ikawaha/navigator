package event

type Touch int

const (
	HoldFinger Touch = iota
	ReleaseFinger
	MoveFinger
)

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
