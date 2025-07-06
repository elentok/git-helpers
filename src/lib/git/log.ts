import { run } from "./index.ts"
import { Commit, Repo } from "./types.ts"

const JSON_FORMAT =
  `{"commit": "%H", "author": "%an", "email": "%ae", "date": "%ad", "message": "%s"}`

export function log(repo: string | Repo, from: string, to: string): Commit[] {
  const result = run(repo, [
    "log",
    `--pretty=format:${JSON_FORMAT}`,
    `${from}..${to}`,
  ])
  console.log("[bazinga] [log.ts] L13", result.stdout)
  return result.stdout.split("\n").map(parseLine) // git log --pretty=format:',' --date=iso | sed '$ s/,$//' | jq .
  // run
}

function parseLine(line: string): Commit | undefined {
  const trimmedLine = line.trim()
  if (trimmedLine.length === 0) {
    return
  }

  try {
    return JSON.parse(line)
  } catch {
    console.error(`Error: Failed to parse line "${line}"`)
  }
}
