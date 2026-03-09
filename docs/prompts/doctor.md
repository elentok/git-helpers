# gx doctor command

Please add a "gx doctor" command that will check the repo for issues:

1. check that the origin remote refs config is defined
2. if the repo uses the ".bare" directory trick:
   1. check that its ".git" file has the correct relative path to ".bare"
   2. check that each worktree's ".git" file has the correct relative path

Also add a "--fix" option that will ask the user for confirmation on every fix
and the fix the issue.
