package presentation

import (
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
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
			encounter.EndReasonTrackCapture: "captura por acercarlo hasta el jugador",
			encounter.EndReasonTrackEscape:  "escape por superar la distancia maxima alcanzable",
			encounter.EndReasonDepthEscape:  "escape por bajar mas alla de la profundidad alcanzable",
			encounter.EndReasonSplashEscape: "escape por chapoteo en superficie",
			encounter.EndReasonDeckCapture:  "captura al agotar la baraja con distancia 2 o menor",
			encounter.EndReasonDeckEscape:   "escape al agotar la baraja con distancia mayor que 2",
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

func (p Presenter) Status(state match.State) StatusView {
	return StatusView{
		RoundNumber:               state.Round + 1,
		FishDistance:              state.Encounter.Distance,
		FishDepth:                 state.Encounter.Depth,
		SurfaceDepth:              state.Encounter.Config.SurfaceDepth,
		MaxDistance:               state.PlayerRig.MaxDistance,
		MaxDepth:                  state.PlayerRig.MaxDepth,
		CaptureDistance:           state.Encounter.Config.CaptureDistance,
		ExhaustionCaptureDistance: state.Encounter.Config.ExhaustionCaptureDistance,
		ActiveCards:               state.Deck.ActiveCards,
		DiscardCards:              state.Deck.DiscardCards,
		RecycleCount:              state.Deck.RecycleCount,
		PlayerWins:                state.Stats.PlayerWins,
		FishWins:                  state.Stats.FishWins,
		Draws:                     state.Stats.Draws,
		MoveOptions:               p.moveOptionsForState(state),
	}
}

func (p Presenter) Round(result match.RoundResult) RoundView {
	return RoundView{
		Status:       p.Status(result.State),
		PlayerMove:   result.PlayerMove,
		FishMove:     result.FishMove,
		PlayerLabel:  p.playerMoveLabel(result.PlayerMove),
		FishLabel:    p.fishMoveLabel(result.FishMove),
		Outcome:      result.Outcome,
		OutcomeLabel: p.roundOutcomeLabel(result.Outcome),
		EventLabel:   p.eventLabel(result.State.Encounter.LastEvent),
	}
}

func (p Presenter) Summary(state match.State) SummaryView {
	return SummaryView{
		TotalRounds:     state.Round,
		FishDistance:    state.Encounter.Distance,
		FishDepth:       state.Encounter.Depth,
		EncounterStatus: state.Encounter.Status,
		OutcomeLabel:    p.encounterOutcomeLabel(state.Encounter.Status),
		EndReasonLabel:  p.endReasonLabel(state.Encounter.EndReason),
		PlayerWins:      state.Stats.PlayerWins,
		FishWins:        state.Stats.FishWins,
		Draws:           state.Stats.Draws,
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

func (p Presenter) moveOptionsForState(state match.State) []MoveOption {
	moveOptions := p.moveOptions()
	for optionIndex := range moveOptions {
		for _, moveState := range state.PlayerMoves.Moves {
			if moveState.Move != moveOptions[optionIndex].Move {
				continue
			}

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
