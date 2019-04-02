// measure distance of RWs of the Tangle

package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

//ExitProb result of simulation
type DistRWsResult struct {
	totalNum int
	// note: need to be careful with  Float mapping because of rounding errors -> if it does not work map int
	Num  map[float64]float64 // # of occurances
	Prob map[float64]float64 // finally calculated probability
}

// variable initialization
func newDistRWsResult() DistRWsResult {
	var result DistRWsResult
	result.Num = make(map[float64]float64)
	result.Prob = make(map[float64]float64)
	// result.countersuccess = append(result.countersuccess, MetricIntFloat64{metric, make(map[int]float64)})
	return result
}

func (sim *Sim) evalTangle_DistRWs(r *DistRWsResult) {
	// Define what we compare to. For now assume the limit lambda->infty
	pRef := make([]float64, 30)
	lamU := 1. // lambda large
	sumpRef := 0.
	pRef[0] = 0.
	for i1 := 1; i1 < len(pRef); i1++ {
		pRef[i1] = math.Exp(-1) / float64(Factorial(float64(i1-1))) // URTS
		if sim.param.TSA == "RW" {
			pRef[i1] = pRef[i1] * gfunc(lamU, i1)
		}
		sumpRef += pRef[i1]
	}
	if sim.param.TSA == "RW" {
		for i1 := 1; i1 < len(pRef); i1++ {
			pRef[i1] = pRef[i1] / sumpRef
		}
	}

	if sim.param.minCut+sim.param.DistRWsSampleLength*2*int(sim.param.Lambda) > sim.param.maxCut {
		fmt.Println("DistRWsSampleLength too large for the available dataset")
		pauseit()
	}

	tsa := BRW{}
	for i1 := 0; i1 < sim.param.DistRWsSampleRWNum; i1++ {
		currentID := 0
		currentSampleID := 0
		sampleSet := make([]int, sim.param.DistRWsSampleLength)
		NumThisSet := make([]float64, len(pRef))
		// 	for i2 := 0; i2 < len(sampleSet); i2++ { // reset SampleSet
		// 	sampleSet[i2] = 0
		// }

		for currentID < sim.param.maxCut {
			currentID, _ = tsa.RandomWalkStep(currentID, sim) // orphaned tips are included here
			if currentID > sim.param.minCut {

				if currentSampleID > sim.param.DistRWsSampleLength {
					// calculate for previous set
					dist := calcDistRWs(&NumThisSet, &pRef)
					intervals := float64(sim.param.DistSlicesResolution)
					r.Num[float64(int(intervals*dist))/intervals]++
					r.Prob[float64(int(intervals*dist))/intervals] = 0 // just to create a cell for this entry
					r.totalNum++
				}

				// update set
				currentSampleIDnow := int(math.Mod(float64(currentSampleID), float64(sim.param.DistRWsSampleLength)))
				if currentSampleID > sim.param.DistRWsSampleLength { //have at least filled sampleSet once
					NumThisSet[sampleSet[currentSampleIDnow]]-- // -1 from old value
				}
				sampleSet[currentSampleIDnow] = len(sim.approvers[currentID])
				NumThisSet[sampleSet[currentSampleIDnow]]++ // +1 to this value
				currentSampleID++

			}
			if len(sim.approvers[currentID]) == 0 { // if we ended at a tip, stop the loop
				currentID = sim.param.TangleSize + 1
			}
		}
	}
}

func gfunc(lamU float64, n int) float64 {
	g := 0.
	a := 1.3
	for i1 := 0; i1 < n+1; i1++ {
		plusval := math.Exp(-0.5*a*lamU) * math.Pow(1+0.5*a, float64(n-i1))
		minusval := math.Exp(0.5*a*lamU) * math.Pow(1-0.5*a, float64(n-i1))
		g += math.Pow(lamU, float64(-i1-1)) * Factorial(float64(n)) / Factorial(float64(n-i1)) * (minusval - plusval) / a
	}
	return g
}

func calcDistRWs(OccurancesVec, pRef *[]float64) float64 {
	num := 0.
	for i1 := 0; i1 < len(*pRef); i1++ {
		num += (*OccurancesVec)[i1]
	}
	l1 := 0.
	for i1 := 0; i1 < len(*pRef); i1++ {
		l1 += math.Abs((*pRef)[i1]-(*OccurancesVec)[i1]/num) / 2. // L1 distance
	}
	return l1
}

// join the results from several simulations
func (r *DistRWsResult) Join(b DistRWsResult) {
	r.totalNum += b.totalNum
	r.Num = joinMapFloat64Float64(r.Num, b.Num)
	r.Prob = joinMapFloat64Float64(r.Prob, b.Prob)
	return
}

// evaluate probabilities
func (r *DistRWsResult) finalprocess() error {
	checknum := 0
	for key, val := range r.Num {
		checknum += int(val)
		r.Prob[key] = val / float64(r.totalNum)
	}
	if checknum != r.totalNum {
		fmt.Println("This should not happen: checknum!=r.totalNum")
		fmt.Println("checknum = ", checknum)
		fmt.Println("r.totalNum = ", r.totalNum)
		pauseit()
	}
	return nil
}

// - - - - - - - - - - - -
// - - - save process - - -
// - - - - - - - - - - - -

// organise save
func (r *DistRWsResult) Save(p Parameters) (err error) {
	if err = r.SaveToFile(p, "Num", r.Num, true); err != nil {
		return err
	}
	if err = r.SaveToFile(p, "Prob", r.Prob, true); err != nil {
		return err
	}
	return err
}

// save a MapFloat64Float64 as a file
func (r *DistRWsResult) SaveToFile(p Parameters, target string, datavec map[float64]float64, normalized bool) error {
	var keys []float64
	// var datapoints int
	for k := range datavec {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.4f", p.Alpha)
	RWsStr := fmt.Sprintf("%d", p.DistRWsSampleLength)
	resStr := fmt.Sprintf("%d", p.DistRWsResolution)
	f, err := os.Create("data/DistRWs__" + target + "_tsa=" + p.TSA + "_lambda" + lambdaStr + "_alpha" + alphaStr + "__RWlength" + RWsStr + "_res" + resStr + ".txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()
	// for i, k := range x {
	for _, k := range keys {
		_, err = f.WriteString(fmt.Sprintf("%f\t%f\n", k, datavec[k])) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
	return nil
}
