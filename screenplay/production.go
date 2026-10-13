package screenplay

import "time"

// DefaultForceAllNarration is the default value of Production.ForceAllNarration.
const DefaultForceAllNarration = false

// Production is an explicit configuration root: build one at the top of a
// program or test suite, then derive every stage from it. It gives the
// "configure once, inherited everywhere" ergonomics of a package-level global
// without the global mutable state.
//
//	prod := NewProduction(WithNarrator(narrator))
//	stage := prod.SetTheStage(CastWhereEveryoneCan(BrowseTheWeb))
type Production struct {
	narrator          *Narrator
	timeout           time.Duration
	polling           time.Duration
	maxBufferedSize   int
	forceAllNarration bool
}

// ProductionOption configures a Production.
type ProductionOption func(*Production)

// WithNarrator sets the narrator every actor of the production narrates through.
func WithNarrator(narrator *Narrator) ProductionOption {
	return func(p *Production) {
		p.narrator = narrator
	}
}

// WithTimeout sets the timeout used by actions that wait on something.
func WithTimeout(timeout time.Duration) ProductionOption {
	return func(p *Production) {
		p.timeout = timeout
	}
}

// WithPolling sets the polling interval used by actions that wait on something.
func WithPolling(polling time.Duration) ProductionOption {
	return func(p *Production) {
		p.polling = polling
	}
}

// WithMaxBufferedSize bounds, in bytes, what one command output, one HTTP body
// or one file read may hold in memory.
func WithMaxBufferedSize(size int) ProductionOption {
	return func(p *Production) {
		p.maxBufferedSize = size
	}
}

// WithForceAllNarration forces every action and question of the production to
// narrate, neutralizing the action.Silently and question.Silently decorators.
// It is handy when debugging a scenario whose interesting step happens to be
// silenced.
func WithForceAllNarration() ProductionOption {
	return func(p *Production) {
		p.forceAllNarration = true
	}
}

// NewProduction creates a production, defaulting the timeout, the polling
// interval and the buffered size to DefaultTimeout, DefaultPolling and
// DefaultMaxBufferedSize.
func NewProduction(options ...ProductionOption) *Production {
	production := &Production{
		timeout:           DefaultTimeout,
		polling:           DefaultPolling,
		maxBufferedSize:   DefaultMaxBufferedSize,
		forceAllNarration: DefaultForceAllNarration,
	}

	for _, option := range options {
		option(production)
	}

	return production
}

// Narrator returns the narrator the production narrates through, or nil when
// none has been set.
func (p *Production) Narrator() *Narrator {
	return p.narrator
}

// Timeout returns the timeout used by actions that wait on something.
func (p *Production) Timeout() time.Duration {
	return p.timeout
}

// Polling returns the polling interval used by actions that wait on something.
func (p *Production) Polling() time.Duration {
	return p.polling
}

// MaxBufferedSize returns the bound, in bytes, on what one command output, one
// HTTP body or one file read may hold in memory.
func (p *Production) MaxBufferedSize() int {
	return p.maxBufferedSize
}

// ForceAllNarration reports whether the production forces every action and
// question to narrate, neutralizing the Silently decorators.
func (p *Production) ForceAllNarration() bool {
	return p.forceAllNarration
}

// SetTheStage builds a stage whose actors all narrate through the production's
// narrator: every actor it casts holds a pointer back to the production and
// reads its narrator from there.
func (p *Production) SetTheStage(cast Cast) *Stage {
	stage := SetTheStage(cast)
	stage.production = p

	return stage
}
