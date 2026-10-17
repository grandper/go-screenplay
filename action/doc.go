// Package action creates an action from a function, without declaring a
// dedicated type, with FromFunc and FromFuncAndQuestions:
//
//	greet := action.FromFunc("greet the world", func(actor *screenplay.Actor) error {
//		fmt.Printf("Hello from %s\n", actor.Name())
//		return nil
//	})
//
// The actions an actor performs live in packages named after their verb, the
// first word of what the actor narrates:
//
//	proceed.To(action).IfThe(question).HappensTo(...)       action/proceed
//	choose.To(action).WhenThe(question).HappensTo(...)      action/choose
//	keep.TryingTo(action).For(5).Seconds()                  action/keep
//	wait.UntilThe(question).HappensTo(...)                  action/wait
//	remember.The(question).As(key)                          action/remember
//	log.The(question)                                       action/log
//	see.The(question).Should(...)                           action/see
package action
