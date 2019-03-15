// measure distance of slices of the Tangle

package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

//ExitProb result of simulation
type DistSlicesResult struct {
	totalNum int
	// note: need to be careful with  Float mapping because of rounding errors -> if it does not work map int
	Num  map[float64]float64 // # of occurances
	Prob map[float64]float64 // finally calculated probability
}

// variable initialization
func newDistSlicesResult() *DistSlicesResult {
	var result DistSlicesResult
	result.Num = make(map[float64]float64)
	result.Prob = make(map[float64]float64)
	// result.countersuccess = append(result.countersuccess, MetricIntFloat64{metric, make(map[int]float64)})
	return &result
}

func (sim *Sim) evalTangle_DistSlices(r *DistSlicesResult) {
	// Define what we compare to
	// URTS  // for now only compare to URTS
	pRef := make([]float64, 30)
	pRef[0] = 0.
	for i1 := 1; i1 < len(pRef); i1++ {
		pRef[i1] = math.Exp(-1) / float64(Factorial(float64(i1-1)))
	}

	// get approver stat for every slice
	sliceID := -1
	NumThisSlice := make([]float64, len(pRef))
	TotalNumThisSlice := 0
	SliceRestTime := ModFloat(sim.tangle[sim.param.minCut].time, sim.param.DistSlicesLength)
	SliceRestTimeNew := 0.
	SliceTime := sim.tangle[sim.param.minCut].time - SliceRestTime
	for i1 := sim.param.minCut; i1 < sim.param.maxCut; i1++ { //note, last slice is not considered because it may be partly above maxCut
		SliceRestTimeNew = sim.tangle[i1].time - SliceTime
		if SliceRestTimeNew > sim.param.DistSlicesLength { //last sliceID finished
			if sliceID > -1 { //discard the first slice because parts of it may be partly below minCut
				dist := calcDist(&NumThisSlice, &pRef, TotalNumThisSlice)
				intervals := float64(sim.param.DistSlicesResolution)
				r.Num[float64(int(intervals*dist))/intervals]++
				r.Prob[float64(int(intervals*dist))/intervals] = 0 // just to create a cell for this entry
				r.totalNum++
				for i2 := 0; i2 < len(NumThisSlice); i2++ { //reset slice to zeros
					NumThisSlice[i2] = 0
				}
				TotalNumThisSlice = 0
				SliceRestTime = ModFloat(sim.tangle[i1].time, sim.param.DistSlicesLength)
				SliceTime = sim.tangle[i1].time - SliceRestTime
			}
			sliceID++
		}
		NumThisSlice[len(sim.approvers[i1])]++
		TotalNumThisSlice++
	}
}

func calcDist(NumThisSlice, pRef *[]float64, num int) float64 {
	l1 := 0.
	for i1 := 0; i1 < len(*pRef); i1++ {
		l1 += math.Abs((*pRef)[i1] - (*NumThisSlice)[i1]/float64(num)) // L1 distance
	}
	return l1
}

// join the results from several simulations
func (r *DistSlicesResult) Join(b DistSlicesResult) {
	r.totalNum += b.totalNum
	r.Num = joinMapFloat64Float64(r.Num, b.Num)
	r.Prob = joinMapFloat64Float64(r.Prob, b.Prob)
	return
}

// evaluate probabilities
func (r *DistSlicesResult) finalprocess() error {
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
func (r DistSlicesResult) Save(p Parameters) (err error) {
	if err = SaveToFile(p, "Num", r.Num, true); err != nil {
		return err
	}
	if err = SaveToFile(p, "Prob", r.Prob, true); err != nil {
		return err
	}
	return err
}

// save a MapFloat64Float64 as a file
func SaveToFile(p Parameters, target string, datavec map[float64]float64, normalized bool) error {
	var keys []float64
	// var datapoints int
	for k := range datavec {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.4f", p.Alpha)
	sliceStr := fmt.Sprintf("%.2f", p.DistSlicesLength)
	resStr := fmt.Sprintf("%d", p.DistSlicesResolution)
	f, err := os.Create("data/DistSlices__" + target + "_tsa=" + p.TSA + "_lambda" + lambdaStr + "_alpha" + alphaStr + "__slicelength" + sliceStr + "_res" + resStr + ".txt")
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
