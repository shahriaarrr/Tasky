package main

import (
	"flag"
	"fmt"
	"github/tasky"
	"os"
)

const (
	taskFile = ".tasky.jsone"
)

func main() {
	add := flag.Bool("add", false, "add new task")
	complete := flag.Int("complete", 0, "task completed✅")
	rm := flag.Int("rm", 0, "task removed successfully :)")
	flag.Parse()

	tasks := &tasky.Todos{}

	if err := tasks.Load(taskFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		tasks.Add("task")
		err := tasks.Store(taskFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *complete > 0:
		err := tasks.Complete(*complete)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = tasks.Store(taskFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *rm > 0:
		err := tasks.Delete(*rm)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = tasks.Store(taskFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(0)
	}
}
