// Package reflected recognizes, at run time, a value that is a
// screenplay.Question[T] for some answer type T. It backs the two entry points
// that receive questions as plain values (actor.Remember and
// action.FromFuncAndQuestions) and thus cannot know T at compile time.
package reflected

import "reflect"

// Question is a value with the methods of a screenplay.Question[T], found by
// reflection: an AnsweredBy method of the form func(actor) (T, error) and a
// String() string method.
type Question struct {
	answeredBy reflect.Value
	str        reflect.Value
}

// QuestionOf reports whether value is a question, where actorType is the type
// of the actor its AnsweredBy method must take, and returns it when it is.
func QuestionOf(value any, actorType reflect.Type) (*Question, bool) {
	receiver := reflect.ValueOf(value)
	if !receiver.IsValid() {
		return nil, false
	}

	answeredBy := receiver.MethodByName("AnsweredBy")
	if !answeredBy.IsValid() || !answers(answeredBy.Type(), actorType) {
		return nil, false
	}

	str := receiver.MethodByName("String")
	if !str.IsValid() || !describes(str.Type()) {
		return nil, false
	}

	return &Question{answeredBy: answeredBy, str: str}, true
}

// answers reports whether a bound method has the signature func(actorType) (T, error).
func answers(method reflect.Type, actorType reflect.Type) bool {
	return method.NumIn() == 1 && method.In(0) == actorType &&
		method.NumOut() == 2 && method.Out(1) == reflect.TypeFor[error]()
}

// describes reports whether a bound method has the signature func() string.
func describes(method reflect.Type) bool {
	return method.NumIn() == 0 && method.NumOut() == 1 && method.Out(0).Kind() == reflect.String
}

// AnsweredBy calls the AnsweredBy method with the given actor and returns its
// answer as an any.
func (q *Question) AnsweredBy(actor any) (any, error) {
	results := q.answeredBy.Call([]reflect.Value{reflect.ValueOf(actor)})

	if err, failed := results[1].Interface().(error); failed && err != nil {
		return nil, err
	}

	return results[0].Interface(), nil
}

// String calls the String method.
func (q *Question) String() string {
	return q.str.Call(nil)[0].String()
}
