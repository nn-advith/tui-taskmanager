package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// TODO-Ception:
// verify edge case when tasks is nil
// delete confirmation in a new prompt
// edit task section
// sqlite3 integration with unique id
// cleanly remove scrollback; hide render changes
// figure out warping

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print("\r")
	}()

	tasks := []string{"Task numero uno", "Task numero dos", "Task numero tres"}
	selected := 0
	entered := -1
	buf := make([]byte, 3)
	for {
		for i := range buf {
			buf[i] = 0
		}
		// fmt.Print("\033[3J")       // remove scrollback
		fmt.Print("\033[2J\033[H") // clear screen ; not optimal but eh
		printHeader()
		printTasks(tasks, selected, entered)
		printCommands()
		fmt.Print("\n\r\033[?25l")
		os.Stdin.Read(buf)
		switch {
		case buf[0] == 'q' || buf[0] == 'Q':
			return
		case buf[0] == 'a' || buf[0] == 'A':
			// try to move to some function

			newtaskname := readInput()
			addTask(&tasks, newtaskname)
		case buf[0] == 13 || buf[0] == 10:
			if entered == -1 {
				entered = selected
			} else {
				entered = -1
			}
		case buf[0] == 27 && buf[1] == 91:
			entered = -1
			switch buf[2] {
			case 65:
				if len(tasks) > 0 {
					selected = (selected - 1 + len(tasks)) % len(tasks)
				}
			case 66:
				if len(tasks) > 0 {
					selected = (selected + 1) % len(tasks)
				}

			}
		case buf[0] == 27:
			entered = -1

		case buf[0] == 'e' || buf[0] == 'E':
			if entered != -1 {
				// TODO: edit selected here

			}
		case buf[0] == 'd' || buf[0] == 'D':
			if entered != -1 {
				// delete selected here
				// TODO: get confirmation here
				if selected == len(tasks)-1 {
					selected = 0
				}
				tasks = append(tasks[:selected], tasks[selected+1:]...)
				entered = -1
				//rotate
			}
		}

	}

}

func readInput() string {
	var input []byte
	buf := make([]byte, 1)
	fmt.Print("\r\n\033[?25h\033[?12hNew task: ")
	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		switch buf[0] {
		case 13, 10:
			return string(input)
		case 127, 8:
			if len(input) > 0 {
				input = input[:len(input)-1]
				fmt.Print("\b \b")
			}
		case 27:
			return ""
		default:
			input = append(input, buf[0])
			fmt.Printf("%c", buf[0])
		}
	}

	return string(input)
}

func printHeader() {
	fmt.Print("\r\033[47m\033[30mTask Manager v9000\033[30m            17-Aug-2025\033[0m")
	fmt.Print("\r\n\n")
}

func printTasks(tasks []string, selected int, entered int) {

	if len(tasks) == 0 {
		fmt.Print("\r\033[3m  All tasks done <=w=>\033[0m")
	}
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
		fmt.Printf("\r%s%s\n", prefix, task)
	}
}

func printCommands() {
	fmt.Print("\r\n\n")
	fmt.Print("\033[47m\033[32m[A]\033[30m Add task    \033[31m[Q]\033[30m  Quit                \033[0m")
}

func addTask(tasks *[]string, newTask string) {
	// show new input for task
	// take imput
	if len(newTask) > 0 {
		*tasks = append(*tasks, newTask)
	}
}
