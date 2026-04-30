package cli

import (
	"fmt"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/habitats"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/rodpresets"
	"pesca/internal/content/watercontexts"
	"pesca/internal/content/waterpools"
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
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

func renderPromptScreen(title string, status presentation.StatusView, options []presentation.MoveOption, opening *encounter.Opening, spawn *fishprofiles.Spawn, lastRound *presentation.RoundView, message string) string {
	var sections []string
	sections = append(sections, renderHeader(title))
	if opening != nil {
		sections = append(sections, renderEncounterOpeningSection(*opening))
	}
	if spawn != nil {
		sections = append(sections, renderFishSpawnSection(*spawn))
	}
	sections = append(sections, renderTrackSection(status))
	sections = append(sections, renderStatsSection(status))
	sections = append(sections, renderFishDiscardSection(status))
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

func renderRodSelectionScreen(title string, presets []rodpresets.Preset, message string) string {
	sections := []string{
		renderHeader(title),
		renderRodSelectionSection(presets),
	}
	if message != "" {
		sections = append(sections, accent("Aviso")+"\n  "+message)
	}

	return clearSequence + strings.Join(sections, "\n\n") + "\n\n"
}

func renderAttachmentSelectionScreen(title string, presets []attachmentpresets.Preset, message string) string {
	sections := []string{
		renderHeader(title),
		renderAttachmentSelectionSection(presets),
	}
	if message != "" {
		sections = append(sections, accent("Aviso")+"\n  "+message)
	}

	return clearSequence + strings.Join(sections, "\n\n") + "\n\n"
}

func renderWaterContextSelectionScreen(title string, presets []watercontexts.Preset, message string) string {
	sections := []string{
		renderHeader(title),
		renderWaterContextSelectionSection(presets),
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

func renderRodConfirmationScreen(title string, preset rodpresets.Preset, message string) string {
	sections := []string{
		renderHeader(title),
		renderRodConfirmationSection(preset),
	}
	if message != "" {
		sections = append(sections, accent("Aviso")+"\n  "+message)
	}

	return clearSequence + strings.Join(sections, "\n\n") + "\n\n"
}

func renderAttachmentConfirmationScreen(title string, preset attachmentpresets.Preset, preview loadout.State, message string) string {
	sections := []string{
		renderHeader(title),
		renderAttachmentConfirmationSection(preset, preview),
	}
	if message != "" {
		sections = append(sections, accent("Aviso")+"\n  "+message)
	}

	return clearSequence + strings.Join(sections, "\n\n") + "\n\n"
}

func renderWaterContextConfirmationScreen(title string, preset watercontexts.Preset, message string) string {
	sections := []string{
		renderHeader(title),
		renderWaterContextConfirmationSection(preset),
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

func renderEncounterOpeningSection(opening encounter.Opening) string {
	return strings.Join([]string{
		accent("Apertura del lance"),
		fmt.Sprintf("  Agua       : %s", opening.WaterContext.Name),
		fmt.Sprintf("  Lance      : %s", opening.CastResult.Band.Label()),
		fmt.Sprintf("  Inicio     : distancia %d | profundidad %d", opening.InitialDistance, opening.InitialDepth),
	}, "\n")
}

func renderFishSpawnSection(spawn fishprofiles.Spawn) string {
	lines := []string{
		accent("Pez en el agua"),
		fmt.Sprintf("  Perfil     : %s", spawn.Profile.Name),
		fmt.Sprintf("  Agua base  : %s", waterpools.Name(spawn.Context.WaterPoolTag)),
		fmt.Sprintf("  Ventana    : distancia %d | profundidad %d", spawn.Context.InitialDistance, spawn.Context.InitialDepth),
		fmt.Sprintf("  Candidatos : %d", spawn.CandidateCount),
	}
	if len(spawn.Context.HabitatTags) > 0 {
		lines = append(lines, "  Habitats   : "+strings.Join(habitats.Names(spawn.Context.HabitatTags), ", "))
	}

	return strings.Join(lines, "\n")
}

func renderFishDiscardSection(status presentation.StatusView) string {
	currentCycle := renderCurrentFishDiscardCycle(status.FishDiscard)
	recycleMode := renderFishRecycleMode(status.FishDiscard)
	previousCycles := renderPreviousFishDiscardCycles(status.FishDiscard)

	return strings.Join([]string{
		accent("Historial del pez"),
		"  Ciclo activo : " + currentCycle,
		fmt.Sprintf("  Reciclado   : %s | retira %d carta%s | %d ciclo%s cerrado%s", recycleMode, status.FishDiscard.CardsToRemove, pluralSuffix(status.FishDiscard.CardsToRemove), status.FishDiscard.RecycleCount, pluralSuffix(status.FishDiscard.RecycleCount), pluralSuffix(status.FishDiscard.RecycleCount)),
		"  Ciclos cerrados: " + previousCycles,
	}, "\n")
}

func renderCurrentFishDiscardCycle(view presentation.FishDiscardView) string {
	if view.CurrentCycleTotalCards == 0 {
		return fmt.Sprintf("C%d sin cartas usadas todavia", maxInt(view.CurrentCycleNumber, 1))
	}

	parts := make([]string, 0, len(view.CurrentCycleEntries)+1)
	for _, entry := range view.CurrentCycleEntries {
		parts = append(parts, entry.Label)
	}

	hiddenCards := view.CurrentCycleTotalCards - len(view.CurrentCycleEntries)
	if hiddenCards > 0 {
		parts = append(parts, fmt.Sprintf("%d oculta%s", hiddenCards, pluralSuffix(hiddenCards)))
	}

	if len(parts) == 0 {
		return fmt.Sprintf("C%d %d carta%s oculta%s", maxInt(view.CurrentCycleNumber, 1), hiddenCards, pluralSuffix(hiddenCards), pluralSuffix(hiddenCards))
	}

	return fmt.Sprintf("C%d %s", maxInt(view.CurrentCycleNumber, 1), strings.Join(parts, " | "))
}

func renderFishRecycleMode(view presentation.FishDiscardView) string {
	if view.ShufflesOnRecycle {
		return "rebaraja"
	}

	return "mantiene orden"
}

func renderPreviousFishDiscardCycles(view presentation.FishDiscardView) string {
	if len(view.PreviousCycles) == 0 {
		return "ninguno"
	}

	parts := make([]string, 0, len(view.PreviousCycles))
	for _, cycle := range view.PreviousCycles {
		parts = append(parts, renderFishDiscardCycleSummary(cycle))
	}

	return strings.Join(parts, " | ")
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

func renderFishDiscardCycleSummary(cycle presentation.FishDiscardCycleSummaryView) string {
	parts := []string{fmt.Sprintf("C%d %d usada%s", cycle.CycleNumber, cycle.TotalCards, pluralSuffix(cycle.TotalCards))}
	if cycle.HiddenCards > 0 {
		parts = append(parts, fmt.Sprintf("%d oculta%s", cycle.HiddenCards, pluralSuffix(cycle.HiddenCards)))
	}

	return strings.Join(parts, ", ")
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

func renderWaterContextSelectionSection(presets []watercontexts.Preset) string {
	lines := []string{
		accent("Situacion de agua"),
		"  Elige el tipo de agua que abrira el encuentro.",
	}

	for index, preset := range presets {
		lines = append(lines,
			fmt.Sprintf("  %d) %s", index+1, preset.Name),
			fmt.Sprintf("     %s", preset.Description),
		)
	}

	lines = append(lines, "  Escribe el numero de la situacion para seleccionarla.")

	return strings.Join(lines, "\n")
}

func renderRodSelectionSection(presets []rodpresets.Preset) string {
	lines := []string{
		accent("Preset de cana"),
		"  Elige la rod base que definira tus limites de apertura y track.",
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

func renderAttachmentSelectionSection(presets []attachmentpresets.Preset) string {
	lines := []string{
		accent("Preset de aditamentos"),
		"  Elige los aditamentos que completaran el loadout sobre tu rod.",
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

func renderWaterContextConfirmationSection(preset watercontexts.Preset) string {
	lines := []string{
		accent("Confirmar situacion de agua"),
		fmt.Sprintf("  Nombre      : %s", preset.Name),
		fmt.Sprintf("  Resumen     : %s", preset.Description),
		fmt.Sprintf("  Profundidad : %d", preset.InitialDepth),
	}
	for _, signal := range preset.Signals {
		lines = append(lines, "  - "+signal)
	}

	return strings.Join(lines, "\n")
}

func renderFishDeckConfirmationSection(preset fishprofiles.FishDeckPreset) string {
	lines := []string{
		accent("Confirmar preset"),
		fmt.Sprintf("  Nombre      : %s", preset.Name),
		fmt.Sprintf("  Arquetipo   : %s", fishprofiles.Name(preset.ArchetypeID)),
		fmt.Sprintf("  Resumen     : %s", preset.Description),
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
		fmt.Sprintf("  Resumen     : %s", preset.Description),
		fmt.Sprintf("  Barajas     : %d colores", len(preset.Config.InitialDecks)),
		fmt.Sprintf("  Recupera    : %d ronda(s)", preset.Config.RecoveryDelayRounds),
	}
	for _, detail := range preset.Details {
		lines = append(lines, "  - "+detail)
	}

	return strings.Join(lines, "\n")
}

func renderRodConfirmationSection(preset rodpresets.Preset) string {
	lines := []string{
		accent("Confirmar cana"),
		fmt.Sprintf("  Nombre          : %s", preset.Name),
		fmt.Sprintf("  Resumen         : %s", preset.Description),
		fmt.Sprintf("  Apertura        : dist %d | prof %d", preset.Config.OpeningMaxDistance, preset.Config.OpeningMaxDepth),
		fmt.Sprintf("  Track           : dist %d | prof %d", preset.Config.TrackMaxDistance, preset.Config.TrackMaxDepth),
	}
	for _, detail := range preset.Details {
		lines = append(lines, "  - "+detail)
	}

	return strings.Join(lines, "\n")
}

func renderAttachmentConfirmationSection(preset attachmentpresets.Preset, preview loadout.State) string {
	effectiveRod, _ := preview.EffectiveRod()
	lines := []string{
		accent("Confirmar aditamentos"),
		fmt.Sprintf("  Nombre          : %s", preset.Name),
		fmt.Sprintf("  Resumen         : %s", preset.Description),
		fmt.Sprintf("  Resultado final : apertura dist %d | prof %d", effectiveRod.OpeningMaxDistance, effectiveRod.OpeningMaxDepth),
		fmt.Sprintf("  Track final     : dist %d | prof %d", effectiveRod.TrackMaxDistance, effectiveRod.TrackMaxDepth),
	}
	if len(preset.Attachments) == 0 {
		lines = append(lines, "  - Sin aditamentos equipados.")
	}
	for _, detail := range preset.Details {
		lines = append(lines, "  - "+detail)
	}
	for _, attachment := range preset.Attachments {
		lines = append(lines, "  * "+attachment.Name+": "+attachment.Description)
	}

	return strings.Join(lines, "\n")
}

func renderFishDeckOrder(preset fishprofiles.FishDeckPreset) string {
	if preset.Shuffle {
		return "barajada"
	}

	return "orden fijo"
}

func renderCastScreen(title string, context encounter.WaterContext, position int, totalSlots int, slotWidth int, message string) string {
	sections := []string{
		renderHeader(title),
		renderCastSection(context, position, totalSlots, slotWidth),
	}
	if message != "" {
		sections = append(sections, accent("Aviso")+"\n  "+message)
	}

	return clearSequence + strings.Join(sections, "\n\n") + "\n\n"
}

func renderCastSection(context encounter.WaterContext, position int, totalSlots int, slotWidth int) string {
	bands := encounter.OrderedCastBands()
	bandParts := make([]string, 0, len(bands))
	for _, band := range bands {
		bandParts = append(bandParts, band.Label())
	}

	lines := []string{
		accent("Lectura del agua"),
		fmt.Sprintf("  Agua       : %s", context.Name),
		fmt.Sprintf("  Resumen    : %s", context.Description),
	}
	for _, signal := range context.VisibleSignals {
		lines = append(lines, "  - "+signal)
	}
	lines = append(lines,
		"",
		accent("Cast"),
		"  Pulsa Enter para detener la barra en la franja deseada.",
		"  Bandas     : "+strings.Join(bandParts, " | "),
		"  Barra      : "+renderCastBar(position, totalSlots, slotWidth),
	)

	return strings.Join(lines, "\n")
}

func renderCastBar(position int, totalSlots int, slotWidth int) string {
	var builder strings.Builder
	_, _ = builder.WriteString("[")
	for slot := 0; slot < totalSlots; slot++ {
		if slot > 0 && slot%slotWidth == 0 {
			_, _ = builder.WriteString("|")
		}
		if slot <= position {
			_, _ = builder.WriteString("=")
			continue
		}
		_, _ = builder.WriteString(".")
	}
	_, _ = builder.WriteString("]")

	return builder.String()
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

func pluralSuffix(count int) string {
	if count == 1 {
		return ""
	}

	return "s"
}

func maxInt(value, minimum int) int {
	if value < minimum {
		return minimum
	}

	return value
}
