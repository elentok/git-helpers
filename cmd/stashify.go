package cmd

import (
	"fmt"
	"os/exec"

	"gx/git"
)

func runStashify(args []string, d deps) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: gx stashify <command> [args...]")
	}

	cwd, err := d.getwd()
	if err != nil {
		return err
	}

	changes, err := git.UncommittedChanges(cwd)
	if err != nil {
		return err
	}

	stashed := false
	if len(changes) > 0 {
		if err := runWithSpinner(d.stdin, d.stderr, "Stashing changes…", func() error {
			_, err := git.Stash(cwd)
			return err
		}); err != nil {
			return fmt.Errorf("stash failed: %w", err)
		}
		stashed = true
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = cwd
	cmd.Stdin = d.stdin
	cmd.Stdout = d.stdout
	cmd.Stderr = d.stderr
	cmdErr := cmd.Run()

	if !stashed {
		return cmdErr
	}

	if cmdErr == nil {
		return runWithSpinner(d.stdin, d.stderr, "Popping stash…", func() error {
			_, err := git.StashPop(cwd)
			return err
		})
	}

	fmt.Fprintf(d.stderr, "\nCommand failed: %v\n", cmdErr)
	confirmed, err := d.confirmForce("Pop stash anyway?")
	if err != nil {
		return err
	}
	if confirmed {
		if _, popErr := git.StashPop(cwd); popErr != nil {
			return fmt.Errorf("stash pop failed: %w", popErr)
		}
	}
	return cmdErr
}
