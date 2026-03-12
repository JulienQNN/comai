package git

import (
	"fmt"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	anytime "github.com/ijt/go-anytime"
)

func Commit(message string, opts CommitOptions) error {
	t := time.Now()

	if opts.Date != "" {
		parsed, err := anytime.Parse(opts.Date, time.Now(), anytime.DefaultToPast)
		if err != nil {
			parsed, err = time.Parse("2006-01-02", opts.Date)
			if err != nil {
				return fmt.Errorf("invalid date %q: %w", opts.Date, err)
			}
		}
		t = parsed
	}

	repo, err := openRepo()
	if err != nil {
		return err
	}

	author, err := resolveAuthor(repo)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	sig := &object.Signature{
		Name:  author.Name,
		Email: author.Email,
		When:  t,
	}

	_, err = w.Commit(message, &gogit.CommitOptions{
		Author:    sig,
		Committer: sig,
	})
	if err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}
	return nil
}
