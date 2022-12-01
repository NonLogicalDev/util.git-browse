package gitutils

import (
	"bytes"
	"golang.org/x/xerrors"
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
		return string(out), xerrors.Errorf(
			"command `git %s` failed with exit code %d and output:\n%s",
			strings.Join(args, " "), cmd.ProcessState.ExitCode(), stderr.String(),
		)
	}

	return string(out), err
}

func GetPathInfo(path string) (relPath string, isDir bool, err error) {
	repoPath, err := GitExec("rev-parse", "--show-toplevel")
	if err != nil {
		return "", false, xerrors.Errorf("failed fetching git top-level dir: %w", err)
	}
	repoPath = strings.Trim(repoPath, "\n")

	absPath, _ := filepath.Abs(path)
	relPath, err = filepath.Rel(repoPath, absPath)
	if err != nil {
		return "", false, xerrors.Errorf("failed computing relative path: %w", err)
	}

	if strings.Contains(relPath, "..") {
		return "", false, xerrors.Errorf("path %#v is not in the repo", relPath)
	}

	pathStat, err := os.Stat(path)
	if err == nil {
		isDir = pathStat.IsDir()
	}

	return filepath.Clean(relPath), isDir, nil
}
