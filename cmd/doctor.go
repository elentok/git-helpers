package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gx/git"
)

func runDoctor(args []string, d deps) error {
	fix := len(args) > 0 && args[0] == "--fix"

	cwd, err := d.getwd()
	if err != nil {
		return err
	}

	repo, err := git.FindRepo(cwd)
	if err != nil {
		// FindRepo can fail when the outer .git file is itself broken.
		// Fall back to checking whether cwd contains a .bare directory.
		repo, err = findRepoWithFallback(cwd)
		if err != nil {
			return err
		}
	}

	issues, err := git.CheckRepo(*repo)
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		fmt.Fprintln(d.stdout, "No issues found.")
		return nil
	}

	for i, issue := range issues {
		fmt.Fprintf(d.stdout, "[%d/%d] %s\n", i+1, len(issues), issue.Description)

		if !issue.CanFix() {
			fmt.Fprintln(d.stdout, "  No automatic fix available.")
			fmt.Fprintln(d.stdout)
			continue
		}

		if !fix {
			fmt.Fprintf(d.stdout, "  Fix: %s\n", issue.FixDescription)
			fmt.Fprintln(d.stdout)
			continue
		}

		confirmed, err := d.confirmForce(issue.FixDescription + "?")
		if err != nil {
			return err
		}
		if confirmed {
			if err := issue.Fix(); err != nil {
				fmt.Fprintf(d.stderr, "  error: %v\n", err)
			} else {
				fmt.Fprintln(d.stdout, "  Fixed.")
			}
		} else {
			fmt.Fprintln(d.stdout, "  Skipped.")
		}
		fmt.Fprintln(d.stdout)
	}

	if !fix {
		fmt.Fprintln(d.stdout, "Run 'gx doctor --fix' to apply fixes.")
	}

	return nil
}

// findRepoWithFallback tries FindRepo first, then checks for a .bare directory
// in dir (used when the outer .git file is itself corrupted).
func findRepoWithFallback(dir string) (*git.Repo, error) {
	bareDir := filepath.Join(dir, ".bare")
	info, err := os.Stat(bareDir)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("no git repo found at %q", dir)
	}
	repo := &git.Repo{
		Root:        bareDir,
		WorktreeDir: dir,
		IsBare:      true,
		MainBranch:  git.RemoteDefaultBranch(bareDir),
	}
	return repo, nil
}

func printDoctorUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: gx doctor [--fix]")
	fmt.Fprintln(w, "  Checks the current repo for common configuration issues.")
	fmt.Fprintln(w, "  --fix  Prompt to apply each fix interactively.")
}
