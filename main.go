package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/exp/constraints"
	"golang.org/x/term"

	_ "modernc.org/sqlite"
)

// TODO-Ception:

type Task struct {
	ID          string
	Name        string
	Description string
	// Entered     bool
}

var opened map[string]struct{} = map[string]struct{}{}

// screen functions
func clearScreen() {
	fmt.Print("\033[2J\033[3J\033[H")
	// clear current sccreen, remove scrollback and set cursor to (1,1)
}

// task functions
func getTasks(db *sql.DB) []Task {
	rows, err := db.Query(`SELECT id, name, description FROM tasks`)
	if err != nil {
		printLog("error reading from db") // this doesnt work. do something here; for notifications in TUI
		return []Task{}
	}
	var tasks []Task
	for rows.Next() {
		var id string
		var name string
		var description string
		if err := rows.Scan(&id, &name, &description); err != nil {
			printLog("error reading row")
		}
		t := Task{ID: id, Name: name, Description: description}
		tasks = append(tasks, t)
	}
	return tasks
}

func deleteTask(db *sql.DB, id string) {
	_, err := db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		printLog("delete error")
	}

}

func addTask(db *sql.DB, newTask Task) {
	// show new input for task
	// take imput
	nt := newTask
	nt.ID = idGen()

	if len(nt.Name) > 0 {

		_, err := db.Exec(`INSERT INTO tasks (id, name, description) VALUES (?,?, ?)`, nt.ID, nt.Name, nt.Description)
		if err != nil {
			printLog("error during insert") // no working :(
		}
	}

}

func editTask(db *sql.DB, updatedTask Task) {

	_, err := db.Exec(`UPDATE tasks SET name = ?, description = ? WHERE id = ?`, updatedTask.Name, updatedTask.Description, updatedTask.ID)
	if err != nil {
		printLog("update error")
	}

}

// utilities
func printLog(msg string) {
	fmt.Print("\r\n\033[33m[LOG]: \033[0m " + msg)
}

func printHeader() {
	fmt.Print("\r\033[47m\033[30mTask Manager v9000\033[30m                                      " + getDate() + "\033[0m")
	fmt.Print("\r\n\n")
}

func printTasks(tasks []Task, selected int) {

	if len(tasks) == 0 {
		fmt.Print("\r\033[3m  No pending tasks \033[0m^_^")
	}

	for i, task := range tasks {
		prefix := "  "
		ntask := "\033[35m" + task.Name + "\033[0m"
		if i == selected {
			prefix = "\033[1m>\033[0m "
		}
		if _, exists := opened[task.ID]; exists {
			if i == selected {
				prefix = "> "
			} else {
				prefix = "  "

			}

			// wrapping: split description into lines ( prepend description )
			maxlength := 40

			newdescription := "\033[3m" + task.Description + "\033[0m"
			var div int
			if len(newdescription)/maxlength == 0 {
				div = 1
			} else {
				div = len(newdescription) / maxlength
			}

			chunksize := len(newdescription) / div
			lines := len(newdescription)/maxlength + 1
			for i := range lines {

				if i == 0 {
					ntask = ntask + "\r\n   \u251C\u2500\033[4mDescription\033[0m: "
				} else {

					ntask = ntask + "\r\n   \u2502              "
				}
				ntask = ntask + newdescription[chunksize*i:min(len(newdescription), chunksize*i+chunksize)]
			}
			// ntask = ntask + "\r\n   \u2502"
			ntask = ntask + "\r\n   \u2514\u2500 \033[34m[E]\033[0m Edit   \033[31m[D]\033[0m Delete"

			// ntask = ntask + "\r\n   \u251C\u2500 Id: \033[3m" + task.ID + "\033[0m"
			// ntask = ntask + "\r\n   \u251C\u2500 Description: \033[3m" + task.Description + "\033[0m"
			// ntask = ntask + "\r\n   \u2502"
			// ntask = ntask + "\r\n   \u2514\u2500 \033[34m[E]\033[0m Edit   \033[31m[D]\033[0m Delete"
			// task = task + "\r\n   \u251C\u2500 [D] Delete"
		}
		fmt.Printf("\r%s%s\n", prefix, ntask)
	}
}

func printCommands() {
	fmt.Print("\r\n\n")
	fmt.Print("\033[47m\033[32m[A]\033[30m Add task    \033[31m[Q]\033[30m  Quit                                          \033[0m")
}

func getDate() string {
	return time.Now().Format("02-Jan-2006")
}

func idGen() string {
	return fmt.Sprintf("%d-%d", rand.Intn(0x100000), time.Now().UnixMilli())
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// input
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
			if string(input) == "Y" || string(input) == "y" {
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

func readInput(width int) Task {

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
				if (len(input)+len("New task: "))%width == width-1 {
					fmt.Printf("\033[A\033[%dC \033[%dD", width, 1)
					fmt.Print("\033[1C")
				} else {
					fmt.Print("\b \b")
				}

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
				if (len(input)+len("Description: "))%width == width-1 {
					fmt.Printf("\033[A\033[%dC \033[%dD", width, 1)
					fmt.Print("\033[1C")
				} else {
					fmt.Print("\b \b")
				}
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

func getUpdatedTask(width int, task Task) Task {

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
				if (len(input)+len("Updated task name: "))%width == width-1 {
					fmt.Printf("\033[A\033[%dC \033[%dD", width, 1)
					fmt.Print("\033[1C")
				} else {
					fmt.Print("\b \b")
				}
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
				if (len(input)+len("Updated Description: "))%width == width-1 {
					fmt.Printf("\033[A\033[%dC \033[%dD", width, 1)
					fmt.Print("\033[1C")
				} else {
					fmt.Print("\b \b")
				}
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

func main() {

	// open sqlite db file
	db, err := sql.Open("sqlite", "file:db/tasks.db")
	if err != nil {
		printLog("unable to open db")
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT
	)`)
	if err != nil {
		printLog("error during table creation")
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print("\r")
	}()

	width, _, _ := term.GetSize(int(os.Stderr.Fd()))

	// tasks := []Task{{ID: idGen(), Name: "Task Numero Uno", Description: "Description Uno"}, {ID: idGen(), Name: "Task Numero Dos", Description: "Description Dos"}, {ID: idGen(), Name: "Task Numero Tres", Description: "Description Tres"}}
	// tasks := []Task{} // replace this with db read
	tasks := getTasks(db)

	selected := 0

	buf := make([]byte, 3)
	for {
		for i := range buf {
			buf[i] = 0
		}
		clearScreen()

		printHeader()
		printTasks(tasks, selected)
		printCommands()
		fmt.Print("\n\r\033[?25l")
		os.Stdin.Read(buf)
		switch {
		case buf[0] == 'q' || buf[0] == 'Q':
			return
		case buf[0] == 'a' || buf[0] == 'A':
			// try to move to some function

			newtask := readInput(width)

			addTask(db, newtask)

			if len(tasks) == 0 {
				selected = 0
			}
			tasks = getTasks(db)

		case buf[0] == 13 || buf[0] == 10:

			if len(tasks) != 0 {

				if _, exists := opened[tasks[selected].ID]; exists {
					delete(opened, tasks[selected].ID)
				} else {
					opened[tasks[selected].ID] = struct{}{}
				}
			}

		case buf[0] == 27 && buf[1] == 91: // arrow cases

			switch buf[2] {
			case 65: // up
				if len(tasks) > 0 {
					selected = (selected - 1 + len(tasks)) % len(tasks)
				}
			case 66: // down
				if len(tasks) > 0 {
					selected = (selected + 1) % len(tasks)
				}

			}
		case buf[0] == 27: // escape
			delete(opened, tasks[selected].ID)
			// tasks[selected].Entered = false

		case buf[0] == 'e' || buf[0] == 'E':

			if _, exists := opened[tasks[selected].ID]; exists {

				// TODO: edit selected here
				// call edit task function here
				updatedTask := getUpdatedTask(width, tasks[selected])
				editTask(db, updatedTask)
				tasks = getTasks(db)

			}

		case buf[0] == 'd' || buf[0] == 'D':
			if _, exists := opened[tasks[selected].ID]; exists {
				// delete selected here
				// TODO: get confirmation here

				doit := getConfirmation()
				if doit {

					deleteTask(db, tasks[selected].ID)

					tasks = getTasks(db)
					if selected > len(tasks)-1 {
						selected = len(tasks) - 1
					}
				} else {
					fmt.Print("\r\n\nAborted.")
					// check how to do this
				}
				//rotate
			}
		}

	}

}
