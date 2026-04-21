package cli

import (
	"fmt"
	"pesca/internal/presentation"
	"strings"
)

const clearSequence = ansiCursorHome + ansiClearScreen

func renderPromptScreen(title string, status presentation.StatusView, options []presentation.MoveOption, lastRound *presentation.RoundView, message string) string {
	var sections []string
	sections = append(sections, renderHeader(title))
	sections = append(sections, renderTrackSection(status))
	sections = append(sections, renderStatsSection(status))
	if lastRound != nil {
		sections = append(sections, renderLastRoundSection(*lastRound))
	}
	sections = append(sections, renderOptionsSection(options))
	if message != "" {
		sections = append(sections, accent("Aviso")+"\n  "+message)
	}

	return clearSequence + strings.Join(sections, "\n\n") + "\n\n"
}

func renderGameOverScreen(title string, summary presentation.SummaryView, lastRound *presentation.RoundView) string {
	var sections []string
	sections = append(sections, renderHeader(title))
	if lastRound != nil {
		sections = append(sections, renderLastRoundSection(*lastRound))
	}
	sections = append(sections, renderGameOverSection(summary))

	return clearSequence + strings.Join(sections, "\n\n") + "\n"
}

func renderHeader(title string) string {
	if title == "" {
		title = "Pesca: duelo contra el pez"
	}

	return strings.Join([]string{
		accent("============================================================"),
		accent(strings.ToUpper(title)),
		dim("Tensa el sedal y arrastra al pez hacia la orilla."),
		accent("============================================================"),
	}, "\n")
}

func renderTrackSection(status presentation.StatusView) string {
	return strings.Join([]string{
		accent("Sedal"),
		"  Orilla  " + renderTrack(status.FishDistance, status.EscapeDistance) + "  Mar abierto",
		fmt.Sprintf("  Distancia actual: %d | Captura <= %d | Escape > %d | Baraja <= %d", status.FishDistance, status.CaptureDistance, status.EscapeDistance, status.ExhaustionCaptureDistance),
	}, "\n")
}

func renderStatsSection(status presentation.StatusView) string {
	return strings.Join([]string{
		accent("Estado del encuentro"),
		fmt.Sprintf("  Ronda %d | Mazo %d | Descarte %d | Rebarajados %d", status.RoundNumber, status.ActiveCards, status.DiscardCards, status.RecycleCount),
		fmt.Sprintf("  Tension del duelo | Jugador %d | Pez %d | Empates %d", status.PlayerWins, status.FishWins, status.Draws),
	}, "\n")
}

func renderLastRoundSection(view presentation.RoundView) string {
	return strings.Join([]string{
		accent("Ultimo lance"),
		"  Tu accion : " + colorizeMove(view.PlayerMove, view.PlayerLabel),
		"  Pez       : " + colorizeMove(view.FishMove, view.FishLabel),
		"  Resultado : " + colorizeRoundOutcome(view.Outcome, view.OutcomeLabel),
		fmt.Sprintf("  Distancia : %d", view.Status.FishDistance),
	}, "\n")
}

func renderOptionsSection(options []presentation.MoveOption) string {
	parts := make([]string, 0, len(options))
	for _, option := range options {
		parts = append(parts, renderMoveOption(option))
	}

	return strings.Join([]string{
		accent("Acciones"),
		"  " + strings.Join(parts, "   "),
		"  Escribe 1, 2 o 3 para actuar.",
	}, "\n")
}

func renderMoveOption(option presentation.MoveOption) string {
	moveLabel := colorizeMove(option.Move, option.Label)
	if option.Available {
		return fmt.Sprintf("%d) %s %s", option.Index, moveLabel, dim(fmt.Sprintf("[%d/%d]", option.RemainingUses, option.MaxUses)))
	}

	if option.RestoresOnRound > 0 {
		return fmt.Sprintf("%d) %s %s", option.Index, moveLabel, dim(fmt.Sprintf("[recarga R%d]", option.RestoresOnRound)))
	}

	return fmt.Sprintf("%d) %s %s", option.Index, moveLabel, dim("[sin usos]"))
}

func renderGameOverSection(summary presentation.SummaryView) string {
	return strings.Join([]string{
		accent("Desenlace"),
		"  Resultado : " + colorizeEncounterStatus(summary.EncounterStatus, summary.OutcomeLabel),
		"  Motivo    : " + summary.EndReasonLabel,
		fmt.Sprintf("  Distancia : %d", summary.FishDistance),
		fmt.Sprintf("  Rondas    : %d | Jugador %d | Pez %d | Empates %d", summary.TotalRounds, summary.PlayerWins, summary.FishWins, summary.Draws),
	}, "\n")
}

func renderTrack(fishDistance, escapeDistance int) string {
	segments := []string{renderPlayerMarker(fishDistance)}
	for trackPosition := 1; trackPosition <= escapeDistance; trackPosition++ {
		segments = append(segments, renderTrackMarker(trackPosition, fishDistance))
	}
	segments = append(segments, renderEscapeMarker())
	if fishDistance > escapeDistance {
		segments = append(segments, accent("F!"))
	}

	return strings.Join(segments, dim("~~~~"))
}

func renderPlayerMarker(fishDistance int) string {
	if fishDistance <= 0 {
		return accent("[J/F]")
	}
	return accent("[J]")
}

func renderTrackMarker(trackPosition, fishDistance int) string {
	if trackPosition == fishDistance {
		return accent("[F]")
	}
	return dim(fmt.Sprintf("[%d]", trackPosition))
}

func renderEscapeMarker() string {
	return accent("[ESC]")
}
