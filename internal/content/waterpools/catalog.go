package waterpools

import "fmt"

type ID string

const (
	Shoreline    ID = "shoreline"
	Offshore     ID = "offshore"
	MixedCurrent ID = "mixed-current"
)

type Pool struct {
	ID          ID
	Name        string
	Description string
}

func DefaultCatalog() []Pool {
	return []Pool{
		{ID: Shoreline, Name: "Costa cercana", Description: "Aguas pegadas a la orilla y primeras franjas utiles del lance."},
		{ID: Offshore, Name: "Mar abierto", Description: "Actividad empujada hacia afuera y ventanas mas largas del agua."},
		{ID: MixedCurrent, Name: "Corriente mixta", Description: "Agua irregular con actividad cambiante y transiciones de lectura."},
	}
}

func (id ID) Validate() error {
	for _, pool := range DefaultCatalog() {
		if pool.ID == id {
			return nil
		}
	}

	return fmt.Errorf("unknown water pool %q", id)
}

func (id ID) String() string {
	return string(id)
}

func Name(id ID) string {
	for _, pool := range DefaultCatalog() {
		if pool.ID == id {
			return pool.Name
		}
	}

	return string(id)
}
