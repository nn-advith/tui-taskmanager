package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nn-advith/cli-taskmanager/task"
)

//functions

func usage() {

	// this is a helper function to show the usage of the cli command
	fmt.Println("Usage:\n1 - Show tasks\n2 - Add task\n3 - Edit task\n4 - Delete task\nh - Show help\nq - Quit")

}

func editHelper() {
	fmt.Println("Usage:\n1 - Update task name\n2 - Update task description\n3 - Update status\n4 - Update due date\n5 - Add comment")
}

func main() {

	// an array to store tasks in memory
	// take input to get the operation : add delete edit
	// file based database

	reader := bufio.NewReader(os.Stdin)

	var op string
	var taskList []task.Task

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

			var comments = make([]string, 0)
			t := task.CreateTask(taskNameInput, pdate, taskDescInput, comments)
			taskList = append(taskList, t)
			fmt.Println("Added task")

		case "3":
			fmt.Println("edit> Editing task. Enter ID of task you want to edit:")
			fmt.Print("edit> ")
			taskIdStr, _ := reader.ReadString('\n')
			taskIdStr = strings.TrimSpace(taskIdStr)
			taskId, _ := strconv.Atoi(taskIdStr)
			fmt.Println(taskId)

			var tempTask *task.Task
			for i := range taskList {
				if taskList[i].ReturnID() == taskId {
					tempTask = &taskList[i]
				}
			}

			fmt.Println("edit> What do you want to edit?")
			editHelper()
			fmt.Print("edit> ")

			eop, _ := reader.ReadString('\n')
			eop = strings.TrimSpace(eop)

			switch eop {
			case "1":
				fmt.Println("edit> Enter task name:")
				newname, _ := reader.ReadString('\n')
				newname = strings.TrimSpace(newname)
				tempTask.UpdateTaskName(newname)
				fmt.Println("edit> Update done")
			case "2":
				fmt.Println("edit> Enter task description:")
				newdesc, _ := reader.ReadString('\n')
				newdesc = strings.TrimSpace(newdesc)
				tempTask.UpdateTaskDesc(newdesc)
				fmt.Println("edit> Update done")
			case "3":
				fmt.Println("edit> Updating status to ", tempTask.GetNextState())
				tempTask.UpdateStatus()
				fmt.Println("edit> Update done")
			case "4":
				fmt.Println("edit> Enter due date:")
				newdue, _ := reader.ReadString('\n')
				newdue = strings.TrimSpace(newdue)
				pdate, err := time.Parse("02/01/2006", newdue)
				if err != nil {
					fmt.Println("error", err)
					continue
				}
				tempTask.UpdateTaskDue(pdate)
				fmt.Println("edit> Update done")
			case "5":
				fmt.Println("edit> Enter comment:")
				comment, _ := reader.ReadString('\n')
				comment = strings.TrimSpace(comment)
				tempTask.AddComment(comment)
				fmt.Println("edit> Update done")
			default:
				fmt.Println("edit> Unknown input")

			}

		case "4":
			fmt.Println("delete> Enter ID of task to delete")
			fmt.Print("delete >")
			delId, _ := reader.ReadString('\n')
			delId = strings.TrimSpace(delId)
			delIdInt, _ := strconv.Atoi(delId)

			{
				var updatedSlice []task.Task

				for i := range taskList {
					if delIdInt != taskList[i].ReturnID() {
						updatedSlice = append(updatedSlice, taskList[i])
					}
				}

				taskList = updatedSlice
			}
			fmt.Println("delete> Done")
		case "h":
			usage()
		case "q":
			{
				fmt.Println("Exit")
				return
			}
		case "":
			continue
		default:
			{
				fmt.Println("Unknown output")
			}
		}

	}

}
