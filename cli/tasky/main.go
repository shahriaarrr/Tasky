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
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(new(bytes.Buffer))

	help := fs.Bool("h", false, "display help")
	helpLong := fs.Bool("help", false, "display help (long flag)")
	add := fs.Bool("a", false, "add new task")
	addLong := fs.Bool("add", false, "add new task (long flag)")
	complete := fs.Int("c", 0, "task completed")
	completeLong := fs.Int("complete", 0, "task completed (long flag)")
	remove := fs.Int("r", 0, "task removed successfully")
	removeLong := fs.Int("remove", 0, "task removed successfully (long flag)")
	list := fs.Bool("l", false, "list all tasks")
	listLong := fs.Bool("list", false, "list all tasks (long flag)")
	edit := fs.Bool("e", false, "edit your task")
	editLong := fs.Bool("edit", false, "edit your task (long flag)")
	
	priority := fs.String("p", "", "Set task priority (Low/Medium/High)")
	priorityLong := fs.String("priority", "", "Set task priority (Low/Medium/High)")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errInvalidUsage.Error())
		return errInvalidUsage
	}

	if *help || *helpLong {
		fs.SetOutput(os.Stdout)
		fs.Usage()
		return nil
	}

	flag.CommandLine = fs

	if fs.NFlag() == 0 && fs.NArg() == 0 {
		displayMenu()
		return nil
	}

	tasks := &tasky.Todos{}

	if err := tasks.Load(taskFile); err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	priorityValue := *priority
	if priorityValue == "" {
		priorityValue = *priorityLong
	}
	
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
		var filteredArgs []string
		var priorityFlag string
		for i := 0; i < len(os.Args[1:]); i++ {
			arg := os.Args[1:][i]
			if arg == "-p" || arg == "--priority" {
				if i+1 < len(os.Args[1:]) {
					priorityFlag = os.Args[1:][i+1]
					i++ 
				}
				continue
			}
			if strings.HasPrefix(arg, "-p=") {
				priorityFlag = strings.TrimPrefix(arg, "-p=")
				continue
			}
			if strings.HasPrefix(arg, "--priority=") {
				priorityFlag = strings.TrimPrefix(arg, "--priority=")
				continue
			}
			filteredArgs = append(filteredArgs, arg)
		}

		if priorityValue == "" {
			priorityValue = priorityFlag
		}
		
		oldArgs := os.Args
		os.Args = append([]string{os.Args[0]}, filteredArgs...)
		defer func() { os.Args = oldArgs }()

		return handleAddTask(tasks, tasky.Priority(priorityValue))
	case completeTask > 0:
		return handleCompleteTask(tasks, completeTask)
	case editTask:
		return handleEditTask(tasks, tasky.Priority(priorityValue))
	case removeTask > 0:
		return handleRemoveTask(tasks, removeTask)
	case listTasks:
		tasks.Print()
		return nil
	default:
		return errInvalidUsage
	}
}

func handleAddTask(tasks *tasky.Todos, priority tasky.Priority) error {
	var filteredArgs []string
	for _, arg := range flag.Args() {
		if arg != "-p" && arg != "--priority" && 
		   !strings.HasPrefix(arg, "-p=") && 
		   !strings.HasPrefix(arg, "--priority=") {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	task := strings.Join(filteredArgs, " ")

	if task == "" {
		return fmt.Errorf("task description cannot be empty")
	}

	if err := tasks.Add(task, priority); err != nil {
		return fmt.Errorf("failed to add task: %w", err)
	}

	if err := tasks.Store(taskFile); err != nil {
		return fmt.Errorf("failed to store tasks: %w", err)
	}

	fmt.Printf("\nBoom! Task added: %s 🤘➕. Priority: %s\nNow go crush it like a boss—or just let it chill like your unread PMs😜! \n\n", task, priority)

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

func handleEditTask(tasks *tasky.Todos, priority tasky.Priority) error {
	if len(flag.Args()) < 2 {
		return fmt.Errorf("usage: -e <index> <new_task> [-p priority]")
	}

	index, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid index: %w", err)
	}

	newTask := strings.Join(flag.Args()[1:], " ")
	if err := tasks.Edit(index, newTask, priority); err != nil {
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

	New Feature: Task Priorities!
	Use -p or --priority flag to set task priority (Low/Medium/High)

	Available commands:
	- Add task: tasky -a "Task description" -p Medium
	- Edit task: tasky -e 1 "New task" -p High
	- Complete task: tasky -c 1
	- Remove task: tasky -r 1
	- List tasks: tasky -l

	You can see more details with -h command.

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
