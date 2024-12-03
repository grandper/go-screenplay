package ability

// SuccessStatusCode is the status code when the command returned successfully.
const SuccessStatusCode = 0

// Result is returned contains the execution result.
type Result struct {
	exitCode int
	stdOut   []byte
	stdErr   []byte
}

// ExitCode returns the exit code of the command.
func (r *Result) ExitCode() int {
	return r.exitCode
}

// StdOut returns the standard output of the command.
func (r *Result) StdOut() []byte {
	return r.stdOut
}

// StdErr returns the standard error of the command.
func (r *Result) StdErr() []byte {
	return r.stdErr
}
