package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print("\r")
	}()

	fmt.Print("\033[12A\033[J")

	// fmt.Print("Tasks:\n")
	// fmt.Print("\r   \033[31mTask numero uno\033[0m\n")
	// fmt.Print("\r   	|\n")
	// fmt.Print("\r   	\u251C\u2500[E]Edit\n")
	// fmt.Print("\r   	\u2514\u2500[D]Delete\n")
	// fmt.Print("\r\n")
	// fmt.Print("\r   Task numero dos\n")
	// fmt.Print("\r   Task numero tres\n")
	// fmt.Print("\r\n")
	// fmt.Print("\r   [A] Add	[X]Clear all\n")

	tasks := []string{"Task numero uno", "Task numero dos", "Task numero tres"}
	selected := 0
	entered := -1
	buf := make([]byte, 3)
	for {
		printTasks(tasks, selected, entered)
		printCommands()
		fmt.Print("\n\r\033[?25l")
		os.Stdin.Read(buf)
		switch {
		case buf[0] == 'q' || buf[0] == 'Q':
			return
		case buf[0] == 'a' || buf[0] == 'A':
			addTask(&tasks, "some new task")
		case buf[0] == 13:
			entered = selected
		case buf[0] == 27 && buf[1] == 91:
			entered = -1
			switch buf[2] {
			case 65:
				selected = (selected - 1 + len(tasks)) % len(tasks)
			case 66:
				selected = (selected + 1) % len(tasks)

			}
		case buf[0] == 27:
			entered = -1

		}

	}

}

func printTasks(tasks []string, selected int, entered int) {
	// if entered != -1 {
	// 	fmt.Printf("\033[%dA", len(tasks))
	// } else {
	// 	fmt.Printf("\033[%dA", len(tasks)+2)

	// }
	fmt.Print("\033[2J\033[H]")
	for i, task := range tasks {
		prefix := "  "
		if i == selected {
			prefix = "> "
		}
		if i == entered {
			prefix = "> "
			task = task + "\r\n   \u251C\u2500 Description: \033[3mSome random description\033[0m"
			task = task + "\r\n   \u2502"
			task = task + "\r\n   \u2514\u2500 \033[34m[E]\033[0m Edit   \033[31m[D]\033[0m Delete"
			// task = task + "\r\n   \u251C\u2500 [D] Delete"
		}
		fmt.Printf("\r\033[K%s%s\n", prefix, task)
	}
}

func printCommands() {
	fmt.Print("\r\n")
	fmt.Print("\033[32m[A]\033[0m Add task    \033[31m[Q]\033[0m  Quit")
}

func addTask(tasks *[]string, newTask string) {
	*tasks = append(*tasks, newTask)
}
