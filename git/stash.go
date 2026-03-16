package git

// Stash saves the dirty state of the working directory onto the stash stack.
func Stash(dir string) (string, error) {
	stdout, stderr, err := run(dir, []string{"stash"})
	return joinOutput(stdout, stderr), err
}

// StashPop applies the most recent stash and removes it from the stash stack.
func StashPop(dir string) (string, error) {
	stdout, stderr, err := run(dir, []string{"stash", "pop"})
	return joinOutput(stdout, stderr), err
}

// Rebase rebases the current branch onto the given ref.
func Rebase(dir string, onto string) (string, error) {
	stdout, stderr, err := run(dir, []string{"rebase", onto})
	return joinOutput(stdout, stderr), err
}
