# .bare trick

I want to support the ".bare" trick for worktrees:

Given a repository "my-repo":

1. Instead of running `git clone --bare git@...:me/my-repo` you run:

   ```
   mkdir my-repo
   cd my-repo
   git clone --bare git@...:me/myrepo .bare
   ```

2. You add a ".git" file in the root that points to the ".bare" directory:

   ```
   gitdir: ./.bare
   ```

This way the root directory is clean, only has ".bare" directory, ".git" file
and the worktrees themselves.

For "gx" to properly support this setup we need these changes:

1. Change `gx clone-wt` to use the new approach
2. Make the "n" mapping to create new worktrees works with the new approach:
   - creates the new directory under `my-repo/` and not `my-repo/.bare`
   - the `.git` file in the worktree directory should have a relative path to the correct directory:
     ```
     gitdir: ../.bare/worktrees/{worktree-name}
     ```

Definition of done: write a test that:

- Uses "gx clone-wt" to clone a repo
- Verify that you can run "gx wt" in its root directory and you see the repos
- Verify that you can run "git log" in one of the repos
