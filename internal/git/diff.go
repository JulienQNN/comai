package git

import (
	"fmt"
	"strings"

	gogit "github.com/go-git/go-git/v5"
)

func getShortStats() (string, error) {
	shortstat, err := runGit("diff", "--staged", "--shortstat")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(shortstat), nil
}

func getFilesWithStatus(repo *gogit.Repository) ([]FileDiff, error) {
	w, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := w.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	statusMap := map[gogit.StatusCode]string{
		gogit.Added:    "NEW",
		gogit.Modified: "MODIFY",
		gogit.Deleted:  "DELETE",
		gogit.Renamed:  "RENAME",
		gogit.Copied:   "MODIFY",
	}

	var files []FileDiff
	for path, s := range status {
		if s.Staging == gogit.Unmodified || s.Staging == gogit.Untracked {
			continue
		}
		label, ok := statusMap[s.Staging]
		if !ok {
			label = "MODIFY"
		}
		files = append(files, FileDiff{Status: label, Path: path})
	}

	return files, nil
}

func GetStagedDiff() (*DiffResult, error) {
	repo, err := openRepo()
	if err != nil {
		return nil, err
	}

	stats, err := getShortStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged diff with --shortstat: %w", err)
	}

	rawDiff, err := runGit("diff", "--staged")
	if err != nil {
		return nil, fmt.Errorf("failed to get staged diff: %w", err)
	}

	if strings.TrimSpace(rawDiff) == "" {
		return nil, fmt.Errorf("no staged changes found. Stage your changes with 'git add' first")
	}

	result := &DiffResult{RawDiff: rawDiff, Stats: stats}
	result.Files, err = getFilesWithStatus(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get file list: %w", err)
	}

	return result, nil
}
