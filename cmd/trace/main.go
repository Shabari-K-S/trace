package main

import (
	"fmt"
	"os"
	"strings"

	"trace/internal/cli"
)

const version = "2.0.0"

const helpText = `Trace - Git-like environment versioning

Usage: trace <command> [options]

Commands:
  init                Initialize trace repository in current directory
  snap <message>      Create a snapshot with the given message
  log [-n <count>]    Show commit history
  status              Show current environment drift from HEAD
  diff [commit]       Compare working environment with a commit
  restore [options]   Restore tracked files to a previous state
  checkout <ref>      Switch to a branch or commit
  branch [name]       List, create, or delete branches

Restore Options:
  --commit <hash>     Restore from specific commit (default: HEAD)
  --no-backup         Don't create backup files before restoring
  <file>...           Restore only specific files

Examples:
  trace init
  trace snap "initial environment setup"
  trace log -n 5
  trace status
  trace diff HEAD~1
  trace restore
  trace restore --commit abc123 .env
  trace branch staging
  trace checkout main

Version: %s
`

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(helpText, version)
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error

	switch command {
	case "init":
		err = cli.Init()

	case "snap":
		message := strings.Join(args, " ")
		err = cli.Snap(message)

	case "log":
		count := 0
		for i, arg := range args {
			if arg == "-n" && i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &count)
				break
			}
		}
		err = cli.Log(count)

	case "status":
		err = cli.Status()

	case "diff":
		target := ""
		if len(args) > 0 {
			target = args[0]
		}
		err = cli.Diff(target)

	case "restore":
		opts := cli.RestoreOptions{}
		var files []string
		for i := 0; i < len(args); i++ {
			switch args[i] {
			case "--commit":
				if i+1 < len(args) {
					opts.CommitRef = args[i+1]
					i++
				}
			case "--no-backup":
				opts.NoBackup = true
			default:
				files = append(files, args[i])
			}
		}
		opts.Files = files
		err = cli.Restore(opts)

	case "checkout":
		if len(args) < 1 {
			err = fmt.Errorf("usage: trace checkout <branch|commit>")
		} else {
			err = cli.Checkout(args[0])
		}

	case "branch":
		name := ""
		delete := false
		for i := 0; i < len(args); i++ {
			if args[i] == "-d" || args[i] == "--delete" {
				delete = true
			} else if name == "" {
				name = args[i]
			}
		}
		err = cli.Branch(name, delete)

	case "help", "--help", "-h":
		fmt.Printf(helpText, version)
		return

	case "version", "--version", "-v":
		fmt.Println("trace version", version)
		return

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'trace help' for usage information.")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
