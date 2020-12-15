package main

import (
	"math"

	"github.com/iotaledger/iota.go/trinary"
)

// Tx defines the data structure of a transaction
type Tx struct {
	id                  int
	time                float64
	timestamp           int64
	attachmentTimestamp int64
	cw                  int
	cw2                 int // TODO: to remove, used only to compare different CW update mechanisms
	ref                 []int
	app                 []int
	firstApproval       float64

	bundle trinary.Hash
}

func (sim *Sim) newGenesis() Tx {

	genesis := Tx{
		id:            0,
		time:          0,
		cw:            1,
		firstApproval: -1,
		cw2:           1,
	}
	sim.tips = append(sim.tips, 0)
	//sim.cw = append(sim.cw, make([]uint64, 1))
	// sim.cw[0] = make([]uint64, 1)
	// sim.cw[0][0] = 1
	return genesis

}

func newTx(sim *Sim, previous Tx) Tx {
	t := Tx{
		id:            previous.id + 1,
		time:          sim.nextTime(previous),
		cw:            1,
		firstApproval: -1,
		cw2:           1,
	}

	return t
}

func (sim Sim) nextTime(t Tx) float64 {
	if sim.param.ConstantRate {
		return t.time + 1/sim.param.Lambda
	}
	//return t.time + rand.ExpFloat64()/sim.param.Lambda
	return t.time + (-1./sim.param.Lambda)*(math.Log(sim.generator.Float64()))

}

//remove txs from tip list that have now a visible approver
func (sim *Sim) removeOldTips(t Tx) {
	var currentTips []int
	for _, tip := range sim.tips {
		if !sim.tangle[tip].hasApprover(t.time, sim.param.H) {
			if sim.param.TSA == "RURTS" {
				if sim.tangle[tip].isTooOld(t.time, sim.param.D) {
					sim.orphanTips = append(sim.orphanTips, tip)
				} else {
					currentTips = append(currentTips, tip)
				}
			} else {
				currentTips = append(currentTips, tip)
			}

		}
	}

	sim.tips = currentTips
}

//given a time "now" and a transaction t, checks that t has (visible) approvers
func (t Tx) hasApprover(now float64, h int) bool {
	return (t.firstApproval > 0 && t.firstApproval+float64(h) < now)
}

//given a time "now" and a transaction t, checks that t is visible
func (t Tx) isVisible(now float64, h int) bool {
	return t.time+float64(h) < now || t.time == 0
}

//given a time "now" and a transaction t, checks that t is visible
func (t Tx) isTooOld(now float64, D int) bool {
	return now-t.time > float64(D)
}

func (t Tx) isGenesis() bool {
	//return t.time == 0
	if len(t.ref) > 0 {
		return false
	}
	return true
}

func (sim *Sim) tipsUpdate(t Tx) []int {
	//defer track(runningtime("udpateTips"))
	sim.removeOldTips(t)
	return sim.revealTips(t)
}

// reveal tips, update CW
func (sim *Sim) revealTips(t Tx) []int {
	var newTips []int
	//fmt.Println("HiddenTips", sim.hiddenTips)

	if len(sim.hiddenTips) == 0 {
		return newTips
	}

	i := 0
	tip := sim.hiddenTips[i]
	for sim.tangle[tip].isVisible(t.time, sim.param.H) {
		sim.updateApprovers(sim.tangle[tip])
		// sim.updateCWOpt(sim.tangle[tip])
		//if sim.param.TSA == "RW" && sim.param.Alpha != 0 {
		//	sim.updateCWOpt(sim.tangle[tip])
		//}
		// if sim.param.TSA == "RW" && sim.param.Alpha != 0 {
		// 	sim.updateCW(sim.tangle[tip])
		//sim.updateCWDFS(sim.tangle[tip])
		// }
		i++
		if i >= len(sim.hiddenTips) {
			break
		}
		tip = sim.hiddenTips[i]
	}

	newTips = sim.hiddenTips[:i]        //set of "new" visible tips
	sim.hiddenTips = sim.hiddenTips[i:] //update hidden tips set
	return newTips
}

func (sim *Sim) updateApprovers(t Tx) {
	//defer track(runningtime("updateApprovers"))
	for _, refID := range t.ref {
		sim.tangle[refID].app = appendUnique(sim.tangle[refID].app, t.id)
	}

}

// func (sim *Sim) updateCW(tip Tx) {
// 	defer sim.b.track(runningtime("updateCW-BitMask"))
// 	sim.cw = append(sim.cw, cwBitMask(sim.tangle[tip.id], sim.cw))
// 	sim.addCW(sim.cw[tip.id], 1)
// }

// func (sim *Sim) updateCWOpt(tip Tx) {
// 	defer sim.b.track(runningtime("updateCW-BitMask-Opt"))
// 	//TODO: use circular array
// 	// sim.cw[tip.id % CWMatrixLen]
// 	sim.cw[tip.id%sim.param.CWMatrixLen] = cwBitMaskOpt(sim.tangle[tip.id], sim) //append(sim.cw, cwBitMask(sim.tangle[tip.id], sim.cw))
// 	sim.addCW(sim.cw[tip.id%sim.param.CWMatrixLen], 1)
// }

// func (sim *Sim) updateCWDFS(tip Tx) {
// 	defer sim.b.track(runningtime("updateCW-DFS"))
// 	set := make(map[int]bool)
// 	dfs(tip, set, sim)
// 	//fmt.Println(tip.id, set)
// 	for k := range set {
// 		sim.tangle[k].cw2++
// 	}

// }

// func (sim *Sim) compareCW() bool {
// 	//fmt.Println(sim.tangle)
// 	//printApprovers(sim.approvers)
// 	for _, t := range sim.tangle {
// 		//fmt.Println(t.id, t.cw, t.cw2)
// 		if t.cw != t.cw2 {
// 			return false
// 		}
// 	}
// 	return true
// }

// func weightedChoose(approvers []int, weights []float64, g *rand.Rand, b Benchmark) (int, int) {
// 	//defer b.track(runningtime("weightedChoose"))
// 	var sum float64
// 	for _, w := range weights {
// 		sum += w
// 	}
// 	rand := g.Float64() * sum

// 	cumSum := weights[0]
// 	for i := 1; i < len(approvers); i++ {
// 		if rand < cumSum {
// 			return approvers[i-1], i - 1
// 		}
// 		cumSum += weights[i]
// 	}
// 	return approvers[len(approvers)-1], len(approvers) - 1
// }

//
// func normalizeWeights(approvers []int, sim *Sim) []float64 {
// 	//defer sim.b.track(runningtime("normalizeWeights"))
// 	cumWeigths := make([]float64, len(approvers))
// 	normalizedWeights := make([]float64, len(approvers))
// 	maxWeigth := 0.
// 	for i, a := range approvers {
// 		cumWeigths[i] = float64(sim.tangle[a].cw)
// 		if cumWeigths[i] > maxWeigth {
// 			maxWeigth = cumWeigths[i]
// 		}
// 	}
// 	for i, w := range cumWeigths {
// 		//w -= maxWeigth + 6.9 // ln(2^10) = 6.9314718056, 10 - number of bits in exponent of double precision
// 		//normalizedWeights[i] = math.Exp(maxFloat64(-700., sim.param.Alpha*w))
// 		w -= maxWeigth
// 		normalizedWeights[i] = math.Exp(sim.param.Alpha * w)
// 	}
// 	return normalizedWeights
// }

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func dfs(t Tx, visited map[int]bool, sim *Sim) {
	if t.id > 0 {
		for _, id := range t.ref {
			if !visited[id] {
				visited[id] = true
				dfs(sim.tangle[id], visited, sim)
			}
		}
	}
}

// func dfsBitMask(t Tx, sim *Sim) []uint64 {
// 	set := make(map[int]bool)
// 	dfs(t, set, sim)
// 	// find max to know size
// 	var size, base uint64
// 	max := 0
// 	for k := range set {
// 		if k > max {
// 			max = k
// 		}
// 	}
// 	base = 64
// 	size = uint64(max) / base
// 	refCW := make([]uint64, int(size)+1)
// 	for k := range set {
// 		l := uint64(k)
// 		//update bit
// 		refCW[int(l/base)] |= (1 << uint64(l%(base)))
// 	}

// 	// for k := range set {
// 	// 	l := uint64(k)
// 	// 	if l/base < size {
// 	// 		//update bit
// 	// 		refCW[int(size)-1] |= (1 << uint64(l%(base)))
// 	// 	} else {
// 	// 		for i = 0; i <= l/base-size; i++ {
// 	// 			refCW = append(refCW, 0) //fill gap
// 	// 		}
// 	// 		size = l / base
// 	// 		refCW[len(refCW)-1] |= 1 << uint64(l%(base)) //add new bit
// 	// 	}
// 	// }
// 	return refCW
// }

// func cwBitMask(t Tx, cw [][]uint64) []uint64 {
// 	refCW := make([][]uint64, len(t.ref))
// 	var i uint64
// 	var base uint64
// 	base = 64
// 	for k, link := range t.ref {
// 		//TODO: add check t.id - link < 50*lambda, if not dfs or additional data structure
// 		l := uint64(link)
// 		size := uint64(len(cw[link]))
// 		refCW[k] = make([]uint64, size)
// 		copy(refCW[k], cw[link])
// 		if l/base < size {
// 			refCW[k][int(size)-1] |= (1 << uint64(l%(base))) //add new link
// 		} else {
// 			for i = 0; i <= l/base-size; i++ {
// 				refCW[k] = append(refCW[k], 0) //fill gap
// 			}
// 			refCW[k][len(refCW[k])-1] |= 1 << uint64(l%(base)) //add new link
// 		}
// 	}

// 	a := []int{}
// 	for _, r := range refCW {
// 		a = append(a, len(r))
// 	}
// 	_, max := max(a)

// 	for i, r := range refCW {
// 		if max-len(r) > 0 {
// 			padding := make([]uint64, max-len(r))
// 			refCW[i] = append(refCW[i], padding...)
// 		}
// 	}

// 	or := refCW[0]
// 	for i := 1; i < len(t.ref); i++ {
// 		for j := 0; j < max; j++ {
// 			or[j] = or[j] | refCW[i][j]
// 		}
// 	}
// 	return or

// }

// func cwBitMaskOpt(t Tx, sim *Sim) []uint64 {
// 	refCW := make([][]uint64, len(t.ref))
// 	var i uint64
// 	var base uint64
// 	base = 64
// 	for k, link := range t.ref {
// 		l := uint64(link)
// 		var size uint64
// 		//check t.id - link < 50*lambda, if not dfs or additional data structure
// 		if t.id-link < sim.param.CWMatrixLen {
// 			size = uint64(len(sim.cw[link%sim.param.CWMatrixLen]))
// 			refCW[k] = make([]uint64, size)
// 			copy(refCW[k], sim.cw[link%sim.param.CWMatrixLen])
// 		} else {
// 			refCW[k] = dfsBitMask(sim.tangle[link], sim)
// 			size = uint64(len(refCW[k]))
// 			//fmt.Println("DFS")
// 		}
// 		if l/base < size {
// 			refCW[k][int(size)-1] |= (1 << uint64(l%(base))) //add new link
// 		} else {
// 			for i = 0; i <= l/base-size; i++ {
// 				refCW[k] = append(refCW[k], 0) //fill gap
// 			}
// 			refCW[k][len(refCW[k])-1] |= 1 << uint64(l%(base)) //add new link
// 		}
// 	}

// 	a := []int{}
// 	for _, r := range refCW {
// 		a = append(a, len(r))
// 	}
// 	_, max := max(a)

// 	for i, r := range refCW {
// 		if max-len(r) > 0 {
// 			padding := make([]uint64, max-len(r))
// 			refCW[i] = append(refCW[i], padding...)
// 		}
// 	}

// 	or := refCW[0]
// 	for i := 1; i < len(t.ref); i++ {
// 		for j := 0; j < max; j++ {
// 			or[j] = or[j] | refCW[i][j]
// 		}
// 	}
// 	return or

// }

// // add BitMask to CW
// func (sim *Sim) addCW(a []uint64, propweight int) {
// 	//defer sim.b.track(runningtime("addCW"))
// 	base := 64
// 	for i, block := range a { // the i iterates through the blocks, block is the unit64 of the block itself
// 		for j := 0; j < base; j++ {
// 			if block&(1<<uint(j)) != 0 {
// 				sim.tangle[j+(i*base)].cw += propweight
// 			}
// 		}
// 	}
// }

// // remove BitMask from CW
// func (sim *Sim) removeCW(a []uint64, propweight int) {
// 	//defer sim.b.track(runningtime("addCW"))
// 	base := 64
// 	for i, block := range a { // the i iterates through the blocks, block is the unit64 of the block itself
// 		for j := 0; j < base; j++ {
// 			if block&(1<<uint(j)) != 0 {
// 				sim.tangle[j+(i*base)].cw -= propweight
// 			}
// 		}
// 	}
// }

// ReferencedByTx checks for particular tx if it is referenced
func ReferencedByTx(a []uint64, ID int) bool {
	base := 64
	baseID := ID / base
	localID := ID - baseID*base
	// if len(a) > baseID && a[baseID]&(1<<uint(localID)) != 0 {
	if len(a) > baseID {
		if a[baseID]&(1<<uint(localID)) != 0 {
			return true
		}
	}
	return false
}

// // setCW set CW bit of a particular tx ID
// func setCW(a []uint64, ID int) {
// 	base := 64
// 	block := ID / base
// 	bitToSet := uint64(ID % base)

// 	if len(a) > block {
// 		a[block] |= (1 << bitToSet)
// 	}
// 	//fmt.Printf("%d \t", ID)
// 	//printCWRef(a)
// }

// ReferencedByRecentTx checks for particular tx if it is referenced by recent tx
func (sim *Sim) ReferencedByRecentTx(searchTx, lastTx, numRecent int) bool {
	for tx := lastTx; tx > lastTx-numRecent && searchTx < tx; tx-- {
		if searchTx < sim.tangle[tx].ref[0] { // only check if at least one of tx's approvers is known
			if ReferencedByTx(sim.cw[tx], searchTx) {
				return true
			}
		}
	}
	return false
}

func (sim *Sim) isLeftBehind(thisTx int) bool {
	recentTx := 3 * int(sim.param.Lambda)
	if thisTx > sim.param.TangleSize-recentTx { // still recent enough to be considered for a root
		// fmt.Println("-")
		return false
	}
	// check if left behind
	// fmt.Println("+")
	// only have cw for len(sim.cw) of the txs
	if sim.ReferencedByRecentTx(thisTx, len(sim.cw)-1, recentTx-1) { // if its not in the most recent cw set abandon the tx.
		return false
	}

	return true

}

//return a ordered list of IDs from a map[int]Tx
func returnIDlist(tanglepart map[int]Tx) []int {
	a := make([]int, len(tanglepart))
	counter := 0
	for id := range tanglepart {
		a[counter] = id
		counter++
	}
	help := 0
	for i1 := 0; i1 < len(a)-1; {
		if a[i1+1] < a[i1] {
			help = a[i1]
			a[i1] = a[i1+1]
			a[i1+1] = help
			if i1 > 0 {
				i1--
			}
		} else {
			i1++
		}
	}
	return a
}
