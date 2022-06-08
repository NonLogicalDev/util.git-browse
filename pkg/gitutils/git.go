package gitutils

import (
	"bytes"
	"fmt"
	"os"
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

func GetPathInfo(path string) (relPath string, isDir bool, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("GetPathInfo(%#v): %w", path, err)
		}
	}()

	repoPath, err := GitExec("rev-parse", "--show-toplevel")
	if err != nil {
		return "", false, nil
	}

	absPath, _ := filepath.Abs(path)
	relPath, err = filepath.Rel(repoPath, absPath)

	fmt.Println(path)
	fmt.Println(absPath)
	fmt.Println(repoPath)
	fmt.Println(relPath)

	if err != nil {
		return "", false, err
	}
	if strings.Contains(relPath, "..") {
		return "", false, fmt.Errorf("path %#v is not in the repo", relPath)
	}

	pathStat, err := os.Stat(path)
	if err == nil {
		isDir = pathStat.IsDir()
	}

	return filepath.Clean(relPath), isDir, nil
}
