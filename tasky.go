package tasky

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
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
	for i, item := range *t {
		i++
		fmt.Printf("%d - %s\n", i, item.Task)
	}
}
