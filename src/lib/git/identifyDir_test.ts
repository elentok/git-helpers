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
        repo: {
          isBare: false,
          root: dir,
        },
        isRepoRoot: true,
        isWorktreeRoot: true,
      })
    })

    it("identifies subdir", () => {
      const dir = createDummyRepo()
      const subdir = path.join(dir, "subdir")
      Deno.mkdirSync(subdir)

      assertEquals(identifyDir(subdir), {
        repo: {
          isBare: false,
          root: dir,
        },
        isRepoRoot: false,
        isWorktreeRoot: false,
      })
    })
  })

  describe("bare repos", () => {
    it("identifies root", () => {
      const dir = createDummyRepo({ bare: true })
      assertEquals(identifyDir(dir), {
        repo: {
          isBare: true,
          root: dir,
        },
        isRepoRoot: true,
        isWorktreeRoot: false,
      })
    })

    it("identifies subdir outside worktree", () => {
      const dir = createDummyRepo({ bare: true })
      const subdir = path.join(dir, "subdir")
      Deno.mkdirSync(subdir)

      assertEquals(identifyDir(subdir), {
        repo: {
          root: dir,
          isBare: true,
        },
        isRepoRoot: false,
        isWorktreeRoot: false,
      })
    })

    it("identifies worktree", () => {
      const dir = createDummyRepo({ bare: true })
      const subdir = path.join(dir, "subdir")
      Deno.mkdirSync(subdir)

      assertEquals(identifyDir(subdir), {
        repo: {
          isBare: true,
          root: dir,
        },
        isRepoRoot: false,
        isWorktreeRoot: false,
      })
    })
  })
})
