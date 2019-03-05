package main

import (
	"fmt"

	"gonum.org/v1/gonum/stat"
)

type pOrphanResult struct {
	op      []float64 // orphanage probability - probability of tx left behind
	top     []float64 // tip orphanage probability - probability of tips left behind
	op2     []float64
	top2    []float64
	tTangle []float64
	tSpine  []float64
	tOrphan []float64
}

func (sim *Sim) runOrphaningP(result *pOrphanResult) {
	// remove txs grater than maxCut from both tangle and spine tangle so to have a comparable cone
	newTangle := sim.tangle[:sim.param.maxCut]
	// calculate spine Tangle up to maxCut (all txs (directly/indirectly) referenced by a GHOST particle (alpha = infinity)
	newSpineTangle := sliceMap(sim.spineTangle, sim.param.maxCut)

	// calculate op
	result.op = append(result.op, 1.-float64(len(newSpineTangle))/float64(len(newTangle)))

	// calculate top by finding all tips left behind and dividing that number over all txs
	top := 0.
	for tx := range newTangle {
		if len(sim.approvers[tx]) == 0 {
			top++
		}
	}
	result.top = append(result.top, top/float64(len(newTangle)))

	result.tTangle = append(result.tTangle, getAverageApprovalTime(sliceToMap(newTangle)))
	result.tSpine = append(result.tSpine, getAverageApprovalTime(newSpineTangle))
	orphanTangle := getOrphanTxs(sim)
	result.tOrphan = append(result.tOrphan, getAverageApprovalTime(orphanTangle))

	sim.runOrphanageRecent(result)
}

func (a pOrphanResult) Join(b pOrphanResult) pOrphanResult {
	if a.op == nil {
		return b
	}
	var result pOrphanResult
	result.op = append(a.op, b.op...)
	result.top = append(a.top, b.top...)
	result.op2 = append(a.op2, b.op2...)
	result.top2 = append(a.top2, b.top2...)
	result.tTangle = append(a.tTangle, b.tTangle...)
	result.tSpine = append(a.tSpine, b.tSpine...)
	result.tOrphan = append(a.tOrphan, b.tOrphan...)
	return result
}

func (a pOrphanResult) String() string {
	result := fmt.Sprintln("Orphanage Probability:", stat.Mean(a.op, nil))
	result += fmt.Sprintln("Tip Orphanage Probability:", stat.Mean(a.top, nil))
	result += fmt.Sprintln("Orphanage Probability 2:", stat.Mean(a.op2, nil))
	result += fmt.Sprintln("Tip Orphanage Probability 2:", stat.Mean(a.top2, nil))
	result += fmt.Sprintln("Avg approval time [Tangle, Spine, Orphan]", stat.Mean(a.tTangle, nil), stat.Mean(a.tSpine, nil), stat.Mean(a.tOrphan, nil))
	return result
}

func newPOrphanResult() *pOrphanResult {
	// variables initialization for pOprhan
	var result pOrphanResult
	return &result
}

func sliceMap(m map[int]Tx, uBound int) map[int]Tx {
	result := make(map[int]Tx)
	for k, v := range m {
		if k < uBound {
			result[k] = v
		}
	}
	return result
}

func getOrphanTxs(sim *Sim) map[int]Tx {
	newTangle := sim.tangle[:sim.param.maxCut]
	newSpineTangle := sliceMap(sim.spineTangle, sim.param.maxCut)
	result := make(map[int]Tx)

	for k, v := range newTangle {
		if _, ok := newSpineTangle[k]; !ok {
			result[k] = v
		}
	}
	return result
}

func getAverageApprovalTime(tangle map[int]Tx) float64 {
	avgTime := 0.
	i := 0
	for _, v := range tangle {
		if v.firstApproval > 0 {
			avgTime += v.firstApproval - v.time
			i++
		}
	}
	return avgTime / float64(i)
}

func sliceToMap(tangle []Tx) map[int]Tx {
	result := make(map[int]Tx)
	for _, t := range tangle {
		result[t.id] = t
	}
	return result
}

func getCWgrowth(tx Tx, sim *Sim) []int {
	cw := 1
	cwGrowth := make([]int, sim.param.TangleSize-tx.id)

	for i := tx.id; i < sim.param.TangleSize; i++ {
		if tx.referencedByTxBitMask(sim.cw[i]) {
			//if tx.referencedByTxDFS(sim.tangle[i], sim) {
			//if ReferencedByTx(sim.cw[i], tx.id) {
			cw++
		}
		cwGrowth[i-tx.id] = cw
	}
	return cwGrowth
}

func (tx Tx) referencedByTxBitMask(txCWbitMask []uint64) bool {
	base := 64
	block := tx.id / base
	bitToCheck := uint64(tx.id % base)

	if len(txCWbitMask) > block && int((txCWbitMask[block]>>bitToCheck)&1) == 1 {
		return true
	}
	return false
}

func (tx Tx) referencedByTxDFS(refTx Tx, sim *Sim) bool {
	_, lastVisibleTx := max(sim.tips)
	if refTx.id <= lastVisibleTx {
		set := make(map[int]bool)
		dfs(refTx, set, sim)
		if _, ok := set[tx.id]; ok {
			return true
		}
	}
	return false
}

func (sim *Sim) runOrphanageRecent(result *pOrphanResult) {
	var base uint
	base = 64
	size := 0
	_, lastVisibleTx := max(sim.tips)

	// finding the size of coneUnionBitMask
	for tx := lastVisibleTx; tx > lastVisibleTx-sim.param.stillrecent; tx-- {
		if len(sim.cw[tx]) > size {
			size = len(sim.cw[tx])
		}
	}
	//data structure to keep information about directly and indirectly referenced txs
	coneUnionBitMask := make([]uint64, size)

	//oring all the cones
	for tx := lastVisibleTx; tx > lastVisibleTx-sim.param.stillrecent; tx-- {
		for block := 0; block < len(sim.cw[tx]); block++ {
			coneUnionBitMask[block] |= sim.cw[tx][block]
		}
	}
	ones := make(map[int]Tx)
	zeros := make(map[int]Tx)
	for block := 0; block < len(coneUnionBitMask); block++ {
		var i uint
		for i = 0; i < base; i++ {
			id := block*int(base) + int(i)
			if int((coneUnionBitMask[block]>>i)&1) == 1 {
				ones[id] = sim.tangle[id]
			} else {
				if id <= lastVisibleTx {
					zeros[id] = sim.tangle[id]
				}
			}
		}
	}

	// set := make(map[int]bool)
	// for tx := lastVisibleTx; tx > lastVisibleTx-sim.param.stillrecent; tx-- {
	// 	dfs(sim.tangle[tx], set, sim)
	// }

	ones = sliceMap(ones, sim.param.maxCut)
	zeros = sliceMap(zeros, sim.param.maxCut)
	top := 0
	for id := range zeros {
		if len(sim.approvers[id]) == 0 {
			top++
		}
	}

	result.op2 = append(result.op2, 1-float64(len(ones))/float64(len(sim.tangle[:sim.param.maxCut])))
	result.top2 = append(result.top2, float64(top)/float64(len(sim.tangle[:sim.param.maxCut])))
}
