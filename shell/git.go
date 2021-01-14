package shell

import (
	"os"
	"path/filepath"
)

// GitInit initializes a git repo in the given directory.
func GitInit(targetDir string, name string) error {
	target := filepath.Join(targetDir, name)
	if err := os.MkdirAll(target, os.ModePerm); err != nil {
		return err
	}
	_, err := gitExec(target, "init")
	return err
}

// GitClone runs 'git clone' with the repoURL and clones the repo contents into
// targetDir. targetDir must be an empty directory.
func GitClone(targetDir, repoURL string) error {
	_, err := gitExec(targetDir, "clone", repoURL, targetDir)
	return err
}

// GitCommit runs `git commit -m message` and returns the output.
func GitCommit(targetDir, message string) ([]byte, error) {
	return gitExec(targetDir, "commit", "-m", message)
}

func gitExec(dir string, subcommand string, args ...string) ([]byte, error) {
	args = append([]string{subcommand}, args...)
	return Exec(dir, "git", args...)
}
