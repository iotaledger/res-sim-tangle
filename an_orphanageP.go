// calculate orphanage probability
// op = sample is past cone of GHOST RW, with maxcut
// op2 = sample is past cone of recent txs, with maxcut
// op3 = linear regression

package main

import (
	"fmt"

	"gonum.org/v1/gonum/stat"
)

type pOrphanResult struct {
	// orphanage probability - probability of tx left behind
	// tip orphanage probability - probability of tips left behind
	op          []float64 // past cone of Ghost RW
	top         []float64
	op2         []float64 // union of past cones of N recent txs
	top2        []float64
	op3         []float64 // linear regression
	top3        []float64
	op4         []float64 // sum Ghost cone + N recent cones
	top4        []float64
	nTipsAtID   []int // for each index record the number of tips
	nOrphanAtID []int // for each index record the number of orphaned txs
	tTangle     []float64
	tSpine      []float64
	tOrphan     []float64
}

func (sim *Sim) runOrphaningP(result *pOrphanResult) {
	// remove txs grater than maxCut from both tangle and spine tangle so to have a comparable cone
	newTangle := sim.tangle[sim.param.minCut:sim.param.maxCut]
	// calculate spine Tangle up to maxCut (all txs (directly/indirectly) referenced by a GHOST particle (alpha = infinity)
	newSpineTangle := sliceMap(sim.spineTangle, sim.param.minCut, sim.param.maxCut)

	result.tTangle = append(result.tTangle, getAverageApprovalTime(sliceToMap(newTangle)))
	result.tSpine = append(result.tSpine, getAverageApprovalTime(newSpineTangle))
	orphanTangle := getOrphanTxs(sim)
	result.tOrphan = append(result.tOrphan, getAverageApprovalTime(orphanTangle))

	sim.runOrphanageGHOST(result, newTangle, newSpineTangle) // calculate op

	recentTipsCones := sim.runOrphanageRecent(result) // calculate op2

	if sim.param.pOrphanLinFitEnabled {
		sim.runOrphanageLinFit(result) // calculate op3
	}

	sim.runOrphanageGHOSTRecent(result, newTangle, newSpineTangle, recentTipsCones)
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
	result.op3 = append(a.op3, b.op3...)
	result.top3 = append(a.top3, b.top3...)
	result.op4 = append(a.op4, b.op4...)
	result.top4 = append(a.top4, b.top4...)
	result.tTangle = append(a.tTangle, b.tTangle...)
	result.tSpine = append(a.tSpine, b.tSpine...)
	result.tOrphan = append(a.tOrphan, b.tOrphan...)
	return result
}

// method String() implements the interface Stringer, so you call this automatically whenever you try to print, e.g., fmt.Println(a)
// it overwrites the functionality of fmt.Println()
func (a pOrphanResult) String() string {
	mean, std := stat.MeanStdDev(a.op, nil)
	result := fmt.Sprintf("Orphanage Probability [mean, StdDev]:\t\t%.5f\t%.5f\n", mean, std)
	mean, std = stat.MeanStdDev(a.top, nil)
	result += fmt.Sprintf("Tip Orphanage Probability [mean, StdDev]:\t%.5f\t%.5f\n", mean, std)
	mean, std = stat.MeanStdDev(a.op2, nil)
	result += fmt.Sprintf("Orphanage Probability 2 [mean, StdDev]:\t\t%.5f\t%.5f\n", mean, std)
	mean, std = stat.MeanStdDev(a.top2, nil)
	result += fmt.Sprintf("Tip Orphanage Probability 2: [mean, StdDev]:\t%.5f\t%.5f\n", mean, std)
	mean, std = stat.MeanStdDev(a.op3, nil)
	result += fmt.Sprintf("Orphanage Probability 3 [mean, StdDev]:\t\t%.5f\t%.5f\n", mean, std)
	mean, std = stat.MeanStdDev(a.top3, nil)
	result += fmt.Sprintf("Tip Orphanage Probability 3: [mean, StdDev]:\t%.5f\t%.5f\n", mean, std)
	mean, std = stat.MeanStdDev(a.op4, nil)
	result += fmt.Sprintf("Orphanage Probability 4 [mean, StdDev]:\t\t%.5f\t%.5f\n", mean, std)
	mean, std = stat.MeanStdDev(a.top4, nil)
	result += fmt.Sprintf("Tip Orphanage Probability 4: [mean, StdDev]:\t%.5f\t%.5f\n", mean, std)
	result += fmt.Sprintf("Avg approval time [Tangle, Spine, Orphan]:\t%.5f\t%.5f\t%.5f\n", stat.Mean(a.tTangle, nil), stat.Mean(a.tSpine, nil), stat.Mean(a.tOrphan, nil))
	return result
}

func newPOrphanResult(p *Parameters) pOrphanResult {
	// variables initialization for pOprhan
	var result pOrphanResult
	if p.pOrphanLinFitEnabled {
		result.nTipsAtID = append(result.nTipsAtID, make([]int, p.TangleSize)...)
		result.nOrphanAtID = append(result.nOrphanAtID, make([]int, p.TangleSize)...)
	}
	return result
}

func sliceMap(m map[int]Tx, lBound, uBound int) map[int]Tx {
	result := make(map[int]Tx)
	for k, v := range m {
		if k < uBound && k >= lBound {
			result[k] = v
		}
	}
	return result
}

func getOrphanTxs(sim *Sim) map[int]Tx {
	newTangle := sim.tangle[sim.param.minCut:sim.param.maxCut]
	newSpineTangle := sliceMap(sim.spineTangle, sim.param.minCut, sim.param.maxCut)
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

func (sim *Sim) runOrphanageGHOST(result *pOrphanResult, newTangle []Tx, newSpineTangle map[int]Tx) {
	// calculate op
	result.op = append(result.op, 1.-float64(len(newSpineTangle))/float64(len(newTangle)))

	// calculate top by finding all tips left behind and dividing that number over all txs
	top := 0.
	for _, tx := range newTangle {
		if len(sim.tangle[tx.id].app) == 0 {
			top++
		}
	}
	result.top = append(result.top, top/float64(len(newTangle)))
}

func (sim *Sim) runOrphanageGHOSTRecent(result *pOrphanResult, newTangle []Tx, newSpineTangle, coneRecent map[int]Tx) {
	// calculate op
	sum := newSpineTangle
	for k, v := range coneRecent {
		sum[k] = v
	}
	result.op4 = append(result.op4, 1.-float64(len(sum))/float64(len(newTangle)))

	// calculate top by finding all tips left behind and dividing that number over all txs
	top := 0.
	for _, tx := range newTangle {
		if len(sim.tangle[tx.id].app) == 0 {
			top++
		}
	}
	result.top4 = append(result.top4, top/float64(len(newTangle)))
}

func (sim *Sim) runOrphanageRecent(result *pOrphanResult) map[int]Tx {
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

	//ORing all the cones
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
			if id < sim.param.maxCut && id >= sim.param.minCut {
				//if coneUnionBitMask[block]&(1<<i) != 0 {
				if (coneUnionBitMask[block]>>i)&1 == 1 {
					ones[id] = sim.tangle[id]
				} else {
					zeros[id] = sim.tangle[id]
				}
			}
		}
	}

	// set := make(map[int]bool)
	// for tx := lastVisibleTx; tx > lastVisibleTx-sim.param.stillrecent; tx-- {
	// 	dfs(sim.tangle[tx], set, sim)
	// }

	//ones = sliceMap(ones, sim.param.minCut, sim.param.maxCut)
	//zeros = sliceMap(zeros, sim.param.minCut, sim.param.maxCut)
	top := 0
	for id := range zeros {
		if len(sim.tangle[id].app) == 0 {
			top++
		}
	}

	result.op2 = append(result.op2, 1-float64(len(ones))/float64(sim.param.maxCut-sim.param.minCut))
	result.top2 = append(result.top2, float64(top)/float64(sim.param.maxCut-sim.param.minCut))
	return ones
}

// how many txs are orphaned at a given index
func (sim *Sim) runAnOPLinfit(tx int, r *pOrphanResult, run int) {
	sim.computeSpine()
	r.nTipsAtID[tx] = len(sim.tips)
	r.nOrphanAtID[tx] = tx - len(sim.spineTangle)
}

// apply linear regression
func (sim *Sim) runOrphanageLinFit(r *pOrphanResult) {
	x := makeRangeInt(1, sim.param.TangleSize)
	y1 := r.nTipsAtID
	m1, _ := linFit(x, y1)
	r.top3 = append(r.top3, m1)
	y2 := r.nOrphanAtID
	m2, _ := linFit(x, y2)
	r.op3 = append(r.op3, m2)
}

func linFit(x, y []int) (float64, float64) {
	xbar := meanInt(x)
	ybar := meanInt(y)
	if len(x) != len(y) {
		fmt.Println("Arrays have not the same length")
		fmt.Println("len(x)", len(x))
		fmt.Println("len(y)", len(y))
		pauseit()
	}
	div := 0.
	m := 0.
	for i1 := 0; i1 < len(x); i1++ {
		div += (float64(x[i1]) - xbar) * (float64(x[i1]) - xbar)
	}
	for i1 := 0; i1 < len(x); i1++ {
		m += (float64(x[i1]) - xbar) * (float64(y[i1]) - ybar) / float64(div)
	}
	b := ybar - m*xbar
	return m, b
}
