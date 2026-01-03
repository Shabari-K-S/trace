package main

import (
	"fmt"
	"log"
	"os"
	"trace/internal/snapshot"
	"trace/internal/store"
	"trace/internal/diff"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: trace [snap | init]")
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		err := store.Init()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("🚀 Initialized empty Trace repository in .trace/")

	case "snap":
		// 1. Collect data
		state, _ := snapshot.Collect()
		
		// 2. Save to .trace
		err := store.SaveSnap(state)
		if err != nil {
			fmt.Println("❌ Failed to save snap:", err)
		} else {
			fmt.Println("📸 Snapshot captured and stored in .trace/snaps/")
		}

	case "diff":
		prev, curr, err := store.GetLatestTwoSnaps()
		if err != nil {
			fmt.Println("❌ Error:", err)
			return
		}
		
		added, removed := diff.Compare(prev, curr)
		fmt.Println("🔍 Comparing last two snapshots...")
		diff.RenderDiff(added, removed)

	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}