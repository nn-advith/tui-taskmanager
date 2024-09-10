package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

//define task

// task {
// 	task name
// 	due date
// 	current state -> [to-do, in-progress, done]
// 	desc
// 	comm
// }

var keyCounter int = 0

// taskstate
type taskState int

const (
	todo = iota
	inprogress
	done
)

var stateName = map[taskState]string{
	todo:       "To-Do",
	inprogress: "In Progress",
	done:       "Done",
}

func (ts taskState) String() string {
	return stateName[ts]
}

func transition(ts taskState) taskState {
	switch ts {
	case todo:
		return inprogress
	case inprogress:
		return done
	case done:
		return done
	default:
		panic(fmt.Errorf("Unkown task state: %s", ts))
	}
}

// task struct
type task struct {
	id     int
	tname  string
	due    time.Time
	desc   string
	comm   []string
	status taskState
}

func createTask(tname string, due time.Time, desc string, comm []string, status taskState) task {
	keyCounter = +1
	return task{id: keyCounter, tname: tname, due: due, desc: desc, comm: comm, status: status}
}

func (t task) String() string {
	return fmt.Sprintf("[---\nID:\t%d\nTask\t:%s\nDue\t:%s\nDescription\t:%s\nStatus\t:%s\n---]\n", t.id, t.tname, t.due, t.desc, t.status)
}

//functions

func usage() {

	// this is a helper function to show the usage of the cli command
	fmt.Println("Usage:\n1 - Show tasks\n2 - Add task\n3 - Edit task\n4 - Delete task\nh - Show help\nq - Quit")

}

func main() {

	// an array to store tasks in memory
	// take input to get the operation : add delete edit
	// file based database

	reader := bufio.NewReader(os.Stdin)

	var op string
	var taskList []task

	//	var tasks = make([]string,0)
	fmt.Println("Launching task CLI v.0.0.1")
	for {

		fmt.Print(">")
		op, _ = reader.ReadString('\n')
		op = strings.TrimSpace(op)

		switch op {

		case "1":
			fmt.Println("Showing tasks")
			fmt.Println(taskList)
		case "2":

			var taskNameInput string
			var taskDescInput string
			var taskDueInput string

			fmt.Println("Adding task. Enter details")
			fmt.Print("Task\t:")
			taskNameInput, _ = reader.ReadString('\n')
			taskNameInput = strings.TrimSpace(taskNameInput)

			fmt.Print("Description\t:")
			taskDescInput, _ = reader.ReadString('\n')
			taskDescInput = strings.TrimSpace(taskDescInput)

			fmt.Print("Due (DD/MM/YYYY)\t:")
			taskDueInput, _ = reader.ReadString('\n')
			taskDueInput = strings.TrimSpace(taskDueInput)

			pdate, err := time.Parse("02/01/2006", taskDueInput)
			if err != nil {
				fmt.Println("error", err)
				continue
			}
			var beginState taskState = todo
			var comments = make([]string, 0)
			t := createTask(taskNameInput, pdate, taskDescInput, comments, beginState)
			taskList = append(taskList, t)
			fmt.Println("Added task")

		case "3":
			fmt.Println("Editing task")
		case "4":
			fmt.Println("Deleting task")
		case "h":
			usage()
		case "q":
			{
				fmt.Println("Exit")
				return
			}
		default:
			{
				fmt.Println("Unknown output")
			}
		}

	}

}
