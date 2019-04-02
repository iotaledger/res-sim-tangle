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
	RandomWalkStep(int, *Sim) (int, int)
	RandomWalkStepBack(int, *Sim) int
	RandomWalkStepInfinity(int, *Sim) (int, int)
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

	AppID := 0
	for i := 0; i < sim.param.K; i++ {
		//URTS with repetition
		if len(sim.tips) == 0 {
			fmt.Println("ERROR: No tips left")
			os.Exit(0)
		}
		j := rand.Intn(len(sim.tips))

		uniqueTip := true
		if sim.param.SingleEdgeEnabled { // consider only SingleEdge Model
			for i1 := 0; i1 < min2Int(i, len(tipsApproved)); i1++ {
				if tipsApproved[i1] == sim.tips[j] {
					uniqueTip = false
				}
			}
		}
		if uniqueTip {
			tipsApproved[AppID] = sim.tips[j] //sim.tangle[sim.vTips[j]].id
			AppID++
		} else { // tip already existed and we are in the SingleEdge Model
			tipsApproved = tipsApproved[:len(tipsApproved)-1]
		}

		//tipsApproved = append(tipsApproved, sim.tangle[sim.vTips[j]].id)
		// tipsApproved[i] = sim.tips[j]
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
func (URW) RandomWalkStep(t int, sim *Sim) (int, int) {
	directApprovers := sim.approvers[t]
	if (len(directApprovers)) == 0 {
		return t, -1
	}
	if (len(directApprovers)) == 1 {
		return directApprovers[0], 0
	}
	j := rand.Intn(len(directApprovers))
	return directApprovers[j], j
}

//RandomWalkStepInfinity returns the choosen tip and its index position when walking over the spine Tangle
func (URW) RandomWalkStepInfinity(t int, sim *Sim) (int, int) {
	directApprovers := sim.spineApprovers[t]
	if (len(directApprovers)) == 0 {
		return t, -1
	}
	if (len(directApprovers)) == 1 {
		return directApprovers[0], 0
	}
	j := rand.Intn(len(directApprovers))
	return directApprovers[j], j
}

//RandomWalkStepBack returns the choosen tx
func (URW) RandomWalkStepBack(t int, sim *Sim) int {
	refs := sim.tangle[t].ref
	if len(refs) == 0 {
		return t
	}
	if (len(refs)) == 1 {
		return refs[0]
	}
	j := rand.Intn(len(refs))
	return refs[j]
}

// ??? this should be called RandoWalkStep, otherwise it is confusing
// ??? we already forward the pointer to the sim, should we not just forward the ID rather than the whole struct?
// ??? Currently each time the RW is called this creates a copy of a tx, while the copy of the int ID is enough
//RandomWalk returns the chosen tip and its index position
func (BRW) RandomWalkStep(t int, sim *Sim) (choosenTip int, approverIndx int) {
	//defer sim.b.track(runningtime("BRW"))
	directApprovers := sim.approvers[t]
	if (len(directApprovers)) == 0 {
		return t, -1
	}
	if (len(directApprovers)) == 1 {
		return directApprovers[0], 0
	}

	nw := normalizeWeights(directApprovers, sim)
	tip, j := weightedChoose(directApprovers, nw, sim.generator, sim.b)
	return tip, j
}

//RandomWalkStepInfinity returns the chosen tip and its index position when walking over the spine Tangle
func (BRW) RandomWalkStepInfinity(t int, sim *Sim) (choosenTip int, approverIndx int) {
	//defer sim.b.track(runningtime("BRW"))
	directApprovers := sim.spineApprovers[t]
	if (len(directApprovers)) == 0 {
		return t, -1
	}
	if (len(directApprovers)) == 1 {
		return directApprovers[0], 0
	}

	nw := normalizeWeights(directApprovers, sim)
	tip, j := weightedChoose(directApprovers, nw, sim.generator, sim.b)
	return tip, j
}

//RandomWalkStepBack returns the chosen tx
func (BRW) RandomWalkStepBack(t int, sim *Sim) (choosenTip int) {
	//defer sim.b.track(runningtime("BRW"))
	refs := sim.tangle[t].ref
	if (len(refs)) == 0 {
		return t
	}
	if (len(refs)) == 1 {
		return refs[0]
	}

	nw := normalizeWeights(refs, sim)
	tip, _ := weightedChoose(refs, nw, sim.generator, sim.b)
	return tip
}

func randomWalk(tsa RandomWalker, t Tx, sim *Sim) []int {
	defer sim.b.track(runningtime("RW"))
	tipsApproved := make([]int, sim.param.K)
	//cache := make(map[int][]float64)

	AppID := 0
	for i := 0; i < sim.param.K; i++ {
		//URTS with repetition  //??? this seems the wrong comment here
		var current int
		for current, _ = tsa.RandomWalkStep(0, sim); len(sim.approvers[current]) > 0; current, _ = tsa.RandomWalkStep(current, sim) {
		}

		uniqueTip := true
		if sim.param.SingleEdgeEnabled { // consider only SingleEdge Model
			for i1 := 0; i1 < min2Int(i, len(tipsApproved)); i1++ {
				if tipsApproved[i1] == current {
					uniqueTip = false
				}
			}
		}
		if uniqueTip {
			tipsApproved[AppID] = current //sim.tangle[sim.vTips[j]].id
			AppID++
		} else { // tip already existed and we are in the SingleEdge Model
			tipsApproved = tipsApproved[:len(tipsApproved)-1]
		}

		//tipsApproved = append(tipsApproved, sim.tangle[sim.vTips[j]].id)
		if sim.tangle[current].firstApproval < 0 {
			sim.tangle[current].firstApproval = t.time
		}
	}
	return tipsApproved
}
