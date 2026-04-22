package cli

import (
	"fmt"
	"pesca/internal/presentation"
	"strings"
)

const clearSequence = ansiCursorHome + ansiClearScreen

const (
	encounterCellWidth = 5
	encounterStride    = encounterCellWidth + 4
)

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
	lines := []string{
		accent("Sedal"),
		renderEncounterAxis(status.MaxDistance),
	}

	for depthLevel := status.SurfaceDepth; depthLevel <= status.MaxDepth; depthLevel++ {
		lines = append(lines, renderEncounterRow(status, depthLevel))
	}

	lines = append(lines,
		renderEncounterEscapeRow(status),
		"  Orilla                                      Mar abierto",
		fmt.Sprintf("  Distancia actual: %d | Captura <= %d | Escape > %d | Baraja <= %d", status.FishDistance, status.CaptureDistance, status.MaxDistance, status.ExhaustionCaptureDistance),
		fmt.Sprintf("  Profundidad actual: %d | Superficie <= %d | Escape > %d", status.FishDepth, status.SurfaceDepth, status.MaxDepth),
	)

	return strings.Join(lines, "\n")
}

func renderStatsSection(status presentation.StatusView) string {
	return strings.Join([]string{
		accent("Estado del encuentro"),
		fmt.Sprintf("  Ronda %d | Mazo %d | Descarte %d | Rebarajados %d", status.RoundNumber, status.ActiveCards, status.DiscardCards, status.RecycleCount),
		fmt.Sprintf("  Tension del duelo | Jugador %d | Pez %d | Empates %d", status.PlayerWins, status.FishWins, status.Draws),
	}, "\n")
}

func renderLastRoundSection(view presentation.RoundView) string {
	lines := []string{
		accent("Ultimo lance"),
		"  Tu accion : " + colorizeMove(view.PlayerMove, view.PlayerLabel),
		"  Pez       : " + colorizeMove(view.FishMove, view.FishLabel),
		"  Resultado : " + colorizeRoundOutcome(view.Outcome, view.OutcomeLabel),
		fmt.Sprintf("  Distancia : %d", view.Status.FishDistance),
		fmt.Sprintf("  Profundidad : %d", view.Status.FishDepth),
	}
	if view.EventLabel != "" {
		lines = append(lines, "  Evento    : "+view.EventLabel)
	}

	return strings.Join(lines, "\n")
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
		fmt.Sprintf("  Profundidad : %d", summary.FishDepth),
		fmt.Sprintf("  Rondas    : %d | Jugador %d | Pez %d | Empates %d", summary.TotalRounds, summary.PlayerWins, summary.FishWins, summary.Draws),
	}, "\n")
}

func renderEncounterAxis(maxDistance int) string {
	labels := make([]string, 0, maxDistance+2)
	for distanceLevel := 0; distanceLevel <= maxDistance; distanceLevel++ {
		labels = append(labels, fmt.Sprintf("%-*d", encounterStride, distanceLevel))
	}
	labels = append(labels, "ESC")

	return strings.Repeat(" ", len(renderEncounterRowPrefix("ESC"))) + strings.Join(labels, "")
}

func renderEncounterRow(status presentation.StatusView, depthLevel int) string {
	rowLabel := renderEncounterRowPrefix(fmt.Sprintf("%d", depthLevel))
	cells := make([]string, 0, status.MaxDistance+2)
	for distanceLevel := 0; distanceLevel <= status.MaxDistance; distanceLevel++ {
		cells = append(cells, renderEncounterCell(status, depthLevel, distanceLevel))
	}
	cells = append(cells, renderEncounterEscapeColumnCell(status, depthLevel))

	return rowLabel + strings.Join(cells, dim("~~~~"))
}

func renderEncounterEscapeRow(status presentation.StatusView) string {
	cells := make([]string, 0, status.MaxDistance+2)
	for distanceLevel := 0; distanceLevel <= status.MaxDistance; distanceLevel++ {
		if status.FishDepth > status.MaxDepth && status.FishDistance == distanceLevel {
			cells = append(cells, accent(padEncounterCell("[F!]")))
			continue
		}

		cells = append(cells, accent(padEncounterCell("[ESC]")))
	}

	if status.FishDepth > status.MaxDepth && status.FishDistance > status.MaxDistance {
		cells = append(cells, accent(padEncounterCell("[F!]")))
	} else {
		cells = append(cells, accent(padEncounterCell("[ESC]")))
	}

	return renderEncounterRowPrefix("ESC") + strings.Join(cells, dim("~~~~"))
}

func renderEncounterCell(status presentation.StatusView, depthLevel, distanceLevel int) string {
	if depthLevel == status.SurfaceDepth && distanceLevel == 0 {
		if status.FishDepth <= status.SurfaceDepth && status.FishDistance <= 0 {
			return accent(padEncounterCell("J/F"))
		}

		return accent(padEncounterCell("J"))
	}

	if status.FishDepth == depthLevel && status.FishDistance == distanceLevel {
		return accent(padEncounterCell("[F]"))
	}

	return dim(padEncounterCell("[ ]"))
}

func renderEncounterEscapeColumnCell(status presentation.StatusView, depthLevel int) string {
	if status.FishDistance > status.MaxDistance && status.FishDepth == depthLevel {
		return accent(padEncounterCell("[F!]"))
	}

	return accent(padEncounterCell("[ESC]"))
}

func renderEncounterRowPrefix(label string) string {
	return fmt.Sprintf("%3s | ", label)
}

func padEncounterCell(content string) string {
	if len(content) >= encounterCellWidth {
		return content
	}

	return content + strings.Repeat(" ", encounterCellWidth-len(content))
}
