package prompt

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/JulienQNN/comai/internal/config"
)

const maxTotalBytes = 6000

// catching low-value files, lockFilePatterns matches auto-generated files by naming convention
var lockFilePatterns = []string{"*.lock", "*.sum", "*-lock.json", "*-lock.yaml"}

func isLockFile(path string) bool {
	base := filepath.Base(path)
	for _, pattern := range lockFilePatterns {
		if matched, _ := filepath.Match(pattern, base); matched {
			return true
		}
	}
	return false
}

// isMetadataLine reports whether a diff line carries no semantic value for the LLM.
func isMetadataLine(line string) bool {
	return strings.HasPrefix(line, "diff --git ") ||
		strings.HasPrefix(line, "index ") ||
		strings.HasPrefix(line, "--- ") ||
		strings.HasPrefix(line, "+++ ") ||
		strings.HasPrefix(line, "\\ No newline")
}

func splitDiffBlocks(raw string) []string {
	if raw == "" {
		return nil
	}
	blocks := strings.Split(raw, "\ndiff --git ")

	for i := 1; i < len(blocks); i++ {
		blocks[i] = "diff --git " + blocks[i]
	}

	return blocks
}

func extractFilePath(block string) string {
	line, _, _ := strings.Cut(block, "\n")
	idx := strings.LastIndex(line, " b/")
	if idx != -1 {
		path := line[idx+3:]
		return strings.Trim(path, `"`)
	}
	return ""
}

func processBlock(block string) string {
	path := extractFilePath(block)

	if strings.Contains(block, "Binary files ") {
		return fmt.Sprintf("[%s: binary file changed]\n", path)
	}

	if isLockFile(path) {
		var added, removed int
		for line := range strings.SplitSeq(block, "\n") {
			switch {
			case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
				added++
			case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
				removed++
			}
		}
		return fmt.Sprintf("[%s: auto-generated file, +%d/-%d lines]\n", path, added, removed)
	}

	var sb strings.Builder
	for line := range strings.SplitSeq(block, "\n") {
		if !isMetadataLine(line) {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// filterDiff processes a raw git diff to remove noise and cap its size.
func filterDiff(raw string) string {
	var sb strings.Builder
	for _, block := range splitDiffBlocks(raw) {
		if sb.Len() >= maxTotalBytes {
			sb.WriteString("...(diff truncated)\n")
			break
		}
		sb.WriteString(processBlock(block))
	}
	return sb.String()
}

func Build(diff string, cfg config.Config) CompletionParams {
	system := fmt.Sprintf(
		"Output ONLY a git commit message in lowercase MaxLength: %v Language: %s without ANY formatting backticks or codeblocks following conventional commit messages <type>(<optional scope>): <description>",
		cfg.MaxLength,
		cfg.Language,
	)

	if cfg.CustomInstructions != "" {
		system = cfg.CustomInstructions
	}

	return CompletionParams{
		SystemPrompt: system,
		UserPrompt:   "Diff:\n" + filterDiff(diff),
		MaxTokens:    128,
	}
}
