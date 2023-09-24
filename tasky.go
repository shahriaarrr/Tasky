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
	Task         string
	done         bool
	created_at   time.Time
	completed_at time.Time
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:         task,
		done:         false,
		created_at:   time.Now(),
		completed_at: time.Time{},
	}

	*t = append(*t, todo)
}

func (t *Todos) Complete(index int) error {
	if index <= 0 || index > len(*t) {
		return errors.New("invalid index")
	}

	(*t)[index-1].completed_at = time.Now()
	(*t)[index-1].done = true

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
	file, err := os.ReadFile(filename)

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

	return os.WriteFile(filename, data, 0644)

}

func (t *Todos) Print() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Tasks"},
			{Align: simpletable.AlignCenter, Text: "state"},
			{Align: simpletable.AlignRight, Text: "created at"},
			{Align: simpletable.AlignRight, Text: "completed at"},
		},
	}

	var cells [][]*simpletable.Cell

	for index, item := range *t {
		index++
		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", index)},
			{Text: item.Task},
			{Text: fmt.Sprintf("%t", item.done)},
			{Text: item.created_at.Format(time.RFC822)},
			{Text: item.completed_at.Format(time.RFC822)},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: "this is your Tasks"},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}
