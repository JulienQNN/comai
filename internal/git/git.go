package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func checkGitRepo() error {
	_, err := runGit("tag")
	if err != nil {
		return fmt.Errorf("you are not in a git repository")
	}
	return nil
}
