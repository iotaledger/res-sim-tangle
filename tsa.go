package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

// TipSelector defines the interface for a TSA
type TipSelector interface {
	TipSelect(Tx, *Sim) []int
}

// URTS implements the uniform random tip selection algorithm
type URTS struct {
	TipSelector
}

// RURTS implements the restricted uniform random tip selection algorithm, where txs are only valid tips up to some age D
type RURTS struct {
	TipSelector
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

		sim.updateFirstApproval(sim.tips[j], t)

	}
	return tipsApproved
}

// TipSelect selects k tips
func (RURTS) TipSelect(t Tx, sim *Sim) []int {
	kNow := sim.param.K
	if sim.param.responseSpamTipsEnabled {
		kNow = increaseK(sim)
	}

	tipsApproved := make([]int, kNow)
	//var j int

	AppID := 0
	for i := 0; i < kNow; i++ {
		//RURTS with repetition
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

		sim.updateFirstApproval(sim.tips[j], t)
	}
	return tipsApproved
}

// dynamically increase K
func increaseK(sim *Sim) int {
	if len(sim.tips) > sim.param.acceptableNumberTips {
		delta := math.Max(0, float64(len(sim.tips)-sim.param.acceptableNumberTips)/(2.*sim.param.Lambda))
		Know := sim.param.K + int(delta*delta*sim.param.responseKIncrease)
		return int(math.Min(float64(sim.param.maxK), float64(Know)))
	}
	return sim.param.K
}

func (sim *Sim) updateFirstApproval(tip int, t Tx) {
	// update first approval
	if sim.tangle[tip].firstApprovalTime < 0 {
		sim.tangle[tip].firstApprovalTime = t.time
	}
	// the current child tx can have a smaller h so we need to check if it gets revealed earlier than other approvers
	if sim.tangle[tip].firstVisibleApprovalTime > t.time+float64(t.h) || sim.tangle[tip].firstVisibleApprovalTime < 0 {
		sim.tangle[tip].firstVisibleApprovalTime = t.time + float64(t.h)
	}

}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// - - - - - - - - - - - - Adversary - - - - - - - - - - - -
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

// TipSelectorAdversary defines the interface for a TSAAdversary
type TipSelectorAdversary interface {
	TipSelectAdversary(Tx, *Sim) []int
}

// SpamGenesis implements a tip selection where txs only attach to the genesis
type SpamGenesis struct {
	TipSelectorAdversary
}

// tx selects the genesis
func (SpamGenesis) TipSelectAdversary(t Tx, sim *Sim) []int {
	tipsApproved := make([]int, 1)
	//var j int
	tipsApproved[0] = 0 // use genesis as only parent

	return tipsApproved
}
