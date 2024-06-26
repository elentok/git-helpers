const decoder = new TextDecoder()

export interface ShellOptions extends Deno.CommandOptions {
  throwError?: boolean
  trim?: boolean
}

export interface ShellResult {
  code: number
  success: boolean
  stdout: string
  stderr: string
}

const DEFAULT_OPTIONS: ShellOptions = {
  trim: true,
  throwError: true,
}

export function shell(cmd: string, options?: ShellOptions): ShellResult {
  options = { ...DEFAULT_OPTIONS, ...options }
  const command = new Deno.Command(cmd, options)
  const { code, success, stdout, stderr } = command.outputSync()
  const result: ShellResult = {
    code,
    success,
    stdout: decoder.decode(stdout),
    stderr: decoder.decode(stderr),
  }
  if (options.trim) {
    result.stdout = result.stdout.trim()
    result.stderr = result.stderr.trim()
  }
  if (success || !options?.throwError) {
    return result
  }

  throw new ShellError(result, cmd, options)
}

export class ShellError extends Error {
  constructor(
    public result: ShellResult,
    public command: string,
    public options?: Deno.CommandOptions,
  ) {
    const fullCommand = [command, ...(options?.args ?? [])].join(" ")
    const output = result.stdout + "\n" + result.stderr
    super(
      `Shell command '${fullCommand}' failed with exitcode ${result.code}:\n\n${output}`,
    )
  }
}
