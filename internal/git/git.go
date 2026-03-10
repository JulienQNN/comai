package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func runGit(args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
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
