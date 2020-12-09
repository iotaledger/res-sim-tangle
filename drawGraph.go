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

	for i := 0; i < len(sim.tangle); i++ {
		keys = append(keys, sim.tangle[i].id)
	}
	// for key := range sim.approvers {
	// 	keys = append(keys, key)
	// }
	// for _, tx := range sim.tips {
	// 	keys = append(keys, tx)
	// }
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
		for _, i := range unique(sim.tangle[tx].app) {
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

//drawNthApprover draws the Tangle highlighting path of random walker transitioning to first or last approver
func drawNthApprover(sim *Sim, nodeMap map[int]int, nApprover int, G *graphviz.Graph) {
	visitMap := make(map[int]int)

	var current int
	if nApprover > 0 { //first approver
		for current = 0; len(sim.tangle[current].app) > 0; current = sim.tangle[current].app[0] {
			visitMap[current]++
		}
	} else { //last approver
		for current = 0; len(sim.tangle[current].app) > 0; current = sim.tangle[current].app[len(sim.tangle[current].app)-1] {
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
