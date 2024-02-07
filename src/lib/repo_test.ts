import { createRepo, CreateRepoOptions, identifyDirectory } from "./repo.ts"
import { assertEquals } from "std/assert/mod.ts"

function createDummyRepo(opts?: CreateRepoOptions): string {
  const dir = Deno.makeTempDirSync()
  createRepo(dir, opts)
  return dir
}

Deno.test("identifyDirectory", () => {
  const dir = createDummyRepo()
  assertEquals(identifyDirectory(dir), {
    repoType: "standard",
    isBare: false,
    isRepoRoot: true,
    isWorktreeRoot: false,
    repoRoot: dir,
    isInsideWorktree: false,
  })
})
