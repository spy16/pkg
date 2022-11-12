package shell

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// Exec executes the given command within workingDir and returns the stdout.
// If there is an error in starting the command or the command exits with non
// zero code, returns the stderr wrapped in Go error.
func Exec(workingDir, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = workingDir

	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	out, err := cmd.Output()
	if errBuf.Len() > 0 {
		errStr := errBuf.String()
		if name == "git" {
			errStr = cleanupGitErr(errStr)
		}
		err = errors.New(errStr)
	}
	return out, err
}

