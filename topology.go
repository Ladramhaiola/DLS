package dls

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"math"
	"sort"
)

type Topology interface {
	// find most suitable processor for given task (in dls that gives the lowest EEST)
	BestProcessor(task *Task) (*Processor, int)
	Processors() []*Processor

	Send(transfer *Transfer)

	// make all transfers required for given task
	SendRequired(proc *Processor, task *Task)

	// get underlying data source
	Src() interface{}
}

type StarBridge struct {
	bridge     *Bridge
	processors []*Processor
	graph      *TaskGraph
}

func (sb *StarBridge) BestProcessor(task *Task) (*Processor, int) {
	est := math.MaxInt32
	processor := sb.processors[0]

	for _, proc := range sb.processors {
		latest := sb.bridge.DataReceived(sb.graph, proc, task)

		if proc.Earliest() > latest {
			latest = proc.Earliest()
		}

		if len(task.Parents) > 0 {
			for _, parent := range task.Parents {
				if parent.start+parent.Cost > latest {
					latest = parent.start + parent.Cost
				}
			}
		}

		if latest < est {
			est = latest
			processor = proc
		}
	}

	return processor, est
}

func (sb *StarBridge) Send(transfer *Transfer) {
	sb.bridge.SetBusy(transfer)
}

func (sb *StarBridge) SendRequired(proc *Processor, task *Task) {
	sort.Slice(task.Parents, func(i, j int) bool {
		return task.Parents[i].start < task.Parents[j].start
	})

	for _, parent := range task.Parents {
		if parent.processor != proc {
			cost := sb.graph.CommCost(parent.ID, task.ID)
			start := sb.bridge.Earliest(parent.start+parent.Cost, cost)

			sb.Send(&Transfer{S: start, E: start + cost, Cost: cost, Src: parent, Dst: task})
		}
	}
}

func (sb *StarBridge) Processors() []*Processor {
	return sb.processors
}

func (sb *StarBridge) Src() interface{} {
	return sb.bridge.busyAt
}

func NewStarBridge(taskGraph *TaskGraph, maxProcCnt int) *StarBridge {
	processors := make([]*Processor, maxProcCnt)
	for i := 0; i < maxProcCnt; i++ {
		processors[i] = &Processor{ID: i}
	}

	return &StarBridge{
		bridge: &Bridge{
			busyAt: rbt.NewWithIntComparator(),
		},
		processors: processors,
		graph:      taskGraph,
	}
}
