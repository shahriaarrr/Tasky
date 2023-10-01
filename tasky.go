package tasky

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
}

func (t *Todos) Complete(index int) error {
	if index <= 0 || index > len(*t) {
		return errors.New("invalid index")
	}

	(*t)[index-1].CompletedAt = time.Now()
	(*t)[index-1].Done = true

	return nil
}

func (t *Todos) Edit(index int, newTask string) error {
	if index <= 0 || index > len(*t) {
		return errors.New("invalid index")
	}

	(*t)[index-1].Task = newTask

	return nil
}

func (t *Todos) Delete(index int) error {
	if index <= 0 || index > len(*t) {
		return errors.New("invalid index")
	}

	*t = append((*t)[:index-1], (*t)[index:]...)

	return nil
}

func (t *Todos) Load(filename string) error {
	// file, err := os.ReadFile(filename)
	homeDir, err := os.UserHomeDir()
	homeDir = homeDir + "/" + filename
	file, err := os.ReadFile(homeDir)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	homeDir = homeDir + "/" + filename

	return os.WriteFile(homeDir, data, 0644)
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
		index++

		task := blue(item.Task)
		done := "❌"
		if item.Done {
			task = green(fmt.Sprintf("%s", item.Task))
			done = green("✅")
		}

		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", index)},
			{Text: task},
			{Text: done},
			{Text: item.CreatedAt.Format(time.RFC822)},
			{Text: item.CompletedAt.Format(time.RFC822)},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("you have %d  pending tasks", t.CountPending()))},
	}}

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
