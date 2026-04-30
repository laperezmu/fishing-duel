package presentation

import (
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
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
