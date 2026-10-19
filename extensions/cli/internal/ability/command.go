package ability

// Command represents a command that is sent.
type Command struct {
	Dir     string
	Program string
	Args    []string
	// Env holds the environment variables added to the environment the child
	// inherits from the test process. A variable named like an inherited one
	// overrides it.
	Env         map[string]string
	Interactive bool
}
