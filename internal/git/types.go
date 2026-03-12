package git

type FileDiff struct {
	Path   string
	Status string
}

type DiffResult struct {
	Files   []FileDiff
	Stats   string
	RawDiff string
}

type CommitOptions struct {
	Date string
}

type AuthorInfo struct {
	Name  string
	Email string
}
