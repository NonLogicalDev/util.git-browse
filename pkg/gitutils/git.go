package gitutils

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func GitExec(args ...string) (string, error) {
	cmd := exec.Command("git")
	cmd.Args = append([]string{"git"}, args...)

	stderr := bytes.NewBuffer(nil)
	cmd.Stderr = stderr

	out, err := cmd.Output()
	if err != nil {
		return string(out), fmt.Errorf(
			"command `git %s` failed with exit code %d and output:\n%s",
			 strings.Join(args, " "), cmd.ProcessState.ExitCode(), stderr.String(),
		)
	}

	return string(out), err
}

func GetRelativePath(path string) (string, error) {
	repoPath, err := GitExec("rev-parse", "--show-toplevel")
	if err != nil {
		return "", nil
	}

	absPath, _ := filepath.Abs(path)
	relPath, err := filepath.Rel(repoPath, absPath)
	if err != nil {
		return "", nil
	}
	if strings.Contains(relPath, "..") {
		return "", fmt.Errorf("path is not in the repo")
	}

	return filepath.Clean(relPath), err
}
