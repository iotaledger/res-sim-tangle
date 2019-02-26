package main

import "fmt"

type pOrphanResult struct {
	p float64
}

func (sim *Sim) runPOrphan(result *pOrphanResult) {
	//remove txs grater than maxCut from both tangle and spine tangle so to have a comparable cone
	newTangle := sim.tangle[:sim.param.maxCut]
	newSpineTangle := sliceSpineTangle(sim, sim.param.maxCut)
	result.p = 1. - float64(len(newSpineTangle))/float64(len(newTangle))
	//fmt.Println("Prob=", result.p)
}

func (a pOrphanResult) Join(b pOrphanResult) pOrphanResult {
	if a.p == 0 {
		return b
	}
	var result pOrphanResult
	result.p = (a.p + b.p) / 2.
	return result
}

func (a pOrphanResult) String() string {
	return fmt.Sprintln("Probability of becoming orphan:", a.p)
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
