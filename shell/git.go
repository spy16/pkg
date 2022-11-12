package shell

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitExec executes a git subcommand with given workingDir and returns the
// stdout of the git command. If there is an error, returns contents of the
// stderr wrapped in Go error.
func GitExec(workingDir, subcommand string, args ...string) ([]byte, error) {
	args = append([]string{subcommand}, args...)

	cmd := exec.Command("git", args...)
	cmd.Dir = workingDir

	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	out, err := cmd.Output()
	if errBuf.Len() > 0 {
		err = errors.New(cleanupGitErr(errBuf.String()))
	}
	return out, err
}

// GitInit initializes a git repo in the given directory.
func GitInit(targetDir string, name string) error {
	target := filepath.Join(targetDir, name)
	if err := os.MkdirAll(target, os.ModePerm); err != nil {
		return err
	}
	_, err := GitExec(target, "init")
	return err
}

// GitClone runs 'git clone' with the repoURL and clones the repo contents into
// targetDir. targetDir must be an empty directory.
func GitClone(targetDir, repoURL string) error {
	_, err := GitExec(targetDir, "clone", repoURL, targetDir)
	return err
}

// GitCommit runs `git commit -m message` and returns the output.
func GitCommit(targetDir, message string) ([]byte, error) {
	return GitExec(targetDir, "commit", "-m", message)
}

func cleanupGitErr(s string) string {
	return strings.Replace(s, "fatal: ", "", -1)
}
