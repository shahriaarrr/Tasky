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

	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(0)
	}
}
