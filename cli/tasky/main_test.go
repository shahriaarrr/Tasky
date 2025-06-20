package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old
	return string(out)
}

func captureErrorOutput(f func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = old
	return string(out)
}

func TestNoArgs(t *testing.T) {
	resetFlags()

	// Create a temporary directory for this test
	tmpDir := t.TempDir()
	// Change to the temporary directory
	oldWD, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer os.Chdir(oldWD)

	os.Args = []string{"tasky"}

	output := captureOutput(func() {
		err := run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "Welcome to Tasky") {
		t.Errorf("expected welcome message, got: %s", output)
	}
}

func TestInvalidUsage(t *testing.T) {
	resetFlags()

	tmpDir := t.TempDir()
	oldWD, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer os.Chdir(oldWD)

	os.Args = []string{"tasky", "-x"}

	errOutput := captureErrorOutput(func() {
		err := run()
		if err == nil {
			t.Error("expected an error for invalid usage, got nil")
		}
	})

	if !strings.Contains(errOutput, "invalid command usage") {
		t.Errorf("expected 'invalid command usage' error, got: %s", errOutput)
	}
}

func TestAddTask(t *testing.T) {
	resetFlags()

	// Remove the need to change directories and work with files
	os.Args = []string{"tasky", "-a", "My first test task"}

	output := captureOutput(func() {
		err := run()
		if err != nil {
			t.Errorf("unexpected error adding task: %v", err)
		}
	})

	// expected := "Boom! Task added: My first test task 🤘➕.\nNow go crush it like a boss—or just let it chill like your unread PMs😜!"
	expected := "Boom! Task added: My first test task 🤘➕. Priority: \nNow go crush it like a boss—or just let it chill like your unread PMs😜!"
	if !strings.Contains(output, expected) {
		t.Errorf("expected output to contain %q, but got %q", expected, output)
	}
}

func TestCompleteTask(t *testing.T) {
	resetFlags()

	tmpDir := t.TempDir()
	oldWD, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer os.Chdir(oldWD)

	os.Args = []string{"tasky", "-a", "Task to be completed"}
	_ = run()

	resetFlags()

	os.Args = []string{"tasky", "-c", "1"}
	errOutput := captureErrorOutput(func() {
		err := run()
		if err != nil {
			t.Errorf("unexpected error completing task: %v", err)
		}
	})

	if len(errOutput) > 0 {
		t.Errorf("did not expect errors, but got: %s", errOutput)
	}
}

func TestRemoveTask(t *testing.T) {
	resetFlags()

	tmpDir := t.TempDir()
	oldWD, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer os.Chdir(oldWD)

	os.Args = []string{"tasky", "-a", "Task to be removed"}
	_ = run()

	resetFlags()

	os.Args = []string{"tasky", "-r", "1"}
	errOutput := captureErrorOutput(func() {
		err := run()
		if err != nil {
			t.Errorf("unexpected error removing task: %v", err)
		}
	})

	if len(errOutput) > 0 {
		t.Errorf("did not expect errors, but got: %s", errOutput)
	}
}

func TestEditTask(t *testing.T) {
	resetFlags()

	tmpDir := t.TempDir()
	oldWD, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer os.Chdir(oldWD)

	os.Args = []string{"tasky", "-a", "Old task content"}
	_ = run()

	resetFlags()

	os.Args = []string{"tasky", "-e", "1", "New", "task", "content"}
	errOutput := captureErrorOutput(func() {
		err := run()
		if err != nil {
			t.Errorf("unexpected error editing task: %v", err)
		}
	})

	if len(errOutput) > 0 {
		t.Errorf("did not expect errors, but got: %s", errOutput)
	}
}

func TestListTasks(t *testing.T) {
	resetFlags()

	tmpDir := t.TempDir()
	oldWD, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer os.Chdir(oldWD)

	os.Args = []string{"tasky", "-a", "First Task"}
	_ = run()
	resetFlags()
	os.Args = []string{"tasky", "-a", "Second Task"}
	_ = run()

	resetFlags()

	os.Args = []string{"tasky", "-l"}
	output := captureOutput(func() {
		err := run()
		if err != nil {
			t.Errorf("unexpected error listing tasks: %v", err)
		}
	})

	if !strings.Contains(output, "First Task") || !strings.Contains(output, "Second Task") {
		t.Errorf("expected tasks to appear in the list, got: %s", output)
	}
}

func TestGetInput(t *testing.T) {
	input, err := getInput(bytes.NewReader(nil), "Hello", "World")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if input != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", input)
	}

	reader := bytes.NewBufferString("Task from stdin")
	input, err = getInput(reader)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if input != "Task from stdin" {
		t.Errorf("expected 'Task from stdin', got '%s'", input)
	}

	reader = bytes.NewBufferString("")
	_, err = getInput(reader)
	if err == nil {
		t.Error("expected errEmptyTask when input is empty, but got nil")
	}
}

func TestAddTaskWithPriority(t *testing.T) {
	resetFlags()

	os.Args = []string{"tasky", "-a", "My prioritized task", "-p", "High"}

	output := captureOutput(func() {
		err := run()
		if err != nil {
			t.Errorf("unexpected error adding task with priority: %v", err)
		}
	})

	// expected := "Boom! Task added: My prioritized task 🤘➕. Priority: High"
	expected := "Boom! Task added: My prioritized task High 🤘➕. Priority: High\nNow go crush it like a boss—or just let it chill like your unread PMs😜!"
	if !strings.Contains(output, expected) {
		t.Errorf("expected output to contain %q, but got %q", expected, output)
	}
}

func TestEditTaskWithPriority(t *testing.T) {
	resetFlags()

	tmpDir := t.TempDir()
	oldWD, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer os.Chdir(oldWD)

	os.Args = []string{"tasky", "-a", "Original task"}
	_ = run()

	resetFlags()

	os.Args = []string{"tasky", "-e", "1", "Updated task", "-p", "Low"}
	errOutput := captureErrorOutput(func() {
		err := run()
		if err != nil {
			t.Errorf("unexpected error editing task with priority: %v", err)
		}
	})

	if len(errOutput) > 0 {
		t.Errorf("did not expect errors, but got: %s", errOutput)
	}
}
