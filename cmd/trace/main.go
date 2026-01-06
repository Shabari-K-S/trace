package main

import (
    "fmt"
    "log"
    "os"

    "trace/internal/config"
    "trace/internal/diff"
    "trace/internal/snapshot"
    "trace/internal/store"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: trace [snap | init | diff]")
        return
    }

    command := os.Args[1]

    switch command {
    case "init":
        if err := store.Init(); err != nil {
            log.Fatal(err)
        }
        if err := config.InitConfig(); err != nil {
            log.Fatal(err)
        }
        fmt.Println("Initialized Trace repo in .trace/ with default config.json")

    case "snap":
        state, err := snapshot.Collect()
        if err != nil {
            log.Fatal(err)
        }

        if err := store.SaveSnap(state); err != nil {
            fmt.Println("Failed to save snap:", err)
        } else {
            fmt.Println("📸 Snapshot captured and stored in .trace/snaps/")
        }

    case "diff":
        prev, curr, err := store.GetLatestStates()
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

		if prev == nil && curr != nil {
			fmt.Println("🔍 Only one snapshot found. Showing everything as newly added compared to an empty baseline.\n")

			addedEnv := curr.EnvKeys
			diff.RenderDiffEnv(addedEnv, nil)

			fileDiff := diff.CompareFiles(nil, curr.Files)
			diff.RenderDiffFiles(fileDiff)
			return
		}

        fmt.Println("🔍 Comparing last two snapshots...")

		fileDiff := diff.CompareFiles(prev.Files, curr.Files)
		diff.RenderDiffFiles(fileDiff)

        added, removed := diff.CompareEnv(prev.EnvKeys, curr.EnvKeys)
        diff.RenderDiffEnv(added, removed)

    default:
        fmt.Printf("Unknown command: %s\n", command)
    }
}
