package snapshot

type PortInfo struct {
    PID  int    `json:"pid"`
    Port int    `json:"port"`
    Proc string `json:"proc"`
}

type DependencyInfo struct {
    Type   string `json:"type"`
    Finger string `json:"finger"`
    Detail string `json:"detail"`
}

type FileSnapshot struct {
    Path string `json:"path"`
    Hash string `json:"hash"`
    // Content string `json:"content,omitempty"` // optional later
}

type State struct {
    Timestamp    string           `json:"timestamp"`
    EnvKeys      []string         `json:"env_keys"`
    Ports        []PortInfo       `json:"ports,omitempty"`
    Dependencies []DependencyInfo `json:"dependencies,omitempty"`
    Files        []FileSnapshot   `json:"files,omitempty"`
}
