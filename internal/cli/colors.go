package cli

import (
	"fmt"
	"pesca/internal/domain"
)

const (
	ansiReset       = "\033[0m"
	ansiBold        = "\033[1m"
	ansiDim         = "\033[2m"
	ansiBlue        = "\033[34m"
	ansiRed         = "\033[31m"
	ansiGreen       = "\033[32m"
	ansiYellow      = "\033[33m"
	ansiCyan        = "\033[36m"
	ansiClearScreen = "\033[2J"
	ansiCursorHome  = "\033[H"
)

func colorizeMove(move domain.Move, text string) string {
	return fmt.Sprintf("%s%s%s%s", ansiBold, ansiCodeForMove(move), text, ansiReset)
}

func accent(text string) string {
	return fmt.Sprintf("%s%s%s%s", ansiBold, ansiCyan, text, ansiReset)
}

func dim(text string) string {
	return fmt.Sprintf("%s%s%s", ansiDim, text, ansiReset)
}

func outcomeColor(text string) string {
	switch text {
	case "gana el jugador", "pez capturado":
		return fmt.Sprintf("%s%s%s%s", ansiBold, ansiGreen, text, ansiReset)
	case "gana el pez", "pez escapado":
		return fmt.Sprintf("%s%s%s%s", ansiBold, ansiRed, text, ansiReset)
	default:
		return fmt.Sprintf("%s%s%s%s", ansiBold, ansiYellow, text, ansiReset)
	}
}

func ansiCodeForMove(move domain.Move) string {
	switch move {
	case domain.Blue:
		return ansiBlue
	case domain.Red:
		return ansiRed
	case domain.Yellow:
		return ansiYellow
	default:
		return ""
	}
}
