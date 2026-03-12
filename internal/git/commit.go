package git

import (
	"context"
	"fmt"
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

func Commit(message string, opts CommitOptions) error {
	args := []string{"commit", "-m", message}

	if opts.Date != "" {
		t, err := ParseDate(opts.Date)
		if err != nil {
			return err
		}
		args = append(args, "--date", t.Format(time.RFC3339))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "git", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("command git commit failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}
