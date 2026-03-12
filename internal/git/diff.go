package git

import (
	"fmt"
	"strings"
)

func getShortStats() (string, error) {
	shortstat, err := runGit("diff", "--staged", "--shortstat")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(shortstat), nil
}

func getFilesWithStatus() ([]FileDiff, error) {
	output, err := runGit("diff", "--staged", "--find-copies", "--name-status")
	if err != nil {
		return nil, err
	}

	statusMap := map[string]string{
		"A": "NEW",
		"M": "MODIFY",
		"D": "DELETE",
		"R": "RENAME",
		"C": "COPY",
	}

	var files []FileDiff
	for line := range strings.SplitSeq(strings.TrimSpace(output), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		label, ok := statusMap[string(fields[0][0])]
		if !ok {
			label = "MODIFY"
		}
		files = append(files, FileDiff{Status: label, Path: fields[len(fields)-1]})
	}
	return files, nil
}

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

	result := &DiffResult{RawDiff: rawDiff, Stats: stats}

	result.Files, err = getFilesWithStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get file list: %w", err)
	}

	return result, nil
}
