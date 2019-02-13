// Analyse the approver distribution with the time difference from a tx y in the past cone of that tx y

package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

//PastCone result of simulation
type PastConeResult struct { //this slices hold the statistics for each approver number mapping over all deltat's
	counter []MetricFloat64Float64
	p       []MetricFloat64Float64
}

//??? use string to create empty value maps to
func newPastConeResult(coneMetrics []string) *PastConeResult {
	// variables initialization for PastCone
	var result PastConeResult
	for _, metric := range coneMetrics {
		result.counter = append(result.counter, MetricFloat64Float64{metric, make(map[float64]float64)})
		result.p = append(result.p, MetricFloat64Float64{metric, make(map[float64]float64)})
	}
	return &result
}

// for each cone, for each member of that cone, count +1 at the particular time
func (sim *Sim) runAnPastCone(result *PastConeResult) {
	base := 64
	// resolution := 4. // how many steps per h
	deltat := 0.

	// initate counters
	// for i2 := 0; i2 <= maxApp; i2++ {
	// 	for i1 := 1. / resolution; i1 < float64(maxt); i1 += 1. / resolution {
	// 		result.counter[i2].v[i1] = 0
	// 		result.p[i2].v[i1] = 0
	// 	}
	// }

	// count occurances
	for i1 := sim.param.minCut; i1 < sim.param.maxCut; i1++ {
		// fmt.Println("i1", i1)
		for i2, block := range sim.cw[i1] { // i2 iterates through the blocks, block is the unit64 of the block itself
			// fmt.Println(" ... i2", i2)
			if i2*base > sim.param.minCut { // only consider tx within this range
				if (i1 - (i2+1)*base) < int(sim.param.AnPastCone.MaxT)*int(sim.param.Lambda) { // only consider tx within this range
					if i1 > 2*sim.param.minCut {
						// fmt.Println("+++++++++++++++++++")
						// fmt.Println(i2 * base)
						// fmt.Println((i1 - i2*base))
						// fmt.Println(int(sim.param.AnPastConeMaxT) * int(sim.param.Lambda))
						// fmt.Println(sim.cw[i1])
						// fmt.Println("   i2", i2)
						// fmt.Println(block)
						// pauseit()
					}
					for i3 := 0; i3 < base; i3++ {
						if block&(1<<uint(i3)) != 0 { // if this is an ancestor of i1 then
							deltat = math.Round((sim.tangle[i1].time-sim.tangle[i2*base+i3].time)*sim.param.AnPastCone.Resolution) / sim.param.AnPastCone.Resolution // need to check that this is picking the correct tx
							result.counter[0].v[deltat]++
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
	}

}

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
	return r
}

func (r PastConeResult) Save(p Parameters) (err error) {
	if err = r.SaveCounter(p); err != nil {
		return err
	}
	if err = r.SaveP(p); err != nil {
		return err
	}
	return err
}

func (r PastConeResult) SaveCounter(p Parameters) error {
	for _, counter := range r.counter {
		counter.SavePastCone(p, "counter", true)
	}
	return nil
}
func (r PastConeResult) SaveP(p Parameters) error {
	for _, prob := range r.p {
		prob.SavePastCone(p, "p", true)
	}
	return nil
}

// Save saves a MetricFloat64Float64 as a file
func (s MetricFloat64Float64) SavePastCone(p Parameters, target string, normalized bool) error {
	var keys []float64
	// var datapoints int
	for k := range s.v {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	// var weigths []float64
	// var x []float64
	// norm := 1.
	// for k := range keys {
	// 	if normalized {
	// 		norm = p.Lambda
	// 	}
	// 	x = append(x, float64(keys[k])/norm)
	// 	weigths = append(weigths, float64(s.v[keys[k]]))
	// 	// datapoints = datapoints + s.v[keys[k]]
	// }

	// save to file for plot

	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.2f", p.Alpha)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/PastCone_" + target + "_" + p.TSA + "_" + rateType + "_" + s.desc +
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

// evaluate probabilities
func (r *PastConeResult) finalprocess(p Parameters) error {
	for i2 := 1; i2 <= len(r.counter)-1; i2++ { // loop over all main options
		for i1 := range r.counter[i2].v {
			r.p[i2].v[i1] = float64(r.counter[i2].v[i1]) / float64(r.counter[0].v[i1])
		}
		// var keys []float64
		// // var datapoints int
		// for k := range s.v {
		// 	keys = append(keys, k)
		// }
		// for k := range
		// for i1 := 1; i1 < int(p.AnPastConeMaxT*p.AnPastConeResolution); i1 += 1 {
		// 	if float64(r.counter[0].v[float64(i1)/p.AnPastConeResolution]) > 0 {
		// r.p[i2].v[float64(i1)/p.AnPastConeResolution] = float64(r.counter[i2].v[float64(i1)/p.AnPastConeResolution]) / float64(r.counter[0].v[float64(i1)/p.AnPastConeResolution])
		// 	}
		// }
	}
	return nil
}
