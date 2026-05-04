package presentation

import (
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/content/anglerprofiles"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/habitats"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/rodpresets"
	"pesca/internal/content/waterpools"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/run"
	"strings"
)

type Catalog struct {
	Title            string
	PlayerMoveLabels map[domain.Move]string
	FishMoveLabels   map[domain.Move]string
	RoundOutcomes    map[domain.RoundOutcome]string
	EncounterEvents  map[encounter.EventKind]string
	EventOutcomes    map[bool]string
	EncounterResults map[encounter.Status]string
	EndReasons       map[encounter.EndReason]string
}

func DefaultCatalog() Catalog {
	return Catalog{
		Title: "Pesca: duelo contra el pez",
		PlayerMoveLabels: map[domain.Move]string{
			domain.Blue:   "Tirar",
			domain.Red:    "Recoger",
			domain.Yellow: "Soltar",
		},
		FishMoveLabels: map[domain.Move]string{
			domain.Blue:   "Embestir",
			domain.Red:    "Aferrarse",
			domain.Yellow: "Zafarse",
		},
		RoundOutcomes: map[domain.RoundOutcome]string{
			domain.Draw:      "empate",
			domain.PlayerWin: "gana el jugador",
			domain.FishWin:   "gana el pez",
		},
		EncounterEvents: map[encounter.EventKind]string{
			encounter.EventKindSplash: "chapotea en la superficie",
		},
		EventOutcomes: map[bool]string{
			false: "permanece sujeto",
			true:  "se suelta del anzuelo",
		},
		EncounterResults: map[encounter.Status]string{
			encounter.StatusCaptured: "pez capturado",
			encounter.StatusEscaped:  "pez escapado",
			encounter.StatusOngoing:  "encuentro en curso",
		},
		EndReasons: map[encounter.EndReason]string{
			encounter.EndReasonTrackCapture: "captura por acercarlo a la orilla y subirlo a la superficie",
			encounter.EndReasonTrackEscape:  "escape por superar la distancia maxima alcanzable",
			encounter.EndReasonDepthEscape:  "escape por bajar mas alla de la profundidad alcanzable",
			encounter.EndReasonSplashEscape: "escape por chapoteo en superficie",
			encounter.EndReasonDeckCapture:  "captura al agotar la baraja con distancia 2 o menor y profundidad 1 o menor",
			encounter.EndReasonDeckEscape:   "escape al agotar la baraja sin cumplir la distancia o profundidad de cierre",
			encounter.EndReasonNone:         "sin resolver",
		},
	}
}

type Presenter struct {
	catalog Catalog
}

func NewPresenter(catalog Catalog) Presenter {
	return Presenter{catalog: catalog}
}

func (p Presenter) Intro() IntroView {
	return IntroView{
		Title:   p.catalog.Title,
		Options: p.moveOptions(),
	}
}

func (p Presenter) Status(snapshot match.StatusSnapshot) StatusView {
	return StatusView{
		RoundNumber:               snapshot.RoundNumber,
		FishDistance:              snapshot.Track.Distance,
		FishDepth:                 snapshot.Track.Depth,
		SurfaceDepth:              snapshot.Track.SurfaceDepth,
		MaxDistance:               snapshot.Track.MaxDistance,
		MaxDepth:                  snapshot.Track.MaxDepth,
		CaptureDistance:           snapshot.Track.CaptureDistance,
		ExhaustionCaptureDistance: snapshot.Track.ExhaustionCaptureDistance,
		ActiveCards:               snapshot.FishDiscard.ActiveCards,
		DiscardCards:              snapshot.FishDiscard.DiscardCards,
		RecycleCount:              snapshot.FishDiscard.RecycleCount,
		PlayerWins:                snapshot.Stats.PlayerWins,
		FishWins:                  snapshot.Stats.FishWins,
		Draws:                     snapshot.Stats.Draws,
		FishDiscard:               p.fishDiscardView(snapshot.FishDiscard),
		MoveOptions:               p.moveOptionsForSnapshot(snapshot.Player),
	}
}

func (p Presenter) Round(snapshot match.RoundSnapshot) RoundView {
	return RoundView{
		Status:       p.Status(snapshot.Status),
		PlayerMove:   snapshot.PlayerMove,
		FishMove:     snapshot.FishMove,
		PlayerLabel:  p.playerMoveLabel(snapshot.PlayerMove),
		FishLabel:    p.fishMoveLabel(snapshot.FishMove),
		Outcome:      snapshot.Outcome,
		OutcomeLabel: p.roundOutcomeLabel(snapshot.Outcome),
		EventLabel:   p.eventLabel(snapshot.Encounter.LastEvent),
	}
}

func (p Presenter) Summary(snapshot match.SummarySnapshot) SummaryView {
	return SummaryView{
		TotalRounds:     snapshot.TotalRounds,
		FishDistance:    snapshot.Encounter.Distance,
		FishDepth:       snapshot.Encounter.Depth,
		EncounterStatus: snapshot.Encounter.Status,
		OutcomeLabel:    p.encounterOutcomeLabel(snapshot.Encounter.Status),
		EndReasonLabel:  p.endReasonLabel(snapshot.Encounter.EndReason),
		PlayerWins:      snapshot.Stats.PlayerWins,
		FishWins:        snapshot.Stats.FishWins,
		Draws:           snapshot.Stats.Draws,
	}
}

func (p Presenter) Opening(opening encounter.Opening) OpeningView {
	return OpeningView{
		WaterLabel:      opening.WaterContext.Name,
		CastLabel:       opening.CastResult.Band.Label(),
		InitialDistance: opening.InitialDistance,
		InitialDepth:    opening.InitialDepth,
	}
}

func (p Presenter) Spawn(spawn fishprofiles.Spawn) SpawnView {
	return SpawnView{
		ProfileLabel:    spawn.Profile.Name,
		WaterBaseLabel:  waterpools.Name(spawn.Context.WaterPoolTag),
		InitialDistance: spawn.Context.InitialDistance,
		InitialDepth:    spawn.Context.InitialDepth,
		CandidateCount:  spawn.CandidateCount,
		HabitatLabels:   habitats.Names(spawn.Context.HabitatTags),
	}
}

func (p Presenter) Cast(context encounter.WaterContext, position, totalSlots, sectionWidth int) CastView {
	return CastView{
		WaterLabel:   context.Name,
		Position:     position,
		TotalSlots:   totalSlots,
		SectionWidth: sectionWidth,
	}
}

func (p Presenter) AnglerProfile(profile anglerprofiles.Profile) AnglerProfileView {
	return AnglerProfileView{
		ProfileID:         profile.ID,
		Name:              profile.Name,
		Description:       profile.Description,
		Details:           append([]string(nil), profile.Details...),
		StartingThread:    profile.StartingThread,
		DeckLabel:         playerDeckPresetName(profile.DeckPresetID),
		RodLabel:          rodPresetName(profile.RodPresetID),
		AttachmentLabel:   attachmentPresetName(profile.AttachmentPresetID),
		UnlockedByDefault: profile.UnlockedByDefault,
	}
}

func (p Presenter) RunIntro(title string, state run.State, route []run.NodeState) RunIntroView {
	routeLabels := make([]string, 0, len(route))
	for _, node := range route {
		routeLabels = append(routeLabels, p.runNodeLabel(node))
	}

	return RunIntroView{
		Title:       title,
		RouteLabels: routeLabels,
		Thread:      state.Thread.Current,
	}
}

func (p Presenter) RunNode(title string, state run.State) RunNodeView {
	return RunNodeView{
		Title:        title,
		ZoneLabel:    p.runZoneLabel(state.Progress.Current.ZoneID),
		NodeLabel:    p.runNodeLabel(state.Progress.Current),
		NodeKind:     state.Progress.Current.Kind,
		Thread:       state.Thread.Current,
		ThreadMax:    state.Thread.Maximum,
		CaptureCount: len(state.Captures),
	}
}

func (p Presenter) RunSummary(title string, state run.State) RunSummaryView {
	return RunSummaryView{
		Title:         title,
		Status:        state.Status,
		StatusLabel:   p.runStatusLabel(state.Status),
		Thread:        state.Thread.Current,
		ThreadMax:     state.Thread.Maximum,
		CaptureCount:  len(state.Captures),
		LastNodeLabel: p.runNodeLabel(state.Progress.Current),
	}
}

func (p Presenter) RunNodeSummary(title string, node run.NodeState, result run.EncounterResult, state run.State, next *run.NodeState) RunNodeSummaryView {
	view := RunNodeSummaryView{
		Title:               title,
		NodeLabel:           p.runNodeLabel(node),
		NodeKind:            node.Kind,
		OutcomeLabel:        p.runEncounterOutcomeLabel(result),
		Thread:              state.Thread.Current,
		ThreadMax:           state.Thread.Maximum,
		ThreadDelta:         -result.ThreadDamage,
		CaptureCount:        len(state.Captures),
		ContinuePromptLabel: "Pulsa Enter para continuar.",
	}
	if result.Capture != nil {
		view.LastCaptureLabel = result.Capture.FishName
	}
	if next != nil {
		view.NextNodeLabel = p.runNodeLabel(*next)
	}

	return view
}

func (p Presenter) eventLabel(event encounter.Event) string {
	if event.Kind == encounter.EventKindNone {
		return ""
	}

	eventLabel := string(event.Kind)
	if configuredLabel, ok := p.catalog.EncounterEvents[event.Kind]; ok {
		eventLabel = configuredLabel
	}

	outcomeLabel := ""
	if configuredOutcomeLabel, ok := p.catalog.EventOutcomes[event.Escaped]; ok {
		outcomeLabel = configuredOutcomeLabel
	}

	if outcomeLabel == "" {
		return eventLabel
	}

	return eventLabel + ": " + outcomeLabel
}

func (p Presenter) moveOptions() []MoveOption {
	return []MoveOption{
		{Index: 1, Move: domain.Blue, Label: p.playerMoveLabel(domain.Blue)},
		{Index: 2, Move: domain.Red, Label: p.playerMoveLabel(domain.Red)},
		{Index: 3, Move: domain.Yellow, Label: p.playerMoveLabel(domain.Yellow)},
	}
}

func (p Presenter) moveOptionsForSnapshot(player match.PlayerOptionsSnapshot) []MoveOption {
	moveOptions := p.moveOptions()
	for optionIndex := range moveOptions {
		for _, moveState := range player.Moves {
			if moveState.Move != moveOptions[optionIndex].Move {
				continue
			}

			moveOptions[optionIndex].CardHint = p.playerCardHint(moveState)
			moveOptions[optionIndex].RemainingUses = moveState.RemainingUses
			moveOptions[optionIndex].MaxUses = moveState.MaxUses
			moveOptions[optionIndex].Available = moveState.RemainingUses > 0
			moveOptions[optionIndex].RestoresOnRound = moveState.RestoresOnRound
			break
		}
	}

	return moveOptions
}

func (p Presenter) playerMoveLabel(move domain.Move) string {
	if label, ok := p.catalog.PlayerMoveLabels[move]; ok {
		return label
	}
	return move.String()
}

func (p Presenter) fishMoveLabel(move domain.Move) string {
	if label, ok := p.catalog.FishMoveLabels[move]; ok {
		return label
	}
	return move.String()
}

func (p Presenter) roundOutcomeLabel(outcome domain.RoundOutcome) string {
	if label, ok := p.catalog.RoundOutcomes[outcome]; ok {
		return label
	}
	return outcome.String()
}

func (p Presenter) encounterOutcomeLabel(status encounter.Status) string {
	if label, ok := p.catalog.EncounterResults[status]; ok {
		return label
	}
	return string(status)
}

func (p Presenter) endReasonLabel(reason encounter.EndReason) string {
	if label, ok := p.catalog.EndReasons[reason]; ok {
		return label
	}
	return string(reason)
}

func (p Presenter) playerCardHint(moveState match.MoveResourceSnapshot) string {
	if !moveState.HasTopCard {
		return ""
	}

	topCard := moveState.TopCard
	if topCard.Name != "" {
		return topCard.Name
	}
	if len(topCard.Effects) == 0 {
		return ""
	}

	parts := make([]string, 0, len(topCard.Effects))
	for _, effect := range topCard.Effects {
		impactParts := make([]string, 0, 5)
		if effect.DistanceShift != 0 {
			impactParts = append(impactParts, fmt.Sprintf("dist %+d", effect.DistanceShift))
		}
		if effect.DepthShift != 0 {
			impactParts = append(impactParts, fmt.Sprintf("prof %+d", effect.DepthShift))
		}
		if effect.CaptureDistanceBonus != 0 {
			impactParts = append(impactParts, fmt.Sprintf("capt %+d", effect.CaptureDistanceBonus))
		}
		if effect.SurfaceDepthBonus != 0 {
			impactParts = append(impactParts, fmt.Sprintf("sup %+d", effect.SurfaceDepthBonus))
		}
		if effect.ExhaustionCaptureDistanceBonus != 0 {
			impactParts = append(impactParts, fmt.Sprintf("baraja %+d", effect.ExhaustionCaptureDistanceBonus))
		}
		if len(impactParts) == 0 {
			continue
		}

		parts = append(parts, triggerLabel(effect.Trigger)+" "+strings.Join(impactParts, ", "))
	}

	return strings.Join(parts, " | ")
}

func (p Presenter) runStatusLabel(status run.Status) string {
	switch status {
	case run.StatusInProgress:
		return "run en curso"
	case run.StatusVictory:
		return "run completada"
	case run.StatusDefeat:
		return "run perdida"
	case run.StatusRetired:
		return "run retirada"
	default:
		return string(status)
	}
}

func (p Presenter) runEncounterOutcomeLabel(result run.EncounterResult) string {
	switch result.Outcome {
	case run.EncounterOutcomeCaptured:
		if result.Capture != nil {
			return "captura confirmada: " + result.Capture.FishName
		}
		return "captura confirmada"
	case run.EncounterOutcomeEscaped:
		if result.ThreadDamage > 0 {
			return fmt.Sprintf("el pez escapa y desgasta %d de hilo", result.ThreadDamage)
		}
		return "el pez escapa sin castigo de hilo"
	case run.EncounterOutcomeDefeated:
		return "run terminada en el nodo"
	case run.EncounterOutcomeRetired:
		return "retirada durante el nodo"
	default:
		return string(result.Outcome)
	}
}

func (p Presenter) runNodeLabel(node run.NodeState) string {
	suffix := strings.ReplaceAll(node.NodeID, "-", " ")
	switch node.Kind {
	case run.NodeKindStart:
		return "Inicio de run"
	case run.NodeKindFishing:
		return "Punto de pesca " + suffix
	case run.NodeKindService:
		return "Servicio " + suffix
	case run.NodeKindCheckpoint:
		return "Checkpoint " + suffix
	case run.NodeKindBoss:
		return "Encuentro final " + suffix
	case run.NodeKindEnd:
		return "Cierre de run"
	default:
		return node.NodeID
	}
}

func (p Presenter) runZoneLabel(zoneID string) string {
	switch zoneID {
	case "shoreline-cove":
		return "Fase 1 - Ensenada cercana"
	case "open-channel":
		return "Fase 2 - Canal abierto"
	case "broken-current":
		return "Fase 3 - Corriente irregular"
	case "reef-shadow":
		return "Fase 4 - Sombra de arrecife"
	case "tidal-gate":
		return "Fase 5 - Paso de marea"
	case "weed-pocket":
		return "Fase 6 - Bolsillo de maleza"
	case "stone-drop":
		return "Fase 7 - Caida de piedra"
	case "deep-ledge":
		return "Fase 8 - Cornisa profunda"
	default:
		return zoneID
	}
}

func playerDeckPresetName(id string) string {
	for _, preset := range playerprofiles.DefaultPresets() {
		if preset.ID == id {
			return preset.Name
		}
	}

	return id
}

func rodPresetName(id string) string {
	for _, preset := range rodpresets.DefaultPresets() {
		if preset.ID == id {
			return preset.Name
		}
	}

	return id
}

func attachmentPresetName(id string) string {
	for _, preset := range attachmentpresets.DefaultPresets() {
		if preset.ID == id {
			return preset.Name
		}
	}

	return id
}

func (p Presenter) fishDiscardView(discard match.FishDiscardSnapshot) FishDiscardView {
	currentCycleEntries := make([]FishDiscardEntryView, 0, len(discard.CurrentCycle.Entries))
	for _, entry := range discard.CurrentCycle.Entries {
		currentCycleEntries = append(currentCycleEntries, FishDiscardEntryView{
			Label: p.fishDiscardEntryLabel(entry),
		})
	}

	previousCycles := make([]FishDiscardCycleSummaryView, 0, len(discard.PreviousCycleStats))
	for _, previousCycle := range discard.PreviousCycleStats {
		previousCycles = append(previousCycles, FishDiscardCycleSummaryView{
			CycleNumber:  previousCycle.Number,
			TotalCards:   previousCycle.TotalCards,
			VisibleCards: previousCycle.VisibleCards,
			HiddenCards:  previousCycle.HiddenCards,
		})
	}

	return FishDiscardView{
		CurrentCycleNumber:     discard.CurrentCycle.Number,
		CurrentCycleTotalCards: discard.CurrentCycle.TotalCards,
		CurrentCycleEntries:    currentCycleEntries,
		PreviousCycles:         previousCycles,
		ShufflesOnRecycle:      discard.ShufflesOnRecycle,
		CardsToRemove:          discard.CardsToRemove,
		RecycleCount:           discard.RecycleCount,
	}
}

func (p Presenter) fishDiscardEntryLabel(entry match.FishDiscardEntryState) string {
	switch entry.Visibility {
	case cards.DiscardVisibilityMasked:
		return "?"
	case cards.DiscardVisibilityMoveOnly:
		return p.fishMoveLabel(entry.Move)
	case cards.DiscardVisibilityHidden:
		return ""
	case cards.DiscardVisibilityFull, "":
		if entry.Name != "" {
			return entry.Name
		}
		return p.fishMoveLabel(entry.Move)
	default:
		if entry.Name != "" {
			return entry.Name
		}
		return p.fishMoveLabel(entry.Move)
	}
}

func triggerLabel(trigger cards.Trigger) string {
	switch trigger {
	case cards.TriggerOnDraw:
		return "draw"
	case cards.TriggerOnOwnerWin:
		return "si gana"
	case cards.TriggerOnOwnerLose:
		return "si pierde"
	case cards.TriggerOnRoundDraw:
		return "empate"
	default:
		return "efecto"
	}
}
