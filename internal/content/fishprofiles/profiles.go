package fishprofiles

import (
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/content/habitats"
	"pesca/internal/content/waterpools"
	"pesca/internal/domain"
)

type CardPattern struct {
	Name              string
	Summary           string
	Move              domain.Move
	Effects           []cards.CardEffect
	DiscardVisibility cards.DiscardVisibility
}

func (pattern CardPattern) BuildCard() cards.FishCard {
	var card cards.FishCard
	if pattern.Name != "" || pattern.Summary != "" {
		card = cards.NewNamedFishCard(pattern.Name, pattern.Summary, pattern.Move, pattern.Effects...)
	} else {
		card = cards.NewFishCard(pattern.Move, pattern.Effects...)
	}

	if pattern.DiscardVisibility != "" {
		card.DiscardVisibility = pattern.DiscardVisibility
	}

	return card
}

type Profile struct {
	ID            string
	ArchetypeID   ArchetypeID
	Name          string
	Description   string
	Details       []string
	Appearance    Appearance
	Cards         []CardPattern
	CardsToRemove int
	Shuffle       bool
}

func (profile Profile) Validate() error {
	if profile.ID == "" {
		return fmt.Errorf("profile id is required")
	}
	if profile.Name == "" {
		return fmt.Errorf("profile name is required")
	}
	if err := profile.ArchetypeID.Validate(); err != nil {
		return err
	}
	if err := profile.Appearance.Validate(); err != nil {
		return fmt.Errorf("appearance: %w", err)
	}

	return nil
}

func (profile Profile) BuildPreset() FishDeckPreset {
	return FishDeckPreset{
		ID:            profile.ID,
		ArchetypeID:   profile.ArchetypeID,
		Name:          profile.Name,
		Description:   profile.Description,
		Details:       append([]string(nil), profile.Details...),
		FishCards:     profile.BuildCards(),
		CardsToRemove: profile.CardsToRemove,
		Shuffle:       profile.Shuffle,
	}
}

func (profile Profile) BuildCards() []cards.FishCard {
	builtCards := make([]cards.FishCard, 0, len(profile.Cards))
	for _, pattern := range profile.Cards {
		builtCards = append(builtCards, pattern.BuildCard())
	}

	return builtCards
}

func DefaultProfiles() []Profile {
	return []Profile{
		{
			ID:          "classic",
			ArchetypeID: ArchetypeBaselineCycle,
			Name:        "Clasico",
			Description: "Baraja base de referencia sin efectos, util para comparar cambios de sistema.",
			Appearance: Appearance{
				WaterPoolTags:      []waterpools.ID{waterpools.Shoreline, waterpools.MixedCurrent},
				MinInitialDistance: 0,
				MaxInitialDistance: 2,
				MinInitialDepth:    0,
				MaxInitialDepth:    2,
			},
			Details: []string{
				"Arquetipo: ciclo base sin presion especializada.",
				"Nueve cartas lisas sin efectos: tres rojas, tres azules y tres amarillas.",
				"Orden: barajada antes de empezar y en cada reciclado.",
				"Reciclado: retira 3 cartas por ciclo.",
			},
			Cards: []CardPattern{
				{Move: domain.Blue},
				{Move: domain.Blue},
				{Move: domain.Blue},
				{Move: domain.Red},
				{Move: domain.Red},
				{Move: domain.Red},
				{Move: domain.Yellow},
				{Move: domain.Yellow},
				{Move: domain.Yellow},
			},
			CardsToRemove: 3,
			Shuffle:       true,
		},
		{
			ID:          "hooked-opening",
			ArchetypeID: ArchetypeDrawTempo,
			Name:        "Apertura con anzuelo",
			Description: "Perfil de tempo que concentra su presion al revelar cartas de apertura.",
			Appearance: Appearance{
				WaterPoolTags:      []waterpools.ID{waterpools.Shoreline, waterpools.Offshore},
				MinInitialDistance: 1,
				MaxInitialDistance: 5,
				MinInitialDepth:    0,
				MaxInitialDepth:    2,
			},
			Details: []string{
				"Arquetipo: draw_tempo.",
				"Rojo - Tiron de apertura: al revelarse permite capturar desde un paso mas lejos ese round.",
				"Azul - Salto de espuma: al revelarse el pez cuenta como un nivel mas cerca de la superficie ese round.",
				"Amarillo - Ultima ventana: al revelarse amplia el margen de captura cuando se agota la baraja ese round.",
				"Orden: fijo para probar cada apertura de forma reproducible.",
			},
			Cards: []CardPattern{
				{Name: "Tiron de apertura", Summary: "Permite capturar desde un paso mas lejos este round.", Move: domain.Red, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}}},
				{Name: "Salto de espuma", Summary: "El pez cuenta como un nivel mas cerca de la superficie este round.", Move: domain.Blue, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnDraw, SurfaceDepthBonus: 1}}},
				{Name: "Ultima ventana", Summary: "Amplia el margen de captura por agotamiento durante este round.", Move: domain.Yellow, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnDraw, ExhaustionCaptureDistanceBonus: 1}}},
				{Name: "Tiron de apertura", Summary: "Permite capturar desde un paso mas lejos este round.", Move: domain.Red, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}}},
			},
			CardsToRemove: 0,
			Shuffle:       false,
		},
		{
			ID:          "horizontal-pressure",
			ArchetypeID: ArchetypeHorizontalPressure,
			Name:        "Presion horizontal",
			Description: "Perfil que prioriza empujar el encuentro hacia mar abierto cuando obtiene ventaja.",
			Appearance: Appearance{
				WaterPoolTags:       []waterpools.ID{waterpools.Offshore},
				MinInitialDistance:  3,
				MaxInitialDistance:  5,
				MinInitialDepth:     0,
				MaxInitialDepth:     2,
				RequiredHabitatTags: []habitats.Tag{habitats.OpenWater},
			},
			Details: []string{
				"Arquetipo: horizontal_pressure.",
				"Azul - Oleaje abierto: si gana, empuja un paso hacia mar abierto.",
				"Rojo - Deriva neutra: en empate gana un paso hacia el escape horizontal.",
				"Amarillo - Oleaje abierto: si gana, empuja un paso hacia mar abierto.",
				"Orden: fijo para leer la presion horizontal de forma controlada.",
			},
			Cards: []CardPattern{
				{Name: "Oleaje abierto", Summary: "Si gana, empuja un paso hacia mar abierto.", Move: domain.Blue, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnOwnerWin, DistanceShift: 1}}},
				{Name: "Deriva neutra", Summary: "En empate gana un paso hacia el escape horizontal.", Move: domain.Red, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnRoundDraw, DistanceShift: 1}}},
				{Name: "Oleaje abierto", Summary: "Si gana, empuja un paso hacia mar abierto.", Move: domain.Yellow, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnOwnerWin, DistanceShift: 1}}},
				{Name: "Deriva neutra", Summary: "En empate gana un paso hacia el escape horizontal.", Move: domain.Red, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnRoundDraw, DistanceShift: 1}}},
			},
			CardsToRemove: 0,
			Shuffle:       false,
		},
		{
			ID:          "vertical-pressure",
			ArchetypeID: ArchetypeVerticalEscape,
			Name:        "Presion vertical",
			Description: "Perfil orientado a hundirse al ganar y a respirar al perder.",
			Appearance: Appearance{
				WaterPoolTags:       []waterpools.ID{waterpools.Offshore, waterpools.MixedCurrent},
				MinInitialDistance:  2,
				MaxInitialDistance:  5,
				MinInitialDepth:     2,
				MaxInitialDepth:     4,
				RequiredHabitatTags: []habitats.Tag{habitats.Bottom, habitats.Channel},
			},
			Details: []string{
				"Arquetipo: vertical_escape.",
				"Azul - Tiron al fondo: si gana, el pez baja un nivel mas profundo.",
				"Rojo - Respiro corto: si pierde, el pez sube un nivel hacia la superficie.",
				"Amarillo - Caida larga: si gana, el pez baja un nivel mas profundo.",
				"Orden: fijo para validar la capa vertical paso a paso.",
			},
			Cards: []CardPattern{
				{Name: "Tiron al fondo", Summary: "Si gana, baja un nivel mas profundo.", Move: domain.Blue, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnOwnerWin, DepthShift: 1}}},
				{Name: "Respiro corto", Summary: "Si pierde, sube un nivel hacia la superficie.", Move: domain.Red, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1}}},
				{Name: "Caida larga", Summary: "Si gana, baja un nivel mas profundo.", Move: domain.Yellow, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnOwnerWin, DepthShift: 1}}},
				{Name: "Respiro corto", Summary: "Si pierde, sube un nivel hacia la superficie.", Move: domain.Blue, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1}}},
			},
			CardsToRemove: 0,
			Shuffle:       false,
		},
		{
			ID:          "surface-control",
			ArchetypeID: ArchetypeSurfaceControl,
			Name:        "Control de superficie",
			Description: "Perfil que gira en torno a mantener el pez cerca de la capa superficial y forzar eventos legibles.",
			Appearance: Appearance{
				WaterPoolTags:       []waterpools.ID{waterpools.Shoreline, waterpools.MixedCurrent},
				MinInitialDistance:  0,
				MaxInitialDistance:  3,
				MinInitialDepth:     0,
				MaxInitialDepth:     1,
				RequiredHabitatTags: []habitats.Tag{habitats.Surface},
			},
			Details: []string{
				"Arquetipo: surface_control.",
				"Azul - Rebote de espuma: si pierde, el pez sube un nivel.",
				"Rojo - Marea calma: en empate el pez cuenta como un nivel mas cerca de la superficie ese round.",
				"Amarillo - Rebote de espuma: si pierde, el pez sube un nivel.",
				"Lectura parcial: una Marea calma solo deja visible el movimiento en el historial del pez.",
				"Orden: fijo para probar lecturas de superficie y cierre visual.",
			},
			Cards: []CardPattern{
				{Name: "Rebote de espuma", Summary: "Si pierde, sube un nivel hacia la superficie.", Move: domain.Blue, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1}}},
				{Name: "Marea calma", Summary: "En empate el pez cuenta como un nivel mas cerca de la superficie.", Move: domain.Red, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnRoundDraw, SurfaceDepthBonus: 1}}, DiscardVisibility: cards.DiscardVisibilityMoveOnly},
				{Name: "Rebote de espuma", Summary: "Si pierde, sube un nivel hacia la superficie.", Move: domain.Yellow, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1}}},
				{Name: "Marea calma", Summary: "En empate el pez cuenta como un nivel mas cerca de la superficie.", Move: domain.Red, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnRoundDraw, SurfaceDepthBonus: 1}}},
			},
			CardsToRemove: 0,
			Shuffle:       false,
		},
		{
			ID:          "deck-exhaustion",
			ArchetypeID: ArchetypeDeckExhaustion,
			Name:        "Agotamiento de mazo",
			Description: "Perfil que concentra su plan en el cierre por agotamiento y en ventanas cortas de cierre.",
			Appearance: Appearance{
				WaterPoolTags:       []waterpools.ID{waterpools.Shoreline, waterpools.MixedCurrent},
				MinInitialDistance:  0,
				MaxInitialDistance:  2,
				MinInitialDepth:     0,
				MaxInitialDepth:     2,
				RequiredHabitatTags: []habitats.Tag{habitats.Weed, habitats.Rock},
			},
			Details: []string{
				"Arquetipo: deck_exhaustion.",
				"Rojo - Ultima ventana: al revelarse amplia la captura por agotamiento ese round.",
				"Azul - Tiron de cierre: al revelarse acerca la captura ese round.",
				"Amarillo - Ultima ventana: repite el empuje al cierre por agotamiento.",
				"Orden: fijo para revisar finales de encuentro y baraja corta.",
			},
			Cards: []CardPattern{
				{Name: "Ultima ventana", Summary: "Amplia el margen de captura por agotamiento durante este round.", Move: domain.Red, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnDraw, ExhaustionCaptureDistanceBonus: 1}}},
				{Name: "Tiron de cierre", Summary: "Permite capturar desde un paso mas lejos este round.", Move: domain.Blue, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}}},
				{Name: "Ultima ventana", Summary: "Amplia el margen de captura por agotamiento durante este round.", Move: domain.Yellow, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnDraw, ExhaustionCaptureDistanceBonus: 1}}},
				{Name: "Tiron de cierre", Summary: "Permite capturar desde un paso mas lejos este round.", Move: domain.Blue, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}}},
			},
			CardsToRemove: 0,
			Shuffle:       false,
		},
		{
			ID:          "mixed-current",
			ArchetypeID: ArchetypeHybridPressure,
			Name:        "Corriente mixta",
			Description: "Perfil mixto que combina ventajas al revelarse con respuestas segun el resultado.",
			Appearance: Appearance{
				WaterPoolTags:      []waterpools.ID{waterpools.MixedCurrent},
				MinInitialDistance: 1,
				MaxInitialDistance: 4,
				MinInitialDepth:    1,
				MaxInitialDepth:    3,
			},
			Details: []string{
				"Arquetipo: hybrid_pressure.",
				"Rojo - Corriente cerrada: al revelarse amplia la captura y, si pierde, sube un nivel hacia la superficie.",
				"Azul - Oleaje abierto: al revelarse se acerca a la superficie y, si gana, empuja un paso hacia mar abierto.",
				"Amarillo - Deriva neutra: en empate gana un paso hacia el escape horizontal.",
				"Orden: fijo para revisar el pipeline completo round a round.",
			},
			Cards: []CardPattern{
				{Name: "Corriente cerrada", Summary: "Amplia la captura al revelarse y sube si pierde.", Move: domain.Red, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}, {Trigger: cards.TriggerOnOwnerLose, DepthShift: -1}}},
				{Name: "Oleaje abierto", Summary: "Se acerca a superficie al revelarse y empuja si gana.", Move: domain.Blue, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnDraw, SurfaceDepthBonus: 1}, {Trigger: cards.TriggerOnOwnerWin, DistanceShift: 1}}},
				{Name: "Deriva neutra", Summary: "En empate gana un paso hacia mar abierto.", Move: domain.Yellow, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnRoundDraw, DistanceShift: 1}}},
				{Name: "Corriente cerrada", Summary: "Amplia la captura al revelarse y sube si pierde.", Move: domain.Red, Effects: []cards.CardEffect{{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1}}},
			},
			CardsToRemove: 0,
			Shuffle:       false,
		},
	}
}
