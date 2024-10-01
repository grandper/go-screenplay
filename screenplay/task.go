package screenplay

// TaskWhere creates a task with the given description and steps.
func TaskWhere(description string, steps ...Performable) Performable {
	return &anonymousTask{
		description: description,
		steps:       steps,
	}
}

type anonymousTask struct {
	description string
	steps       []Performable
}

// String returns the description of the task.
func (at *anonymousTask) String() string {
	return at.description
}

// PerformAs performs the task as the provided actor.
func (at *anonymousTask) PerformAs(actor *Actor) error {
	for _, step := range at.steps {
		err := actor.AttemptsTo(step)
		if err != nil {
			return err
		}
	}

	return nil
}
