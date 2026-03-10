package git

import (
	"fmt"
	"strings"
)

// get short stats for printing it
func getShortStats() (string, error) {
	shortstat, err := runGit("diff", "--staged", "--shortstat")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(shortstat), nil
}

// get reliable file listing for better visibility
func getFilesWithStatus() ([]FileDiff, error) {
	output, err := runGit("diff", "--staged", "--name-status")
	if err != nil {
		return nil, err
	}

	var files []FileDiff
	statusMap := map[string]string{
		"A": "NEW",
		"M": "MODIFY",
		"D": "DELETE",
		"R": "RENAME",
	}

	for line := range strings.SplitSeq(strings.TrimSpace(output), "\n") {
		// Avoid empty lines and malformed output
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		status := statusMap[string(fields[0][0])]
		if status == "" {
			status = "MODIFY"
		}

		files = append(files, FileDiff{
			Status: status,
			Path:   fields[len(fields)-1],
		})
	}

	return files, nil
}

// get the staged diff
func GetStagedDiff() (*DiffResult, error) {
	if err := checkGitRepo(); err != nil {
		return nil, err
	}
	stats, err := getShortStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged diff: %w", err)
	}

	rawDiff, err := runGit("diff", "--staged")
	if err != nil {
		return nil, fmt.Errorf("failed to get staged diff: %w", err)
	}

	if strings.TrimSpace(rawDiff) == "" {
		return nil, fmt.Errorf("no staged changes found. Stage your changes with 'git add' first")
	}

	result := &DiffResult{
		RawDiff: rawDiff,
		Stats:   stats,
	}

	// Get file list with each status
	result.Files, err = getFilesWithStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get file list: %w", err)
	}

	return result, nil
}
