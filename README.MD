# chasm — Character Assembly CLI  
> my personal, interactive command-line string / command composer.

Written in pure Go and delivered as a single binary that works out-of-the-box on **macOS** and **Windows**.  
Core mission: **assemble the most complex commands with the fewest keystrokes**.
[中文](README_CN.MD)
---

## ✨ Key Features

1. **Smart Auto-Completion**  
   Logic lives in `components/insidewords.go`; swap in an external `insidewords.toml` for **hot-reload without recompiling**.

2. **Built-in File Navigator**  
   Mini file browser with fuzzy search, multi-select, and `Tab` path-completion.

3. **Macro Command Vault**  
   Store once, reuse forever—macros are fully customizable.

4. **Pipeline Friendly**  
   Pipe straight to any shell: `chasm ... | sh`.

5. **Clipboard Auto-Copy**  
   Generated string is **instantly copied** to the system clipboard—`Ctrl+V` to paste & run.

6. **Zero-Dependency & Cross-Platform**  
   Single static binary built with Go 1.20+. Runs on macOS & Windows with no extras.

---

## 🚀 Installation

```bash
# Run source directly
go run main.go

# Build locally
go build -o chasm

# Install globally
go install

🎮 Interactive Mode Cheat-Sheet

| Key          | Mode         | Purpose                              |
| ------------ | ------------ | ------------------------------------ |
| `Tab` / `↑↓` | **CMDMOD**   | Default command input + completion   |
| `Ctrl+P`     | **BATMODE**  | Quick-insert / search macro commands |
| `Ctrl+O`     | **FILEMODE** | Pick files/folders (multi-select)    |
| `Ctrl+N`     | **ARGMODE**  | Select from upstream pipe list       |


Universal Shortcuts

Shift+←/→: switch categories (BATMODE / FILEMODE)
Enter: confirm & return to CMDMOD
Real-time fuzzy search box always available


📌 Typical Workflow

Goal: run seed download <filepath>
Launch chasm → CMDMOD.
Type seed → Tab to accept completion.
At <filepath> hit Ctrl+O → FILEMODE to choose file.
Tab to select → Enter inject path → Enter exit.
Final command is now in clipboard; paste with Ctrl+V or run directly:

chasm seed download /path/to/file | sh
⚙️ Customization
| Method       | File                        | Notes                               |
| ------------ | --------------------------- | ----------------------------------- |
| **Embedded** | `components/insidewords.go` | edit source & rebuild               |
| **External** | `insidewords.toml`          | drop-in replacement, zero recompile |

🛠️ Tech Stack

Language: Go 1.20+
Deps: stdlib only + tiny clipboard helper
Terminal: ANSI & Windows Console compatible