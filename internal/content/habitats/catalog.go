package habitats

import "fmt"

type Tag string

const (
	Bottom    Tag = "bottom"
	Channel   Tag = "channel"
	OpenWater Tag = "open-water"
	Surface   Tag = "surface"
	Weed      Tag = "weed"
	Rock      Tag = "rock"
)

type Habitat struct {
	Tag         Tag
	Name        string
	Description string
}

func DefaultCatalog() []Habitat {
	return []Habitat{
		{Tag: Bottom, Name: "Fondo", Description: "Capas profundas o fondos pesados del agua."},
		{Tag: Channel, Name: "Canal", Description: "Canales o pasos de corriente donde la profundidad manda."},
		{Tag: OpenWater, Name: "Agua abierta", Description: "Ventanas expuestas y recorridos largos del lance."},
		{Tag: Surface, Name: "Superficie", Description: "Capas altas y lectura cercana a espuma o reflejos."},
		{Tag: Weed, Name: "Maleza", Description: "Cobertura de vegetacion o agua sucia donde se esconde el pez."},
		{Tag: Rock, Name: "Roca", Description: "Fondos de piedra y zonas de roce o refugio duro."},
	}
}

func (tag Tag) Validate() error {
	for _, habitat := range DefaultCatalog() {
		if habitat.Tag == tag {
			return nil
		}
	}

	return fmt.Errorf("unknown habitat tag %q", tag)
}

func Strings(tags []Tag) []string {
	values := make([]string, 0, len(tags))
	for _, tag := range tags {
		values = append(values, string(tag))
	}

	return values
}

func Names(tags []Tag) []string {
	values := make([]string, 0, len(tags))
	for _, tag := range tags {
		values = append(values, Name(tag))
	}

	return values
}

func Name(tag Tag) string {
	for _, habitat := range DefaultCatalog() {
		if habitat.Tag == tag {
			return habitat.Name
		}
	}

	return string(tag)
}
