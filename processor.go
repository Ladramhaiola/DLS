package dls

import (
	"fmt"
	"strings"
)

type Processor struct {
	ID        int
	scheduled []*Task
}

func (p *Processor) Earliest() int {
	if len(p.scheduled) == 0 {
		return 0
	}
	last := p.scheduled[len(p.scheduled)-1]
	return last.start + last.Cost
}

// schedule task at specified time
func (p *Processor) Insert(at int, task *Task) {
	task.start = at
	task.processor = p
	p.scheduled = append(p.scheduled, task)
}

func (p *Processor) String() string {
	var res strings.Builder
	res.WriteString(fmt.Sprintf("P%d {\n", p.ID))
	for _, task := range p.scheduled {
		res.WriteString("\t" + task.String() + "\n")
	}
	res.WriteString("}")
	return res.String()
}
