// Analyse the approver distribution with the time difference from a tx y in the past cone of that tx y
// time is measured backwards in relative indexes

package main

import (
	"fmt"
)

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//Data structure for results
type FutureConeResult struct { //these slices hold the statistics for each approver number mapping over all t
	counter        []MetricFloat64Float64
	p              []MetricFloat64Float64
	kappa          []MetricFloat64Float64
	counterInttime []MetricFloat64Float64
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//??? use string to create empty value maps to
func newFutureConeResult(coneMetrics []string) FutureConeResult {
	// variables initialization for FutureCone
	var result FutureConeResult
	for _, metric := range coneMetrics {
		result.counter = append(result.counter, MetricFloat64Float64{metric, make(map[float64]float64)})
		result.p = append(result.p, MetricFloat64Float64{metric, make(map[float64]float64)})
		result.counterInttime = append(result.counterInttime, MetricFloat64Float64{metric, make(map[float64]float64)})
		result.kappa = append(result.kappa, MetricFloat64Float64{metric, make(map[float64]float64)})
	}
	return result
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// for each cone, for each member of that cone, count +1 at the particular time
func (sim *Sim) runAnFutureCone(result *FutureConeResult) {
	maxchildID := 0
	_ = maxchildID

	if len(sim.cw) < sim.param.TangleSize-sim.param.minCut {
		fmt.Println(".\n\nNeed to check that the CWMatrix is not limited too much for this analysis.")
		pauseit()
	}

	if sim.param.maxCutrange < int(sim.param.AnFutureCone.MaxT*sim.param.Lambda) {
		fmt.Println("maxCutrange < MaxT !")
		pauseit()
	}

	// count occurances
	for i1 := sim.param.minCut; i1 < sim.param.maxCut; i1++ { //only consider roots that are within this cut ranges
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// evaluate measurements
func (r *FutureConeResult) finalprocess(p Parameters) error {
	for i2 := 1; i2 < len(r.counter); i2++ { // loop over all main options
		for _, i1 := range r.counter[i2].getkeys() {
			r.p[i2].v[i1] = float64(r.counter[i2].v[i1]) / float64(r.counter[0].v[i1])
		}
	}
	for _, i1 := range r.kappa[0].getkeys() {
		r.kappa[0].v[i1] = r.kappa[0].v[i1] / r.counterInttime[0].v[i1]
	}
	return nil
}

//Join joins FutureConeResult
func (r *FutureConeResult) Join(b FutureConeResult) (res FutureConeResult) {
	if r.counter == nil {
		return b
	}
	for i := range b.counter {
		res.counter = append(res.counter, joinMapMetricFloat64Float64(r.counter[i], b.counter[i]))
	}
	for i := range b.p {
		res.p = append(res.p, joinMapMetricFloat64Float64(r.p[i], b.p[i]))
	}
	for i := range b.counterInttime {
		res.counterInttime = append(res.counterInttime, joinMapMetricFloat64Float64(r.counterInttime[i], b.counterInttime[i]))
	}
	for i := range b.kappa {
		res.kappa = append(res.kappa, joinMapMetricFloat64Float64(r.kappa[i], b.kappa[i]))
	}
	return res
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// organise saving
func (r FutureConeResult) Save(p Parameters) (err error) {
	if err = SaveArrayMetricFloat64Float64(p, "AnFutureCone_counterInttime", r.counterInttime); err != nil {
		return err
	}
	if err = SaveArrayMetricFloat64Float64(p, "AnFutureCone_kappa", r.kappa); err != nil {
		return err
	}
	if err = SaveArrayMetricFloat64Float64(p, "AnFutureCone_counter", r.counter); err != nil {
		return err
	}
	if err = SaveArrayMetricFloat64Float64(p, "AnFutureCone_p", r.p); err != nil {
		return err
	}
	return err
}
