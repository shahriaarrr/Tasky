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

const (
	taskFile = ".tasky.json"
	exitSuccess = 0
	exitError = 1
)

var (
	errEmptyTask = errors.New("your task is empty")
	errInvalidUsage = errors.New("invalid command usage")
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(exitError)
	}
}

func run() error {
	// Define command-line flags
	add := flag.Bool("a", false, "add new task")
	addLong := flag.Bool("add", false, "add new task (long flag)")
	complete := flag.Int("c", 0, "task completed")
	completeLong := flag.Int("complete", 0, "task completed (long flag)")
	rm := flag.Int("r", 0, "task removed successfully")
	rmLong := flag.Int("rm", 0, "task removed successfully (long flag)")
	list := flag.Bool("l", false, "list all tasks")
	listLong := flag.Bool("list", false, "list all tasks (long flag)")
	edit := flag.Bool("e", false, "edit your task")
	editLong := flag.Bool("edit", false, "edit your task (long flag)")
	flag.Parse()

	// Display the welcome menu if no commands are provided
	if len(os.Args) == 1 {
		displayMenu()
		return nil
	}

	tasks := &tasky.Todos{}

	// Load tasks from file
	if err := tasks.Load(taskFile); err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	// Use the long flag if it's provided, otherwise use the short flag
	addTask := *addLong || *add
	completeTask := *completeLong
	if completeTask == 0 {
		completeTask = *complete
	}
	removeTask := *rmLong
	if removeTask == 0 {
		removeTask = *rm
	}
	listTasks := *listLong || *list
	editTask := *editLong || *edit

	switch {
	case addTask:
		return handleAddTask(tasks)
	case completeTask > 0:
		return handleCompleteTask(tasks, completeTask)
	case editTask:
		return handleEditTask(tasks)
	case removeTask > 0:
		return handleRemoveTask(tasks, removeTask)
	case listTasks:
		tasks.Print()
		return nil
	default:
		return errInvalidUsage
	}
}

func handleAddTask(tasks *tasky.Todos) error {
	task, err := getInput(os.Stdin, flag.Args()...)
	if err != nil {
		return fmt.Errorf("failed to get input: %w", err)
	}

	tasks.Add(task)
	if err := tasks.Store(taskFile); err != nil {
		return fmt.Errorf("failed to store tasks: %w", err)
	}
	return nil
}

func handleCompleteTask(tasks *tasky.Todos, index int) error {
	if err := tasks.Complete(index); err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}

	if err := tasks.Store(taskFile); err != nil {
		return fmt.Errorf("failed to store tasks: %w", err)
	}
	return nil
}

func handleEditTask(tasks *tasky.Todos) error {
	if len(flag.Args()) < 2 {
		return fmt.Errorf("usage: -e <index> <new_task>")
	}

	index, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid index: %w", err)
	}

	newTask := strings.Join(flag.Args()[1:], " ")
	if err := tasks.Edit(index, newTask); err != nil {
		return fmt.Errorf("failed to edit task: %w", err)
	}

	if err := tasks.Store(taskFile); err != nil {
		return fmt.Errorf("failed to store tasks: %w", err)
	}
	return nil
}

func handleRemoveTask(tasks *tasky.Todos, index int) error {
	if err := tasks.Delete(index); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if err := tasks.Store(taskFile); err != nil {
		return fmt.Errorf("failed to store tasks: %w", err)
	}
	return nil
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
		return "", errEmptyTask
	}

	return text, nil
}

