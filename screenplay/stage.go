package screenplay

import (
	"errors"
	"strings"
)

var (
	ErrNoActorInTheSpotlight = errors.New("no actor is in the spotlight")
)

// Stage manages actors and tracks which one is in the spotlight.
type Stage struct {
	cast      Cast
	actors    map[string]*Actor
	spotlight *Actor
}

// SetTheStage creates a new stage with the given cast.
func SetTheStage(cast Cast) *Stage {
	return &Stage{
		cast:   cast,
		actors: map[string]*Actor{},
	}
}

// TheActorCalled returns an actor by name, creating and preparing
// one through the cast if needed. The actor is placed in the spotlight.
func (s *Stage) TheActorCalled(name string) *Actor {
	key := strings.ToLower(name)

	actor, exists := s.actors[key]
	if !exists {
		actor = ActorNamed(name)
		s.cast.Prepare(actor)
		s.actors[key] = actor
	}

	s.spotlight = actor

	return actor
}

// TheActorInTheSpotlight returns the actor currently in the spotlight.
func (s *Stage) TheActorInTheSpotlight() (*Actor, error) {
	if s.spotlight == nil {
		return nil, ErrNoActorInTheSpotlight
	}

	return s.spotlight, nil
}

// AnActorIsOnStage reports whether an actor is currently in the spotlight.
func (s *Stage) AnActorIsOnStage() bool {
	return s.spotlight != nil
}

// DrawTheCurtain makes all actors exit the stage.
func (s *Stage) DrawTheCurtain() error {
	var errs []error

	for _, actor := range s.actors {
		if err := actor.Exit(); err != nil {
			errs = append(errs, err)
		}
	}

	s.actors = map[string]*Actor{}
	s.spotlight = nil

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
