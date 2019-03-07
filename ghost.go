package main

func (sim *Sim) computeSpine() {
	sim.spineTangle = make(map[int]Tx)
	_, spineTip := ghostWalk(sim.tangle[0], sim)
	set := make(map[int]bool)
	dfs(spineTip, set, sim)
	//fmt.Println(len(set))
	// add genesis
	sim.spineTangle[0] = sim.tangle[0]
	// add directly and indirectly referenced txs
	for key := range set {
		sim.spineTangle[key] = sim.tangle[key]
		for _, v := range sim.approvers[key] {
			if _, exist := set[v]; exist {
				// update approvers
				sim.spineApprovers[key] = append(sim.spineApprovers[key], v)
			}
		}
	}
}

func ghostWalk(t Tx, sim *Sim) (path []int, tip Tx) {
	defer sim.b.track(runningtime("Ghost RW"))

	var current Tx
	for current = ghostStep(sim.tangle[0], sim); len(sim.approvers[current.id]) > 0; current = ghostStep(current, sim) {
		path = append(path, current.id)
		//fmt.Println(current.id, "\t", len(sim.approvers[current.id]), sim.tangle[current.id].cw)
	}
	return path, current
}

func ghostStep(t Tx, sim *Sim) Tx {
	directApprovers := sim.approvers[t.id]
	if (len(directApprovers)) == 0 {
		return t
	}
	if (len(directApprovers)) == 1 {
		return sim.tangle[directApprovers[0]]
	}

	var cws []int
	for _, approver := range directApprovers {
		cws = append(cws, sim.tangle[approver].cw)
	}
	//fmt.Println(cws)
	maxCW, _ := max(cws)
	return sim.tangle[directApprovers[maxCW]]
}
