package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github/tasky"
	"io"
	"os"
	"strconv"
	"strings"
)

const taskFile = ".tasky.json"

func main() {
	// Define command-line flags
	add := flag.Bool("add", false, "add new task")
	complete := flag.Int("complete", 0, "task completed")
	rm := flag.Int("rm", 0, "task removed successfully")
	list := flag.Bool("list", false, "list all tasks")
	edit := flag.Bool("edit", false, "edit your task")
	flag.Parse()

	// Display the welcome menu if no commands are provided
	if len(os.Args) == 1 {
		displayMenu()
		return
	}

	tasks := &tasky.Todos{}

	// Load tasks from file
	if err := tasks.Load(taskFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			handleError(err)
		}

		tasks.Add(task)
		err = tasks.Store(taskFile)
		handleError(err)

	case *complete > 0:
		err := tasks.Complete(*complete)
		handleError(err)

		err = tasks.Store(taskFile)
		handleError(err)

	case *edit:
		if len(flag.Args()) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: -edit <index> <new_task>")
			os.Exit(1)
		}

		index := flag.Arg(0)
		newTask := strings.Join(flag.Args()[1:], " ")

		indexInt, err := strconv.Atoi(index)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid index:", err)
			os.Exit(1)
		}

		err = tasks.Edit(indexInt, newTask)
		handleError(err)

		err = tasks.Store(taskFile)
		handleError(err)

	case *rm > 0:
		err := tasks.Delete(*rm)
		handleError(err)

		err = tasks.Store(taskFile)
		handleError(err)

	case *list:
		tasks.Print()

	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(0)
	}
}

// displayMenu displays the welcome message and available commands
func displayMenu() {
	menu := `
████████╗ █████╗ ███████╗██╗  ██╗██╗   ██╗
╚══██╔══╝██╔══██╗██╔════╝██║ ██╔╝╚██╗ ██╔╝
   ██║   ███████║███████╗█████╔╝  ╚████╔╝ 
   ██║   ██╔══██║╚════██║██╔═██╗   ╚██╔╝  
   ██║   ██║  ██║███████║██║  ██╗   ██║   
   ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝   ╚═╝   
                                          
	Welcome to Tasky👋
	Your personal command-line task manager🧑‍💼

	Tasky helps you efficiently manage your to-do list directly from the terminal.
	Whether you're tracking daily tasks, marking items as complete, or editing existing tasks,
	Tasky provides a simple yet powerful interface to keep your tasks organized.

	You can see Available commands with -h command.

	Stay on top of your tasks with Tasky!

	for more details: https://github.com/shahriaarrr/Tasky

	© Developed with ❤️  and ☕ By Shahriar Ghasempour.
`
	fmt.Println(menu)
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

	if len(text) == 0 {
		return "", errors.New("your task is empty :)")
	}

	return text, nil
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
