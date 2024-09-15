package task

import (
	"fmt"
	"strings"
	"time"
)

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
		panic(fmt.Errorf("unkown task state: %s", ts))
	}
}

// task struct
type Task struct {
	id     int
	tname  string
	due    time.Time
	desc   string
	comm   []string
	status taskState
}

func CreateTask(tname string, due time.Time, desc string, comm []string) Task {
	keyCounter = keyCounter + 1
	var beginState taskState = todo
	return Task{id: keyCounter, tname: tname, due: due, desc: desc, comm: comm, status: beginState}
}

func (t Task) String() string {
	commentList := strings.Join(t.comm, "\n")
	return fmt.Sprintf("[\nID:\t%d\nTask\t:%s\nDue\t:%s\nDescription\t:%s\nStatus\t:%s\nComments\t:\n%s\n]\n", t.id, t.tname, t.due, t.desc, t.status, commentList)
}

func (t *Task) AddComment(c string) {
	t.comm = append(t.comm, c)
}

func (t *Task) UpdateTaskName(n string) {
	t.tname = n
}

func (t *Task) UpdateTaskDesc(n string) {
	t.desc = n
}

func (t *Task) UpdateTaskDue(n time.Time) {
	t.due = n
}

func (t *Task) UpdateStatus() {
	t.status = transition(t.status)
}

func (t Task) GetNextState() string {
	return fmt.Sprint(transition(t.status))
}

func (t Task) ReturnID() int {
	return t.id
}
