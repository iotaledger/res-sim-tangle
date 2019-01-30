package main

import (
	"math"
	"math/rand"
)

// Tx defines the data structure of a transaction
type Tx struct {
	id            int
	time          float64
	cw            int
	cw2           int // TODO: to remove, used only to compare different CW update mechanisms
	ref           []int
	firstApproval float64
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
	sim.cw = append(sim.cw, make([]uint64, 1))
	sim.cw[0][0] = 1
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

func (sim *Sim) removeOldTips(t Tx) {
	var currentTips []int
	for _, tip := range sim.tips {
		if !sim.tangle[tip].hasApprover(t.time, sim.param.H) {
			currentTips = append(currentTips, tip)
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

func (t Tx) isGenesis() bool {
	return t.time == 0
}

func (sim *Sim) tipsUpdate(t Tx) []int {
	//defer track(runningtime("udpateTips"))
	sim.removeOldTips(t)
	return sim.revealTips(t)
}

func (sim *Sim) revealTips(t Tx) []int {
	var newTips []int
	//fmt.Println("HiddenTips", sim.hiddenTips)

	if len(sim.hiddenTips) == 0 {
		return newTips
	}

	i := 0
	tip := sim.hiddenTips[i]
	for sim.tangle[tip].isVisible(t.time, sim.param.H) {
		sim.approvers = updateApprovers(sim.approvers, sim.tangle[tip])
		if sim.param.TSA == "RW" || sim.param.VelocityEnabled == true {
			sim.updateCW(sim.tangle[tip])
			//sim.updateCWDFS(sim.tangle[tip])
		}
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

func updateApprovers(a map[int][]int, t Tx) map[int][]int {
	//defer track(runningtime("updateApprovers"))
	for _, ref := range t.ref {
		a[ref] = append(a[ref], t.id)
	}
	return a
}

func (sim *Sim) updateCW(tip Tx) {
	//defer sim.b.track(runningtime("updateCW-BitMask"))
	sim.cw = append(sim.cw, cwBitMask(sim.tangle[tip.id], sim.cw))
	sim.addCW(sim.cw[tip.id])

}

func (sim *Sim) updateCWDFS(tip Tx) {
	//defer sim.b.track(runningtime("updateCW-DFS"))
	set := make(map[int]bool)
	//dfs(tip, set, sim)
	//fmt.Println(tip.id, set)
	for k := range set {
		sim.tangle[k].cw2++
	}

}

func (sim *Sim) compareCW() bool {
	//fmt.Println(sim.tangle)
	//printApprovers(sim.approvers)
	for _, t := range sim.tangle {
		//fmt.Println(t.id, t.cw, t.cw2)
		if t.cw != t.cw2 {
			return false
		}
	}
	return true
}

func weightedChoose(approvers []int, weights []float64, g *rand.Rand, b Benchmark) (int, int) {
	//defer b.track(runningtime("weightedChoose"))
	var sum float64
	for _, w := range weights {
		sum += w
	}
	rand := g.Float64() * sum

	cumSum := weights[0]
	for i := 1; i < len(approvers); i++ {
		if rand < cumSum {
			return approvers[i-1], i - 1
		}
		cumSum += weights[i]
	}
	return approvers[len(approvers)-1], len(approvers) - 1
}

func normalizeWeights(approvers []int, sim *Sim) []float64 {
	//defer sim.b.track(runningtime("normalizeWeights"))
	cumWeigths := make([]float64, len(approvers))
	normalizedWeights := make([]float64, len(approvers))
	maxWeigth := 0.
	for i, a := range approvers {
		cumWeigths[i] = float64(sim.tangle[a].cw)
		if cumWeigths[i] > maxWeigth {
			maxWeigth = cumWeigths[i]
		}
	}
	for i, w := range cumWeigths {
		//w -= maxWeigth + 6.9 // ln(2^10) = 6.9314718056, 10 - number of bits in exponent of double precision
		//normalizedWeights[i] = math.Exp(maxFloat64(-700., sim.param.Alpha*w))
		w -= maxWeigth
		normalizedWeights[i] = math.Exp(sim.param.Alpha * w)
	}
	return normalizedWeights
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func cwBitMask(t Tx, cw [][]uint64) []uint64 {
	refCW := make([][]uint64, len(t.ref))
	var i uint64
	var base uint64
	base = 64
	for k, link := range t.ref {
		//TODO: add check t.id - link < 50*lambda, if not dfs or additional data structure
		l := uint64(link)
		size := uint64(len(cw[link]))
		refCW[k] = make([]uint64, size)
		copy(refCW[k], cw[link])
		if l/base < size {
			refCW[k][int(size)-1] |= (1 << uint64(l%(base))) //add new link
		} else {
			for i = 0; i <= l/base-size; i++ {
				refCW[k] = append(refCW[k], 0) //fill gap
			}
			refCW[k][len(refCW[k])-1] |= 1 << uint64(l%(base)) //add new link
		}
	}

	a := []int{}
	for _, r := range refCW {
		a = append(a, len(r))
	}
	_, max := max(a)

	for i, r := range refCW {
		if max-len(r) > 0 {
			padding := make([]uint64, max-len(r))
			refCW[i] = append(refCW[i], padding...)
		}
	}

	or := refCW[0]
	for i := 1; i < len(t.ref); i++ {
		for j := 0; j < max; j++ {
			or[j] = or[j] | refCW[i][j]
		}
	}
	return or

}

func (sim *Sim) addCW(a []uint64) {
	//defer sim.b.track(runningtime("addCW"))
	base := 64
	for i, block := range a {
		for j := 0; j < base; j++ {
			if block&(1<<uint(j)) != 0 {
				sim.tangle[j+(i*base)].cw++
			}
		}
	}
}
