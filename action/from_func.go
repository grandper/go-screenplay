package action

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/screenplay"
)

// PerformableFn represents is function that makes an actor performs an action or a task.
type PerformableFn func(theActor *screenplay.Actor) error

// FromFunc creates a new performable from a function.
func FromFunc(description string, fn PerformableFn) screenplay.Performable {
	return &funcPerformable{
		performAS:   fn,
		description: description,
	}
}

type funcPerformable struct {
	performAS   PerformableFn
	description string
}

// String describes the action.
func (p *funcPerformable) String() string {
	return p.description
}

// PerformAs performs the task or the action as the provided actor.
func (p *funcPerformable) PerformAs(actor *screenplay.Actor) error {
	return p.performAS(actor)
}

// PauseAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*funcPerformable)(nil)

// Sentinel errors returned by FromFuncAndQuestions.
var (
	ErrTaskFuncNotFunction    = errors.New("taskFunc must be a function")
	ErrNoQuestions            = errors.New("at least one question must be provided")
	ErrFirstParamNotActor     = errors.New("first parameter must be *Actor")
	ErrQuestionCountMismatch  = errors.New("number of questions does not match function parameters")
	ErrTaskFuncMustReturnOne  = errors.New("taskFunc must return one value (Performable)")
	ErrTaskFuncMustReturnTask = errors.New("taskFunc must return a Performable")
	ErrAnswerNotAssignable    = errors.New("cannot assign or convert answer to expected type")
)

// FromFuncAndQuestions creates a performable from a description format string,
// a function and questions. When String is called, the answers to the
// questions are asked to the actor that performed the action (if any) and
// passed as arguments to fmt.Sprintf with the format string.
func FromFuncAndQuestions(
	description string,
	taskFunc any,
	questions ...screenplay.Question,
) screenplay.Performable {
	return &funcWithQuestionsPerformable{
		descriptionFormat: description,
		taskFunc:          taskFunc,
		questions:         questions,
	}
}

type funcWithQuestionsPerformable struct {
	descriptionFormat string
	taskFunc          any
	questions         []screenplay.Question
	actor             *screenplay.Actor
}

// String describes the action by asking the questions to the actor that
// performed (or is performing) the action, and formatting descriptionFormat
// with their answers using fmt.Sprintf. If PerformAs has not been called yet,
// or if a question fails to answer, the format string is returned unchanged.
func (p *funcWithQuestionsPerformable) String() string {
	if p.actor == nil {
		return p.descriptionFormat
	}
	answers := make([]any, 0, len(p.questions))
	for _, q := range p.questions {
		answer, err := q.AnsweredBy(p.actor)
		if err != nil {
			return p.descriptionFormat
		}
		answers = append(answers, answer)
	}

	return fmt.Sprintf(p.descriptionFormat, answers...)
}

// PerformAs performs the task or the action as the provided actor.
func (p *funcWithQuestionsPerformable) PerformAs(actor *screenplay.Actor) error {
	p.actor = actor
	funcValue := reflect.ValueOf(p.taskFunc)
	funcType := funcValue.Type()
	if funcType.Kind() != reflect.Func {
		return ErrTaskFuncNotFunction
	}
	if len(p.questions) == 0 {
		return ErrNoQuestions
	}
	if funcType.NumIn() < 1 || funcType.In(0) != reflect.TypeOf(&screenplay.Actor{}) {
		return ErrFirstParamNotActor
	}
	if funcType.NumIn()-1 != len(p.questions) {
		return ErrQuestionCountMismatch
	}
	if funcType.NumOut() != 1 {
		return ErrTaskFuncMustReturnOne
	}
	taskType := reflect.TypeOf((*screenplay.Performable)(nil)).Elem()
	if !funcType.Out(0).Implements(taskType) {
		return ErrTaskFuncMustReturnTask
	}
	args := []reflect.Value{reflect.ValueOf(actor)}
	for i, q := range p.questions {
		answer, err := q.AnsweredBy(actor)
		if err != nil {
			return err
		}
		expectedType := funcType.In(i + 1)
		answerValue := reflect.ValueOf(answer)
		if !answerValue.Type().AssignableTo(expectedType) {
			if answerValue.Type().ConvertibleTo(expectedType) {
				answerValue = answerValue.Convert(expectedType)
			} else {
				return ErrAnswerNotAssignable
			}
		}
		args = append(args, answerValue)
	}
	results := funcValue.Call(args)
	task, ok := results[0].Interface().(screenplay.Performable)
	if !ok {
		return ErrTaskFuncMustReturnTask
	}
	return task.PerformAs(actor)
}

// funcWithQuestionsPerformable implements the screenplay.Performable interface.
var _ screenplay.Performable = (*funcWithQuestionsPerformable)(nil)
