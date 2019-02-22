package main

import (
	"fmt"
	"math/rand"
	"os"
)

// TipSelector defines the interface for a TSA
type TipSelector interface {
	TipSelect(Tx, *Sim) []int
}

// RandomWalker defines the interface for a random walk
type RandomWalker interface {
	RandomWalk(Tx, *Sim) (Tx, int)
	RandomWalkBack(Tx, *Sim) Tx
}

// URTS defines a concrete type of a TSA
type URTS struct {
	TipSelector
}

// URW defines a concrete type of a TSA
type URW struct {
	TipSelector
	RandomWalker
}

// BRW defines a concrete type of a TSA
type BRW struct {
	TipSelector
	RandomWalker
}

// TipSelect selects k tips
func (URTS) TipSelect(t Tx, sim *Sim) []int {
	tipsApproved := make([]int, sim.param.K)
	//var j int

	for i := 0; i < sim.param.K; i++ {
		//URTS with repetition

		if len(sim.tips) == 0 {
			fmt.Println("ERROR: No tips left")
			os.Exit(0)
		}
		j := rand.Intn(len(sim.tips))

		//tipsApproved = append(tipsApproved, sim.tangle[sim.vTips[j]].id)
		tipsApproved[i] = sim.tips[j]
		if sim.tangle[sim.tips[j]].firstApproval < 0 {
			sim.tangle[sim.tips[j]].firstApproval = t.time
		}
	}
	return tipsApproved
}

// TipSelect selects k tips
func (tsa URW) TipSelect(t Tx, sim *Sim) []int {
	return randomWalk(tsa, t, sim)
}

// TipSelect selects k tips
func (tsa BRW) TipSelect(t Tx, sim *Sim) []int {
	return randomWalk(tsa, t, sim)
}

//RandomWalk returns the choosen tip and its index position
func (URW) RandomWalk(t Tx, sim *Sim) (Tx, int) {
	directApprovers := sim.approvers[t.id]
	if (len(directApprovers)) == 0 {
		return t, -1
	}
	if (len(directApprovers)) == 1 {
		return sim.tangle[directApprovers[0]], 0
	}
	j := rand.Intn(len(directApprovers))
	return sim.tangle[directApprovers[j]], j
}

//RandomWalkBack returns the choosen tx
func (URW) RandomWalkBack(t Tx, sim *Sim) Tx {
	refs := t.ref
	if len(refs) == 0 {
		return t
	}
	if (len(refs)) == 1 {
		return sim.tangle[refs[0]]
	}
	j := rand.Intn(len(refs))
	return sim.tangle[refs[j]]
}

//RandomWalk returns the chosen tip and its index position
func (BRW) RandomWalk(t Tx, sim *Sim) (choosenTip Tx, approverIndx int) {
	//defer sim.b.track(runningtime("BRW"))
	directApprovers := sim.approvers[t.id]
	if (len(directApprovers)) == 0 {
		return t, -1
	}
	if (len(directApprovers)) == 1 {
		return sim.tangle[directApprovers[0]], 0
	}
	// nw, ok := cache[t.id]
	// if !ok {
	// 	nw = normalizeWeights(directApprovers, sim)
	// 	cache[t.id] = nw
	// } else {
	// 	//fmt.Println("HIT")
	// }

	nw := normalizeWeights(directApprovers, sim)
	tip, j := weightedChoose(directApprovers, nw, sim.generator, sim.b)
	return sim.tangle[tip], j
}

//RandomWalkBack returns the chosen tx
func (BRW) RandomWalkBack(t Tx, sim *Sim) (choosenTip Tx) {
	//defer sim.b.track(runningtime("BRW"))
	refs := t.ref
	if (len(refs)) == 0 {
		return t
	}
	if (len(refs)) == 1 {
		return sim.tangle[t.ref[0]]
	}

	nw := normalizeWeights(refs, sim)
	tip, _ := weightedChoose(refs, nw, sim.generator, sim.b)
	return sim.tangle[tip]
}

func randomWalk(tsa RandomWalker, t Tx, sim *Sim) []int {
	defer sim.b.track(runningtime("RW"))
	tipsApproved := make([]int, sim.param.K)
	//cache := make(map[int][]float64)

	for i := 0; i < sim.param.K; i++ {
		//URTS with repetition  //??? this seems the wrong comment here
		var current Tx
		for current, _ = tsa.RandomWalk(sim.tangle[0], sim); len(sim.approvers[current.id]) > 0; current, _ = tsa.RandomWalk(current, sim) {
		}

		//tipsApproved = append(tipsApproved, sim.tangle[sim.vTips[j]].id)
		tipsApproved[i] = current.id //sim.tangle[sim.vTips[j]].id
		if sim.tangle[current.id].firstApproval < 0 {
			sim.tangle[current.id].firstApproval = t.time
		}
	}
	return tipsApproved
}
