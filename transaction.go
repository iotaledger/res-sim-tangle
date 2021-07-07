package main

import (
	"math"
)

// Tx defines the data structure of a transaction
type Tx struct {
	id                       int // id of the tx
	nodeID                   int // if of the issuing node
	time                     float64
	h                        int
	timestamp                int64
	attachmentTimestamp      int64
	cw                       int
	aw                       float64 // Approval weight
	ref                      []int   // parents
	app                      []int   // approvers
	firstApprovalTime        float64 // first time transaction is approved
	confirmationTime         int     // first time transaction reached confirmation status
	firstVisibleApprovalTime float64 // what is this?

	//bundle trinary.Hash
}

func (sim *Sim) newGenesis() Tx {

	genesis := Tx{
		id:                       0,
		nodeID:                   0,
		time:                     0,
		aw:                       1,
		firstApprovalTime:        -1,
		firstVisibleApprovalTime: -1,
		confirmationTime:         -1,
	}
	sim.tips = append(sim.tips, 0)
	//sim.cw = append(sim.cw, make([]uint64, 1))
	// sim.cw[0] = make([]uint64, 1)
	// sim.cw[0][0] = 1
	return genesis

}

func newTx(sim *Sim, previous Tx, nodeID int) Tx {
	t := Tx{
		id:                       previous.id + 1,
		nodeID:                   nodeID,
		time:                     sim.nextTime(previous),
		h:                        sim.setDelay(),
		aw:                       0,
		firstApprovalTime:        -1,
		firstVisibleApprovalTime: -1,
		confirmationTime:         -1,
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

func (sim Sim) setDelay() int {
	if sim.param.p < sim.generator.Float64() {
		return sim.param.Hsmall
	}
	return sim.param.Hlarge
}

//remove txs from tip list that have now a visible approver
func (sim *Sim) removeOldTips(t Tx) {
	var currentTips []int
	for _, tip := range sim.tips {
		if !sim.tangle[tip].hasVisibleApprover(t.time) {
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
func (t Tx) hasVisibleApprover(now float64) bool {
	return (t.firstApprovalTime > 0 && t.firstVisibleApprovalTime < now)
}

//given a time "now",  a transaction t and a max window h, checks that t is visible
func (t Tx) isVisible(now float64) bool {
	return t.time+float64(t.h) < now || t.time == 0
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

// reveal tips, update CW
func (sim *Sim) revealTips(t Tx) []int {
	var newTips []int
	var newHiddenTips []int
	//fmt.Println("HiddenTips", sim.hiddenTips)

	if len(sim.hiddenTips) == 0 {
		return newTips
	}

	for _, tip := range sim.hiddenTips { //go through all hidden tips
		if sim.tangle[tip].isVisible(t.time) { // if tip is revealed now
			sim.updateApprovers(sim.tangle[tip])
			newTips = append(newTips, tip)
		} else {
			newHiddenTips = append(newHiddenTips, tip)
		}
	}

	sim.hiddenTips = newHiddenTips //update hidden tips set
	return newTips
}

func (sim *Sim) updateApprovers(t Tx) {
	//defer track(runningtime("updateApprovers"))
	for _, refID := range t.ref {
		sim.tangle[refID].app = appendUnique(sim.tangle[refID].app, t.id)
	}

}

func (sim *Sim) updateAW(t Tx, nodeID int) {
	if sim.nodeApprover[nodeID][t.id] {
		return // we stop exploring the pas cone
	} else {
		sim.nodeApprover[nodeID][t.id] = true
		sim.tangle[t.id].aw += sim.mana[nodeID]
		//fmt.Println("updating past cone of", t.id)
		//fmt.Println("List of parents", t.ref)
		if sim.tangle[t.id].aw > 0.5 {
			sim.tangle[t.id].confirmationTime = sim.tangleAge - sim.tangle[t.id].id //
			// we stop exploring the pas cone
		}
		if sim.tangle[t.id].aw > 0.8 {
			return
		}
		for _, refID := range t.ref {
			//fmt.Println("Checking parents", refID)
			sim.updateAW(sim.tangle[refID], nodeID)
		}
		return
	}
}

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
