package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github/tasky"
)

const taskFile = ".tasky.json"

type menuChoice int

type menuItem struct {
	title string
	desc  string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

const (
	menuAdd menuChoice = iota
	menuComplete
	menuEdit
	menuDelete
	menuList
	menuExit
)

type state int

const (
	stateMenu state = iota
	stateAdd
	stateComplete
	stateEdit
	stateDelete
	stateList
)

type model struct {
	todos      *tasky.Todos
	filename   string
	state      state
	menu       list.Model
	textinput  textinput.Model
	indexinput textinput.Model
	err        error

	editIndex int

	listOutput string

	completeItems  []completeItem
	completeCursor int
}

type completeItem struct {
	title    string
	index    int
	selected bool
}

func main() {
	items := []list.Item{
		menuItem{"Add Task", "Add a new task"},
		menuItem{"Complete Task", "Mark a task as done"},
		menuItem{"Edit Task", "Edit an existing task"},
		menuItem{"Remove Task", "Delete a task"},
		menuItem{"List Tasks", "Show all tasks"},
		menuItem{"Exit", "Quit the app"},
	}

	menuList := list.New(items, list.NewDefaultDelegate(), 40, 14)
	menuList.Title = "Tasky Menu"

	txtInput := textinput.New()
	txtInput.Placeholder = "Enter task name"

	idxInput := textinput.New()
	idxInput.Placeholder = "Enter task index (e.g. 1)"

	t := &tasky.Todos{}
	if err := t.Load(taskFile); err != nil {

		fmt.Fprintln(os.Stderr, "warning: failed to load tasks:", err)
	}

	taskModel := model{
		todos:      t,
		filename:   taskFile,
		state:      stateMenu,
		menu:       menuList,
		textinput:  txtInput,
		indexinput: idxInput,
		err:        nil,
		editIndex:  0,
		listOutput: "",
	}

	if _, err := tea.NewProgram(taskModel, tea.WithMouseAllMotion()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateMenu:
		return m.updateMenu(msg)
	case stateAdd:
		return m.updateAdd(msg)
	case stateComplete:
		return m.updateComplete(msg)
	case stateEdit:
		return m.updateEdit(msg)
	case stateDelete:
		return m.updateDelete(msg)
	case stateList:
		return m.updateList(msg)
	}
	return m, nil
}

func (m model) View() string {

	errLine := ""
	if m.err != nil {
		errLine = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("Error: "+m.err.Error()) + "\n\n"
	}

	switch m.state {
	case stateMenu:
		return errLine + lipgloss.NewStyle().Margin(1, 2).Render(m.menu.View())
	case stateAdd:
		return errLine + "Add Task (ENTER to save, ESC to cancel):\n\n" + m.textinput.View()
	case stateComplete:
		var b strings.Builder
		header := "Select tasks to complete (Space/Enter to toggle). Use ↑/↓, j/k or mouse wheel.\n"
		header += "We'll only return when you select 'main menu'.\n\n"
		b.WriteString(header)

		green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
		dim := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

		for i, it := range m.completeItems {
			cursor := " "
			if m.completeCursor == i {
				cursor = ">"
			}
			check := "[ ]"
			line := it.title
			if it.index > 0 && it.selected {
				check = "[x]"
				line = green.Render(line)
			}
			if it.index == 0 {
				line = dim.Render(line)
			}
			b.WriteString(fmt.Sprintf("%s %s %s\n", cursor, check, line))
		}
		return errLine + b.String()
	case stateEdit:

		if m.editIndex == 0 {
			return errLine + "Edit Task - enter number to edit (ENTER to confirm):\n\n" + m.indexinput.View()
		}
		return errLine + fmt.Sprintf("Edit Task #%d - enter new text (ENTER to save):\n\n", m.editIndex) + m.textinput.View()
	case stateDelete:
		return errLine + "Delete Task - enter number (ENTER to confirm, ESC to cancel):\n\n" + m.indexinput.View()
	case stateList:
		out := m.listOutput
		if strings.TrimSpace(out) == "" {
			out = "No tasks.\n"
		}
		out += "\n\nPress any key to return to menu."
		return errLine + out
	}
	return errLine + ""
}

func (m model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":

			m.err = nil
			switch m.menu.Index() {
			case int(menuAdd):
				m.state = stateAdd
				m.textinput.SetValue("")
				m.textinput.Focus()
				return m, nil
			case int(menuComplete):

				m = m.withBuiltCompleteItems()
				m.completeCursor = 0
				m.state = stateComplete
				return m, nil
			case int(menuEdit):
				m.state = stateEdit
				m.editIndex = 0
				m.indexinput.SetValue("")
				m.indexinput.Focus()
				return m, nil
			case int(menuDelete):
				m.state = stateDelete
				m.indexinput.SetValue("")
				m.indexinput.Focus()
				return m, nil
			case int(menuList):

				m.listOutput = buildTableString(*m.todos)
				m.state = stateList
				return m, nil
			case int(menuExit):
				return m, tea.Quit
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)
	return m, cmd
}

func (m model) updateAdd(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.err = nil
			task := strings.TrimSpace(m.textinput.Value())
			if task == "" {
				m.err = fmt.Errorf("task cannot be empty")
			} else {
				if err := m.todos.Add(task); err != nil {
					m.err = err
				} else {
					if err := m.todos.Store(m.filename); err != nil {
						m.err = err
					}
				}
			}
			m.state = stateMenu
			return m, nil
		case "esc":
			m.err = nil
			m.state = stateMenu
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m model) updateComplete(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.completeCursor > 0 {
				m.completeCursor--
			}
			return m, nil
		case "down", "j":
			if m.completeCursor < len(m.completeItems)-1 {
				m.completeCursor++
			}
			return m, nil
		case "home":
			m.completeCursor = 0
			return m, nil
		case "end":
			if len(m.completeItems) > 0 {
				m.completeCursor = len(m.completeItems) - 1
			}
			return m, nil
		case "enter", " ":
			if len(m.completeItems) == 0 {
				return m, nil
			}
			cur := m.completeItems[m.completeCursor]
			if cur.index == 0 {
				m.applyCompletions()
				return m, nil
			}
			m.completeItems[m.completeCursor].selected = !m.completeItems[m.completeCursor].selected
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.MouseMsg:

		switch msg.Type {
		case tea.MouseWheelUp:
			if m.completeCursor > 0 {
				m.completeCursor--
			}
			return m, nil
		case tea.MouseWheelDown:
			if m.completeCursor < len(m.completeItems)-1 {
				m.completeCursor++
			}
			return m, nil
		case tea.MouseLeft:

			if len(m.completeItems) == 0 {
				return m, nil
			}
			cur := m.completeItems[m.completeCursor]
			if cur.index == 0 {
				m.applyCompletions()
				return m, nil
			}
			m.completeItems[m.completeCursor].selected = !m.completeItems[m.completeCursor].selected
			return m, nil
		}
	}
	return m, nil
}

func (m model) updateEdit(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.err = nil

			if m.editIndex == 0 {
				idx, err := strconv.Atoi(strings.TrimSpace(m.indexinput.Value()))
				if err != nil || idx <= 0 || idx > len(*m.todos) {
					m.err = fmt.Errorf("invalid index")
					m.state = stateMenu
					return m, nil
				}
				m.editIndex = idx

				m.textinput.SetValue((*m.todos)[idx-1].Task)
				m.textinput.Focus()
				return m, nil
			}

			newText := strings.TrimSpace(m.textinput.Value())
			if newText == "" {
				m.err = fmt.Errorf("task cannot be empty")
			} else {
				if err := m.todos.Edit(m.editIndex, newText); err != nil {
					m.err = err
				} else {
					if err := m.todos.Store(m.filename); err != nil {
						m.err = err
					}
				}
			}

			m.editIndex = 0
			m.state = stateMenu
			return m, nil
		case "esc":
			m.err = nil
			m.editIndex = 0
			m.state = stateMenu
			return m, nil
		}
	}
	var cmd tea.Cmd

	if m.editIndex == 0 {
		m.indexinput, cmd = m.indexinput.Update(msg)
	} else {
		m.textinput, cmd = m.textinput.Update(msg)
	}
	return m, cmd
}

func (m model) updateDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.err = nil
			idx, err := strconv.Atoi(strings.TrimSpace(m.indexinput.Value()))
			if err != nil {
				m.err = fmt.Errorf("invalid index")
			} else {
				if err2 := m.todos.Delete(idx); err2 != nil {
					m.err = err2
				} else {
					if err3 := m.todos.Store(m.filename); err3 != nil {
						m.err = err3
					}
				}
			}
			m.state = stateMenu
			return m, nil
		case "esc":
			m.err = nil
			m.state = stateMenu
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.indexinput, cmd = m.indexinput.Update(msg)
	return m, cmd
}

func (m model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {

	if _, ok := msg.(tea.KeyMsg); ok {
		m.state = stateMenu
	}
	return m, nil
}

func (m model) withBuiltCompleteItems() model {
	items := make([]completeItem, 0, len(*m.todos)+1)
	for i, it := range *m.todos {
		if !it.Done {
			items = append(items, completeItem{
				title:    fmt.Sprintf("%d. %s", i+1, it.Task),
				index:    i + 1,
				selected: false,
			})
		}
	}

	items = append(items, completeItem{title: "main menu", index: 0})
	m.completeItems = items
	return m
}

func (m *model) applyCompletions() {
	m.err = nil
	for _, it := range m.completeItems {
		if it.index > 0 && it.selected {
			if err := m.todos.Complete(it.index); err != nil {
				m.err = err
			}
		}
	}
	if err := m.todos.Store(m.filename); err != nil {
		m.err = err
	}
	m.state = stateMenu
}

func buildTableString(tasks tasky.Todos) string {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Tasks"},
			{Align: simpletable.AlignCenter, Text: "State"},
			{Align: simpletable.AlignRight, Text: "Created At"},
			{Align: simpletable.AlignRight, Text: "Completed At"},
		},
	}

	var cells [][]*simpletable.Cell
	for i, it := range tasks {
		done := "❌"
		completedAt := "-"
		taskText := it.Task
		if it.Done {
			done = "✅"
			completedAt = it.CompletedAt.Format(time.RFC822)
		}
		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", i+1)},
			{Text: taskText},
			{Text: done},
			{Text: it.CreatedAt.Format(time.RFC822)},
			{Text: completedAt},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}
	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{
				Align: simpletable.AlignCenter,
				Span:  5,
				Text:  fmt.Sprintf("You have %d pending tasks", tasks.CountPending()),
			},
		},
	}
	table.SetStyle(simpletable.StyleUnicode)

	return table.String()
}
