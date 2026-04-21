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
		return "empate"
	case PlayerWin:
		return "gana el jugador"
	case FishWin:
		return "gana el pez"
	default:
		return "desconocido"
	}
}
