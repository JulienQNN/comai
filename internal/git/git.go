package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	gogit "github.com/go-git/go-git/v5"
	gogitconfig "github.com/go-git/go-git/v5/config"
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

// openRepo opens the git repository from the current working directory.
// DetectDotGit walks up the directory tree like git does, so it works from any subdirectory.
func openRepo() (*gogit.Repository, error) {
	repo, err := gogit.PlainOpenWithOptions(".", &gogit.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, fmt.Errorf("you are not in a git repository")
	}
	return repo, nil
}

// resolveAuthor reads user name/email from the local repo config,
// falling back to the global git config (~/.gitconfig) if not set locally.
func resolveAuthor(repo *gogit.Repository) (AuthorInfo, error) {
	local, err := repo.Config()
	if err == nil && local.User.Name != "" {
		return AuthorInfo{Name: local.User.Name, Email: local.User.Email}, nil
	}

	global, err := gogitconfig.LoadConfig(gogitconfig.GlobalScope)
	if err != nil {
		return AuthorInfo{}, fmt.Errorf("failed to read git config: %w", err)
	}
	return AuthorInfo{Name: global.User.Name, Email: global.User.Email}, nil
}

// GetAuthorInfo returns the resolved git user identity (local config → global fallback).
func GetAuthorInfo() (AuthorInfo, error) {
	repo, err := openRepo()
	if err != nil {
		return AuthorInfo{}, err
	}
	return resolveAuthor(repo)
}
