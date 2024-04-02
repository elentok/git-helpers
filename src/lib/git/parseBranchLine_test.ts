import { parseBranchLine } from "./parseBranchLine.ts"
import { assertEquals } from "std/assert/mod.ts"
import { Branch } from "./types.ts"

Deno.test(parseBranchLine.name, () => {
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
    assertEquals(parseBranchLine(line), branch)
  }
})
