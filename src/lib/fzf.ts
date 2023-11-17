import { isPresent } from "./helpers.ts"

export interface FzfFlags {
  allowMultiple?: boolean
  selectOne?: boolean
  exitZero?: boolean
  prompt?: string
}

export interface FzfOptions extends FzfFlags {
  items: string[]
}

const FZF_CODE_NO_MATCH = 1
const FZF_CODE_TERMINATED_BY_USER = 130

export async function fzf({ items, selectOne, exitZero, ...rest }: FzfOptions): Promise<string[]> {
  if (items.length === 1 && selectOne) {
    return items
  }

  if (items.length === 0 && exitZero) {
    return items
  }

  const args = buildArgs(rest)
  const command = new Deno.Command("fzf", { stdin: "piped", stdout: "piped", args })
  const proc = command.spawn()

  console.log("[elentok] [fzf.ts] fzf", items)
  const encoder = new TextEncoder()
  const writer = proc.stdin.getWriter()
  writer.write(encoder.encode(`${items.join("\n")}\n`))

  const { code, success, stdout } = await proc.output()
  if (!success) {
    switch (code) {
      case FZF_CODE_NO_MATCH:
      case FZF_CODE_TERMINATED_BY_USER:
        return []
      default:
        throw new Error(`fzf exited with status ${code}`)
    }
  }

  const outputLines = new TextDecoder().decode(stdout).trim().split("\n")
  return outputLines
}

function buildArgs({ allowMultiple, prompt }: FzfFlags): string[] {
  return [allowMultiple ? "--multi" : null, prompt ? ["--prompt", prompt] : null]
    .filter(isPresent)
    .flat()
}
