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
	j := rand.Intn(len(directApprovers))
	return sim.tangle[directApprovers[j]], j

}

//RandomWalk returns the choosen tip and its index position
func (BRW) RandomWalk(t Tx, sim *Sim) (choosenTip Tx, approverIndx int) {
	//defer sim.b.track(runningtime("BRW"))
	directApprovers := sim.approvers[t.id]
	if (len(directApprovers)) == 0 {
		return t, -1
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

func randomWalk(tsa RandomWalker, t Tx, sim *Sim) []int {
	defer sim.b.track(runningtime("RW"))
	tipsApproved := make([]int, sim.param.K)
	//cache := make(map[int][]float64)

	for i := 0; i < sim.param.K; i++ {
		//URTS with repetition
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

// //remove (visible)approved transactions from tips set. Return []tx of approved tips
// func (tsa URTS) TipsUpdate(t tx, sim *Sim) []int {
// 	//defer track(runningtime("udpateTips"))
// 	sim.removeOldTips(t)
// 	return tsa.revealTips(t, sim)
// }

// //remove (visible)approved transactions from tips set. Return []tx of approved tips
// func (tsa URW) TipsUpdate(t tx, sim *Sim) []int {
// 	//defer track(runningtime("udpateTips"))
// 	sim.removeOldTips(t)
// 	return tsa.revealTips(t, sim)
// }

// //remove (visible) approved transactions from tips set. Return []tx of approved tips
// func (tsa BRW) TipsUpdate(t tx, sim *Sim) []int {
// 	//defer track(runningtime("udpateTips"))
// 	sim.removeOldTips(t)
// 	return sim.revealTips(t)
// }

// //return set of new visible tips from the hidden tips set and update hidden tips set
// func (URW) revealTips(t tx, sim *Sim) []int {
// 	//var i, tip int
// 	var newTips []int
// 	//fmt.Println("HiddenTips", sim.hiddenTips)
// 	for i, tip := range sim.hiddenTips { //iterate hidden tips set until it finds the first NOT visible tip

// 		if sim.tangle[tip].isVisible(t.time, sim.param.H) {
// 			sim.approvers = updateApprovers(sim.approvers, sim.tangle[tip])
// 		} else {
// 			newTips = sim.hiddenTips[:i]        //set of "new" visible tips
// 			sim.hiddenTips = sim.hiddenTips[i:] //update hidden tips set
// 			return newTips
// 		}
// 	}

// 	//All the hidden tips are visible
// 	newTips = sim.hiddenTips
// 	sim.hiddenTips = sim.hiddenTips[:0]
// 	return newTips
// }

// func (URTS) revealTips(t tx, sim *Sim) []int {
// 	//var i, tip int
// 	var newTips []int
// 	//fmt.Println("HiddenTips", sim.hiddenTips)
// 	for i, tip := range sim.hiddenTips { //iterate hidden tips set until it finds the first NOT visible tip

// 		if sim.tangle[tip].isVisible(t.time, sim.param.H) {
// 			sim.approvers = updateApprovers(sim.approvers, sim.tangle[tip])
// 		} else {
// 			newTips = sim.hiddenTips[:i]        //set of "new" visible tips
// 			sim.hiddenTips = sim.hiddenTips[i:] //update hidden tips set
// 			return newTips
// 		}
// 	}

// 	//All the hidden tips are visible
// 	newTips = sim.hiddenTips
// 	sim.hiddenTips = sim.hiddenTips[:0]
// 	return newTips
// }

// func (BRW) revealTips(t tx, sim *Sim) []int {
// 	//var i, tip int
// 	var newTips []int
// 	//fmt.Println("HiddenTips", sim.hiddenTips)

// 	if len(sim.hiddenTips) == 0 {
// 		return newTips
// 	}

// 	i := 0
// 	tip := sim.hiddenTips[i]
// 	for sim.tangle[tip].isVisible(t.time, sim.param.H) {
// 		sim.approvers = updateApprovers(sim.approvers, sim.tangle[tip])
// 		if sim.param.TSA == "rw" {
// 			sim.updateCW(sim.tangle[tip])
// 			//sim.updateCWDFS(sim.tangle[tip])
// 		}
// 		i++
// 		if i >= len(sim.hiddenTips) {
// 			break
// 		}
// 		tip = sim.hiddenTips[i]
// 	}

// 	newTips = sim.hiddenTips[:i]        //set of "new" visible tips
// 	sim.hiddenTips = sim.hiddenTips[i:] //update hidden tips set
// 	return newTips
// }
