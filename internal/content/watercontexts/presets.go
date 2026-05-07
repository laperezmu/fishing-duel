package watercontexts

import (
	"fmt"
	"pesca/internal/content/waterpools"
	"pesca/internal/encounter"
)

type ID string

const (
	ShorelineCove ID = "shoreline-cove"
	OpenChannel   ID = "open-channel"
	BrokenCurrent ID = "broken-current"
	ReefShadow    ID = "reef-shadow"
	TidalGate     ID = "tidal-gate"
	WeedPocket    ID = "weed-pocket"
	StoneDrop     ID = "stone-drop"
	WindLane      ID = "wind-lane"
	DeepLedge     ID = "deep-ledge"
)

type Preset struct {
	ID           ID
	Name         string
	Description  string
	Signals      []string
	PoolTag      waterpools.ID
	Distances    map[encounter.CastBand]int
	InitialDepth int
}

var (
	defaultPresets = []Preset{
		{
			ID:          ShorelineCove,
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
			ID:          OpenChannel,
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
			ID:          BrokenCurrent,
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
		{
			ID:          ReefShadow,
			Name:        "Sombra de arrecife",
			Description: "La actividad se pega a estructuras medias y castiga los lances demasiado cortos con ventanas menos limpias.",
			Signals: []string{
				"Manchas oscuras cortan el brillo del agua a media distancia.",
				"Hay vibracion irregular alrededor de piedra sumergida.",
			},
			PoolTag:      waterpools.MixedCurrent,
			Distances:    consecutiveDistances(1),
			InitialDepth: 2,
		},
		{
			ID:          TidalGate,
			Name:        "Paso de marea",
			Description: "La marea abre un pasillo estrecho donde los peces empujan entre profundidad media y corriente viva.",
			Signals: []string{
				"La corriente se acelera entre dos franjas mas calmas.",
				"Burbujas largas marcan un corredor hacia fuera.",
			},
			PoolTag:      waterpools.Offshore,
			Distances:    consecutiveDistances(2),
			InitialDepth: 2,
		},
		{
			ID:          WeedPocket,
			Name:        "Bolsillo de maleza",
			Description: "El agua parece retenida por vegetacion sumergida y favorece lances cortos con pelea mas vertical.",
			Signals: []string{
				"La superficie respira despacio entre parches verdes.",
				"Las ondulaciones cortas se frenan antes del canal central.",
			},
			PoolTag:      waterpools.Shoreline,
			Distances:    consecutiveDistances(0),
			InitialDepth: 2,
		},
		{
			ID:          StoneDrop,
			Name:        "Caida de piedra",
			Description: "El fondo rompe en un escalon brusco y la lectura premia sostener profundidad sin ceder toda la distancia.",
			Signals: []string{
				"El color del agua se oscurece de golpe despues del primer tercio.",
				"Pequenos golpes en superficie delatan actividad sobre la caida.",
			},
			PoolTag:      waterpools.MixedCurrent,
			Distances:    consecutiveDistances(1),
			InitialDepth: 3,
		},
		{
			ID:          WindLane,
			Name:        "Calle de viento",
			Description: "El viento ordena una via larga de espuma y la actividad buena aparece donde la deriva se estabiliza mar adentro.",
			Signals: []string{
				"Una linea continua de espuma apunta fuera de la costa.",
				"Las crestas rompen siempre en la misma direccion.",
			},
			PoolTag:      waterpools.Offshore,
			Distances:    consecutiveDistances(2),
			InitialDepth: 1,
		},
		{
			ID:          DeepLedge,
			Name:        "Cornisa profunda",
			Description: "La actividad rota por una pared lejana y obliga a probar aperturas largas con margen vertical desde el inicio.",
			Signals: []string{
				"El agua plana se corta en un borde oscuro mar adentro.",
				"Las burbujas suben aisladas desde una misma linea profunda.",
			},
			PoolTag:      waterpools.Offshore,
			Distances:    consecutiveDistances(3),
			InitialDepth: 3,
		},
	}
	defaultPresetByID = buildDefaultPresetIndex(defaultPresets)
	defaultPhaseOrder = []ID{ShorelineCove, OpenChannel, BrokenCurrent, ReefShadow, TidalGate, WeedPocket, StoneDrop, DeepLedge}
	defaultPhaseIndex = buildDefaultPhaseIndex(defaultPhaseOrder)
)

func (preset Preset) BuildContext() encounter.WaterContext {
	bandDistances := make(map[encounter.CastBand]int, len(preset.Distances))
	for band, initialDistance := range preset.Distances {
		bandDistances[band] = initialDistance
	}

	return encounter.WaterContext{
		ID:                  string(preset.ID),
		Name:                preset.Name,
		Description:         preset.Description,
		VisibleSignals:      append([]string(nil), preset.Signals...),
		PoolTag:             preset.PoolTag,
		BandInitialDistance: bandDistances,
		BaseInitialDepth:    preset.InitialDepth,
	}
}

func DefaultPresets() []Preset {
	presets := make([]Preset, 0, len(defaultPresets))
	for _, preset := range defaultPresets {
		presets = append(presets, clonePreset(preset))
	}

	return presets
}

func ResolveDefaultPreset(id ID) (Preset, error) {
	preset, ok := defaultPresetByID[id]
	if !ok {
		return Preset{}, fmt.Errorf("unknown water context preset %q", id)
	}

	return clonePreset(preset), nil
}

func DefaultPhaseLabel(id ID) string {
	preset, ok := defaultPresetByID[id]
	if !ok {
		return string(id)
	}
	phaseIndex, ok := defaultPhaseIndex[id]
	if !ok {
		return preset.Name
	}

	return fmt.Sprintf("Fase %d - %s", phaseIndex+1, preset.Name)
}

func consecutiveDistances(start int) map[encounter.CastBand]int {
	bands := encounter.OrderedCastBands()
	distances := make(map[encounter.CastBand]int, len(bands))
	for offset, band := range bands {
		distances[band] = start + offset
	}

	return distances
}

func clonePreset(preset Preset) Preset {
	clonedPreset := preset
	clonedPreset.Signals = append([]string(nil), preset.Signals...)
	clonedPreset.Distances = make(map[encounter.CastBand]int, len(preset.Distances))
	for band, distance := range preset.Distances {
		clonedPreset.Distances[band] = distance
	}

	return clonedPreset
}

func buildDefaultPresetIndex(presets []Preset) map[ID]Preset {
	index := make(map[ID]Preset, len(presets))
	for _, preset := range presets {
		index[preset.ID] = clonePreset(preset)
	}

	return index
}

func buildDefaultPhaseIndex(order []ID) map[ID]int {
	index := make(map[ID]int, len(order))
	for i, id := range order {
		index[id] = i
	}

	return index
}
