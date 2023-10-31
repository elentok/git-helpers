import { _parseBranchLine } from "./branch.ts"
import { assertEquals } from "https://deno.land/std@0.204.0/assert/mod.ts"
import { Branch } from "./types.ts"

Deno.test(_parseBranchLine.name, () => {
  _parseBranchLine("")

  const examples: Record<string, Branch | null> = {
    "  *main": { type: "local", name: "main", gitName: "main" },
    "  bob": { type: "local", name: "bob", gitName: "bob" },
    bob: { type: "local", name: "bob", gitName: "bob" },
    "remotes/origin/my-branch": {
      type: "remote",
      name: "my-branch",
      gitName: "origin/my-branch",
      remoteName: "origin",
    },
  }

  for (const [line, branch] of Object.entries(examples)) {
    assertEquals(_parseBranchLine(line), branch)
  }
})
