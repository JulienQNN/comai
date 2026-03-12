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
	if _, err := runGit("rev-parse", "--git-dir"); err != nil {
		return fmt.Errorf("you are not in a git repository")
	}
	return nil
}

func GetAuthorInfo() (AuthorInfo, error) {
	name, err := runGit("config", "user.name")
	if err != nil {
		return AuthorInfo{}, fmt.Errorf("failed to read git user name: %w", err)
	}
	email, err := runGit("config", "user.email")
	if err != nil {
		return AuthorInfo{}, fmt.Errorf("failed to read git user email: %w", err)
	}
	return AuthorInfo{
		Name:  strings.TrimSpace(name),
		Email: strings.TrimSpace(email),
	}, nil
}
