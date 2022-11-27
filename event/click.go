package event

type Click int

const (
	SingleClick Click = iota
	HoldClick
	ReleaseClick
)

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
