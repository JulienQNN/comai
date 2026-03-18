package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	anytime "github.com/ijt/go-anytime"
)

func ParseDate(s string) (time.Time, error) {
	before := time.Now()
	t, err := anytime.Parse(s, before, anytime.DefaultToPast)
	if err != nil || t.Equal(before.Truncate(time.Second)) {
		parsed, err2 := time.Parse("2006-01-02", s)
		if err2 != nil {
			return time.Time{}, fmt.Errorf(
				"invalid date %q: not recognized (try: yesterday, last friday, 2024-01-01)",
				s,
			)
		}
		return parsed, nil
	}
	return t, nil
}

func FormatDate(date string) (string, error) {
	t, err := ParseDate(date)
	if err != nil {
		return "", err
	}
	return t.Format(time.RFC3339), nil
}

func Commit(message string, opts CommitOptions) error {
	args := []string{"commit", "-m", message}
	var committerDateEnv []string
	
	if opts.Date != "" {
		args = append(args, "--date", opts.Date)
	}
	committerDateEnv = append(os.Environ(), "GIT_COMMITTER_DATE="+opts.Date)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", args...)

	if committerDateEnv != nil {
		cmd.Env = committerDateEnv
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command git commit failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}
