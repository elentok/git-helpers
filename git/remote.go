package git

import "strings"

// ListRemotes returns the names of all configured remotes.
func ListRemotes(repo Repo) ([]string, error) {
	out, err := run(repo.Root, []string{"remote"})
	if err != nil {
		return nil, err
	}
	var remotes []string
	for _, r := range strings.Split(out, "\n") {
		if r != "" {
			remotes = append(remotes, r)
		}
	}
	return remotes, nil
}

// UpdateRemotes fetches updates from all remotes.
func UpdateRemotes(repo Repo) error {
	_, err := run(repo.Root, []string{"remote", "update"})
	return err
}

// PruneRemote removes remote-tracking references for deleted remote branches.
func PruneRemote(repo Repo, remote string) error {
	_, err := run(repo.Root, []string{"remote", "prune", remote})
	return err
}

// PruneAllRemotes prunes all configured remotes.
func PruneAllRemotes(repo Repo) error {
	remotes, err := ListRemotes(repo)
	if err != nil {
		return err
	}
	for _, remote := range remotes {
		if err := PruneRemote(repo, remote); err != nil {
			return err
		}
	}
	return nil
}
