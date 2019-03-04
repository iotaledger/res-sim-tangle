package main

import "fmt"

type pOrphanResult struct {
	op  float64 // orphanage probability - probability of tx left behind
	top float64 // tip orphanage probability - probability of tips left behind
}

func (sim *Sim) runOrphaningP(result *pOrphanResult) {
	// remove txs grater than maxCut from both tangle and spine tangle so to have a comparable cone
	newTangle := sim.tangle[:sim.param.maxCut]
	// calculate spine Tangle up to maxCut (all txs (directly/indirectly) referenced by a GHOST particle (alpha = infinity)
	newSpineTangle := sliceSpineTangle(sim, sim.param.maxCut)

	// calculate op
	result.op = 1. - float64(len(newSpineTangle))/float64(len(newTangle))

	// calculate top by finding all tips left behind and dividing that number over all txs
	for tx := range newTangle {
		if len(sim.approvers[tx]) == 0 {
			result.top++
		}
	}
	result.top /= float64(len(newTangle))
}

func (a pOrphanResult) Join(b pOrphanResult) pOrphanResult {
	if a.op == 0 {
		return b
	}
	var result pOrphanResult
	result.op = (a.op + b.op) / 2.
	result.top = (a.top + b.top) / 2.
	return result
}

func (a pOrphanResult) String() string {
	result := fmt.Sprintln("Orphanage Probability:", a.op)
	result += fmt.Sprintln("Tip Orphanage Probability:", a.top)
	return result
}

func newPOrphanResult() *pOrphanResult {
	// variables initialization for pOprhan
	var result pOrphanResult
	return &result
}

func sliceSpineTangle(sim *Sim, uBound int) map[int]Tx {
	result := make(map[int]Tx)
	for k, v := range sim.spineTangle {
		if k < uBound {
			result[k] = v
		}
	}
	return result
}
