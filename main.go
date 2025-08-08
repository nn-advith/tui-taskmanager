package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/term"
)

// TODO-Ception:
// verify edge case when tasks is nil
// delete confirmation in a new prompt
// edit task section
// sqlite3 integration with unique id
// cleanly remove scrollback; hide render changes
// figure out warping
// multiple tasks open

type Task struct {
	ID          string
	Name        string
	Description string
}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print("\r")
	}()

	tasks := []Task{{ID: idGen(), Name: "Task Numero Uno", Description: "Description Uno"}, {ID: idGen(), Name: "Task Numero Dos", Description: "Description Dos"}, {ID: idGen(), Name: "Task Numero Tres", Description: "Description Tres"}}
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

			newtask := readInput()

			addTask(&tasks, newtask)

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
				// call edit task function here
				updatedTask := getUpdatedTask(tasks[selected])
				editTask(&tasks, updatedTask)

			}
		case buf[0] == 'd' || buf[0] == 'D':
			if entered != -1 {
				// delete selected here
				// TODO: get confirmation here

				doit := getConfirmation()
				if doit {

					if selected == len(tasks)-1 {
						selected = 0
					}
					tasks = append(tasks[:selected], tasks[selected+1:]...)
					entered = -1
				} else {
					// fmt.Print("\r\n\033[?12hAborted.")
					// check how to do this
				}
				//rotate
			}
		}

	}

}

func getConfirmation() bool {
	var input []byte
	buf := make([]byte, 1)
	fmt.Print("\r\n\033[?25h\033[?12hConfirm Delete [\033[32mY\033[0m/\033[31mother\033[0m]: ")
	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		switch buf[0] {
		case 13, 10:
			// return string(input)
			if string(input) == "Y" {
				return true
			}
			return false

		case 127, 8:
			if len(input) > 0 {
				input = input[:len(input)-1]
				fmt.Print("\b \b")
			}
		case 27:
			// return ""
			return false
		default:
			input = append(input, buf[0])
			fmt.Printf("%c", buf[0])
		}
	}
	return false
}

func readInput() Task {

	var input []byte
	buf := make([]byte, 1)
	fmt.Print("\r\n\033[?25h\033[?12hNew task: ")
	var newTask Task
	namedone := false
	descdone := false
	for !namedone {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		switch buf[0] {
		case 13, 10:
			// return string(input)
			newTask.Name = string(input)
			namedone = true
		case 127, 8:
			if len(input) > 0 {
				input = input[:len(input)-1]
				fmt.Print("\b \b")
			}
		case 27:
			// return ""
			newTask = Task{}
			namedone = true
			descdone = true
		default:
			input = append(input, buf[0])
			fmt.Printf("%c", buf[0])
		}
	}

	input = []byte{}
	fmt.Print("\r\n\033[?25h\033[?12hDescription: ")
	for !descdone {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		switch buf[0] {
		case 13, 10:
			// return string(input)
			newTask.Description = string(input)
			descdone = true
		case 127, 8:
			if len(input) > 0 {
				input = input[:len(input)-1]
				fmt.Print("\b \b")
			}
		case 27:
			// return ""
			newTask = Task{}
			descdone = true
		default:
			input = append(input, buf[0])
			fmt.Printf("%c", buf[0])
		}
	}

	return newTask

}

func printHeader() {
	fmt.Print("\r\033[47m\033[30mTask Manager v9000\033[30m            " + getDate() + "\033[0m")
	fmt.Print("\r\n\n")
}

func printTasks(tasks []Task, selected int, entered int) {

	if len(tasks) == 0 {
		fmt.Print("\r\033[3m  All tasks done <=w=>\033[0m")
	}

	for i, task := range tasks {
		prefix := "  "
		ntask := task.Name
		if i == selected {
			prefix = "> "
		}
		if i == entered {
			prefix = "> "
			// ntask = ntask + "\r\n   \u251C\u2500 Id: \033[3m" + task.ID + "\033[0m"
			ntask = ntask + "\r\n   \u251C\u2500 Description: \033[3m" + task.Description + "\033[0m"
			ntask = ntask + "\r\n   \u2502"
			ntask = ntask + "\r\n   \u2514\u2500 \033[34m[E]\033[0m Edit   \033[31m[D]\033[0m Delete"
			// task = task + "\r\n   \u251C\u2500 [D] Delete"
		}
		fmt.Printf("\r%s%s\n", prefix, ntask)
	}
}

func printCommands() {
	fmt.Print("\r\n\n")
	fmt.Print("\033[47m\033[32m[A]\033[30m Add task    \033[31m[Q]\033[30m  Quit                \033[0m")
}

func addTask(tasks *[]Task, newTask Task) {
	// show new input for task
	// take imput
	nt := newTask
	nt.ID = idGen()
	if len(nt.Name) > 0 {
		*tasks = append(*tasks, nt)
	}
}

func editTask(tasks *[]Task, updatedTask Task) {
	for i := range *tasks {
		if (*tasks)[i].ID == updatedTask.ID {
			(*tasks)[i] = updatedTask
		}
	}
}

// func editTask(tas)

func getDate() string {
	return time.Now().Format("02-Jan-2006")
}

func idGen() string {
	return fmt.Sprintf("%d-%d", rand.Intn(0x100000), time.Now().UnixMilli())
}

func getUpdatedTask(task Task) Task {

	input := []byte(task.Name)
	buf := make([]byte, 1)
	fmt.Print("\r\n\033[?25h\033[?12hUpdated task name: " + string(input))
	var newTask Task
	newTask.ID = task.ID

	namedone := false
	descdone := false
	for !namedone {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		switch buf[0] {
		case 13, 10:
			// return string(input)
			newTask.Name = string(input)
			namedone = true
		case 127, 8:
			if len(input) > 0 {
				input = input[:len(input)-1]
				fmt.Print("\b \b")
			}
		case 27:
			// return ""
			newTask = Task{}
			namedone = true
			// descdone = true
		default:
			input = append(input, buf[0])
			fmt.Printf("%c", buf[0])
		}
	}

	input = []byte(task.Description)
	fmt.Print("\r\n\033[?25h\033[?12hUpdated Description: " + string(input))
	for !descdone {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		switch buf[0] {
		case 13, 10:
			// return string(input)
			newTask.Description = string(input)
			descdone = true
		case 127, 8:
			if len(input) > 0 {
				input = input[:len(input)-1]
				fmt.Print("\b \b")
			}
		case 27:
			// return ""
			newTask = Task{}
			descdone = true
		default:
			input = append(input, buf[0])
			fmt.Printf("%c", buf[0])
		}
	}

	return newTask
}
