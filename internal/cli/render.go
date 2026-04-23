package cli

import (
	"fmt"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/presentation"
	"strings"
)

const clearSequence = ansiCursorHome + ansiClearScreen

const (
	encounterAxisIndent = 6
	encounterCellWidth  = 5
	encounterSeparator  = "~~~~"
	encounterColumnSpan = encounterCellWidth + len(encounterSeparator)
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

func renderFishDeckSelectionScreen(title string, presets []fishprofiles.FishDeckPreset, message string) string {
	sections := []string{
		renderHeader(title),
		renderFishDeckSelectionSection(presets),
	}
	if message != "" {
		sections = append(sections, accent("Aviso")+"\n  "+message)
	}

	return clearSequence + strings.Join(sections, "\n\n") + "\n\n"
}

func renderPlayerDeckSelectionScreen(title string, presets []playerprofiles.DeckPreset, message string) string {
	sections := []string{
		renderHeader(title),
		renderPlayerDeckSelectionSection(presets),
	}
	if message != "" {
		sections = append(sections, accent("Aviso")+"\n  "+message)
	}

	return clearSequence + strings.Join(sections, "\n\n") + "\n\n"
}

func renderFishDeckConfirmationScreen(title string, preset fishprofiles.FishDeckPreset, message string) string {
	sections := []string{
		renderHeader(title),
		renderFishDeckConfirmationSection(preset),
	}
	if message != "" {
		sections = append(sections, accent("Aviso")+"\n  "+message)
	}

	return clearSequence + strings.Join(sections, "\n\n") + "\n\n"
}

func renderPlayerDeckConfirmationScreen(title string, preset playerprofiles.DeckPreset, message string) string {
	sections := []string{
		renderHeader(title),
		renderPlayerDeckConfirmationSection(preset),
	}
	if message != "" {
		sections = append(sections, accent("Aviso")+"\n  "+message)
	}

	return clearSequence + strings.Join(sections, "\n\n") + "\n\n"
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
		renderEncounterAxisDivider(status.MaxDistance),
	}

	for depthLevel := status.SurfaceDepth; depthLevel <= status.MaxDepth; depthLevel++ {
		lines = append(lines, renderEncounterRow(status, depthLevel))
	}

	lines = append(lines,
		renderEncounterEscapeRow(status),
		renderEncounterShoreLabels(status.MaxDistance),
		fmt.Sprintf("  Distancia actual: %d | Captura <= %d | Escape > %d", status.FishDistance, status.CaptureDistance, status.MaxDistance),
		fmt.Sprintf("  Profundidad actual: %d | Superficie <= %d | Escape > %d", status.FishDepth, status.SurfaceDepth, status.MaxDepth),
		fmt.Sprintf("  Baraja agotada: captura con distancia <= %d y profundidad <= %d", status.ExhaustionCaptureDistance, status.SurfaceDepth+1),
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
	hint := ""
	if option.CardHint != "" {
		hint = " " + dim("{"+option.CardHint+"}")
	}
	if option.Available {
		return fmt.Sprintf("%d) %s %s%s", option.Index, moveLabel, dim(fmt.Sprintf("[%d/%d]", option.RemainingUses, option.MaxUses)), hint)
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

func renderFishDeckSelectionSection(presets []fishprofiles.FishDeckPreset) string {
	lines := []string{
		accent("Preset del pez"),
		"  Elige una baraja del pez para el duelo.",
	}

	for index, preset := range presets {
		lines = append(lines,
			fmt.Sprintf("  %d) %s", index+1, preset.Name),
			fmt.Sprintf("     %s", preset.Description),
		)
	}

	lines = append(lines, "  Escribe el numero del preset para seleccionarlo.")

	return strings.Join(lines, "\n")
}

func renderPlayerDeckSelectionSection(presets []playerprofiles.DeckPreset) string {
	lines := []string{
		accent("Preset del jugador"),
		"  Elige una baraja del jugador para probar el duelo.",
	}

	for index, preset := range presets {
		lines = append(lines,
			fmt.Sprintf("  %d) %s", index+1, preset.Name),
			fmt.Sprintf("     %s", preset.Description),
		)
	}

	lines = append(lines, "  Escribe el numero del preset para seleccionarlo.")

	return strings.Join(lines, "\n")
}

func renderFishDeckConfirmationSection(preset fishprofiles.FishDeckPreset) string {
	lines := []string{
		accent("Confirmar preset"),
		fmt.Sprintf("  Nombre      : %s", preset.Name),
		fmt.Sprintf("  Arquetipo   : %s", preset.ArchetypeID),
		fmt.Sprintf("  Descripcion : %s", preset.Description),
		fmt.Sprintf("  Cartas      : %d", len(preset.FishCards)),
		fmt.Sprintf("  Orden       : %s", renderFishDeckOrder(preset)),
		fmt.Sprintf("  Reciclado   : retira %d cartas por ciclo", preset.CardsToRemove),
	}
	for _, detail := range preset.Details {
		lines = append(lines, "  - "+detail)
	}

	return strings.Join(lines, "\n")
}

func renderPlayerDeckConfirmationSection(preset playerprofiles.DeckPreset) string {
	lines := []string{
		accent("Confirmar preset"),
		fmt.Sprintf("  Nombre      : %s", preset.Name),
		fmt.Sprintf("  Descripcion : %s", preset.Description),
		fmt.Sprintf("  Barajas     : %d colores", len(preset.Config.InitialDecks)),
		fmt.Sprintf("  Recupera    : %d ronda(s)", preset.Config.RecoveryDelayRounds),
	}
	for _, detail := range preset.Details {
		lines = append(lines, "  - "+detail)
	}

	return strings.Join(lines, "\n")
}

func renderFishDeckOrder(preset fishprofiles.FishDeckPreset) string {
	if preset.Shuffle {
		return "barajada"
	}

	return "orden fijo"
}

func renderEncounterAxis(maxDistance int) string {
	labels := make([]string, 0, maxDistance+2)
	for distanceLevel := 0; distanceLevel <= maxDistance; distanceLevel++ {
		labels = append(labels, centerEncounterText(fmt.Sprintf("%d", distanceLevel), encounterColumnSpan))
	}
	labels = append(labels, centerEncounterText("ESC", encounterCellWidth))

	return strings.Repeat(" ", encounterAxisIndent) + strings.Join(labels, "")
}

func renderEncounterAxisDivider(maxDistance int) string {
	lineWidth := encounterRenderableWidth(maxDistance)
	return strings.Repeat(" ", encounterAxisIndent) + strings.Repeat("-", lineWidth)
}

func renderEncounterRow(status presentation.StatusView, depthLevel int) string {
	rowLabel := renderEncounterRowPrefix(fmt.Sprintf("%d", depthLevel))
	cells := make([]string, 0, status.MaxDistance+2)
	for distanceLevel := 0; distanceLevel <= status.MaxDistance; distanceLevel++ {
		cells = append(cells, renderEncounterCell(status, depthLevel, distanceLevel))
	}
	cells = append(cells, renderEncounterEscapeColumnCell(status, depthLevel))

	return rowLabel + strings.Join(cells, dim(encounterSeparator))
}

func renderEncounterEscapeRow(status presentation.StatusView) string {
	cells := make([]string, 0, status.MaxDistance+2)
	for distanceLevel := 0; distanceLevel <= status.MaxDistance; distanceLevel++ {
		if status.FishDepth > status.MaxDepth && status.FishDistance == distanceLevel {
			cells = append(cells, accent(renderEncounterToken("[F!]")))
			continue
		}

		cells = append(cells, accent(renderEncounterToken("[ESC]")))
	}

	if status.FishDepth > status.MaxDepth && status.FishDistance > status.MaxDistance {
		cells = append(cells, accent(renderEncounterToken("[F!]")))
	} else {
		cells = append(cells, accent(renderEncounterToken("[ESC]")))
	}

	return renderEncounterRowPrefix("ESC") + strings.Join(cells, dim(encounterSeparator))
}

func renderEncounterCell(status presentation.StatusView, depthLevel, distanceLevel int) string {
	if depthLevel == status.SurfaceDepth && distanceLevel == 0 {
		if status.FishDepth <= status.SurfaceDepth && status.FishDistance <= 0 {
			return accent(renderEncounterToken("J/F"))
		}

		return accent(renderEncounterToken("J"))
	}

	if status.FishDepth == depthLevel && status.FishDistance == distanceLevel {
		return accent(renderEncounterToken("[F]"))
	}

	return dim(renderEncounterToken("[ ]"))
}

func renderEncounterEscapeColumnCell(status presentation.StatusView, depthLevel int) string {
	if status.FishDistance > status.MaxDistance && status.FishDepth == depthLevel {
		return accent(renderEncounterToken("[F!]"))
	}

	return accent(renderEncounterToken("[ESC]"))
}

func renderEncounterRowPrefix(label string) string {
	return fmt.Sprintf("%3s | ", label)
}

func renderEncounterShoreLabels(maxDistance int) string {
	const (
		leftLabel  = "Orilla"
		rightLabel = "Mar abierto"
	)

	padding := encounterRenderableWidth(maxDistance) - len(leftLabel) - len(rightLabel)
	if padding < 1 {
		padding = 1
	}

	return "  " + leftLabel + strings.Repeat(" ", padding) + rightLabel
}

func encounterRenderableWidth(maxDistance int) int {
	return (maxDistance+1)*encounterColumnSpan + encounterCellWidth
}

func centerEncounterText(label string, width int) string {
	if len(label) >= width {
		return label
	}

	leftPadding := (width - len(label)) / 2
	rightPadding := width - len(label) - leftPadding

	return strings.Repeat(" ", leftPadding) + label + strings.Repeat(" ", rightPadding)
}

func renderEncounterToken(label string) string {
	return fmt.Sprintf("%-*s", encounterCellWidth, label)
}
