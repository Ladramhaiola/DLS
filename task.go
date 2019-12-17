package dls

import "fmt"

type TaskState int

const (
	Unexamined TaskState = iota
	Examined
)

type Task struct {
	ID int

	SL   int
	DL   int
	Cost int

	// start time on processor
	start int

	Parents  []*Task
	Children []*Task

	marked bool
	state  TaskState

	processor *Processor
}

func (t *Task) String() string {
	return fmt.Sprintf("T%d {start:%d cost:%d}", t.ID, t.start, t.Cost)
}
