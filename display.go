package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/capossele/GoGraphviz/graphviz"
)

//func printCWRef(a map[int][]uint64) {
func printCWRef(a [][]uint64) {
	for i := 0; i < len(a); i++ {
		fmt.Printf("%d: ", i)
		for j := len(a[i]) - 1; j >= 0; j-- {
			fmt.Printf("\t%064b\n", a[i][j])
		}
		fmt.Println()
	}
}

func printApprovers(a map[int][]int) {
	var keys []int
	for key := range a {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, t := range keys {
		fmt.Println(t, "<-", unique(a[t]))
	}
}

func saveApprovers(a map[int][]int) {
	var keys []int
	for key := range a {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	f, err := os.Create("data/approvers.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
	}
	defer f.Close()

	for _, t := range keys {
		_, err = f.WriteString(fmt.Sprintln(t, "<-", unique(a[t]))) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
}

func saveTangle(tangle []Tx) {
	f, err := os.Create("data/tangle.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
	}
	defer f.Close()

	for _, t := range tangle {
		_, err = f.WriteString(fmt.Sprintln(t)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
}

func printCW(sim Sim) {
	for _, t := range sim.tangle {
		fmt.Println(t.id, t.cw)
	}
}

func printTips(a map[int]bool) {
	var keys []int
	for key := range a {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	fmt.Println(keys)
}

func (sim *Sim) visualizeTangle() {
	G := createTangleGraph(0, sim)
	G.GenerateDOT(os.Stdout)
}

func createTangleGraph(tx int, sim *Sim) *graphviz.Graph {
	G := &graphviz.Graph{}
	visited := make(map[int]bool)
	rootNode := make(map[int]int)
	addSubGraph(tx, sim, visited, rootNode, G)
	G.DefaultNodeAttribute(graphviz.Shape, graphviz.ShapeCircle)
	G.GraphAttribute(graphviz.NodeSep, "0.3")
	G.MakeDirected()
	return G
}

func addSubGraph(tx int, sim *Sim, visited map[int]bool, rootNode map[int]int, G *graphviz.Graph) int {
	if _, ok := visited[tx]; !ok {
		//add new node if does not exist yet
		rootNode[tx] = G.AddNode(fmt.Sprint(tx))
		visited[tx] = true
	}

	for _, i := range unique(sim.approvers[tx]) {
		node := addSubGraph(i, sim, visited, rootNode, G)
		G.AddEdge(node, rootNode[tx], "")
	}
	return rootNode[tx]
}
