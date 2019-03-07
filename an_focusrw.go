// Analyse the approver distribution with the time difference from a tx y in the past cone of that tx y

package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

//FocusRWResult result of simulation
type FocusRWResult struct { //this slices hold the statistics for each approver number mapping over all deltat's
	countersuccess []MetricIntFloat64
	countertotal   []MetricIntFloat64
	prob           []MetricIntFloat64
}

//??? use string to create empty value maps to
func newFocusRWResult(Metrics []string) *FocusRWResult {
	// variables initialization for FocusRW
	var result FocusRWResult
	for _, metric := range Metrics {
		result.countersuccess = append(result.countersuccess, MetricIntFloat64{metric, make(map[int]float64)})
		result.countertotal = append(result.countertotal, MetricIntFloat64{metric, make(map[int]float64)})
		result.prob = append(result.prob, MetricIntFloat64{metric, make(map[int]float64)})
	}
	return &result
}

// add PCs and get probabilities
func (sim *Sim) runAnFocusRW(r *FocusRWResult) {
	// base := 64
	// deltat := 0.
	pAn := sim.param.AnFocusRW
	var current Tx
	var accepttx bool
	_ = accepttx
	PCweight := 0
	// var tsa RandomWalker
	tsa := BRW{}

	// apply PC for each tx, reveiled PC propto revealed txs since i1
	counter := 0
	// fmt.Println((sim.param.maxCut-sim.param.minCut)/int(sim.param.Lambda), "h data")
	// hidden tips do not have a propagation vector yet ->  i1 < len(sim.cw)
	// for i1 := sim.param.minCut; i1 < (sim.param.TangleSize-int(sim.param.Lambda)*200) && i1 < len(sim.cw); i1++ {
	for i1 := sim.param.minCut; i1 < len(sim.cw); i1++ { //consider only tx that are revealed
		// fmt.Println("-+-+--", i1, "-+-+--")
		// fmt.Println(sim.param.TangleSize - sim.param.AnFocusRW.acceptalwayslastN)
		// pauseit()
		if (i1-sim.param.minCut)%(int(math.Max(100, float64(sim.param.maxCut-sim.param.minCut))/100)) == 0 {
			counter++
			// fmt.Println(counter, "/", sim.param.maxCut-sim.param.minCut)
		}

		if sim.param.Lambda*sim.param.Alpha < 0.1 {
			accepttx = true
		} else {
			accepttx = !(sim.isLeftBehind(i1))
		}

		if accepttx {
			PCweight = int(float64(sim.param.TangleSize-i1) * sim.param.AnFocusRW.murel)
			sim.addCW(sim.cw[i1], PCweight)                                               // add cw from PC to all ancestors
			sim.tangle[i1].cw += PCweight                                                 // add also to root
			r.countertotal[0].v[int(float64(sim.param.TangleSize-i1)/sim.param.Lambda)]++ // add +1 to this particular PC time
			for i3 := 0; i3 < pAn.nRWs; i3++ {
				for current, _ = tsa.RandomWalk(sim.tangle[0], sim); (len(sim.approvers[current.id]) > 0) && (current.id < i1); current, _ = tsa.RandomWalk(current, sim) {
				}
				if current.id == i1 {
					r.countersuccess[0].v[int(float64(sim.param.TangleSize-i1)/sim.param.Lambda)] += 1. / float64(pAn.nRWs)
				}
			}
			sim.removeCW(sim.cw[i1], PCweight) // remove cw from PC to all ancestors
			sim.tangle[i1].cw -= PCweight      // remove also from root
		}
	}
	// fmt.Println("\nFocus RW Data gathered.")
}

//Join joins FocusRWResult
func (r *FocusRWResult) Join(b FocusRWResult) (res FocusRWResult) {
	if r.countersuccess == nil {
		return b
	}
	for i := range b.countersuccess {
		res.countersuccess = append(res.countersuccess, joinMapMetricIntFloat64(r.countersuccess[i], b.countersuccess[i]))
	}
	for i := range b.countertotal {
		res.countertotal = append(res.countertotal, joinMapMetricIntFloat64(r.countertotal[i], b.countertotal[i]))
	}
	for i := range b.prob {
		res.prob = append(res.prob, joinMapMetricIntFloat64(r.prob[i], b.prob[i]))
	}
	return res
}

// evaluate probabilities
func (r *FocusRWResult) finalprocess(p Parameters) error {
	for i1 := range r.countersuccess[0].v {
		r.prob[0].v[i1] = r.countersuccess[0].v[i1] / r.countertotal[0].v[i1]
	}
	return nil
}

// - - - - - - - - - - - -
// - - - save process - - -
// - - - - - - - - - - - -

// Save saves FocusRWResult
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

// SaveCountersuccess saves countersuccess
func (r FocusRWResult) SaveCountersuccess(p Parameters) error {
	for _, val := range r.countersuccess {
		val.SaveFocusRW(p, "countersuccess", true)
	}
	return nil
}

// SaveCountertotal saves countertotal
func (r FocusRWResult) SaveCountertotal(p Parameters) error {
	for _, val := range r.countertotal {
		val.SaveFocusRW(p, "Countertotal", true)
	}
	return nil
}

// SaveProb saves prob
func (r FocusRWResult) SaveProb(p Parameters) error {
	for _, val := range r.prob {
		val.SaveFocusRW(p, "prob", true)
	}
	return nil
}

// SaveFocusRW saves a MetricIntFloat64 as a file
func (s MetricIntFloat64) SaveFocusRW(p Parameters, target string, normalized bool) error {
	var keys []int
	// var datapoints int
	for k := range s.v {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.4f", p.Alpha)
	murelStr := fmt.Sprintf("%.1f", p.AnFocusRW.murel)
	f, err := os.Create("data/FocusRW_lambda" + lambdaStr + "_alpha" + alphaStr + "_murel" + murelStr + ".txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()
	// for i, k := range x {
	for _, k := range keys {
		_, err = f.WriteString(fmt.Sprintf("%f\t%f\n", float64(k), s.v[k])) // writing...
		// _, err = f.WriteString(fmt.Sprintf("%f\t%f\n", k, weigths[i]/float64(datapoints)*norm)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
	return nil
}
