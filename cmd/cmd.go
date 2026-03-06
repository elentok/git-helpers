package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gx/config"
	"gx/git"
	"gx/ui/confirm"
	"gx/ui/worktrees"

	tea "github.com/charmbracelet/bubbletea"
)

type deps struct {
	stdin        io.Reader
	stdout       io.Writer
	stderr       io.Writer
	getwd        func() (string, error)
	runWorktrees func(string) error
	confirmForce func(string) (bool, error)
	initConfig   func() (string, error)
	getenv       func(string) string
	runEditor    func(editor, path string, in io.Reader, out, err io.Writer) error
}

func defaultDeps() deps {
	return deps{
		stdin:        os.Stdin,
		stdout:       os.Stdout,
		stderr:       os.Stderr,
		getwd:        os.Getwd,
		runWorktrees: runWorktrees,
		confirmForce: confirm.Run,
		initConfig:   config.Init,
		getenv:       os.Getenv,
		runEditor:    runEditorCommand,
	}
}

// Execute runs gx with the provided arguments.
func Execute(args []string) error {
	return execute(args, defaultDeps())
}

func execute(args []string, d deps) error {
	if len(args) == 0 {
		return d.runWorktrees("")
	}

	switch args[0] {
	case "worktrees", "wt":
		return d.runWorktrees("")
	case "clone-wt":
		return runCloneWT(args[1:], d)
	case "push":
		return runPush(d)
	case "init":
		return runInit(d)
	case "edit-config":
		return runEditConfig(d)
	case "-h", "--help", "help":
		printUsage(d.stdout)
		return nil
	default:
		printUsage(d.stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  gx worktrees|wt")
	fmt.Fprintln(w, "  gx clone-wt <repo-url> [directory]")
	fmt.Fprintln(w, "  gx push")
	fmt.Fprintln(w, "  gx init")
	fmt.Fprintln(w, "  gx edit-config")
}

func runWorktrees(_ string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	repo, err := git.FindRepo(cwd)
	if err != nil {
		return err
	}

	// Detect which worktree the user launched from, if any.
	var activeWorktreePath string
	if info, err := git.IdentifyDir(cwd); err == nil {
		activeWorktreePath = info.WorktreeRoot
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	settings := worktrees.Settings{
		UseNerdFontIcons: cfg.UseNerdFontIcons,
	}
	m := worktrees.NewWithSettings(*repo, activeWorktreePath, settings)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func runCloneWT(args []string, d deps) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("clone-wt expects <repo-url> [directory]")
	}

	cwd, err := d.getwd()
	if err != nil {
		return err
	}

	repoURL := args[0]
	target := ""
	if len(args) == 2 {
		target = args[1]
	}

	repoRoot, err := git.CloneBare(repoURL, target, cwd)
	if err != nil {
		return err
	}

	repo := git.Repo{Root: repoRoot, IsBare: true, MainBranch: git.RemoteDefaultBranch(repoRoot)}
	branch := repo.MainBranch
	if branch == "" {
		return fmt.Errorf("unable to determine default branch for %s", repoRoot)
	}

	wtPath := filepath.Join(repo.Root, branch)
	if err := git.AddWorktreeFromRemote(repo, wtPath, branch, "origin/"+branch); err != nil {
		return fmt.Errorf("clone succeeded but initial worktree creation failed: %w", err)
	}

	fmt.Fprintf(d.stdout, "Cloned bare repo to %s and created worktree %s\n", repoRoot, wtPath)
	return nil
}

func runPush(d deps) error {
	cwd, err := d.getwd()
	if err != nil {
		return err
	}

	info, err := git.IdentifyDir(cwd)
	if err != nil {
		return err
	}
	if info.Repo.IsBare && info.WorktreeRoot == "" {
		return fmt.Errorf("gx push must be run from a regular repo or linked worktree")
	}

	pushDir := cwd
	if info.WorktreeRoot != "" {
		pushDir = info.WorktreeRoot
	}

	branch, err := git.CurrentBranch(pushDir)
	if err != nil {
		return err
	}
	if branch == "HEAD" {
		return fmt.Errorf("cannot push from detached HEAD")
	}

	remote := git.BranchRemote(info.Repo, branch)
	pushLabel := fmt.Sprintf("Pushing %s to %s...", branch, remote)
	if err := runWithSpinner(d.stdin, d.stderr, pushLabel, func() error {
		return git.PushBranch(pushDir, remote, branch)
	}); err != nil {
		if !git.IsNonFastForwardPushError(err) {
			return err
		}

		prompt := fmt.Sprintf("Push rejected for %s/%s. Force push with lease?", remote, branch)
		confirmed, confirmErr := d.confirmForce(prompt)
		if confirmErr != nil {
			return confirmErr
		}
		if !confirmed {
			return fmt.Errorf("push aborted")
		}
		forceLabel := fmt.Sprintf("Force-pushing %s to %s with lease...", branch, remote)
		if forceErr := runWithSpinner(d.stdin, d.stderr, forceLabel, func() error {
			return git.PushBranchForceWithLease(pushDir, remote, branch)
		}); forceErr != nil {
			prompt := fmt.Sprintf("--force-with-lease failed: %v\nRun plain --force for %s/%s?", forceErr, remote, branch)
			confirmedForce, confirmErr := d.confirmForce(prompt)
			if confirmErr != nil {
				return confirmErr
			}
			if !confirmedForce {
				return fmt.Errorf("push aborted after --force-with-lease failure")
			}
			forceLabel = fmt.Sprintf("Force-pushing %s to %s...", branch, remote)
			if err := runWithSpinner(d.stdin, d.stderr, forceLabel, func() error {
				return git.PushBranchForce(pushDir, remote, branch)
			}); err != nil {
				return err
			}
			fmt.Fprintf(d.stdout, "Force-pushed %s to %s with --force\n", branch, remote)
			return nil
		}
		fmt.Fprintf(d.stdout, "Force-pushed %s to %s with --force-with-lease\n", branch, remote)
		return nil
	}

	fmt.Fprintf(d.stdout, "Pushed %s to %s\n", branch, remote)
	return nil
}

func runInit(d deps) error {
	path, err := d.initConfig()
	if err != nil {
		return err
	}
	fmt.Fprintf(d.stdout, "Created config file at %s\n", path)
	return nil
}

func runEditConfig(d deps) error {
	path, err := config.FilePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		createdPath, initErr := d.initConfig()
		if initErr != nil {
			return initErr
		}
		fmt.Fprintf(d.stdout, "Created config file at %s\n", createdPath)
	} else if err != nil {
		return err
	}

	editor := d.getenv("EDITOR")
	if strings.TrimSpace(editor) == "" {
		return fmt.Errorf("$EDITOR is not set")
	}
	return d.runEditor(editor, path, d.stdin, d.stdout, d.stderr)
}

func runEditorCommand(editor, path string, in io.Reader, out, errOut io.Writer) error {
	parts := strings.Fields(editor)
	if len(parts) == 0 {
		return fmt.Errorf("$EDITOR is empty")
	}
	args := append(parts[1:], path)
	cmd := exec.Command(parts[0], args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = errOut
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run editor %q: %w", editor, err)
	}
	return nil
}
