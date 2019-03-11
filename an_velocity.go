package main

import (
	"fmt"
	"math"
	"os"
	"sort"

	"gonum.org/v1/gonum/stat"
)

//Velocity result of simulation
type velocityResult struct {
	vID        []MetricIntInt //???creates a map[int]int with a keyword
	vTime      []MetricFloat64Int
	dApprovers []MetricIntInt
	vCW        []MetricIntInt
	vCWfirst   []MetricIntInt
}

//??? use string to create empty value maps to vID, vTime, dApprovers
func newVelocityResult(veloMetrics []string, param Parameters) *velocityResult {
	// variables initialization for velocity
	var result velocityResult
	for _, metric := range veloMetrics {
		if metric != "backG" {
			result.vID = append(result.vID, MetricIntInt{metric, make(map[int]int)})
			result.vTime = append(result.vTime, MetricFloat64Int{metric, make(map[float64]int)})
			if metric != "backU" && metric != "backB" {
				result.vCW = append(result.vCW, MetricIntInt{metric, make(map[int]int)})
				if metric != "only-1" {
					result.vCWfirst = append(result.vCWfirst, MetricIntInt{metric, make(map[int]int)})
				}
				if metric == "rw" || metric == "all" || metric == "URW" {
					result.dApprovers = append(result.dApprovers, MetricIntInt{metric, make(map[int]int)})
				}
			}
		}
	}
	if param.SpineEnabled {
		for _, metric := range veloMetrics {
			if metric != "backU" && metric != "backB" {
				result.vID = append(result.vID, MetricIntInt{"S" + metric, make(map[int]int)})
				result.vTime = append(result.vTime, MetricFloat64Int{"S" + metric, make(map[float64]int)})
				if metric != "backG" {
					result.vCW = append(result.vCW, MetricIntInt{"S" + metric, make(map[int]int)})
					if metric != "only-1" {
						result.vCWfirst = append(result.vCWfirst, MetricIntInt{"S" + metric, make(map[int]int)})
					}
					if metric == "rw" || metric == "all" || metric == "URW" {
						result.dApprovers = append(result.dApprovers, MetricIntInt{"S" + metric, make(map[int]int)})
					}
				}
			}
		}
	}
	return &result
}

func (sim *Sim) runVelocityStat(result *velocityResult) {
	numParticles := 1000
	if sim.param.TSA != "RW" {
		sim.velocityURTS(result.vID[0].v, result.vTime[0].v, result.dApprovers[0].v, result.vCW[0].v, result.vCWfirst[0].v)
		sim.velocityAll(result.vID[1].v, result.vTime[1].v, result.dApprovers[1].v, result.vCW[1].v, result.vCWfirst[1].v)
		sim.velocityOfIndex(result.vID[2].v, result.vTime[2].v, result.vCW[2].v, result.vCWfirst[2].v, 1)
		sim.velocityOfIndex(result.vID[3].v, result.vTime[3].v, result.vCW[3].v, result.vCWfirst[3].v, -1)
		sim.velocityOfIndex(result.vID[4].v, result.vTime[4].v, result.vCW[4].v, result.vCWfirst[4].v, 2)
		sim.velocityOfIndex(result.vID[5].v, result.vTime[5].v, result.vCW[5].v, result.vCWfirst[5].v, 3)
		sim.velocityOfIndex(result.vID[6].v, result.vTime[6].v, result.vCW[6].v, result.vCWfirst[6].v, 4)
		sim.velocityOfOnlyIndex(result.vID[7].v, result.vTime[7].v, result.vCW[7].v, 1)
		sim.velocityCWLowUpBound(result.vID[8].v, result.vTime[8].v, result.vCW[8].v, result.vCWfirst[7].v, 1)
		sim.velocityCWLowUpBound(result.vID[9].v, result.vTime[9].v, result.vCW[9].v, result.vCWfirst[8].v, -1)
		sim.velocityCWLowUpBoundRW(result.vID[10].v, result.vTime[10].v, result.vCW[10].v, result.vCWfirst[9].v, 1)
		sim.velocityCWLowUpBoundRW(result.vID[11].v, result.vTime[11].v, result.vCW[11].v, result.vCWfirst[10].v, -1)
		sim.velocityBackURTS(result.vID[12].v, result.vTime[12].v)
	} else {
		//fmt.Println(len(result.vID), len(result.vTime), len(result.dApprovers), len(result.vCW), len(result.vCWfirst))
		sim.velocityParticleRW(result.vID[0].v, result.vTime[0].v, result.dApprovers[0].v, result.vCW[0].v, result.vCWfirst[0].v, numParticles, "default")
		sim.velocityAll(result.vID[1].v, result.vTime[1].v, result.dApprovers[1].v, result.vCW[1].v, result.vCWfirst[1].v)
		sim.velocityOfIndexRW(result.vID[2].v, result.vTime[2].v, result.vCW[2].v, result.vCWfirst[2].v, 1, numParticles)
		sim.velocityOfIndexRW(result.vID[3].v, result.vTime[3].v, result.vCW[3].v, result.vCWfirst[3].v, -1, numParticles)
		sim.velocityCWLowUpBoundRW(result.vID[4].v, result.vTime[4].v, result.vCW[4].v, result.vCWfirst[4].v, 1)
		sim.velocityCWLowUpBoundRW(result.vID[5].v, result.vTime[5].v, result.vCW[5].v, result.vCWfirst[5].v, -1)
		sim.velocityParticleBackRW(result.vID[6].v, result.vTime[6].v, URW{}, numParticles)
		sim.velocityParticleBackRW(result.vID[7].v, result.vTime[7].v, BRW{}, numParticles)
		sim.velocityParticleRW(result.vID[8].v, result.vTime[8].v, result.dApprovers[2].v, result.vCW[6].v, result.vCWfirst[6].v, numParticles, "URW")
		if sim.param.SpineEnabled {
			sim.velocityParticleRWSpine(result.vID[9].v, result.vTime[9].v, result.dApprovers[3].v, result.vCW[7].v, result.vCWfirst[7].v, numParticles, "default")
			sim.velocityAllSpine(result.vID[10].v, result.vTime[10].v, result.dApprovers[4].v, result.vCW[8].v, result.vCWfirst[8].v)
			sim.velocityOfIndexRWSpine(result.vID[11].v, result.vTime[11].v, result.vCW[9].v, result.vCWfirst[9].v, 1, numParticles)
			sim.velocityOfIndexRWSpine(result.vID[12].v, result.vTime[12].v, result.vCW[10].v, result.vCWfirst[10].v, -1, numParticles)
			sim.velocityCWLowUpBoundSpine(result.vID[13].v, result.vTime[13].v, result.vCW[11].v, result.vCWfirst[11].v, 1)
			sim.velocityCWLowUpBoundSpine(result.vID[14].v, result.vTime[14].v, result.vCW[12].v, result.vCWfirst[12].v, -1)
			sim.velocityParticleRWSpine(result.vID[15].v, result.vTime[15].v, result.dApprovers[5].v, result.vCW[13].v, result.vCWfirst[13].v, numParticles, "URW")
			sim.velocityParticleGhostBack(result.vID[16].v, result.vTime[16].v)
		}
	}

}

func (sim Sim) velocityURTS(v map[int]int, t map[float64]int, d map[int]int, w, wFirst map[int]int) {
	for i := sim.param.minCut; i < sim.param.maxCut; i++ {
		if len(sim.approvers[i]) > 0 {
			l := sim.generator.Intn(len(sim.approvers[i]))
			delta := sim.approvers[i][l] - i
			deltaTime := math.Round((sim.tangle[sim.approvers[i][l]].time-sim.tangle[i].time)*100) / 100
			v[delta]++
			d[l+1]++
			t[deltaTime]++
			deltaCW := sim.tangle[i].cw - sim.tangle[sim.approvers[i][l]].cw
			w[deltaCW]++
			if len(sim.approvers[i]) > 1 {
				wFirst[deltaCW]++
			}
			// if float64(delta)/sim.param.Lambda != deltaTime {
			// 	fmt.Println(sim.approvers[i][l], "-", i, float64(delta)/sim.param.Lambda, "|", deltaTime, sim.tangle[sim.approvers[i][l]].time, "-", sim.tangle[i].time)
			// }
		}
	}
	//fmt.Println(t)
}

func (sim Sim) velocityBackURTS(v map[int]int, t map[float64]int) {
	for i := sim.param.maxCut; i > sim.param.minCut; i-- {

		l := sim.generator.Intn(len(sim.tangle[i].ref))
		delta := i - sim.tangle[i].ref[l]
		deltaTime := math.Round((sim.tangle[i].time-sim.tangle[sim.tangle[i].ref[l]].time)*100) / 100
		v[delta]++
		t[deltaTime]++

	}
}

func (sim *Sim) velocityParticleRW(v map[int]int, t map[float64]int, d map[int]int, w, wFirst map[int]int, nParticles int, mode string) {
	for i := 0; i < nParticles; i++ {
		prev := sim.tangle[0]
		//start := sim.generator.Intn(sim.param.minCut)
		//prev := sim.tangle[start]
		//tsa := mode
		var tsa RandomWalker
		if mode == "default" && sim.param.Alpha != 0 {
			tsa = BRW{}
		} else {
			tsa = URW{}
		}

		for current, currentIdx := tsa.RandomWalk(prev, sim); len(sim.approvers[current.id]) > 0; current, currentIdx = tsa.RandomWalk(current, sim) {
			if current.id > sim.param.minCut && current.id < sim.param.maxCut {
				delta := current.id - prev.id
				v[delta]++
				d[currentIdx+1]++
				deltaTime := math.Round((current.time-prev.time)*100) / 100
				t[deltaTime]++
				deltaCW := prev.cw - current.cw
				w[deltaCW]++
				if len(sim.approvers[prev.id]) > 1 {
					wFirst[deltaCW]++
				}
			}
			prev = current
		}
	}
}

func (sim *Sim) velocityParticleRWSpine(v map[int]int, t map[float64]int, d map[int]int, w, wFirst map[int]int, nParticles int, mode string) {
	for i := 0; i < nParticles; i++ {
		prev := sim.spineTangle[0]
		//start := sim.generator.Intn(sim.param.minCut)
		//prev := sim.tangle[start]
		var tsa RandomWalker
		if mode == "default" && sim.param.Alpha != 0 {
			tsa = BRW{}
			//fmt.Println("Doing BRW")
		} else {
			tsa = URW{}
		}

		for current, currentIdx := tsa.RandomWalkSpine(prev, sim); len(sim.spineApprovers[current.id]) > 0; current, currentIdx = tsa.RandomWalkSpine(current, sim) {
			if current.id > sim.param.minCut && current.id < sim.param.maxCut {
				delta := current.id - prev.id
				v[delta]++
				d[currentIdx+1]++
				deltaTime := math.Round((current.time-prev.time)*100) / 100
				t[deltaTime]++
				deltaCW := prev.cw - current.cw
				w[deltaCW]++
				if len(sim.spineApprovers[prev.id]) > 1 {
					wFirst[deltaCW]++
				}
			}
			prev = current
		}
	}
}

func (sim *Sim) velocityParticleBackRW(v map[int]int, t map[float64]int, mode RandomWalker, nParticles int) {
	for i := 0; i < nParticles; i++ {
		//start := sim.generator.Intn(sim.param.maxCutrange) + sim.param.minCut - 1 // -1 just to be sure start is larger than TangleSize
		_, start := ghostWalk(sim.tangle[0], sim)
		prev := start

		tsa := mode

		//var tsa RandomWalker
		// if sim.param.Alpha != 0 {
		// 	tsa = BRW{}
		// } else {
		// 	tsa = URW{}
		// }

		for current := tsa.RandomWalkBack(prev, sim); current.id > sim.param.minCut; current = tsa.RandomWalkBack(current, sim) {
			if current.id > sim.param.minCut && current.id < sim.param.maxCut {
				delta := prev.id - current.id
				v[delta]++
				deltaTime := math.Round((prev.time-current.time)*100) / 100
				t[deltaTime]++
			}
			prev = current
		}
	}
}

func (sim *Sim) velocityParticleGhostBack(v map[int]int, t map[float64]int) {

	//fmt.Println("Doing backG")
	//start := sim.generator.Intn(sim.param.maxCutrange) + sim.param.minCut - 1 // -1 just to be sure start is larger than TangleSize
	_, start := ghostWalk(sim.tangle[0], sim)
	prev := start

	for current := ghostWalkBack(prev, sim); current.id > sim.param.minCut; current = ghostWalkBack(current, sim) {
		if current.id > sim.param.minCut && current.id < sim.param.maxCut {
			delta := prev.id - current.id
			v[delta]++
			deltaTime := math.Round((prev.time-current.time)*100) / 100
			t[deltaTime]++
		}
		prev = current
	}
}

func (sim Sim) velocityOfIndex(v map[int]int, t map[float64]int, w, wFirst map[int]int, index int) {
	for i := sim.param.minCut; i < sim.param.maxCut; i++ {
		if index > 0 && len(sim.approvers[i]) > index-1 {
			delta := sim.approvers[i][index-1] - i
			v[delta]++
			deltaTime := math.Round((sim.tangle[sim.approvers[i][index-1]].time-sim.tangle[i].time)*100) / 100
			t[deltaTime]++
			deltaCW := sim.tangle[i].cw - sim.tangle[sim.approvers[i][index-1]].cw
			w[deltaCW]++
			if len(sim.approvers[i]) > 1 {
				wFirst[deltaCW]++
			}
		} else if index < 0 && len(sim.approvers[i]) > 1 {
			delta := sim.approvers[i][len(sim.approvers[i])-1] - i
			v[delta]++
			deltaTime := math.Round((sim.tangle[sim.approvers[i][len(sim.approvers[i])-1]].time-sim.tangle[i].time)*100) / 100
			t[deltaTime]++
			deltaCW := sim.tangle[i].cw - sim.tangle[sim.approvers[i][len(sim.approvers[i])-1]].cw
			w[deltaCW]++
			if len(sim.approvers[i]) > 1 {
				wFirst[deltaCW]++
			}
		}
	}
}

func (sim Sim) velocityOfIndexSpine(v map[int]int, t map[float64]int, w, wFirst map[int]int, index int) {
	for i := sim.param.minCut; i < sim.param.maxCut; i++ {
		if index > 0 && len(sim.spineApprovers[i]) > index-1 {
			delta := sim.spineApprovers[i][index-1] - i
			v[delta]++
			deltaTime := math.Round((sim.spineTangle[sim.spineApprovers[i][index-1]].time-sim.spineTangle[i].time)*100) / 100
			t[deltaTime]++
			deltaCW := sim.spineTangle[i].cw - sim.spineTangle[sim.spineApprovers[i][index-1]].cw
			w[deltaCW]++
			if len(sim.spineApprovers[i]) > 1 {
				wFirst[deltaCW]++
			}
		} else if index < 0 && len(sim.spineApprovers[i]) > 1 {
			delta := sim.spineApprovers[i][len(sim.spineApprovers[i])-1] - i
			v[delta]++
			deltaTime := math.Round((sim.spineTangle[sim.spineApprovers[i][len(sim.spineApprovers[i])-1]].time-sim.spineTangle[i].time)*100) / 100
			t[deltaTime]++
			deltaCW := sim.spineTangle[i].cw - sim.spineTangle[sim.spineApprovers[i][len(sim.spineApprovers[i])-1]].cw
			w[deltaCW]++
			if len(sim.spineApprovers[i]) > 1 {
				wFirst[deltaCW]++
			}
		}
	}
}

func (sim *Sim) velocityOfIndexRW(v map[int]int, t map[float64]int, w, wFirst map[int]int, index int, nParticles int) {

	//for i := 0; i < nParticles; i++ {
	start := 0 //sim.generator.Intn(sim.param.minCut)

	for current := sim.tangle[start]; len(sim.approvers[current.id]) > 0 && current.id < sim.param.maxCut; {
		if index > 0 && len(sim.approvers[current.id]) > index-1 {
			delta := sim.approvers[current.id][index-1] - current.id
			deltaTime := math.Round((sim.tangle[sim.approvers[current.id][index-1]].time-sim.tangle[current.id].time)*100) / 100
			deltaCW := sim.tangle[current.id].cw - sim.tangle[sim.approvers[current.id][index-1]].cw
			if current.id > sim.param.minCut {
				v[delta]++
				t[deltaTime]++
				w[deltaCW]++
				if len(sim.approvers[current.id]) > 1 {
					wFirst[deltaCW]++
				}
			}
			current = sim.tangle[sim.approvers[current.id][index-1]]
		} else if index < 0 && len(sim.approvers[current.id]) > 0 {

			delta := sim.approvers[current.id][len(sim.approvers[current.id])-1] - current.id
			deltaTime := math.Round((sim.tangle[sim.approvers[current.id][len(sim.approvers[current.id])-1]].time-sim.tangle[current.id].time)*100) / 100
			deltaCW := sim.tangle[current.id].cw - sim.tangle[sim.approvers[current.id][len(sim.approvers[current.id])-1]].cw
			if current.id > sim.param.minCut {
				v[delta]++
				t[deltaTime]++
				w[deltaCW]++
				if len(sim.approvers[current.id]) > 1 {
					wFirst[deltaCW]++
				}
			}
			current = sim.tangle[sim.approvers[current.id][len(sim.approvers[current.id])-1]]
		} else {
			break
		}
	}
	//}
}

func (sim *Sim) velocityOfIndexRWSpine(v map[int]int, t map[float64]int, w, wFirst map[int]int, index int, nParticles int) {

	//for i := 0; i < nParticles; i++ {
	start := 0 //sim.generator.Intn(sim.param.minCut)

	for current := sim.spineTangle[start]; len(sim.spineApprovers[current.id]) > 0 && current.id < sim.param.maxCut; {
		if index > 0 && len(sim.spineApprovers[current.id]) > index-1 { //lower bound for index = 1
			delta := sim.spineApprovers[current.id][index-1] - current.id
			deltaTime := math.Round((sim.spineTangle[sim.spineApprovers[current.id][index-1]].time-sim.spineTangle[current.id].time)*100) / 100
			deltaCW := sim.spineTangle[current.id].cw - sim.spineTangle[sim.spineApprovers[current.id][index-1]].cw
			if current.id > sim.param.minCut && current.id < sim.param.maxCut {
				v[delta]++
				t[deltaTime]++
				w[deltaCW]++
				if len(sim.spineApprovers[current.id]) > 1 {
					wFirst[deltaCW]++
				}
			}
			current = sim.spineTangle[sim.spineApprovers[current.id][index-1]]
		} else if index < 0 && len(sim.spineApprovers[current.id]) > 0 { //upper bound
			delta := sim.spineApprovers[current.id][len(sim.spineApprovers[current.id])-1] - current.id
			deltaTime := math.Round((sim.spineTangle[sim.spineApprovers[current.id][len(sim.spineApprovers[current.id])-1]].time-sim.spineTangle[current.id].time)*100) / 100
			deltaCW := sim.spineTangle[current.id].cw - sim.spineTangle[sim.spineApprovers[current.id][len(sim.spineApprovers[current.id])-1]].cw
			if current.id > sim.param.minCut && current.id < sim.param.maxCut {
				v[delta]++
				t[deltaTime]++
				w[deltaCW]++
				if len(sim.spineApprovers[current.id]) > 1 {
					wFirst[deltaCW]++
				}
			}
			current = sim.spineTangle[sim.spineApprovers[current.id][len(sim.spineApprovers[current.id])-1]]
		} else {
			break
		}
	}
	//}
}

func (sim Sim) velocityCWLowUpBound(v map[int]int, t map[float64]int, w, wFirst map[int]int, index int) {
	var current, prev Tx
	for i := sim.param.minCut; i < sim.param.maxCut; i++ {
		prev = sim.tangle[i]
		if len(sim.approvers[i]) > 0 {
			var cws []int
			for _, approver := range sim.approvers[i] {
				cws = append(cws, sim.tangle[approver].cw)
			}
			if index > 0 {
				//select approvers with max cw
				maxCW, _ := max(cws)
				current = sim.tangle[sim.approvers[i][maxCW]]
			} else if index < 0 {
				//select approvers with min cw
				minCW, _ := min(cws)
				current = sim.tangle[sim.approvers[i][minCW]]
			}
		}
		delta := current.id - prev.id
		deltaTime := math.Round((current.time-prev.time)*100) / 100
		deltaCW := prev.cw - current.cw
		if current.id > sim.param.minCut && current.id < sim.param.maxCut {
			v[delta]++
			t[deltaTime]++
			w[deltaCW]++
			if len(sim.approvers[current.id]) > 1 {
				wFirst[deltaCW]++
			}
		}
	}
}

func (sim *Sim) velocityCWLowUpBoundRW(v map[int]int, t map[float64]int, w, wFirst map[int]int, index int) {

	start := 0 //sim.generator.Intn(sim.param.minCut)

	for current := sim.tangle[start]; len(sim.approvers[current.id]) > 0 && current.id < sim.param.maxCut; {
		prev := current
		if index > 0 && len(sim.approvers[current.id]) > index-1 { //lower bound for index = 1
			//select approvers with max cw
			var cws []int
			for _, approver := range sim.approvers[current.id] {
				cws = append(cws, sim.tangle[approver].cw)
			}
			maxCW, _ := max(cws)
			current = sim.tangle[sim.approvers[current.id][maxCW]]
		} else if index < 0 && len(sim.approvers[current.id]) > 0 { //upper bound
			//select approvers with min cw
			var cws []int
			for _, approver := range sim.approvers[current.id] {
				cws = append(cws, sim.tangle[approver].cw)
			}
			minCW, _ := min(cws)
			current = sim.tangle[sim.approvers[current.id][minCW]]
		}
		delta := current.id - prev.id
		deltaTime := math.Round((current.time-prev.time)*100) / 100
		deltaCW := prev.cw - current.cw
		if current.id > sim.param.minCut && current.id < sim.param.maxCut {
			v[delta]++
			t[deltaTime]++
			w[deltaCW]++
			if len(sim.approvers[current.id]) > 1 {
				wFirst[deltaCW]++
			}
		}
	}
}

func (sim *Sim) velocityCWLowUpBoundSpine(v map[int]int, t map[float64]int, w, wFirst map[int]int, index int) {

	start := 0 //sim.generator.Intn(sim.param.minCut)

	for current := sim.spineTangle[start]; len(sim.spineApprovers[current.id]) > 0 && current.id < sim.param.maxCut; {
		prev := current
		if index > 0 && len(sim.spineApprovers[current.id]) > index-1 { //lower bound for index = 1
			//select approvers with max cw
			var cws []int
			for _, approver := range sim.spineApprovers[current.id] {
				cws = append(cws, sim.spineTangle[approver].cw)
			}
			maxCW, _ := max(cws)
			current = sim.spineTangle[sim.spineApprovers[current.id][maxCW]]
		} else if index < 0 && len(sim.spineApprovers[current.id]) > 0 { //upper bound
			//select approvers with min cw
			var cws []int
			for _, approver := range sim.spineApprovers[current.id] {
				cws = append(cws, sim.spineTangle[approver].cw)
			}
			minCW, _ := min(cws)
			current = sim.spineTangle[sim.spineApprovers[current.id][minCW]]
		}
		delta := current.id - prev.id
		deltaTime := math.Round((current.time-prev.time)*100) / 100
		deltaCW := prev.cw - current.cw
		if current.id > sim.param.minCut && current.id < sim.param.maxCut {
			v[delta]++
			t[deltaTime]++
			w[deltaCW]++
			if len(sim.spineApprovers[current.id]) > 1 {
				wFirst[deltaCW]++
			}
		}
	}
}

func (sim Sim) velocityOfOnlyIndex(v map[int]int, t map[float64]int, w map[int]int, index int) {
	for i := sim.param.minCut; i < sim.param.maxCut; i++ {
		if index > 0 && len(sim.approvers[i]) == index {
			delta := sim.approvers[i][index-1] - i
			v[delta]++
			deltaTime := math.Round((sim.tangle[sim.approvers[i][index-1]].time-sim.tangle[i].time)*100) / 100
			t[deltaTime]++
			deltaCW := sim.tangle[i].cw - sim.tangle[sim.approvers[i][index-1]].cw
			w[deltaCW]++
		}
	}
}

func (sim Sim) velocityAll(v map[int]int, t map[float64]int, d map[int]int, w, wFirst map[int]int) {
	for i := sim.param.minCut; i < sim.param.maxCut; i++ {

		d[len(sim.approvers[i])]++
		for _, a := range sim.approvers[i] {
			delta := a - i
			v[delta]++
			deltaTime := math.Round((sim.tangle[a].time-sim.tangle[i].time)*100) / 100
			t[deltaTime]++
			deltaCW := sim.tangle[i].cw - sim.tangle[a].cw
			w[deltaCW]++
			if len(sim.approvers[i]) > 1 {
				wFirst[deltaCW]++
			}
		}
	}
}

func (sim Sim) velocityAllSpine(v map[int]int, t map[float64]int, d map[int]int, w, wFirst map[int]int) {
	//for i := sim.param.minCut; i < sim.param.maxCut; i++ {
	for i := range sliceMap(sim.spineTangle, sim.param.minCut, sim.param.maxCut) {
		d[len(sim.spineApprovers[i])]++
		for _, a := range sim.spineApprovers[i] {
			delta := a - i
			v[delta]++
			deltaTime := math.Round((sim.spineTangle[a].time-sim.spineTangle[i].time)*100) / 100
			t[deltaTime]++
			deltaCW := sim.spineTangle[i].cw - sim.spineTangle[a].cw
			w[deltaCW]++
			if len(sim.spineApprovers[i]) > 1 {
				wFirst[deltaCW]++
			}
		}
	}
}

func (velo *velocityResult) Join(b velocityResult) (r velocityResult) {
	if velo.vID == nil {
		return b
	}
	for i := range b.vID {
		r.vID = append(r.vID, joinMapMetricIntInt(velo.vID[i], b.vID[i]))
	}
	for i := range b.dApprovers {
		r.dApprovers = append(r.dApprovers, joinMapMetricIntInt(velo.dApprovers[i], b.dApprovers[i]))
	}
	for i := range b.vTime {
		r.vTime = append(r.vTime, joinMapMetricFloat64Int(velo.vTime[i], b.vTime[i]))
	}

	for i := range b.vCW {
		// r.vCW = append(r.vCW, joinMapStatInt(velo.vCW[i], b.vCW[i]))
		r.vCW = append(r.vCW, joinMapMetricIntInt(velo.vCW[i], b.vCW[i]))
	}

	for i := range b.vCWfirst {
		// r.vCWfirst = append(r.vCWfirst, joinMapStatInt(velo.vCWfirst[i], b.vCWfirst[i]))
		r.vCWfirst = append(r.vCWfirst, joinMapMetricIntInt(velo.vCWfirst[i], b.vCWfirst[i]))
	}

	return r
}

func (velo velocityResult) Save(p Parameters) (err error) {
	if err = velo.SaveVID(p); err != nil {
		return err
	}
	if err = velo.SaveVTime(p); err != nil {
		return err
	}
	if err = velo.saveApprovers(p); err != nil {
		return err
	}
	if err = velo.saveCW(p); err != nil {
		return err
	}
	if err = velo.saveCWfirst(p); err != nil {
		return err
	}

	return err
}

func (velo velocityResult) SaveStat(p Parameters) (err error) {
	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.4f", p.Alpha)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/velocity_stat_" + p.TSA + "_" + rateType +
		"_lambda_" + lambdaStr +
		"_alpha_" + alphaStr + "_.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(velo.Stat(p)) // writing...

	if err != nil {
		fmt.Printf("error writing string: %v", err)
		return err
	}

	return nil
}

func (velo velocityResult) Stat(p Parameters) (result string) {
	result = velo.StatVID(p)
	result += "\n"
	result += velo.StatVTime(p)
	result += "\n"
	result += velo.StatCW(p)
	result += "\n"
	result += velo.StatCWfirst(p)
	result += "\n"
	result += velo.StatApprovers(p)

	return result
}

// ToString converts a MetricIntInt to a string
func (s MetricIntInt) ToString(p Parameters, normalized bool) (result string) {
	if len(s.v) == 0 {
		//fmt.Println("Empty " + s.desc + " slice")
		return ""
	}
	var keys []int
	var datapoints int
	for k := range s.v {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// calculate statistics
	var weigths []float64
	var x []float64
	for k := range keys {
		norm := 1.
		if normalized {
			norm = p.Lambda
		}
		x = append(x, float64(keys[k])/norm)
		weigths = append(weigths, float64(s.v[keys[k]]))
		datapoints = datapoints + s.v[keys[k]]
	}

	//fmt.Println(s.desc, x, weigths)

	var avg, std = stat.MeanStdDev(x, weigths)
	_, variance := stat.MeanVariance(x, weigths)
	skew := stat.Skew(x, weigths)
	mode, _ := stat.Mode(x, weigths)
	//median := median(x, weigths)
	median := stat.Quantile(0.5, stat.Empirical, x, weigths)
	//median := 0.

	//result += fmt.Sprintf("%s\n", s.desc)

	if variance > 10000 {
		result += fmt.Sprintf("%s\t\t%.2f\t\t%.4f\t\t%.2f\t\t%.2f\t\t%.2f\t%.2f\t\t%.2f\t\t%.3f\t\t%.2f\t\t%.2f\t\t%d\n", s.desc, p.Lambda, p.Alpha, avg, std, variance, median, mode, skew, x[0], x[len(x)-1], datapoints)
	} else {
		result += fmt.Sprintf("%s\t\t%.2f\t\t%.4f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.3f\t\t%.2f\t\t%.2f\t\t%d\n", s.desc, p.Lambda, p.Alpha, avg, std, variance, median, mode, skew, x[0], x[len(x)-1], datapoints)
	}
	return result
}

// ToString converts a MetricFloat64Int to a string
func (s MetricFloat64Int) ToString(p Parameters, normalized bool) (result string) {
	if len(s.v) == 0 {
		//fmt.Println("Empty " + s.desc + " slice")
		return ""
	}
	var keys []float64
	var datapoints int
	for k := range s.v {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	// calculate statistics
	var weigths []float64
	var x []float64
	for k := range keys {
		norm := 1.
		if normalized {
			norm = p.Lambda
		}
		x = append(x, float64(keys[k])/norm)
		weigths = append(weigths, float64(s.v[keys[k]]))
		datapoints = datapoints + s.v[keys[k]]
	}

	var avg, std = stat.MeanStdDev(x, weigths)
	_, variance := stat.MeanVariance(x, weigths)
	skew := stat.Skew(x, weigths)
	mode, _ := stat.Mode(x, weigths)
	//median := median(x, weigths)
	median := stat.Quantile(0.5, stat.Empirical, x, weigths)
	//median := 0.0
	//result += fmt.Sprintf("%s\n", s.desc)

	if variance > 10000 {
		result += fmt.Sprintf("%s\t\t%.2f\t\t%.4f\t\t%.2f\t\t%.2f\t\t%.2f\t%.2f\t\t%.2f\t\t%.3f\t\t%.2f\t\t%.2f\t\t%d\n", s.desc, p.Lambda, p.Alpha, avg, std, variance, median, mode, skew, x[0], x[len(x)-1], datapoints)
	} else {
		result += fmt.Sprintf("%s\t\t%.2f\t\t%.4f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.3f\t\t%.2f\t\t%.2f\t\t%d\n", s.desc, p.Lambda, p.Alpha, avg, std, variance, median, mode, skew, x[0], x[len(x)-1], datapoints)
	}
	return result
}

// Save saves a MetricIntInt on a file
func (s MetricIntInt) Save(p Parameters, aType, target string, normalized bool) error {
	var keys []int
	var datapoints int
	for k := range s.v {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// calculate statistics
	var weigths []float64
	var x []float64
	norm := 1.
	for k := range keys {
		if normalized {
			norm = p.Lambda
		}
		x = append(x, float64(keys[k])/norm)
		weigths = append(weigths, float64(s.v[keys[k]]))
		datapoints = datapoints + s.v[keys[k]]
	}
	// save to file for plot

	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.4f", p.Alpha)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/" + aType + "_" + target + "_" + p.TSA + "_" + rateType + "_" + s.desc +
		"_lambda_" + lambdaStr +
		"_alpha_" + alphaStr + "_.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()
	for i, k := range x {
		//fmt.Println("Key:", k, "Value:", m[k])
		if target == "approvers" || !normalized {
			_, err = f.WriteString(fmt.Sprintf("%d\t%f\n", int(k), weigths[i]/float64(datapoints)*norm)) // writing...
		} else {
			_, err = f.WriteString(fmt.Sprintf("%f\t%f\n", k, weigths[i]/float64(datapoints)*norm)) // writing...
		}
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
	return nil
}

// Save saves a MetricFloat64Int as a file
func (s MetricFloat64Int) Save(p Parameters, target string, normalized bool) error {
	var keys []float64
	var datapoints int
	for k := range s.v {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	var weigths []float64
	var x []float64
	norm := 1.
	for k := range keys {
		if normalized {
			norm = p.Lambda
		}
		x = append(x, float64(keys[k])/norm)
		weigths = append(weigths, float64(s.v[keys[k]]))
		datapoints = datapoints + s.v[keys[k]]
	}
	// save to file for plot

	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.4f", p.Alpha)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/velocity_" + target + "_" + p.TSA + "_" + rateType + "_" + s.desc +
		"_lambda_" + lambdaStr +
		"_alpha_" + alphaStr + "_.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()
	for i, k := range x {
		_, err = f.WriteString(fmt.Sprintf("%f\t%f\n", k, weigths[i]/float64(datapoints)*norm)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
	return nil
}

func (velo velocityResult) SaveVID(p Parameters) error {
	for _, velocity := range velo.vID {
		velocity.Save(p, "velocity", "ID", true)
	}
	return nil
}

func (velo velocityResult) SaveVTime(p Parameters) error {
	for _, velocity := range velo.vTime {
		velocity.Save(p, "time", false)
	}
	return nil
}

func (velo velocityResult) saveApprovers(p Parameters) error {
	for _, velocity := range velo.dApprovers {
		velocity.Save(p, "velocity", "approvers", false)
	}
	return nil
}

func (velo velocityResult) saveCW(p Parameters) error {
	for _, velocity := range velo.vCW {
		velocity.Save(p, "velocity", "cw", true)
	}
	return nil
}

func (velo velocityResult) saveCWfirst(p Parameters) error {
	for _, velocity := range velo.vCWfirst {
		velocity.Save(p, "velocity", "cw-first", true)
	}
	return nil
}

func (velo velocityResult) StatVID(p Parameters) (result string) {
	result += "#Velocity Stats [ID]\n"
	result += "#Stat Type\tLambda\t\tAlpha\t\tMean\t\tStdDev\t\tVariance\tMedian\t\tMode\t\tSkew\t\tMinVal\t\tMaxVal\t\tN\n"
	for _, velocity := range velo.vID {
		result += velocity.ToString(p, true)
	}
	return result
}

func (velo velocityResult) StatApprovers(p Parameters) (result string) {
	result += "#Direct Approvers Stats\n"
	result += "#Stat Type\tLambda\t\tAlpha\t\tMean\t\tStdDev\t\tVariance\tMedian\t\tMode\t\tSkew\t\tMinVal\t\tMaxVal\t\tN\n"
	for _, velocity := range velo.dApprovers {
		result += velocity.ToString(p, false)
	}
	return result

}

func (velo velocityResult) StatVTime(p Parameters) (result string) {
	result += "#Velocity Stats [Time]\n"
	result += "#Stat Type\tLambda\t\tAlpha\t\tMean\t\tStdDev\t\tVariance\tMedian\t\tMode\t\tSkew\t\tMinVal\t\tMaxVal\t\tN\n"
	for _, velocity := range velo.vTime {
		result += velocity.ToString(p, false)
	}
	return result

}

func (velo velocityResult) StatCW(p Parameters) (result string) {
	result += "#Velocity Stats CW\n"
	result += "#Stat Type\tLambda\t\tAlpha\t\tMean\t\tStdDev\t\tVariance\tMedian\t\tMode\t\tSkew\t\tMinVal\t\tMaxVal\t\tN\n"
	for _, velocity := range velo.vCW {
		result += velocity.ToString(p, true)
	}
	return result
}

func (velo velocityResult) StatCWfirst(p Parameters) (result string) {
	result += "#Velocity Stats CW-first\n"
	result += "#Stat Type\tLambda\t\tAlpha\t\tMean\t\tStdDev\t\tVariance\tMedian\t\tMode\t\tSkew\t\tMinVal\t\tMaxVal\t\tN\n"
	for _, velocity := range velo.vCWfirst {
		result += velocity.ToString(p, true)
	}
	return result
}

func (velo velocityResult) getTimeMetric(s string) (x, y []float64) {
	for _, velocity := range velo.vTime {
		if velocity.desc == s {
			return mapFloatIntToSlices(velocity.v, 1)
		}
	}
	return nil, nil
}

func mapIntIntToSlices(a map[int]int, normFactor float64) (x, y []float64) {
	if normFactor == 0 {
		panic("Normalization Factor cannot be 0")
	}
	var keys []int
	for k := range a {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for k := range keys {
		x = append(x, float64(keys[k])/normFactor)
		y = append(y, float64(a[keys[k]]))
	}
	return x, y
}

func mapFloatIntToSlices(a map[float64]int, normFactor float64) (x, y []float64) {
	if normFactor == 0 {
		panic("Normalization Factor cannot be 0")
	}
	var keys []float64
	for k := range a {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	for k := range keys {
		x = append(x, float64(keys[k])/normFactor)
		y = append(y, float64(a[keys[k]]))
	}
	return x, y
}
