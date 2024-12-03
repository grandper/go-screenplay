package ability

// Command represents a command that is sent.
type Command struct {
	Dir         string
	Program     string
	Args        []string
	Env         map[string]string
	Interactive bool
}
