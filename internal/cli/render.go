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
		"  Orilla  " + renderTrack(status.Distance, status.EscapeDistance) + "  Mar abierto",
		fmt.Sprintf("  Distancia actual: %d | Captura <= %d | Escape > %d | Baraja <= %d", status.Distance, status.CaptureDistance, status.EscapeDistance, status.ExhaustionCaptureDistance),
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
		"  Resultado : " + outcomeColor(view.Outcome),
		fmt.Sprintf("  Distancia : %d", view.Status.Distance),
	}, "\n")
}

func renderOptionsSection(options []presentation.MoveOption) string {
	parts := make([]string, 0, len(options))
	for _, option := range options {
		parts = append(parts, fmt.Sprintf("%d) %s", option.Index, colorizeMove(option.Move, option.Label)))
	}

	return strings.Join([]string{
		accent("Acciones"),
		"  " + strings.Join(parts, "   "),
		"  Escribe 1, 2 o 3 para actuar.",
	}, "\n")
}

func renderGameOverSection(summary presentation.SummaryView) string {
	return strings.Join([]string{
		accent("Desenlace"),
		"  Resultado : " + outcomeColor(summary.Outcome),
		"  Motivo    : " + summary.EndReason,
		fmt.Sprintf("  Distancia : %d", summary.Distance),
		fmt.Sprintf("  Rondas    : %d | Jugador %d | Pez %d | Empates %d", summary.TotalRounds, summary.PlayerWins, summary.FishWins, summary.Draws),
	}, "\n")
}

func renderTrack(distance, escapeDistance int) string {
	segments := []string{renderPlayerMarker(distance)}
	for position := 1; position <= escapeDistance; position++ {
		segments = append(segments, renderTrackMarker(position, distance))
	}
	segments = append(segments, renderEscapeMarker())
	if distance > escapeDistance {
		segments = append(segments, accent("F!"))
	}

	return strings.Join(segments, dim("~~~~"))
}

func renderPlayerMarker(distance int) string {
	if distance <= 0 {
		return accent("[J/F]")
	}
	return accent("[J]")
}

func renderTrackMarker(position, distance int) string {
	if position == distance {
		return accent("[F]")
	}
	return dim(fmt.Sprintf("[%d]", position))
}

func renderEscapeMarker() string {
	return accent("[ESC]")
}
