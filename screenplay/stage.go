package screenplay

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrNoActorInTheSpotlight = errors.New("no actor is in the spotlight")
)

// Stage manages actors and tracks which one is in the spotlight.
// It is safe for concurrent use: actors can be asked for from several
// goroutines without racing.
type Stage struct {
	mutex      sync.Mutex
	cast       Cast
	production *Production
	actors     map[string]*Actor
	spotlight  *Actor
}

// SetTheStage creates a new stage with the given cast.
func SetTheStage(cast Cast) *Stage {
	return &Stage{
		cast:   cast,
		actors: map[string]*Actor{},
	}
}

// ActorNamed returns an actor by name, creating and preparing
// one through the cast if needed. The actor is placed in the spotlight.
func (s *Stage) ActorNamed(name string) *Actor {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := strings.ToLower(name)

	actor, exists := s.actors[key]
	if !exists {
		actor = ActorNamed(name)
		actor.production = s.production
		s.cast.Prepare(actor)
		s.actors[key] = actor
	}

	s.spotlight = actor

	return actor
}

// TheActorInTheSpotlight returns the actor currently in the spotlight.
func (s *Stage) TheActorInTheSpotlight() (*Actor, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.spotlight == nil {
		return nil, ErrNoActorInTheSpotlight
	}

	return s.spotlight, nil
}

// AnActorIsOnStage reports whether an actor is currently in the spotlight.
func (s *Stage) AnActorIsOnStage() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.spotlight != nil
}

// DrawTheCurtain makes all actors exit the stage.
func (s *Stage) DrawTheCurtain() error {
	// The actors are taken under the lock and exit outside of it: an exit may
	// run cleanup tasks that ask the stage for an actor.
	s.mutex.Lock()
	actors := s.actors
	s.actors = map[string]*Actor{}
	s.spotlight = nil
	s.mutex.Unlock()

	var errs []error

	for _, actor := range actors {
		if err := actor.Exit(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
