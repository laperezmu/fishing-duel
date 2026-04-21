package cli

import (
	"fmt"
	"pesca/internal/domain"
	"pesca/internal/encounter"
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

func colorizeRoundOutcome(outcome domain.RoundOutcome, outcomeLabel string) string {
	switch outcome {
	case domain.PlayerWin:
		return fmt.Sprintf("%s%s%s%s", ansiBold, ansiGreen, outcomeLabel, ansiReset)
	case domain.FishWin:
		return fmt.Sprintf("%s%s%s%s", ansiBold, ansiRed, outcomeLabel, ansiReset)
	default:
		return fmt.Sprintf("%s%s%s%s", ansiBold, ansiYellow, outcomeLabel, ansiReset)
	}
}

func colorizeEncounterStatus(status encounter.Status, outcomeLabel string) string {
	switch status {
	case encounter.StatusCaptured:
		return fmt.Sprintf("%s%s%s%s", ansiBold, ansiGreen, outcomeLabel, ansiReset)
	case encounter.StatusEscaped:
		return fmt.Sprintf("%s%s%s%s", ansiBold, ansiRed, outcomeLabel, ansiReset)
	default:
		return fmt.Sprintf("%s%s%s%s", ansiBold, ansiYellow, outcomeLabel, ansiReset)
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
