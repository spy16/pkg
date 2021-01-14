package shell

import (
	"os"
	"path/filepath"
)

// Init initializes a git repo in the given directory.
func Init(targetDir string, name string) error {
	target := filepath.Join(targetDir, name)
	if err := os.MkdirAll(target, os.ModePerm); err != nil {
		return err
	}
	_, err := gitExec(target, "init")
	return err
}

// Clone runs 'git clone' with the repoURL and clones the repo contents into
// targetDir. targetDir must be an empty directory.
func Clone(targetDir, repoURL string) error {
	_, err := gitExec(targetDir, "clone", repoURL, targetDir)
	return err
}

// Commit runs `git commit -m message` and returns the output.
func Commit(targetDir, message string) ([]byte, error) {
	return gitExec(targetDir, "commit", "-m", message)
}

func gitExec(dir string, subcommand string, args ...string) ([]byte, error) {
	args = append([]string{subcommand}, args...)
	return Exec(dir, "git", args...)
}
