package main

import (
	"fmt"
	"os"
	"sort"

	"gonum.org/v1/gonum/stat"
)

//Entropy result of simulation
type entropyResult struct {
	tips    MetricIntInt // number of particles reaching specific tips
	ep      [][]float64  // exit probabilities (number of rows = number of Tangles )
	mean    []float64    // exit probability means over Tangles
	median  []float64    // exit probability medians over Tangles
	std     []float64    // exit probability std dev over Tangles
	ep2     [][]float64  // exit probabilities (number of rows = number of Tangles), ep2 maps onto 2lambda intervals
	mean2   []float64    // exit probability means over Tangles
	median2 []float64    // exit probability medians over Tangles
	std2    []float64    // exit probability std dev over Tangles
}

func newEntropyResult() entropyResult {
	// variables initialization for entropy
	var result entropyResult
	return result
}

func (sim *Sim) runEntropyStat(result *entropyResult) {
	Nparticle := 100000
	result.tips = MetricIntInt{"Tips", make(map[int]int)}
	if sim.param.TSA != "RW" {
		sim.entropyURTS(result.tips.v, Nparticle)
		result.ep = append(result.ep, sortEntropy(result.tips.v))
		// fmt.Println(result.ep)
		result.ep2 = append(result.ep2, mapEntropy01(sortEntropy(result.tips.v), sim.param))
		// fmt.Println(result.ep2)
		// fmt.Println("in here")
		// pauseit()
	} else {
		sim.entropyParticleRW(result.tips.v, Nparticle)
		result.ep = append(result.ep, sortEntropy(result.tips.v))
		result.ep2 = append(result.ep2, mapEntropy01(sortEntropy(result.tips.v), sim.param))
	}
}

//sort the exit probabilities from max to min
func sortEntropy(v map[int]int) (r []float64) {
	r = make([]float64, len(v))
	var values []int
	var datapoints float64
	for _, val := range v {
		values = append(values, val)
		datapoints += float64(val)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	for i, val := range values {
		r[i] = float64(val) / datapoints
	}
	return r
}

// map the exit probabilites from 0-1
func mapEntropy01(r []float64, p Parameters) (m []float64) {
	m = make([]float64, int(p.Lambda)*2)
	fillboxes := float64(len(m)) / float64(len(r))
	leftboxes := 0.
	remainer := 0.
	vID := 0
	for i1 := 0; i1 < len(r); i1++ {
		if remainer+fillboxes < 1 { // next r[i1] fits still into the same box
			m[vID] += r[i1]
			remainer += fillboxes
		} else {
			m[vID] += r[i1] * (1 - remainer) / fillboxes
			vID++
			leftboxes = fillboxes - (1 - remainer)
			for ; leftboxes > 1; leftboxes-- {
				m[vID] += r[i1] / fillboxes
				vID++
			}
			if vID == 2*int(p.Lambda) && leftboxes < 1e-6 { // this is just some overflow due to some precission error
				vID--
			}
			if leftboxes > 0 { //could be zero now
				m[vID] += leftboxes * r[i1] / fillboxes
				remainer = leftboxes
			}
		}
	}
	if vID < 2*int(p.Lambda)-1 || vID > 2*int(p.Lambda) { //brief check if we are one of those indices
		fmt.Println("Error: vID=", vID, " != ", 2*int(p.Lambda))
		pauseit()
	}
	return m
}

// calculate exit probabilities and save in v
func (sim *Sim) entropyURTS(v map[int]int, nParticles int) {
	for i := 0; i < nParticles; i++ {
		tip := sim.generator.Intn(len(sim.tips))
		v[sim.tips[tip]]++
	}
	// for i := 0; i < len(sim.tips); i++ {
	// 	v[i]++
	// }
}

// calculate exit probabilities and save in v
func (sim *Sim) entropyParticleRW(v map[int]int, nParticles int) {
	for i := 0; i < nParticles; i++ {
		var tsa RandomWalker
		if sim.param.Alpha > 0 {
			tsa = BRW{}
		} else {
			tsa = URW{}
		}
		var current int
		for current, _ = tsa.RandomWalkStep(0, sim); len(sim.approvers[current]) > 0; current, _ = tsa.RandomWalkStep(current, sim) {
		}
		v[current]++
	}
}

func (a *entropyResult) Join(b entropyResult) (r entropyResult) {
	if a.ep == nil {
		return b
	}
	r.ep = append(r.ep, a.ep...)
	r.ep = append(r.ep, b.ep...)
	r.ep2 = append(r.ep2, a.ep2...)
	r.ep2 = append(r.ep2, b.ep2...)
	return r
}

func (e *entropyResult) Save(p Parameters, str string) (err error) {
	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.4f", p.Alpha)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/entropy_stat_" + p.TSA + "_" + rateType +
		"_lambda_" + lambdaStr +
		"_alpha_" + alphaStr + "_" + str + ".txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(e.Stat(p)) // writing...

	if err != nil {
		fmt.Printf("error writing string: %v", err)
		return err
	}

	return nil

}

func (e *entropyResult) Stat(p Parameters) (result string) {
	result += "#Entropy Stats [exit probabilities]\n"
	result += "#Index\t\tSingle\t\tMean\t\tStdDev\t\tMedian\n"

	// find len of each row
	var lenRows []int
	for _, row := range e.ep {
		lenRows = append(lenRows, len(row))
	}

	// make same len for all the rows and fill with 0s if smaller
	_, nColumns := max(lenRows)
	for i, row := range e.ep {
		for j := 0; j < nColumns-len(row); j++ {
			e.ep[i] = append(e.ep[i], 0.0)
		}
	}

	// calculate means and std deviations by extrapolating columns first
	for i := 0; i < nColumns; i++ {
		var c []float64
		for _, row := range e.ep {
			c = append(c, row[i])
		}
		sort.Float64s(c)
		mean, std := stat.MeanStdDev(c, nil)
		median := stat.Quantile(0.5, stat.Empirical, c, nil)
		e.mean = append(e.mean, mean)
		e.median = append(e.median, median)
		e.std = append(e.std, std)
	}

	for i := 0; i < nColumns; i++ {
		result += fmt.Sprintf("%d\t\t%.4f\t\t%.4f\t\t%.4f\t\t%.4f\n", i, e.ep[0][i], e.mean[i], e.std[i], e.median[i])

	}
	return result
}
