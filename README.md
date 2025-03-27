<div align="center">

# 🚀 Tasky

*A sleek, powerful CLI task management tool*

[![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg?style=for-the-badge)](https://www.gnu.org/licenses/gpl-3.0)
[![GitHub stars](https://img.shields.io/github/stars/shahriaarrr/Tasky?style=for-the-badge)](https://github.com/shahriaarrr/Tasky/stargazers)

*Manage your tasks with lightning speed, right from your terminal*

</div>

## ✨ Features

- 📋 Simple and intuitive task management
- 💾 Persistent local storage of tasks
- 🌈 Color-coded task tracking
- 🏆 Task Priority Support (Low/Medium/High)
- 🖥️ Cross-platform compatibility
- 🚦 Easy-to-use command-line interface
- 🏗️ Minimal configuration required

## 📸 Terminal Experience

<div align="center">
  <p><strong>Task Management Made Easy</strong></p>
  
  ```bash
  # Add a task with priority
  $ tasky -a "Prepare project presentation" -p High
  Boom! Task added: Prepare project presentation 🤘➕. Priority: High

  # List tasks
  $ tasky -l
  ┌───┬────────────────────────┬──────────┬───────┬─────────────┬──────────────┐
  │ # │ Tasks                  │ Priority │ State │ Created At  │ Completed At │
  ├───┼────────────────────────┼──────────┼───────┼─────────────┼──────────────┤
  │ 1 │ Prepare presentation   │   High   │  ❌   │ Mar 27 2025 │      -       │
  └───┴────────────────────────┴──────────┴───────┴─────────────┴──────────────┘
  ```
</div>

## 📦 Installation

<details>
<summary><b>Binary Release (Recommended)</b></summary>

1. Visit [Releases](https://github.com/shahriaarrr/Tasky/releases)
2. Download the binary for your operating system
3. Add to your system PATH
</details>

<details>
<summary><b>Build from Source</b></summary>

```bash
# Clone the repository
git clone https://github.com/shahriaarrr/Tasky.git

# Navigate to project directory
cd Tasky

# Build the project
go build ./cli/tasky

# Optional: Install system-wide
go install ./cli/tasky
```
</details>

## 🎮 Quick Start

### Basic Commands

| Command | Description |
|---------|-------------|
| `tasky -a "Task description" -p Priority` | Add a new task with priority |
| `tasky -e 1 "New description" -p Priority` | Edit task with priority |
| `tasky -l` | List all tasks |
| `tasky -c 1` | Complete task #1 |
| `tasky -r 1` | Remove task #1 |

### Priority Levels
- `-p Low`: Low priority tasks
- `-p Medium`: Medium priority tasks (default)
- `-p High`: High priority tasks

## 🔧 Advanced Configuration

Tasky works out of the box with sensible defaults:

```bash
# Tasks are stored in ~/.tasky.json
# No additional configuration needed!
```

## 🌟 Pro Tips

- Use short flags for quick actions (`-a`, `-l`, `-c`, `-p`)
- Long flags also work (`--add`, `--list`, `--complete`, `--priority`)
- Tasks are automatically saved after each operation
- Priorities are color-coded in the task list
- Tasks default to Medium priority if not specified

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 🐛 Issues & Feedback

Found a bug? Have a suggestion? 
[Open an issue](https://github.com/shahriaarrr/Tasky/issues) and help improve Tasky!

## 💖 Support

If you love Tasky, consider:
- ⭐ Starring the repository
- 📣 Sharing with your network
- 💡 Contributing to the project

---

<div align="center">
  <p>Crafted with ❤️ and ☕ by Shahriar Ghasempour</p>
</div>
