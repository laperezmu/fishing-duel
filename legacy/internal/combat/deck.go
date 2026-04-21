package combat

import "math/rand"

type FatigueResult struct {
	Triggered    bool
	TotalFatigue bool
	Removed      FamilyCounts
}

type FishDeck struct {
	draw    []FishFamily
	discard []FishFamily
}

func NewShuffledFishDeck(rng *rand.Rand) FishDeck {
	cards := []FishFamily{
		Embiste, Embiste, Embiste,
		Aguante, Aguante, Aguante,
		Quiebre, Quiebre, Quiebre,
	}
	rng.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return FishDeck{draw: cards}
}

func NewFishDeckFromDraw(draw []FishFamily) FishDeck {
	cloned := make([]FishFamily, len(draw))
	copy(cloned, draw)
	return FishDeck{draw: cloned}
}

func (d FishDeck) Clone() FishDeck {
	clonedDraw := make([]FishFamily, len(d.draw))
	copy(clonedDraw, d.draw)
	clonedDiscard := make([]FishFamily, len(d.discard))
	copy(clonedDiscard, d.discard)
	return FishDeck{draw: clonedDraw, discard: clonedDiscard}
}

func (d *FishDeck) Draw() (FishFamily, bool) {
	if len(d.draw) == 0 {
		return 0, false
	}
	card := d.draw[0]Lo
	d.draw = d.draw[1:]
	return card, true
}

func (d *FishDeck) DiscardCard(card FishFamily) {
	d.discard = append(d.discard, card)
}

func (d *FishDeck) ReshuffleForFatigue(rng *rand.Rand) FatigueResult {
	result := FatigueResult{Triggered: true}
	if len(d.discard) == 0 {
		result.TotalFatigue = true
		return result
	}

	d.draw = append(d.draw[:0], d.discard...)
	d.discard = d.discard[:0]
	rng.Shuffle(len(d.draw), func(i, j int) {
		d.draw[i], d.draw[j] = d.draw[j], d.draw[i]
	})

	if len(d.draw) <= 3 {
		for _, card := range d.draw {
			result.Removed.Increment(card)
		}
		d.draw = d.draw[:0]
		result.TotalFatigue = true
		return result
	}

	removed := d.draw[:3]
	for _, card := range removed {
		result.Removed.Increment(card)
	}
	d.draw = append([]FishFamily(nil), d.draw[3:]...)
	return result
}

func (d FishDeck) DrawCounts() FamilyCounts {
	var counts FamilyCounts
	for _, card := range d.draw {
		counts.Increment(card)
	}
	return counts
}

func (d FishDeck) DiscardCounts() FamilyCounts {
	var counts FamilyCounts
	for _, card := range d.discard {
		counts.Increment(card)
	}
	return counts
}

func (d FishDeck) DrawLen() int {
	return len(d.draw)
}

func (d FishDeck) DiscardLen() int {
	return len(d.discard)
}
