package fishprofiles

import _ "embed"

var (
	//go:embed data/default_profiles.json
	defaultProfilesJSON []byte
	//go:embed data/default_pools.json
	defaultPoolsJSON []byte
	defaultCatalog   = loadDefaultCatalog()
)

func DefaultCatalog() Catalog {
	return defaultCatalog
}

const DefaultEncounterFishPoolID PoolID = "all-default-fish"

func ResolveDefaultProfile(id ProfileID) (Profile, error) {
	return DefaultCatalog().ProfileByID(id)
}

func ResolveDefaultPool(id PoolID) (Pool, error) {
	return DefaultCatalog().PoolByID(id)
}

func loadDefaultCatalog() Catalog {
	catalog, err := LoadCatalog(defaultProfilesJSON, defaultPoolsJSON)
	if err != nil {
		panic(err)
	}

	return catalog
}
