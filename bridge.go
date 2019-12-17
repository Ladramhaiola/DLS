package dls

import (
	"sort"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

type Bridge struct {
	busyAt *rbt.Tree
}

type Transfer struct {
	Src *Task
	Dst *Task

	// start, end, transaction cost
	S, E, Cost int
}

// Earliest free interval
func (b *Bridge) Earliest(start, duration int) int {
	it := b.busyAt.Iterator()
	for it.Next() {
		transfer := it.Value().(*Transfer)
		if transfer.S-start >= duration {
			return start
		}
		if transfer.E > start {
			start = transfer.E
		}
	}
	return start
}

// load bridge with given data transfer
func (b *Bridge) SetBusy(transfer *Transfer) {
	b.busyAt.Put(transfer.S, transfer)
}

// time when all required data sent to this task
func (b *Bridge) DataReceived(graph *TaskGraph, proc *Processor, task *Task) (latest int) {
	changes := make([]int, 0)

	sort.Slice(task.Parents, func(i, j int) bool {
		return task.Parents[i].start < task.Parents[j].start
	})

	for _, parent := range task.Parents {
		if parent.processor != proc {
			cost := graph.CommCost(parent.ID, task.ID)
			start := b.Earliest(parent.start+parent.Cost, cost)

			if start+cost > latest {
				latest = start + cost
			}

			b.SetBusy(&Transfer{S: start, E: start + cost, Cost: cost, Src: parent, Dst: task})
			changes = append(changes, start)
		}
	}

	for _, change := range changes {
		b.busyAt.Remove(change)
	}
	return
}
