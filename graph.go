package dls

import (
	"github.com/emirpasic/gods/sets/hashset"
)

type TaskGraph struct {
	nodes map[int]*Task
	edges map[[2]int]int
}

func (graph *TaskGraph) AddTasks(tasks ...[2]int) {
	for _, t := range tasks {
		graph.nodes[t[0]] = &Task{ID: t[0], Cost: t[1]}
	}
}

func (graph *TaskGraph) AddConns(conns ...[3]int) {
	for _, conn := range conns {
		src := graph.nodes[conn[0]]
		dst := graph.nodes[conn[1]]

		if src == nil || dst == nil {
			continue
		}

		key := [2]int{conn[0], conn[1]}
		graph.edges[key] = conn[2]

		src.Children = append(src.Children, dst)
		dst.Parents = append(dst.Parents, src)
	}
}

func (graph *TaskGraph) CommCost(src, dst int) int {
	if cost, ok := graph.edges[[2]int{src, dst}]; ok {
		return cost
	}
	return -1
}

// Topologically sorted graph nodes
func (graph *TaskGraph) TopSort() (L []*Task) {
	S := hashset.New()

	for _, task := range graph.nodes {
		task.marked = false
		if len(task.Parents) == 0 {
			S.Add(task)
		}
	}

	for S.Size() > 0 {
		n := S.Values()[0].(*Task)
		S.Remove(n)
		L = append(L, n)

		n.marked = true
		for _, child := range n.Children {
			if checkMarked(child) {
				S.Add(child)
			}
		}
	}

	return L
}

func StaticLevel(graph *TaskGraph) {
	for _, task := range graph.nodes {
		max := 0
		for _, parent := range task.Parents {
			s := parent.SL + graph.CommCost(parent.ID, task.ID) + parent.Cost
			if s > max {
				max = s
			}
		}
		task.SL = max
	}
}

func checkMarked(task *Task) bool {
	for _, parent := range task.Parents {
		if !parent.marked {
			return false
		}
	}
	return true
}

func Graph() *TaskGraph {
	return &TaskGraph{
		nodes: make(map[int]*Task),
		edges: make(map[[2]int]int),
	}
}
