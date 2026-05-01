package fishprofiles

import (
	"fmt"
	"pesca/internal/content/habitats"
	"pesca/internal/content/waterpools"
	"pesca/internal/encounter"
	"sort"
	"strings"
)

type SpawnRandomizer interface {
	Intn(n int) int
}

type Appearance struct {
	WaterPoolTags       []waterpools.ID
	MinInitialDistance  int
	MaxInitialDistance  int
	MinInitialDepth     int
	MaxInitialDepth     int
	RequiredHabitatTags []habitats.Tag
}

type SpawnContext struct {
	WaterPoolTag    waterpools.ID
	InitialDistance int
	InitialDepth    int
	HabitatTags     []habitats.Tag
}

type Spawn struct {
	Profile        Profile
	Context        SpawnContext
	CandidateCount int
}

func NewSpawnContext(opening encounter.Opening, habitatTags []habitats.Tag) (SpawnContext, error) {
	context := SpawnContext{
		WaterPoolTag:    opening.WaterContext.PoolTag,
		InitialDistance: opening.InitialDistance,
		InitialDepth:    opening.InitialDepth,
		HabitatTags:     append([]habitats.Tag(nil), habitatTags...),
	}

	if err := context.Validate(); err != nil {
		return SpawnContext{}, err
	}

	return context, nil
}

func (context SpawnContext) Validate() error {
	if context.WaterPoolTag == "" {
		return fmt.Errorf("spawn context water pool tag is required")
	}
	if err := context.WaterPoolTag.Validate(); err != nil {
		return err
	}
	if context.InitialDistance < 0 {
		return fmt.Errorf("spawn context initial distance must be greater than or equal to 0")
	}
	if context.InitialDepth < 0 {
		return fmt.Errorf("spawn context initial depth must be greater than or equal to 0")
	}
	for _, habitatTag := range context.HabitatTags {
		if err := habitatTag.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (appearance Appearance) Validate() error {
	if len(appearance.WaterPoolTags) == 0 {
		return fmt.Errorf("appearance water pool tags are required")
	}
	for _, waterPoolTag := range appearance.WaterPoolTags {
		if err := waterPoolTag.Validate(); err != nil {
			return err
		}
	}
	if appearance.MinInitialDistance < 0 {
		return fmt.Errorf("appearance min initial distance must be greater than or equal to 0")
	}
	if appearance.MaxInitialDistance < appearance.MinInitialDistance {
		return fmt.Errorf("appearance max initial distance must be greater than or equal to min initial distance")
	}
	if appearance.MinInitialDepth < 0 {
		return fmt.Errorf("appearance min initial depth must be greater than or equal to 0")
	}
	if appearance.MaxInitialDepth < appearance.MinInitialDepth {
		return fmt.Errorf("appearance max initial depth must be greater than or equal to min initial depth")
	}
	for _, habitatTag := range appearance.RequiredHabitatTags {
		if err := habitatTag.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (appearance Appearance) Matches(context SpawnContext) bool {
	if err := appearance.Validate(); err != nil {
		return false
	}
	if err := context.Validate(); err != nil {
		return false
	}
	if !containsWaterPool(appearance.WaterPoolTags, context.WaterPoolTag) {
		return false
	}
	if context.InitialDistance < appearance.MinInitialDistance || context.InitialDistance > appearance.MaxInitialDistance {
		return false
	}
	if context.InitialDepth < appearance.MinInitialDepth || context.InitialDepth > appearance.MaxInitialDepth {
		return false
	}
	if len(appearance.RequiredHabitatTags) > 0 && !sharesAnyHabitatTag(appearance.RequiredHabitatTags, context.HabitatTags) {
		return false
	}

	return true
}

func (appearance Appearance) MatchScore(context SpawnContext) int {
	score := 0
	score += 6 - minInt(5, appearance.MaxInitialDistance-appearance.MinInitialDistance)
	score += 6 - minInt(5, appearance.MaxInitialDepth-appearance.MinInitialDepth)
	if len(appearance.RequiredHabitatTags) > 0 && sharesAnyHabitatTag(appearance.RequiredHabitatTags, context.HabitatTags) {
		score += 10
	}

	return score
}

func ResolveSpawn(profiles []Profile, context SpawnContext) (Spawn, error) {
	return resolveSpawn(profiles, context, nil)
}

func ResolveSpawnWithRandomizer(profiles []Profile, context SpawnContext, randomizer SpawnRandomizer) (Spawn, error) {
	return resolveSpawn(profiles, context, randomizer)
}

func resolveSpawn(profiles []Profile, context SpawnContext, randomizer SpawnRandomizer) (Spawn, error) {
	if err := context.Validate(); err != nil {
		return Spawn{}, err
	}
	if len(profiles) == 0 {
		return Spawn{}, fmt.Errorf("at least one fish profile is required")
	}

	type candidate struct {
		profile Profile
		score   int
		index   int
	}

	candidates := make([]candidate, 0, len(profiles))
	for index, profile := range profiles {
		if err := profile.Validate(); err != nil {
			return Spawn{}, fmt.Errorf("profile %s: %w", profile.ID, err)
		}
		if !profile.Appearance.Matches(context) {
			continue
		}

		candidates = append(candidates, candidate{
			profile: profile,
			score:   profile.Appearance.MatchScore(context),
			index:   index,
		})
	}
	if len(candidates) == 0 {
		return Spawn{}, fmt.Errorf("no fish profile matches water %s at distance %d and depth %d with habitats [%s]", context.WaterPoolTag, context.InitialDistance, context.InitialDepth, strings.Join(habitats.Strings(context.HabitatTags), ", "))
	}

	sort.SliceStable(candidates, func(left int, right int) bool {
		if candidates[left].score == candidates[right].score {
			return candidates[left].index < candidates[right].index
		}
		return candidates[left].score > candidates[right].score
	})

	selectedCandidate := candidates[0]
	if randomizer != nil {
		topScoreCount := 1
		for topScoreCount < len(candidates) && candidates[topScoreCount].score == selectedCandidate.score {
			topScoreCount++
		}
		selectedCandidate = candidates[randomizer.Intn(topScoreCount)]
	}

	return Spawn{
		Profile:        selectedCandidate.profile,
		Context:        context,
		CandidateCount: len(candidates),
	}, nil
}

func containsWaterPool(values []waterpools.ID, target waterpools.ID) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

func sharesAnyHabitatTag(left []habitats.Tag, right []habitats.Tag) bool {
	seenTags := make(map[habitats.Tag]struct{}, len(left))
	for _, value := range left {
		seenTags[value] = struct{}{}
	}
	for _, value := range right {
		if _, exists := seenTags[value]; exists {
			return true
		}
	}

	return false
}

func minInt(left int, right int) int {
	if left < right {
		return left
	}

	return right
}
