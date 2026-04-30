package watercontexts

import (
	"pesca/internal/content/waterpools"
	"pesca/internal/encounter"
)

type Preset struct {
	ID           string
	Name         string
	Description  string
	Signals      []string
	PoolTag      waterpools.ID
	Distances    map[encounter.CastBand]int
	InitialDepth int
}

func (preset Preset) BuildContext() encounter.WaterContext {
	bandDistances := make(map[encounter.CastBand]int, len(preset.Distances))
	for band, initialDistance := range preset.Distances {
		bandDistances[band] = initialDistance
	}

	return encounter.WaterContext{
		ID:                  preset.ID,
		Name:                preset.Name,
		Description:         preset.Description,
		VisibleSignals:      append([]string(nil), preset.Signals...),
		PoolTag:             preset.PoolTag,
		BandInitialDistance: bandDistances,
		BaseInitialDepth:    preset.InitialDepth,
	}
}

func DefaultPresets() []Preset {
	return []Preset{
		{
			ID:          "shoreline-cove",
			Name:        "Ensenada cercana",
			Description: "Las marcas del agua se rompen cerca de la orilla y la actividad util aparece antes de forzar un lance largo.",
			Signals: []string{
				"Espuma corta pegada a la costa.",
				"Destellos y remolinos en la primera mitad del agua.",
			},
			PoolTag:      waterpools.Shoreline,
			Distances:    consecutiveDistances(0),
			InitialDepth: 1,
		},
		{
			ID:          "open-channel",
			Name:        "Canal abierto",
			Description: "La corriente empuja la actividad hacia fuera y conviene sostener el lance hasta la parte media o larga del agua.",
			Signals: []string{
				"La corriente principal tira mar adentro.",
				"Las ondas vivas aparecen lejos de la orilla.",
			},
			PoolTag:      waterpools.Offshore,
			Distances:    consecutiveDistances(1),
			InitialDepth: 1,
		},
		{
			ID:          "broken-current",
			Name:        "Corriente irregular",
			Description: "Las marcas del agua cambian rapido y el lance mas extremo pierde consistencia aunque el centro del canal siga siendo util.",
			Signals: []string{
				"La espuma cruza en diagonales y corta el agua en dos.",
				"Hay actividad intermitente a media distancia.",
			},
			PoolTag:      waterpools.MixedCurrent,
			Distances:    consecutiveDistances(0),
			InitialDepth: 1,
		},
	}
}

func consecutiveDistances(start int) map[encounter.CastBand]int {
	bands := encounter.OrderedCastBands()
	distances := make(map[encounter.CastBand]int, len(bands))
	for offset, band := range bands {
		distances[band] = start + offset
	}

	return distances
}
