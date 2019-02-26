package main

import "fmt"

type pOrphanResult struct {
	op float64
}

func (sim *Sim) runOrphaningP(result *pOrphanResult) {
	//remove txs grater than maxCut from both tangle and spine tangle so to have a comparable cone
	newTangle := sim.tangle[:sim.param.maxCut]
	newSpineTangle := sliceSpineTangle(sim, sim.param.maxCut)
	result.op = 1. - float64(len(newSpineTangle))/float64(len(newTangle))
	//fmt.Println("Prob=", result.p)
}

func (a pOrphanResult) Join(b pOrphanResult) pOrphanResult {
	if a.op == 0 {
		return b
	}
	var result pOrphanResult
	result.op = (a.op + b.op) / 2.
	return result
}

func (a pOrphanResult) String() string {
	return fmt.Sprintln("Orphaning Probability:", a.op)
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
