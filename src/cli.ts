import { Command } from "npm:commander@11.1.0"
import { destroy } from "./commands/destroy.ts"

const program = new Command()
program.command("destroy").description("Destroys a local and remote branch").action(destroy)

program.parse()
