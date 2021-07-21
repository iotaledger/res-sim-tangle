// measure app stats during RW

package main

import (
	"fmt"
	"os"
	"sort"
)

//ExitProb result of simulation
type AppStatsAllResult struct {
	totalNum int
	// note: need to be careful with  Float mapping because of rounding errors -> if it does not work map int
	Num  map[int]float64 // # of occurances
	Prob map[int]float64 // finally calculated probability
}

// variable initialization
func newAppStatsAllResult() AppStatsAllResult {
	var result AppStatsAllResult
	result.Num = make(map[int]float64)
	result.Prob = make(map[int]float64)
	return result
}

func (sim *Sim) evalTangle_AppStatsAll(r *AppStatsAllResult) {
	for i1 := sim.param.MinCut; i1 < sim.param.MaxCut; i1++ {
		r.totalNum++
		r.Num[len(sim.tangle[i1].app)]++
	}
}

// join the results from several simulations
func (r *AppStatsAllResult) Join(b AppStatsAllResult) {
	r.totalNum += b.totalNum
	r.Num = joinMapIntFloat64(r.Num, b.Num)
	r.Prob = joinMapIntFloat64(r.Prob, b.Prob)
	return
}

// evaluate probabilities
func (r *AppStatsAllResult) finalprocess() error {
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
func (r *AppStatsAllResult) Save(p Parameters) (err error) {
	if err = r.SaveToFile(p, "Num", r.Num, true); err != nil {
		return err
	}
	if err = r.SaveToFile(p, "Prob", r.Prob, true); err != nil {
		return err
	}
	return err
}

// save a MapIntFloat64 as a file
func (r *AppStatsAllResult) SaveToFile(p Parameters, target string, datavec map[int]float64, normalized bool) error {
	var keys []int
	// var datapoints int
	for k := range datavec {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	f, err := os.Create("data/AppStatsAll__" + target + "_tsa=" + p.TSA + "_lambda" + lambdaStr + ".txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()
	// for i, k := range x {
	for _, k := range keys {
		_, err = f.WriteString(fmt.Sprintf("%d\t%f\n", k, datavec[k])) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
	return nil
}
