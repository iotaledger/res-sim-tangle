// Functions within this file uses Graphviz to define a graph.
// Output file is in dot format, thus, it is recommended to install Graphviz
// and use the dot command to produce a png/pdf of the graph.
// e.g., cat TangleGraph.dot | dot  -Tpng -o graph.png
//
// VS code has a nice plugin to preview a .dot file called Graphviz Preview

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/capossele/GoGraphviz/graphviz"
)

// visualizeTangle draws the Tangle Graph by generating a dot file.
// input: path, if != nil, is the path to highlight
// input: mode, is the type of graph:
//			1: simple Tangle with/without highlighed path
//			2: Ghost path, Ghost cone, Orphans + tips (TODO: clustering needs to be done manually)
//			3: Tangle with tx visiting probability in red gradients
//			4: Tangle with highlighted path of random walker transitioning to first approver
//			5: Tangle with highlighted path of random walker transitioning to last approver
func (sim *Sim) visualizeTangle(path map[int]int, mode int) {
	G := createTangleGraph(0, sim, path, mode)
	fmt.Println("\n")
	f, err := os.Create("graph/TangleGraph.dot")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		//return err
	}
	defer f.Close()
	G.GenerateDOT(f)
}

// createTangleGraph returns a graphviz graph based on the drawing mode selected.
// General properties of the graph can be modified here
func createTangleGraph(tx int, sim *Sim, path map[int]int, mode int) *graphviz.Graph {
	G := &graphviz.Graph{}
	nodeTxs := make(map[int]int)
	switch mode {
	case 1:
		drawTangle(sim, nodeTxs, path, "green", G)
	case 2:
		drawGhostConeOrphans(sim, nodeTxs, G)
	case 3:
		drawVisitingProbability(sim, nodeTxs, G)
	case 4:
		drawNthApprover(sim, nodeTxs, 1, G)
	case 5:
		drawNthApprover(sim, nodeTxs, -1, G)
	default:
		drawTangle(sim, nodeTxs, path, "green", G)
	}
	sortRankTransactons(sim, nodeTxs, G)
	G.DefaultNodeAttribute(graphviz.Shape, graphviz.ShapeCircle)
	G.GraphAttribute(graphviz.NodeSep, "0.3")
	if sim.param.horizontalOrientation {
		G.GraphAttribute("rankdir", "RL")
	}
	//G.GraphAttribute("bgcolor", "transparent")
	G.MakeDirected()
	return G
}

//drawTangle draws the Tangle with/without highlighing with a given color a given path
func drawTangle(sim *Sim, nodeMap, path map[int]int, color string, G *graphviz.Graph) {
	//sort txs and tips
	var keys []int
	for key := range sim.approvers {
		keys = append(keys, key)
	}
	for _, tx := range sim.tips {
		keys = append(keys, tx)
	}
	sort.Ints(keys)

	//all txs
	for _, tx := range keys {
		nodeMap[tx] = G.AddNode(fmt.Sprint(tx))
		if path != nil {
			if _, ok := path[tx]; ok {
				G.NodeAttribute(nodeMap[tx], "style", "filled")
				G.NodeAttribute(nodeMap[tx], "fillcolor", color)
			}
		}
	}

	//edges
	for _, tx := range keys {
		for _, i := range unique(sim.approvers[tx]) {
			edge := G.AddEdge(nodeMap[i], nodeMap[tx], "")
			G.EdgeAttribute(edge, "color", "#E0E0E0")
			if _, ok := path[tx]; ok {
				if _, ok := path[i]; ok {
					G.EdgeAttribute(edge, "color", color)
					G.EdgeAttribute(edge, "penwidth", "2")
					//G.EdgeAttribute(edge, "label", fmt.Sprintf("%.2f", sim.tangle[i].time-sim.tangle[tx].time))
				}
			}
		}
	}
}

//drawGhostConeOrphans draws the Tangle highlighting Ghost path, Ghost cone, Orphans + tips.
//TODO: clustering needs to be done manually
func drawGhostConeOrphans(sim *Sim, nodeMap map[int]int, G *graphviz.Graph) {
	//ghost path
	path, tip := ghostWalk(sim.tangle[0], sim)
	for _, tx := range path {
		nodeMap[tx] = G.AddNode(fmt.Sprint(tx))
		G.NodeAttribute(nodeMap[tx], "style", "filled")
		G.NodeAttribute(nodeMap[tx], "fillcolor", "red")
	}
	nodeMap[tip.id] = G.AddNode(fmt.Sprint(tip.id))
	G.NodeAttribute(nodeMap[tip.id], "style", "filled")
	G.NodeAttribute(nodeMap[tip.id], "fillcolor", "red")

	//ghost cone
	sim.computeSpine()
	var keys []int
	for key := range sim.spineApprovers {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, tx := range keys {
		if _, ok := nodeMap[tx]; !ok {
			nodeMap[tx] = G.AddNode(fmt.Sprint(tx))
			G.NodeAttribute(nodeMap[tx], "style", "filled")
			G.NodeAttribute(nodeMap[tx], "fillcolor", "turquoise")
		}
	}

	//orphans and tips
	keys = []int{}
	for key := range sim.approvers {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, tx := range keys {
		if _, ok := nodeMap[tx]; !ok {
			nodeMap[tx] = G.AddNode(fmt.Sprint(tx))
			//G.NodeAttribute(nodeMap[tx], "xlabel", fmt.Sprintf("%.3f", sim.tangle[tx].time))
		}
	}

	//edges
	for _, tx := range keys {
		for _, i := range unique(sim.approvers[tx]) {
			if _, ok := nodeMap[i]; !ok {
				nodeMap[i] = G.AddNode(fmt.Sprint(i))
			}
			//G.AddEdge(nodeMap[i], nodeMap[tx], fmt.Sprintf("%.3f", sim.tangle[i].time))
			G.AddEdge(nodeMap[i], nodeMap[tx], "")
		}
	}
}

//drawVisitingProbability draws the Tangle with tx visiting probability in red gradients
func drawVisitingProbability(sim *Sim, nodeMap map[int]int, G *graphviz.Graph) {
	visitMap := make(map[int]int)
	for i := 0; i < 100000; i++ {
		var current Tx
		//current := sim.tangle[0]
		//visitMap[current.id]++
		var tsa RandomWalker
		if sim.param.Alpha != 0 {
			tsa = BRW{}
		} else {
			tsa = URW{}
		}
		for current = sim.tangle[0]; len(sim.approvers[current.id]) > 0; current, _ = tsa.RandomWalk(current, sim) {
			visitMap[current.id]++
		}
		visitMap[current.id]++
	}

	//all txs
	keys := []int{}
	for key := range sim.approvers {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, tx := range keys {
		nodeMap[tx] = G.AddNode(fmt.Sprint(tx))
		G.NodeAttribute(nodeMap[tx], "style", "filled")
		c := int((255. * float64(visitMap[tx])) / float64(100000))
		c = 255 - c
		color := fmt.Sprintf("#FF%02x%02x", c, c)
		G.NodeAttribute(nodeMap[tx], "fillcolor", strings.ToUpper(color))
	}

	//edges
	for _, tx := range keys {
		for _, i := range unique(sim.approvers[tx]) {
			if _, ok := nodeMap[i]; !ok {
				nodeMap[i] = G.AddNode(fmt.Sprint(i))
				G.NodeAttribute(nodeMap[i], "style", "filled")
				c := int((255. * float64(visitMap[i])) / float64(100000))
				c = 255 - c
				color := fmt.Sprintf("#FF%02x%02x", c, c)
				G.NodeAttribute(nodeMap[i], "fillcolor", strings.ToUpper(color))
			}
			//G.AddEdge(nodeMap[i], nodeMap[tx], fmt.Sprintf("%.3f", sim.tangle[i].time))
			G.AddEdge(nodeMap[i], nodeMap[tx], "")
		}
	}
}

//drawNthApprover draws the Tangle highlighting path of random walker transitioning to first or last approver
func drawNthApprover(sim *Sim, nodeMap map[int]int, nApprover int, G *graphviz.Graph) {
	visitMap := make(map[int]int)

	var current int
	if nApprover > 0 { //first approver
		for current = 0; len(sim.approvers[current]) > 0; current = sim.approvers[current][0] {
			visitMap[current]++
		}
	} else { //last approver
		for current = 0; len(sim.approvers[current]) > 0; current = sim.approvers[current][len(sim.approvers[current])-1] {
			visitMap[current]++
		}
	}
	visitMap[current]++

	switch nApprover {
	case 1:
		drawTangle(sim, nodeMap, visitMap, "green", G)
	case -1:
		drawTangle(sim, nodeMap, visitMap, "red", G)
	}
}

//sortRankTransaction assigns the same rank to txs belonging to the same epoch h
func sortRankTransactons(sim *Sim, nodeMap map[int]int, G *graphviz.Graph) {
	endT := sim.param.TangleSize / int(sim.param.Lambda)
	for h := 1; h < endT; h++ {
		sameRank := []int{}
		for tx := range nodeMap {
			if sim.tangle[tx].time > float64(h-1) && sim.tangle[tx].time < float64(h) {
				sameRank = append(sameRank, tx)
			}
		}
		if len(sameRank) > 1 {
			for i := 1; i < len(sameRank); i++ {
				G.MakeSameRank(nodeMap[sameRank[0]], nodeMap[sameRank[i]])
			}
		}
	}
}

// visualizeRW performs 10 random walk and draws the Tangle at each step.
// Useful for generating GIF or video.
func (sim *Sim) visualizeRW() {
	for i := 0; i < 10; i++ {
		path := make(map[int]int)
		var current Tx
		var tsa RandomWalker
		if sim.param.Alpha != 0 {
			tsa = BRW{}
		} else {
			tsa = URW{}
		}
		for current = sim.tangle[0]; len(sim.approvers[current.id]) > 0; current, _ = tsa.RandomWalk(current, sim) {
			path[current.id]++
			sim.visualizeTangle(path, 1)
			time.Sleep(100 * time.Millisecond)
		}
		path[current.id]++
		sim.visualizeTangle(path, 1)
		time.Sleep(100 * time.Millisecond)
	}
}
