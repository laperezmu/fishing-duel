package fishprofiles

import _ "embed"

var (
	//go:embed data/default_profiles.json
	defaultProfilesJSON []byte
	//go:embed data/default_pools.json
	defaultPoolsJSON []byte
)

func DefaultCatalog() Catalog {
	catalog, err := LoadCatalog(defaultProfilesJSON, defaultPoolsJSON)
	if err != nil {
		panic(err)
	}

	return catalog
}
