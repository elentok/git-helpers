When pushing a branch to GitHub for the first time, it shows the following message:

```
remote: Create a pull request for '{branch-name}' on GitHub by visiting:
remote:      https://github.com/elentok/gx/pull/new/{branch-name}
```

Update the workspace view's push command to look for this in the output,
if there's a URL, show it to the user in a modal and ask them for confirmation
to open it and if they approve open it.
