package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

type ctResult struct {
	confirmationTime [][]int   // confirmation time of each tx
	mean             []float64 // avg of confirmation time in descending order over different Tangles
	variance         []float64
	totalMean        float64 // avg confirmation time over all tx and tangles
	totalVariance    float64 // variance of confirmation time over all tx and tangles
}

func newCTResult(p Parameters) ctResult {
	var result ctResult
	result.confirmationTime = make([][]int, p.nRun)
	for i := 0; i < p.nRun; i++ {
		result.confirmationTime[i] = make([]int, p.TangleSize)
	}
	result.mean = make([]float64, p.TangleSize)
	result.variance = make([]float64, p.TangleSize)
	return result
}

func (sim *Sim) fillCT(run int, r *ctResult) {
	for _, tx := range sim.tangle {
		r.confirmationTime[run][tx.id] = tx.confirmationTime
	}
	//order row
	sort.Sort(sort.Reverse(sort.IntSlice(r.confirmationTime[run])))
}

func (r *ctResult) Statistics(p Parameters) {
	var confirmed bool
	var counter int
	var mean float64
	var variance float64
	for j := 0; j < p.TangleSize; j++ {
		confirmed = true
		counter = 0
		mean = 0
		variance = 0
		for i := 0; i < p.nRun; i++ {
			if r.confirmationTime[i][j] == -1 {
				r.mean[j] = -1
				r.variance[j] = -1
				confirmed = false
				break
			}
			mean += float64(r.confirmationTime[i][j])
			counter++
		}
		if confirmed {
			mean = mean / float64(counter)
			counter = 0
			for i := 0; i < p.nRun; i++ {
				variance += math.Pow((float64(r.confirmationTime[i][j]) - mean), 2)
				counter++
			}
			r.mean[j] = mean
			r.variance[j] = variance / float64(counter)
		}
	}
	counter = 0
	mean = 0
	for j := 0; j < p.TangleSize; j++ {
		for i := 0; i < p.nRun; i++ {
			if r.confirmationTime[i][j] > -1 {
				counter += 1
				mean += float64(r.confirmationTime[i][j])
			}
		}
	}
	mean = mean / float64(counter)
	counter = 0
	variance = 0
	for j := 0; j < p.TangleSize; j++ {
		for i := 0; i < p.nRun; i++ {
			if r.confirmationTime[i][j] > -1 {
				counter += 1
				variance += math.Pow((mean - float64(r.confirmationTime[i][j])), 2)
			}
		}
	}
	r.totalMean = mean
	r.totalVariance = variance / float64(counter)
	fmt.Println("Stopping Statistics")
	// total pdf
	// r.tPDF = MetricIntInt{"pdf", make(map[int]int)}
	// for _, row := range r.pdf {
	// 	r.tPDF = joinMapMetricIntInt(r.tPDF, row)
	// 	//fmt.Println(r.tPDF)
	// }
}

func (a ctResult) Join(b ctResult) ctResult {
	if a.mean == nil {
		return b
	}
	var result ctResult
	result.confirmationTime = append(a.confirmationTime, b.confirmationTime...)
	//result.mean = a.mean   // placeholder to be checked
	//result.variance = a.variance
	return result
}

func (a ctResult) ctToString(p Parameters, sample int) string {
	fmt.Println("Starting ctToString")
	result := "# CT of each tx (descendent order)\n"
	result += "#Index\t\tsample\t\tavg\t\tvar\t\tstd\n"
	for j := range a.confirmationTime[0] {
		result += fmt.Sprintf("%d\t\t%d\t\t%.2f\t\t%.2f\t\t%.4f\n", j, a.confirmationTime[sample][j], a.mean[j], a.variance[j], math.Sqrt(a.variance[j]))
	}
	fmt.Println("Stopping ctToString")
	return result
}

func (a ctResult) Save(p Parameters, sample int) error {
	fmt.Println("Starting Save CT")
	err := a.SaveCT(p)
	if err != nil {
		fmt.Println("error Saving CT", err)
		return err
	}
	fmt.Println("Stopping Save CT")
	return err
}

func (a ctResult) SaveCT(p Parameters) (err error) {

	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/cw_" + p.TSA + "_" + rateType +
		"_lambda_" + lambdaStr + "_.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(a.ctToString(p, 0)) // writing...

	if err != nil {
		fmt.Printf("error writing string: %v", err)
		return err
	}
	return nil

}
