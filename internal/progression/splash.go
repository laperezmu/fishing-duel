package progression

type SplashEscapeDeciderFunc func(chance float64) bool

func (decider SplashEscapeDeciderFunc) ShouldEscape(chance float64) bool {
	return decider(chance)
}
