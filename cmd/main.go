package main

import "dls"

func main() {
	taskGraph := dls.Graph()
	taskGraph.AddTasks(
		[2]int{1, 3},
		[2]int{2, 5},
		[2]int{3, 4},
		[2]int{4, 6},
		[2]int{5, 2},
		[2]int{6, 2},
		[2]int{7, 4},
		[2]int{8, 3},
		[2]int{9, 6},
		[2]int{10, 2},
	)
	taskGraph.AddConns(
		[3]int{1, 3, 8},
		[3]int{1, 4, 7},
		[3]int{1, 6, 6},

		[3]int{2, 3, 10},
		[3]int{2, 4, 14},
		[3]int{2, 5, 10},

		[3]int{3, 6, 5},
		[3]int{3, 8, 8},

		[3]int{4, 7, 12},
		[3]int{4, 9, 7},

		[3]int{5, 6, 4},
		[3]int{5, 7, 10},
		[3]int{5, 10, 14},

		[3]int{6, 8, 4},
		[3]int{6, 10, 8},

		[3]int{7, 9, 5},
	)

	topology := dls.NewStarBridge(taskGraph, 4)
	dls.DLS(taskGraph, topology)
}
