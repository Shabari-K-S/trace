# Trace: Project-Level State Versioning (Still in Development)

**Trace** is a developer tool that versions your **local development environment** at the project level. It ensures that if a project "worked yesterday," you can return to that exact environment state today.

Trace focuses on the **Project Boundary**: it ignores your system-wide clutter and tracks only the ports, variables, and configs that belong to *this* codebase.

---

## 🔍 The Goal: Stop "Environmental Drift"

Code is versioned by Git. But the **environment** (the DB connection, the port, the `.env` file, the local dependencies) is usually invisible. Trace makes it visible.

### The "Project-Level" Pillars:

1. **Project Ports:** What ports is *this specific project* trying to use? (e.g., Is your Docker container or local Node server actually listening?)
2. **Local Context:** Changes in `.env`, `.env.local`, or `config/` files within the project root.
3. **Dependency State:** Tracking changes in `node_modules`, `venv`, or `go.mod` without committing them to Git.

---

## 🛠 Project Roadmap (Revised)

### Phase 1: The "Scope" (Project Detection)

* [ ] Logic to identify the "Project Root" (look for `.git` or `go.mod`).
* [ ] **Task:** Create a scanner that looks for `.env` files and parses keys (not values!) to track what config variables are required.

### Phase 2: The "Process Link"

* [ ] **Task:** Filter system processes to find only those running *inside* this project directory.
* [ ] Match those processes to active ports.

### Phase 3: The "Project Snap"

* [ ] Command: `trace snap "before database migration"`.
* [ ] This saves a JSON manifest of the project's "Health" at that moment.

---

## 🚀 Getting Started

### Installation (Go)

```bash
git clone https://github.com/Shabari-K-S/trace.git
cd trace
go build -o trace
```

### Usage Concept

```bash
# Initialize trace in your project folder
trace init

# Capture the current state
trace snap "working setup"

# After a break or update, check for drift
trace status
```

---

### 🔄 The Trace Lifecycle

Trace works alongside your development flow. Here is how it monitors your project's health:

1. **`trace init`**: Establish the project boundary. Trace identifies your root directory and ignores your global system noise.
2. **`trace snap`**: Captures a "clean" state. It records your active project ports, `.env` keys, and current dependency versions.
3. **`trace watch`**: A background mode that stays silent until something breaks. It notices if a port you need is suddenly taken by another process or if a config key is missing.
4. **`trace diff`**: The "Time Machine." Compare your current broken state against your last successful `snap` to see the exact environmental drift.

---

### 🎯 Key Use Cases

* **"The Morning After":** You start work on Monday, but the dev server won't boot. Run `trace diff` to see if a background update changed your environment over the weekend.
* **"The New Contributor":** A new dev joins the team. They run `trace status` to see exactly which ports and env variables they are missing compared to the "Master Snap."
* **"The Dependency Ghost":** You ran `npm install` and now everything is slow. Trace shows you which new background processes were spawned.

---

### 🛡 Privacy & Security

Trace is built with a **Security-First** mindset:

* **Zero-Value Storage:** Trace records that `API_KEY` exists, but it *never* records the actual value of the key.
* **Local Only:** No data ever leaves your machine. Your environment signatures are stored in a hidden `.trace` folder within your project.
* **Process Isolation:** It only tracks processes that originate from or interact with your project directory.

