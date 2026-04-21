package domain

type Move int

const (
	Blue Move = iota
	Red
	Yellow
)

func (m Move) String() string {
	switch m {
	case Blue:
		return "blue"
	case Red:
		return "red"
	case Yellow:
		return "yellow"
	default:
		return "unknown"
	}
}

type RoundOutcome int

const (
	Draw RoundOutcome = iota
	PlayerWin
	FishWin
)

func (o RoundOutcome) String() string {
	switch o {
	case Draw:
		return "draw"
	case PlayerWin:
		return "player win"
	case FishWin:
		return "fish win"
	default:
		return "unknown"
	}
}
