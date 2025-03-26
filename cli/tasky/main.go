package main

import (
	"bufio"
	"bytes"
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
	taskFile    = ".tasky.json"
	exitSuccess = 0
	exitError   = 1
)

var (
	errEmptyTask    = errors.New("your task is empty")
	errInvalidUsage = errors.New("invalid command usage")
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(exitError)
	}
}

func run() error {
	// Create a new FlagSet with ContinueOnError to allow manual error handling.
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	// override the default behavior by manually checking help flags.
	fs.SetOutput(new(bytes.Buffer))

	// Define command-line flags using the custom flag set.
	help := fs.Bool("h", false, "display help")
	helpLong := fs.Bool("help", false, "display help (long flag)\n")
	add := fs.Bool("a", false, "add new task")
	addLong := fs.Bool("add", false, "add new task (long flag)\n")
	complete := fs.Int("c", 0, "task completed")
	completeLong := fs.Int("complete", 0, "task completed (long flag)\n")
	remove := fs.Int("r", 0, "task removed successfully")
	removeLong := fs.Int("remove", 0, "task removed successfully (long flag)\n")
	list := fs.Bool("l", false, "list all tasks")
	listLong := fs.Bool("list", false, "list all tasks (long flag)\n")
	edit := fs.Bool("e", false, "edit your task")
	editLong := fs.Bool("edit", false, "edit your task (long flag)\n")

	// Parse command-line arguments using the custom flag set.
	err := fs.Parse(os.Args[1:])
	if err != nil {
		// Print the custom error message to stderr.
		fmt.Fprintln(os.Stderr, errInvalidUsage.Error())
		return errInvalidUsage
	}

	// If the help flag is provided, print the usage message and exit.
	if *help || *helpLong {
		// Change output to stdout before printing the usage message.
		fs.SetOutput(os.Stdout)
		fs.Usage()
		return nil
	}

	// Make the parsed flag set available for flag.Args() calls in other functions.
	flag.CommandLine = fs

	// Display the welcome menu if no flags and positional arguments are provided.
	if fs.NFlag() == 0 && fs.NArg() == 0 {
		displayMenu()
		return nil
	}

	tasks := &tasky.Todos{}

	// Load tasks from the file.
	if err := tasks.Load(taskFile); err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	// Choose between short and long flags.
	addTask := *addLong || *add
	completeTask := *completeLong
	if completeTask == 0 {
		completeTask = *complete
	}
	removeTask := *removeLong
	if removeTask == 0 {
		removeTask = *remove
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

	fmt.Printf("\nBoom! Task added: %s 🤘➕.\nNow go crush it like a boss—or just let it chill like your unread PMs😜! \n\n", task)

	return nil
}

func handleCompleteTask(tasks *tasky.Todos, index int) error {
	if err := tasks.Complete(index); err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}

	if err := tasks.Store(taskFile); err != nil {
		return fmt.Errorf("failed to store tasks: %w", err)
	}

	fmt.Printf("\nBoom! Task %d got smashed like your weekend plans! 🤘💥✅\n\n", index)
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

	fmt.Printf("\nLook at that! Task %d got a facelift – even your mom approves! 😎📝✨\n\n", index)
	return nil
}

func handleRemoveTask(tasks *tasky.Todos, index int) error {
	if err := tasks.Delete(index); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if err := tasks.Store(taskFile); err != nil {
		return fmt.Errorf("failed to store tasks: %w", err)
	}

	fmt.Printf("\nAdios! Task %d vanished faster than your last paycheck! 😂🗑️🚀\n\n", index)
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
