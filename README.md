> [!WARNING] 
> Development is currently happening on the `feature/bubbletea` branch, not on `main`.  
> If you wish to contribute, please make your changes on the `feature/bubbletea` branch and open your Pull Requests against it.

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
- 🖥️ Cross-platform compatibility
- 🚦 Easy-to-use command-line interface
- 🏗️ Minimal configuration required

## 📸 Terminal Experience

<div align="center">
  <p><strong>Task Management Made Easy</strong></p>
  
  ```bash
  # Add a task
  $ tasky -a "Prepare project presentation"
  Boom! Task added: Prepare project presentation 🤘➕

  # List tasks
  $ tasky -l
  ┌───┬────────────────────────┬───────┬─────────────┬──────────────┐
  │ # │ Tasks                  │ State │ Created At  │ Completed At │
  ├───┼────────────────────────┼───────┼─────────────┼──────────────┤
  │ 1 │ Prepare presentation   │  ❌   │ Mar 27 2025 │      -       │
  └───┴────────────────────────┴───────┴─────────────┴──────────────┘
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

# get dependency packages
go get

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
| `tasky -a "Task description"` | Add a new task |
| `tasky -l` | List all tasks |
| `tasky -c 1` | Complete task #1 |
| `tasky -e 1 "New description"` | Edit task #1 |
| `tasky -r 1` | Remove task #1 |

## 🔧 Advanced Configuration

Tasky works out of the box with sensible defaults:

```bash
# Tasks are stored in ~/.tasky.json
# No additional configuration needed!
```

## 🌟 Pro Tips

- Use short flags for quick actions (`-a`, `-l`, `-c`)
- Long flags also work (`--add`, `--list`, `--complete`)
- Tasks are automatically saved after each operation

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
