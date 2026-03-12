package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	anytime "github.com/ijt/go-anytime"
)

func Commit(message string, opts CommitOptions) error {
	args := []string{"commit", "-m", message}

	if opts.Date != "" {
		t, err := anytime.Parse(opts.Date, time.Now(), anytime.DefaultToPast)
		if err != nil {
			t, err = time.Parse("2006-01-02", opts.Date)
			if err != nil {
				return fmt.Errorf("invalid date %q: %w", opts.Date, err)
			}
		}
		args = append(args, "--date", t.Format(time.RFC3339))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "git", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}
