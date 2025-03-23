package tasky

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alexeyco/simpletable"
)

var (
	errInvalidIndex = errors.New("invalid index")
	errEmptyTask = errors.New("task cannot be empty")
)

type item struct {
	Task        string    `json:"task"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

type Todos []item

func (t *Todos) Add(task string) error {
	if len(task) == 0 {
		return errEmptyTask
	}

	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
	return nil
}

func (t *Todos) Complete(index int) error {
	if !t.isValidIndex(index) {
		return errInvalidIndex
	}

	(*t)[index-1].CompletedAt = time.Now()
	(*t)[index-1].Done = true

	return nil
}

func (t *Todos) Edit(index int, newTask string) error {
	if !t.isValidIndex(index) {
		return errInvalidIndex
	}

	if len(newTask) == 0 {
		return errEmptyTask
	}

	(*t)[index-1].Task = newTask
	return nil
}

func (t *Todos) Delete(index int) error {
	if !t.isValidIndex(index) {
		return errInvalidIndex
	}

	*t = append((*t)[:index-1], (*t)[index:]...)
	return nil
}

func (t *Todos) Load(filename string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	filePath := filepath.Join(homeDir, filename)
	file, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	if len(file) == 0 {
		return nil
	}

	if err := json.Unmarshal(file, t); err != nil {
		return fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	return nil
}

func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	filePath := filepath.Join(homeDir, filename)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (t *Todos) Print() {
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
	for index, item := range *t {
		task := blue(item.Task)
		done := "❌"
		completedAt := "-"

		if item.Done {
			task = green(item.Task)
			done = green("✅")
			completedAt = item.CompletedAt.Format(time.RFC822)
		}

		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", index+1)},
			{Text: task},
			{Text: done},
			{Text: item.CreatedAt.Format(time.RFC822)},
			{Text: completedAt},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}
	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{
				Align: simpletable.AlignCenter,
				Span: 5,
				Text: red(fmt.Sprintf("You have %d pending tasks", t.CountPending())),
			},
		},
	}

	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}

func (t *Todos) CountPending() int {
	total := 0
	for _, item := range *t {
		if !item.Done {
			total++
		}
	}
	return total
}

func (t *Todos) isValidIndex(index int) bool {
	return index > 0 && index <= len(*t)
}

