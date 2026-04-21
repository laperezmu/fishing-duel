package combat

import (
	"fmt"
	"strings"
)

type Action uint8

const (
	Forzar Action = iota
	Tensar
	Soltar
	actionCount
)

var allActions = [...]Action{Forzar, Tensar, Soltar}

func AllActions() []Action {
	return allActions[:]
}

func (a Action) String() string {
	switch a {
	case Forzar:
		return "forzar"
	case Tensar:
		return "tensar"
	case Soltar:
		return "soltar"
	default:
		return fmt.Sprintf("action(%d)", a)
	}
}

type FishFamily uint8

const (
	Embiste FishFamily = iota
	Aguante
	Quiebre
	fishFamilyCount
)

func (f FishFamily) String() string {
	switch f {
	case Embiste:
		return "embiste"
	case Aguante:
		return "aguante"
	case Quiebre:
		return "quiebre"
	default:
		return fmt.Sprintf("fish_family(%d)", f)
	}
}

type Color uint8

const (
	ColorNone Color = iota
	ColorRed
	ColorBlue
	ColorYellow
)

func (c Color) String() string {
	switch c {
	case ColorRed:
		return "red"
	case ColorBlue:
		return "blue"
	case ColorYellow:
		return "yellow"
	default:
		return "none"
	}
}

func ParseColor(raw string) (Color, error) {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case "red", "rojo":
		return ColorRed, nil
	case "blue", "azul":
		return ColorBlue, nil
	case "yellow", "amarillo":
		return ColorYellow, nil
	case "", "none", "ninguno":
		return ColorNone, nil
	default:
		return ColorNone, fmt.Errorf("unknown color %q", raw)
	}
}

type FishKind uint8

const (
	FishBlack FishKind = iota
	FishMono
	FishBiColor
)

func (k FishKind) String() string {
	switch k {
	case FishBlack:
		return "black"
	case FishMono:
		return "mono"
	case FishBiColor:
		return "bicolor"
	default:
		return fmt.Sprintf("fish_kind(%d)", k)
	}
}

type Outcome uint8

const (
	OutcomeOngoing Outcome = iota
	OutcomeCapture
	OutcomeEscape
)

func (o Outcome) String() string {
	switch o {
	case OutcomeCapture:
		return "capture"
	case OutcomeEscape:
		return "escape"
	default:
		return "ongoing"
	}
}

type TerminalReason uint8

const (
	TerminalReasonNone TerminalReason = iota
	TerminalReasonTrackCapture
	TerminalReasonTrackEscape
	TerminalReasonTotalFatigue
)

func (r TerminalReason) String() string {
	switch r {
	case TerminalReasonTrackCapture:
		return "track_capture"
	case TerminalReasonTrackEscape:
		return "track_escape"
	case TerminalReasonTotalFatigue:
		return "total_fatigue"
	default:
		return "none"
	}
}

type ColorMask uint8

func colorMaskFor(c Color) ColorMask {
	switch c {
	case ColorRed:
		return 1 << 0
	case ColorBlue:
		return 1 << 1
	case ColorYellow:
		return 1 << 2
	default:
		return 0
	}
}

type RodModifiers struct {
	Offensive ColorMask
}

func (r RodModifiers) HasOffensive(c Color) bool {
	return r.Offensive&colorMaskFor(c) != 0
}

func (r RodModifiers) WithOffensive(c Color) RodModifiers {
	r.Offensive |= colorMaskFor(c)
	return r
}

type FishProfile struct {
	Kind   FishKind
	Colors [2]Color
}

func NewBlackFish() FishProfile {
	return FishProfile{Kind: FishBlack}
}

func NewMonoFish(color Color) FishProfile {
	return FishProfile{Kind: FishMono, Colors: [2]Color{color, ColorNone}}
}

func NewBiColorFish(a, b Color) FishProfile {
	if a == ColorNone || b == ColorNone || a == b {
		panic("bicolor fish requires two different non-none colors")
	}
	return FishProfile{Kind: FishBiColor, Colors: normalizeBiColors(a, b)}
}

func normalizeBiColors(a, b Color) [2]Color {
	switch {
	case (a == ColorRed && b == ColorYellow) || (a == ColorYellow && b == ColorRed):
		return [2]Color{ColorRed, ColorYellow}
	case (a == ColorBlue && b == ColorRed) || (a == ColorRed && b == ColorBlue):
		return [2]Color{ColorBlue, ColorRed}
	case (a == ColorBlue && b == ColorYellow) || (a == ColorYellow && b == ColorBlue):
		return [2]Color{ColorBlue, ColorYellow}
	default:
		if a < b {
			return [2]Color{a, b}
		}
		return [2]Color{b, a}
	}
}

func (f FishProfile) HasColor(c Color) bool {
	return f.Colors[0] == c || f.Colors[1] == c
}

func (f FishProfile) String() string {
	switch f.Kind {
	case FishBlack:
		return "black"
	case FishMono:
		return f.Colors[0].String()
	case FishBiColor:
		if f.Colors[0] == ColorBlue && f.Colors[1] == ColorRed {
			return "blue-red"
		}
		return f.Colors[0].String() + "-" + f.Colors[1].String()
	default:
		return "unknown"
	}
}

func ParseFishProfile(raw string) (FishProfile, error) {
	clean := strings.TrimSpace(strings.ToLower(raw))
	switch clean {
	case "black", "negro":
		return NewBlackFish(), nil
	case "red", "rojo":
		return NewMonoFish(ColorRed), nil
	case "blue", "azul":
		return NewMonoFish(ColorBlue), nil
	case "yellow", "amarillo":
		return NewMonoFish(ColorYellow), nil
	case "red-yellow", "yellow-red", "rojo-amarillo", "amarillo-rojo":
		return NewBiColorFish(ColorRed, ColorYellow), nil
	case "blue-red", "red-blue", "azul-rojo", "rojo-azul":
		return NewBiColorFish(ColorBlue, ColorRed), nil
	case "blue-yellow", "yellow-blue", "azul-amarillo", "amarillo-azul":
		return NewBiColorFish(ColorBlue, ColorYellow), nil
	default:
		return FishProfile{}, fmt.Errorf("unknown fish profile %q", raw)
	}
}

type FamilyCounts [fishFamilyCount]uint8

func NewFamilyCounts(embiste, aguante, quiebre uint8) FamilyCounts {
	return FamilyCounts{embiste, aguante, quiebre}
}

func (c FamilyCounts) Count(f FishFamily) uint8 {
	return c[f]
}

func (c FamilyCounts) Total() int {
	return int(c[Embiste] + c[Aguante] + c[Quiebre])
}

func (c *FamilyCounts) Increment(f FishFamily) {
	c[f]++
}

func (c *FamilyCounts) Decrement(f FishFamily) {
	if c[f] == 0 {
		panic("cannot decrement empty family count")
	}
	c[f]--
}

type BeliefState struct {
	TrackPos           uint8
	Fish               FishProfile
	Draw               FamilyCounts
	Discard            FamilyCounts
	FatigueCount       uint8
	BaitGuardAvailable bool
	RodMods            RodModifiers
}

type CombatConfig struct {
	InitialTrackPos    int
	Fish               FishProfile
	BaitGuardAvailable bool
	RodMods            RodModifiers
}

type CombatResult struct {
	Outcome              Outcome
	TerminalReason       TerminalReason
	Rounds               int
	Fatigues             int
	ActionUsage          [actionCount]int
	BaitSaves            int
	OffensiveModTriggers int
	InitialTrackPos      int
	FinalTrackPos        int
}

type ActionValues struct {
	Values     [actionCount]float64
	BestAction Action
	BestValue  float64
}

type weightedBeliefState struct {
	State       BeliefState
	Probability float64
}
