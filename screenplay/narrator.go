package screenplay

import "sync"

// Kind classifies a narrated line, mirroring ScreenPy's act/scene/beat/aside.
type Kind int

const (
	KindAct   Kind = iota // suite-level grouping
	KindScene             // smaller grouping
	KindBeat              // a single step: an action performed or a question asked
	KindAside             // an ad-hoc comment injected into the narrative
)

// String describes the kind.
func (k Kind) String() string {
	switch k {
	case KindAct:
		return "act"
	case KindScene:
		return "scene"
	case KindBeat:
		return "beat"
	case KindAside:
		return "aside"
	default:
		return "unknown"
	}
}

// Level is the "gravitas" of a line: how important it is. Adapters typically
// map it onto their logging library's severity levels.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo        // the default for beats
	LevelWarn
	LevelError // used automatically when a beat fails
)

// Phase distinguishes the moment a step is announced from the moment it
// finishes. Simple adapters may ignore PhaseEnd unless Event.Err is set.
type Phase int

const (
	PhaseBegin Phase = iota // the step is about to run (or a one-shot aside/act/scene)
	PhaseEnd                // the step finished; Event.Err reports its outcome
)

// Event is a single unit of narration handed to every attached Adapter.
type Event struct {
	Kind    Kind
	Level   Level
	Phase   Phase
	Depth   int    // nesting depth, for indentation
	Actor   string // name of the actor in the spotlight, if any
	Message string // human-readable description (usually a Performable/Question String())
	Answer  any    // for question beats: the answer that was given
	Err     error  // non-nil when a beat failed
}

// Adapter is the customization seam. Implement it on top of any logging
// library to decide where and how the narrative is rendered.
//
// The library ships two reference implementations: one over the standard "log"
// package (narrator/logadapter) and one over "log/slog" (narrator/slogadapter).
type Adapter interface {
	// Narrate renders a single narration event.
	Narrate(event Event)
}

// Narrator receives narration from actors and fans it out to its adapters.
// The zero value is a valid, silent narrator (no adapters attached).
//
// A Narrator is safe for concurrent use: the actions that run performables from
// several goroutines (action.Concurrently, action.Asynchronously) all narrate
// through the same microphone.
type Narrator struct {
	mutex    sync.Mutex
	adapters []Adapter
	depth    int
}

// NewNarrator creates a Narrator with the given adapters already attached.
func NewNarrator(adapters ...Adapter) *Narrator {
	return &Narrator{adapters: adapters}
}

// AttachAdapter registers one or more adapters. Every event is delivered to
// every attached adapter, in registration order.
func (n *Narrator) AttachAdapter(adapters ...Adapter) *Narrator {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.adapters = append(n.adapters, adapters...)

	return n
}

// active reports whether the narrator has at least one adapter attached.
// A silent narrator (the default) skips narration entirely so that, with no
// adapter attached, the library behaves exactly as it does without narration.
func (n *Narrator) active() bool {
	if n == nil {
		return false
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	return len(n.adapters) > 0
}

// emit is the single choke point that reaches the adapters.
func (n *Narrator) emit(event Event) {
	if n == nil {
		return
	}

	n.mutex.Lock()
	event.Depth = n.depth
	adapters := n.adapters
	n.mutex.Unlock()

	for _, adapter := range adapters {
		adapter.Narrate(event)
	}
}

// WhispersTheAside injects an ad-hoc, one-shot comment at the current depth
// (@aside).
func (n *Narrator) WhispersTheAside(actor, message string) {
	n.emit(Event{
		Kind:    KindAside,
		Level:   LevelInfo,
		Phase:   PhaseBegin,
		Actor:   actor,
		Message: message,
	})
}

// AnnouncesTheAct announces a suite-level grouping boundary (@act).
func (n *Narrator) AnnouncesTheAct(name string) {
	n.emit(Event{Kind: KindAct, Level: LevelInfo, Phase: PhaseBegin, Message: name})
}

// SetsTheScene announces a smaller grouping boundary (@scene).
func (n *Narrator) SetsTheScene(name string) {
	n.emit(Event{Kind: KindScene, Level: LevelInfo, Phase: PhaseBegin, Message: name})
}

// StatesTheFact announces a step (@beat), runs it, then reports its outcome.
// Nested steps (a task made of steps) automatically increase Depth, producing
// the indented, hierarchical output ScreenPy is known for.
func (n *Narrator) StatesTheFact(actor, message string, do func() error) error {
	if !n.active() {
		return do()
	}

	n.emit(Event{
		Kind:    KindBeat,
		Level:   LevelInfo,
		Phase:   PhaseBegin,
		Actor:   actor,
		Message: message,
	})

	n.mutex.Lock()
	n.depth++
	n.mutex.Unlock()

	err := do()

	n.mutex.Lock()
	n.depth--
	n.mutex.Unlock()

	level := LevelInfo
	if err != nil {
		level = LevelError
	}

	n.emit(Event{
		Kind:    KindBeat,
		Level:   level,
		Phase:   PhaseEnd,
		Actor:   actor,
		Message: message,
		Err:     err,
	})

	return err
}

// RevealsTheAnswer narrates the answer a question was given.
func (n *Narrator) RevealsTheAnswer(actor, message string, answer any, err error) {
	level := LevelInfo
	if err != nil {
		level = LevelError
	}

	n.emit(Event{
		Kind:    KindBeat,
		Level:   level,
		Phase:   PhaseEnd,
		Actor:   actor,
		Message: message,
		Answer:  answer,
		Err:     err,
	})
}
