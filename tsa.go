package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/willf/bitset"
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

// HPS implements the "Heaviest Pair Selection" algorithm, where the tips are selected in such a way
// that the number of referenced transactions is maximized.
type HPS struct{}

// getReferences returns the set of all the transactions directly or indirectly referenced by t.
// The references are computed using recursion and dynamic programming.
func getReferences(t int, tangle []Tx, cache []*bitset.BitSet) *bitset.BitSet {
	if cache[t] != nil {
		return cache[t]
	}

	result := bitset.New(uint(t))
	for _, r := range tangle[t].ref {
		result.InPlaceUnion(getReferences(r, tangle, cache))
		result.Set(uint(r))
	}
	cache[t] = result
	return result
}

// heaviestPairs finds the tip pairs the reference the most transactions.
func heaviestPairs(sim *Sim) [][2]int {
	// cache the references of all the nodes
	cache := make([]*bitset.BitSet, len(sim.tangle))
	cache[0] = bitset.New(0) // genesis has no referenced txs

	var bestWeight uint
	var bestResults [][2]int

	// loop through all pairs of tips and find the pair with the most referenced txs.
	for _, t1 := range sim.tips {
		ref1 := getReferences(t1, sim.tangle, cache)

		for _, t2 := range sim.tips {
			if t1 >= t2 {
				continue // we don't care about the order in the pair
			}

			ref2 := getReferences(t2, sim.tangle, cache)

			// the weight are all the txs referenced by t1,t2 plus t1 and t2 themselves
			weight := ref1.UnionCardinality(ref2) + 2
			if weight > bestWeight {
				bestWeight = weight
				bestResults = [][2]int{[2]int{t1, t2}}
			} else if weight == bestWeight {
				bestResults = append(bestResults, [2]int{t1, t2})
			}
		}
	}

	return bestResults
}

func (HPS) TipSelect(t Tx, sim *Sim) (result []int) {
	if len(sim.tips) <= sim.param.K {
		result = append([]int{}, sim.tips...)
	} else {
		if sim.param.K != 2 {
			panic("Only 2 references supported.")
		}
		heaviest := heaviestPairs(sim)

		// select a random pair from the set of best pairs
		result = heaviest[rand.Intn(len(heaviest))][:]
	}

	// mark the approval time if needed
	for _, x := range result {
		if sim.tangle[x].firstApproval < 0 {
			sim.tangle[x].firstApproval = t.time
		}
	}
	return result
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
func (RURTS) TipSelect(t Tx, sim *Sim) []int {
	tipsApproved := make([]int, sim.param.K)
	//var j int

	AppID := 0
	for i := 0; i < sim.param.K; i++ {
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

		//tipsApproved = append(tipsApproved, sim.tangle[sim.vTips[j]].id)
		// tipsApproved[i] = sim.tips[j]
		if sim.tangle[sim.tips[j]].firstApproval < 0 {
			sim.tangle[sim.tips[j]].firstApproval = t.time
		}
	}
	return tipsApproved
}
