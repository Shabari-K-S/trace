# Trace: Project-Level State Versioning (Phase 1 Complete)

**Trace** is a developer tool that versions your **local development environment** at the project level. It ensures that if a project "worked yesterday," you can return to that exact environment state today.

Trace focuses on the **Project Boundary**: it ignores your system-wide clutter and tracks only the ports, variables, and configs that belong to *this* codebase.

***

## 🔍 The Goal: Stop "Environmental Drift"

Code is versioned by Git. But the **environment** (the DB connection, the port, the `.env` file, the local dependencies) is usually invisible. Trace makes it visible.

### The "Project-Level" Pillars:

1. **Project Ports:** What ports is *this specific project* trying to use? (e.g., Is your Docker container or local Node server actually listening?)
2. **Local Context:** Changes in `.env`, `.env.local`, or `config/` files within the project root.
3. **Dependency State:** Tracking changes in `node_modules`, `venv`, or `go.mod` without committing them to Git.

***

## ✅ Phase 1: Config Tracking (Complete)

- [x] **Project Root:** Uses current directory (`.git`/`go.mod` detection coming soon).
- [x] **Config Scanner:** Parses `.env` keys (values never stored) + tracks any file content via `.trace/config.json`.
- [x] **`trace init`:** Creates `.trace/` with `config.json` listing files to track (`.env` by default).
- [x] **`trace snap`:** Captures env keys + file content hashes.
- [x] **`trace diff`:** Compares snapshots (handles single snapshot as "all new").

## 🛠 Phase 2: Process & Port Detection

* [ ] Filter system processes running inside project directory.
* [ ] Match processes to active ports.

## 🛠 Phase 3: Full Snapshots & Watch Mode

* [ ] `trace snap "msg"` with descriptive labels.
* [ ] `trace status` for current drift.
* [ ] `trace watch` for background monitoring.

***

## 🚀 Getting Started

### Installation

```bash
git clone https://github.com/Shabari-K-S/trace.git
cd trace
go build -o trace
```

### Usage

```bash
# Initialize in your project folder
./trace init

# Customize .trace/config.json to add files like:
# "docker-compose.yml", "config/database.yml"

# Capture current state
./trace snap

# Check for drift
./trace diff
```

**Example `trace diff` output:**
```
🔍 Only one snapshot found. Showing everything as newly added...
 + [ENV ADDED]   DATABASE_URL
 + [FILE ADDED]    .env
```

***

### 🔄 Trace Lifecycle

1. **`./trace init`**: Creates `.trace/config.json` defining tracked files.
2. **`./trace snap`**: Records env keys + file content hashes.
3. **`./trace diff`**: Compares snapshots, shows added/removed/modified env vars & files.
4. **`./trace watch`** (coming soon): Background drift detection.

***

### 🎯 Key Use Cases

- **"Morning After"**: `trace diff` shows what changed in `.env` or config files.
- **"New Contributor"**: Run `trace diff` vs master snap to see missing setup.
- **"Config Drift"**: Track changes in `docker-compose.yml` or local configs.

***

### 🛡 Privacy & Security

- **Zero-Value Storage**: Only env **keys** and file **hashes** stored, never values or full content.
- **Local Only**: Data stays in `.trace/` folder.
- **Project Scoped**: Only tracks files listed in your config.

***

**Next**: Port detection and `watch` mode. Star/follow for updates! 🚀