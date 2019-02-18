// Analyse the approver distribution with the time difference from a tx y in the past cone of that tx y

package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

//FocusRW result of simulation
type FocusRWResult struct { //this slices hold the statistics for each approver number mapping over all deltat's
	countersuccess []MetricFloat64Float64
	countertotal   []MetricFloat64Float64
	prob           []MetricFloat64Float64
}

//??? use string to create empty value maps to
func newFocusRWResult(Metrics []string) *FocusRWResult {
	// variables initialization for FocusRW
	var result FocusRWResult
	for _, metric := range Metrics {
		result.countersuccess = append(result.countersuccess, MetricFloat64Float64{metric, make(map[float64]float64)})
		result.countertotal = append(result.countertotal, MetricFloat64Float64{metric, make(map[float64]float64)})
		result.prob = append(result.prob, MetricFloat64Float64{metric, make(map[float64]float64)})
	}
	return &result
}

// add PCs and get probabilities
func (sim *Sim) runAnFocusRW(r *FocusRWResult) {
	// base := 64
	// deltat := 0.
	pAn := sim.param.AnFocusRW
	var current Tx
	// var tsa RandomWalker
	tsa := BRW{}

	// apply PC for each tx
	counter := 0
	fmt.Println((sim.param.maxCut-sim.param.minCut)/int(sim.param.Lambda), "h data")
	for i1 := sim.param.minCut; i1 < sim.param.maxCut; i1++ {
		fmt.Println(sim.param.maxCut - sim.param.minCut)
		if (i1-sim.param.minCut)%(int(math.Max(100, float64(sim.param.maxCut-sim.param.minCut))/100)) == 0 {
			counter++
			fmt.Println(counter, "/", sim.param.maxCut-sim.param.minCut)
		}
		// fmt.Println("Tx", i1)
		// assess for a range of PC lengths
		for i2 := 1; i2 < int(float64(pAn.maxiMT)*pAn.murel); i2++ {
			r.countertotal[0].v[float64(i2)/pAn.murel/sim.param.Lambda]++
			// add weight for each added PC tx
			sim.addCW(sim.cw[i1]) // add +1 to all ancestors
			sim.tangle[i1].cw++   // add also +1 to root

			for i3 := 0; i3 < pAn.nRWs; i3++ {
				for current, _ = tsa.RandomWalk(sim.tangle[0], sim); (len(sim.approvers[current.id]) > 0) && (current.id < i1); current, _ = tsa.RandomWalk(current, sim) {
				}
				if current.id == i1 {
					r.countersuccess[0].v[float64(i2)/pAn.murel/sim.param.Lambda] += 1. / float64(pAn.nRWs)
				}
			}
		}
		// remove all added weights
		for i2 := 1; i2 < int(float64(pAn.maxiMT)*pAn.murel); i2++ {
			sim.removeCW(sim.cw[i1])
			sim.tangle[i1].cw--
		}
	}
	fmt.Println("\nFocus RW Data gathered.")
}

func (a *FocusRWResult) Join(b FocusRWResult) (r FocusRWResult) {
	if a.countersuccess == nil {
		return b
	}
	for i := range b.countersuccess {
		r.countersuccess = append(r.countersuccess, joinMapMetricFloat64Float64(a.countersuccess[i], b.countersuccess[i]))
	}
	for i := range b.countertotal {
		r.countertotal = append(r.countertotal, joinMapMetricFloat64Float64(a.countertotal[i], b.countertotal[i]))
	}
	for i := range b.prob {
		r.prob = append(r.prob, joinMapMetricFloat64Float64(a.prob[i], b.prob[i]))
	}
	return r
}

// evaluate probabilities
func (r *FocusRWResult) finalprocess(p Parameters) error {
	for i1 := range r.countersuccess[0].v {
		r.prob[0].v[i1] = float64(r.countersuccess[0].v[i1]) / float64(r.countertotal[0].v[i1])
	}
	return nil
}

// - - - - - - - - - - - -
// - - - save process - - -
// - - - - - - - - - - - -

func (r FocusRWResult) Save(p Parameters) (err error) {
	if err = r.SaveCountersuccess(p); err != nil {
		return err
	}
	if err = r.SaveCountertotal(p); err != nil {
		return err
	}
	if err = r.SaveProb(p); err != nil {
		return err
	}
	return err
}

func (r FocusRWResult) SaveCountersuccess(p Parameters) error {
	for _, val := range r.countersuccess {
		val.SaveFocusRW(p, "countersuccess", true)
	}
	return nil
}
func (r FocusRWResult) SaveCountertotal(p Parameters) error {
	for _, val := range r.countertotal {
		val.SaveFocusRW(p, "Countertotal", true)
	}
	return nil
}
func (r FocusRWResult) SaveProb(p Parameters) error {
	for _, val := range r.prob {
		val.SaveFocusRW(p, "prob", true)
	}
	return nil
}

// Save saves a MetricFloat64Float64 as a file
func (s MetricFloat64Float64) SaveFocusRW(p Parameters, target string, normalized bool) error {
	var keys []float64
	// var datapoints int
	for k := range s.v {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.2f", p.Alpha)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/FocusRW_" + target + "_" + p.TSA + "_" + rateType + "_" + s.desc +
		"_lambda_" + lambdaStr +
		"_alpha_" + alphaStr + "_.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()
	// for i, k := range x {
	for _, k := range keys {
		_, err = f.WriteString(fmt.Sprintf("%f\t%f\n", k, s.v[k])) // writing...
		// _, err = f.WriteString(fmt.Sprintf("%f\t%f\n", k, weigths[i]/float64(datapoints)*norm)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
	return nil
}
