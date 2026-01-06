package store

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "time"

    "trace/internal/snapshot"
)

const TraceDir = ".trace"
const SnapsDir = ".trace/snaps"

func Init() error {
    return os.MkdirAll(SnapsDir, 0755)
}

func SaveSnap(state snapshot.State) error {
    if err := os.MkdirAll(SnapsDir, 0755); err != nil {
        return err
    }

    filename := fmt.Sprintf("%d.json", time.Now().Unix())
    path := filepath.Join(SnapsDir, filename)

    data, err := json.MarshalIndent(state, "", "  ")
    if err != nil {
        return err
    }

    if err := os.WriteFile(path, data, 0644); err != nil {
        return err
    }

    // Update history log
    logPath := filepath.Join(TraceDir, "history.log")
    logEntry := fmt.Sprintf("[%s] Saved snap to %s\n", time.Now().Format(time.Kitchen), filename)

    f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err == nil {
        defer f.Close()
        _, _ = f.WriteString(logEntry)
    }

    return nil
}

func GetLatestStates() (prev *snapshot.State, curr *snapshot.State, err error) {
    files, err := os.ReadDir(SnapsDir)
    if err != nil {
        return nil, nil, err
    }
    n := len(files)
    if n == 0 {
        return nil, nil, fmt.Errorf("no snapshots found")
    }

    // helper to read a file into State
    readState := func(path string) (*snapshot.State, error) {
        data, err := os.ReadFile(path)
        if err != nil {
            return nil, err
        }
        var s snapshot.State
        if err := json.Unmarshal(data, &s); err != nil {
            return nil, err
        }
        return &s, nil
    }

    if n == 1 {
        f := filepath.Join(SnapsDir, files[0].Name())
        s, err := readState(f)
        if err != nil {
            return nil, nil, err
        }
        // prev is nil, curr is the only snapshot
        return nil, s, nil
    }

    // n >= 2
    f1 := filepath.Join(SnapsDir, files[n-2].Name())
    f2 := filepath.Join(SnapsDir, files[n-1].Name())

    s1, err := readState(f1)
    if err != nil {
        return nil, nil, err
    }
    s2, err := readState(f2)
    if err != nil {
        return nil, nil, err
    }

    return s1, s2, nil
}


// GetLatestTwoStates returns the two most recent snapshot States.
func GetLatestTwoStates() (prev, current snapshot.State, err error) {
    files, err := os.ReadDir(SnapsDir)
    if err != nil || len(files) < 2 {
        return snapshot.State{}, snapshot.State{}, fmt.Errorf("not enough snapshots to compare (need at least 2)")
    }

    f1 := filepath.Join(SnapsDir, files[len(files)-2].Name())
    f2 := filepath.Join(SnapsDir, files[len(files)-1].Name())

    prevData, err := os.ReadFile(f1)
    if err != nil {
        return snapshot.State{}, snapshot.State{}, err
    }
    currData, err := os.ReadFile(f2)
    if err != nil {
        return snapshot.State{}, snapshot.State{}, err
    }

    var s1, s2 snapshot.State
    if err := json.Unmarshal(prevData, &s1); err != nil {
        return snapshot.State{}, snapshot.State{}, err
    }
    if err := json.Unmarshal(currData, &s2); err != nil {
        return snapshot.State{}, snapshot.State{}, err
    }

    return s1, s2, nil
}
