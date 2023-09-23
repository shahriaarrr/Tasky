package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github/tasky"
	"io"
	"os"
	"strings"
)

const (
	taskFile = ".tasky.jsone"
)

func main() {
	add := flag.Bool("add", false, "add new task")
	complete := flag.Int("complete", 0, "task completed✅")
	rm := flag.Int("rm", 0, "task removed successfully :)")
	list := flag.Bool("list", false, "all Tasks")
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

	case *list:
		tasks.Print()

	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(0)
	}
}

func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(scanner.Text()) == 0 {
		return "", errors.New("your task is empty :)")
	}

	return text, nil
}
