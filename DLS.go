package dls

import (
	"fmt"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"sort"
)

func DLS(graph *TaskGraph, topology Topology) {
	// Set initial SL for each node
	StaticLevel(graph)

	readyList := make([]*Task, 0)
	for _, task := range graph.nodes {
		if len(task.Parents) == 0 {
			readyList = append(readyList, task)
		}
	}

	for len(readyList) > 0 {
		for _, task := range readyList {
			proc, est := topology.BestProcessor(task)
			task.DL = task.SL - est
			task.processor = proc
		}

		// sort ready tasks by dynamic level
		sort.Slice(readyList, func(i, j int) bool {
			return readyList[i].DL > readyList[j].DL
		})

		task := readyList[0]
		readyList = readyList[1:]

		topology.SendRequired(task.processor, task)
		task.processor.Insert(task.SL-task.DL, task)

		task.state = Examined

		for _, task := range graph.nodes {
			if Ready(task) && task.state == Unexamined && !Contains(readyList, task) {
				readyList = append(readyList, task)
			}
		}
	}

	// todo: normal output
	for _, proc := range topology.Processors() {
		fmt.Println(proc)
	}

	tree := topology.Src().(*rbt.Tree)
	it := tree.Iterator()
	for it.Next() {
		t := it.Value().(*Transfer)
		fmt.Printf("%+v\n", t)
	}
}

// ready checks if all task parents are scheduled
func Ready(task *Task) bool {
	for _, parent := range task.Parents {
		if parent.state != Examined {
			return false
		}
	}
	return true
}

func Contains(a []*Task, b *Task) bool {
	for _, it := range a {
		if it.ID == b.ID {
			return true
		}
	}
	return false
}
