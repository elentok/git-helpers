import { identifyDir } from "./identifyDir.ts"
import { assertEquals } from "std/assert/mod.ts"
import * as path from "https://deno.land/std@0.214.0/path/mod.ts"
import { describe, it } from "https://deno.land/std@0.214.0/testing/bdd.ts"
import { createDummyRepo } from "../test-helpers.ts"

describe("identifyDir", () => {
  describe("standard repo", () => {
    it("identifies root", () => {
      const dir = createDummyRepo()
      assertEquals(identifyDir(dir), {
        isBare: false,
        isRepoRoot: true,
        isWorktreeRoot: true,
        repoRoot: dir,
      })
    })

    it("identifies subdir", () => {
      const dir = createDummyRepo()
      const subdir = path.join(dir, "subdir")
      Deno.mkdirSync(subdir)

      assertEquals(identifyDir(subdir), {
        isBare: false,
        isRepoRoot: false,
        isWorktreeRoot: false,
        repoRoot: dir,
      })
    })
  })

  describe("bare repos", () => {
    it("identifies root", () => {
      const dir = createDummyRepo({ bare: true })
      assertEquals(identifyDir(dir), {
        isBare: true,
        isRepoRoot: true,
        isWorktreeRoot: false,
        repoRoot: dir,
      })
    })

    it("identifies subdir outside worktree", () => {
      const dir = createDummyRepo({ bare: true })
      const subdir = path.join(dir, "subdir")
      Deno.mkdirSync(subdir)

      assertEquals(identifyDir(subdir), {
        isBare: true,
        isRepoRoot: false,
        isWorktreeRoot: false,
        repoRoot: dir,
      })
    })

    it("identifies worktree", () => {
      const dir = createDummyRepo({ bare: true })
      const subdir = path.join(dir, "subdir")
      Deno.mkdirSync(subdir)

      assertEquals(identifyDir(subdir), {
        isBare: true,
        isRepoRoot: false,
        isWorktreeRoot: false,
        repoRoot: dir,
      })
    })
  })
})
