// Analyse the approver distribution with the time difference from a tx y in the past cone of that tx y
// time is measured backwards in relative indexes

package main

import (
	"fmt"
	"math"
)

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//Data structure for results
type PastConeResult struct { //these slices hold the statistics for each approver number mapping over all t
	counter        []MetricFloat64Float64
	p              []MetricFloat64Float64
	kappa          []MetricFloat64Float64
	counterInttime []MetricFloat64Float64
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//??? use string to create empty value maps to
func newPastConeResult(coneMetrics []string) *PastConeResult {
	// variables initialization for PastCone
	var result PastConeResult
	for _, metric := range coneMetrics {
		result.counter = append(result.counter, MetricFloat64Float64{metric, make(map[float64]float64)})
		result.p = append(result.p, MetricFloat64Float64{metric, make(map[float64]float64)})
		result.counterInttime = append(result.counterInttime, MetricFloat64Float64{metric, make(map[float64]float64)})
		result.kappa = append(result.kappa, MetricFloat64Float64{metric, make(map[float64]float64)})
	}
	return &result
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// for each cone, for each member of that cone, count +1 at the particular time
func (sim *Sim) runAnPastCone(result *PastConeResult) {
	base := 64
	// resolution := 4. // how many steps per h
	deltat := 0.
	maxchildID := 0
	_ = maxchildID

	// count occurances
	for i1 := sim.param.minCut; i1 < sim.param.maxCut; i1++ { //only consider roots that are within this cut ranges
		maxchildID = 0
		// fmt.Println("i1", i1)
		for i2, block := range sim.cw[i1] { // i2 iterates through the blocks, block is the unit64 of the block itself
			// fmt.Println(" ... i2", i2)
			if (i2+1)*base > sim.param.minCut { // only consider blocks above min cut
				if (i1 - (i2+1)*base) < int(sim.param.AnPastCone.MaxT)*int(sim.param.Lambda) { // only consider block if its within maxT
					for i3 := 0; i3 < base; i3++ {
						deltat = math.Round((sim.tangle[i1].time-sim.tangle[i2*base+i3].time)*sim.param.AnPastCone.Resolution) / sim.param.AnPastCone.Resolution // need to check that this is picking the correct tx
						if i3+i2*base < i1 {                                                                                                                     // only count if we are in the past of i1
							maxchildID = max2(i3+i2*base, maxchildID)
							result.counterInttime[0].v[float64(int(deltat))]++
						} else { // if we counted the above correctly we should not get here.
							if block&(1<<uint(i3)) != 0 { // c
								fmt.Println("This should not happen.")
								pauseit()
							}
						}
						if block&(1<<uint(i3)) != 0 { // if this is an ancestor of i1 then
							result.counter[0].v[deltat]++
							result.kappa[0].v[float64(int(deltat))]++
							if len(sim.approvers[i2*base+i3]) < sim.param.AnPastCone.MaxApp { //if smaller than maximum considered add +1 to maxApp
								result.counter[len(sim.approvers[i2*base+i3])].v[deltat]++
							} else { //if larger than maximum considered add +1 to maxApp
								result.counter[sim.param.AnPastCone.MaxApp].v[deltat]++
							}
						}
					}
				}
			}
		}
		if maxchildID < i1-1 && maxchildID > 0 { // potentially the cw array was smaller than i1, therefore add the values for those as well.
			for i2 := maxchildID + 1; i2 < i1; i2++ {
				deltat = math.Round(sim.tangle[i1].time - sim.tangle[i2].time)
				result.counterInttime[0].v[float64(int(deltat))]++
			}
		}
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// evaluate measurements
func (r *PastConeResult) finalprocess(p Parameters) error {
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

// - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// join Simulations
func (PER *PastConeResult) Join(b PastConeResult) (r PastConeResult) {
	if PER.counter == nil {
		return b
	}
	for i := range b.counter {
		r.counter = append(r.counter, joinMapMetricFloat64Float64(PER.counter[i], b.counter[i]))
	}
	for i := range b.p {
		r.p = append(r.p, joinMapMetricFloat64Float64(PER.p[i], b.p[i]))
	}
	for i := range b.counterInttime {
		r.counterInttime = append(r.counterInttime, joinMapMetricFloat64Float64(PER.counterInttime[i], b.counterInttime[i]))
	}
	for i := range b.kappa {
		r.kappa = append(r.kappa, joinMapMetricFloat64Float64(PER.kappa[i], b.kappa[i]))
	}
	return r
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// organise saving
func (r PastConeResult) Save(p Parameters) (err error) {
	if err = SaveArrayMetricFloat64Float64(p, "AnPastCone_counterInttime", r.counterInttime); err != nil {
		return err
	}
	if err = SaveArrayMetricFloat64Float64(p, "AnPastCone_kappa", r.kappa); err != nil {
		return err
	}
	if err = SaveArrayMetricFloat64Float64(p, "AnPastCone_counter", r.counter); err != nil {
		return err
	}
	if err = SaveArrayMetricFloat64Float64(p, "AnPastCone_p", r.p); err != nil {
		return err
	}
	return err
}
