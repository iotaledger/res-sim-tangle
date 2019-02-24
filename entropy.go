package main

import (
	"fmt"
	"os"
	"sort"

	"gonum.org/v1/gonum/stat"
)

//Velocity result of simulation
type entropyResult struct {
	tips MetricIntInt
	ep   [][]float64
	mean []float64
	std  []float64
}

func newEntropyResult() *entropyResult {
	// variables initialization for entropy
	var result entropyResult
	return &result
}

func (sim *Sim) runEntropyStat(result *entropyResult) {
	result.tips = MetricIntInt{"Tips", make(map[int]int)}
	if sim.param.TSA != "RW" {
		//
	} else {
		sim.entropyParticleRW(result.tips.v, 100000)
		result.ep = append(result.ep, sortEntropy(result.tips.v))
	}
}

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

func (sim *Sim) entropyParticleRW(v map[int]int, nParticles int) {
	for i := 0; i < nParticles; i++ {
		prev := sim.tangle[0]
		tsa := URW{}

		var current Tx
		for current, _ = tsa.RandomWalk(prev, sim); len(sim.approvers[current.id]) > 0; current, _ = tsa.RandomWalk(current, sim) {
		}

		v[current.id]++
	}
}

func (a *entropyResult) Join(b entropyResult) (r entropyResult) {
	if a.ep == nil {
		return b
	}
	r.ep = append(r.ep, a.ep...)
	r.ep = append(r.ep, b.ep...)
	return r
}

func (e *entropyResult) Save(p Parameters) (err error) {
	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.2f", p.Alpha)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/entropy_stat_" + p.TSA + "_" + rateType +
		"_lambda_" + lambdaStr +
		"_alpha_" + alphaStr + "_.txt")
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
	result += "#Index\t\tSingle\t\tMean\t\tStdDev\n"

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

	// calculate means by extrapolating columns first
	for i := 0; i < nColumns; i++ {
		var c []float64
		for _, row := range e.ep {
			c = append(c, row[i])
		}
		mean, std := stat.MeanStdDev(c, nil)
		e.mean = append(e.mean, mean)
		e.std = append(e.std, std)
	}

	for i := 0; i < nColumns; i++ {
		result += fmt.Sprintf("%d\t\t%.4f\t\t%.4f\t\t%.4f\n", i, e.ep[0][i], e.mean[i], e.std[i])

	}
	return result
}
