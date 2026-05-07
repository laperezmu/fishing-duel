package fishprofiles

import (
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/content/habitats"
	"pesca/internal/content/waterpools"
)

type PoolID string

func (id PoolID) Validate() error {
	if id == "" {
		return fmt.Errorf("pool id is required")
	}

	return nil
}

type Pool struct {
	ID          PoolID
	Name        string
	Description string
	Entries     []PoolEntry
}

type PoolEntry struct {
	ProfileID ProfileID
	Weight    int
}

func (pool Pool) Validate() error {
	if pool.ID == "" {
		return pool.ID.Validate()
	}
	if pool.Name == "" {
		return fmt.Errorf("pool name is required")
	}
	if len(pool.Entries) == 0 {
		return fmt.Errorf("pool %s must reference at least one profile id", pool.ID)
	}
	seenProfileIDs := make(map[ProfileID]struct{}, len(pool.Entries))
	for _, entry := range pool.Entries {
		if entry.ProfileID == "" {
			return fmt.Errorf("pool %s has an empty profile id reference", pool.ID)
		}
		if entry.Weight <= 0 {
			return fmt.Errorf("pool %s has invalid weight %d for profile id %s", pool.ID, entry.Weight, entry.ProfileID)
		}
		if _, exists := seenProfileIDs[entry.ProfileID]; exists {
			return fmt.Errorf("pool %s references duplicated profile id %s", pool.ID, entry.ProfileID)
		}
		seenProfileIDs[entry.ProfileID] = struct{}{}
	}

	return nil
}

type Catalog struct {
	profiles     []Profile
	pools        []Pool
	profileIndex map[ProfileID]Profile
	poolIndex    map[PoolID]Pool
}

func NewCatalog(profiles []Profile, pools []Pool) (Catalog, error) {
	if len(profiles) == 0 {
		return Catalog{}, fmt.Errorf("catalog requires at least one fish profile")
	}

	seenProfileIDs := make(map[ProfileID]struct{}, len(profiles))
	clonedProfiles := make([]Profile, 0, len(profiles))
	profileIndex := make(map[ProfileID]Profile, len(profiles))
	for _, profile := range profiles {
		if err := profile.Validate(); err != nil {
			return Catalog{}, fmt.Errorf("profile %s: %w", profile.ID, err)
		}
		if _, exists := seenProfileIDs[profile.ID]; exists {
			return Catalog{}, fmt.Errorf("duplicated fish profile id %s", profile.ID)
		}
		seenProfileIDs[profile.ID] = struct{}{}
		clonedProfile := cloneProfile(profile)
		clonedProfiles = append(clonedProfiles, clonedProfile)
		profileIndex[profile.ID] = clonedProfile
	}

	seenPoolIDs := make(map[PoolID]struct{}, len(pools))
	clonedPools := make([]Pool, 0, len(pools))
	poolIndex := make(map[PoolID]Pool, len(pools))
	for _, pool := range pools {
		if err := pool.Validate(); err != nil {
			return Catalog{}, err
		}
		if _, exists := seenPoolIDs[pool.ID]; exists {
			return Catalog{}, fmt.Errorf("duplicated fish pool id %s", pool.ID)
		}
		for _, entry := range pool.Entries {
			if _, exists := seenProfileIDs[entry.ProfileID]; !exists {
				return Catalog{}, fmt.Errorf("fish pool %s references unknown profile id %s", pool.ID, entry.ProfileID)
			}
		}
		seenPoolIDs[pool.ID] = struct{}{}
		clonedPool := clonePool(pool)
		clonedPools = append(clonedPools, clonedPool)
		poolIndex[pool.ID] = clonedPool
	}

	return Catalog{profiles: clonedProfiles, pools: clonedPools, profileIndex: profileIndex, poolIndex: poolIndex}, nil
}

func (catalog Catalog) Profiles() []Profile {
	profiles := make([]Profile, 0, len(catalog.profiles))
	for _, profile := range catalog.profiles {
		profiles = append(profiles, cloneProfile(profile))
	}

	return profiles
}

func (catalog Catalog) Pools() []Pool {
	pools := make([]Pool, 0, len(catalog.pools))
	for _, pool := range catalog.pools {
		pools = append(pools, clonePool(pool))
	}

	return pools
}

func (catalog Catalog) ResolvePool(poolID PoolID) ([]Profile, error) {
	pool, err := catalog.PoolByID(poolID)
	if err != nil {
		return nil, err
	}

	profiles := make([]Profile, 0, resolvedPoolProfileCount(pool))
	for _, entry := range pool.Entries {
		profile, err := catalog.ProfileByID(entry.ProfileID)
		if err != nil {
			return nil, fmt.Errorf("resolve fish pool %s: %w", poolID, err)
		}
		for range entry.Weight {
			profiles = append(profiles, cloneProfile(profile))
		}
	}

	return profiles, nil
}

func (catalog Catalog) ProfileByID(profileID ProfileID) (Profile, error) {
	profile, ok := catalog.profileIndex[profileID]
	if !ok {
		return Profile{}, fmt.Errorf("unknown fish profile id %s", profileID)
	}

	return cloneProfile(profile), nil
}

func (catalog Catalog) PoolByID(poolID PoolID) (Pool, error) {
	pool, ok := catalog.poolIndex[poolID]
	if !ok {
		return Pool{}, fmt.Errorf("unknown fish pool id %s", poolID)
	}

	return clonePool(pool), nil
}

func cloneProfile(profile Profile) Profile {
	clonedProfile := profile
	clonedProfile.Details = append([]string(nil), profile.Details...)
	clonedProfile.Appearance.WaterPoolTags = append([]waterpools.ID(nil), profile.Appearance.WaterPoolTags...)
	clonedProfile.Appearance.RequiredHabitatTags = append([]habitats.Tag(nil), profile.Appearance.RequiredHabitatTags...)
	clonedProfile.Cards = append([]CardPattern(nil), profile.Cards...)
	for index := range clonedProfile.Cards {
		clonedProfile.Cards[index].Effects = append([]cards.CardEffect(nil), clonedProfile.Cards[index].Effects...)
	}

	return clonedProfile
}

func clonePool(pool Pool) Pool {
	clonedPool := pool
	clonedPool.Entries = append([]PoolEntry(nil), pool.Entries...)

	return clonedPool
}

func resolvedPoolProfileCount(pool Pool) int {
	total := 0
	for _, entry := range pool.Entries {
		total += entry.Weight
	}

	return total
}
