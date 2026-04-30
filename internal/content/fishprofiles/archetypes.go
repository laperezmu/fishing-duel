package fishprofiles

import "fmt"

type ArchetypeID string

const (
	ArchetypeBaselineCycle      ArchetypeID = "baseline_cycle"
	ArchetypeHorizontalPressure ArchetypeID = "horizontal_pressure"
	ArchetypeVerticalEscape     ArchetypeID = "vertical_escape"
	ArchetypeSurfaceControl     ArchetypeID = "surface_control"
	ArchetypeDrawTempo          ArchetypeID = "draw_tempo"
	ArchetypeDeckExhaustion     ArchetypeID = "deck_exhaustion"
	ArchetypeHybridPressure     ArchetypeID = "hybrid_pressure"
)

type Archetype struct {
	ID          ArchetypeID
	Name        string
	Description string
}

func DefaultArchetypes() []Archetype {
	return []Archetype{
		{ID: ArchetypeBaselineCycle, Name: "Ciclo base", Description: "Perfil generico sin presion especializada."},
		{ID: ArchetypeHorizontalPressure, Name: "Presion horizontal", Description: "Favorece escape o control hacia mar abierto."},
		{ID: ArchetypeVerticalEscape, Name: "Presion vertical", Description: "Trabaja la profundidad como eje principal del duelo."},
		{ID: ArchetypeSurfaceControl, Name: "Control de superficie", Description: "Mantiene al pez cerca de capas altas y eventos visibles."},
		{ID: ArchetypeDrawTempo, Name: "Tempo de apertura", Description: "Concentra valor al revelar cartas de apertura."},
		{ID: ArchetypeDeckExhaustion, Name: "Agotamiento", Description: "Juega alrededor del cierre por agotamiento del mazo."},
		{ID: ArchetypeHybridPressure, Name: "Presion hibrida", Description: "Combina ventajas de apertura y respuestas por outcome."},
	}
}

func (id ArchetypeID) Validate() error {
	for _, archetype := range DefaultArchetypes() {
		if archetype.ID == id {
			return nil
		}
	}

	return fmt.Errorf("unknown fish archetype %q", id)
}

func Name(id ArchetypeID) string {
	for _, archetype := range DefaultArchetypes() {
		if archetype.ID == id {
			return archetype.Name
		}
	}

	return string(id)
}
