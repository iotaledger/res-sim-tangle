package main

import (
	"fmt"
	"math"
	"os"
	"sort"

	"gonum.org/v1/gonum/stat"
)

type ctResult struct {
	confirmationTime [][]int   // confirmation time of each tx
	mean             []float64 // avg of confirmation time in descending order over different Tangles
	variance         []float64
	totalMean        float64 // avg confirmation time over all tx and tangles
	totalVariance    float64 // variance of confirmation time over all tx and tangles
}

func newCTResult(p Parameters) ctResult {
	// variables initialization for entropy
	var result ctResult
	result.confirmationTime = make([][]int, p.nRun)
	for i := range result.confirmationTime {
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
	for j := range r.mean {
		confirmed = true
		var col []float64
		for i := range r.confirmationTime {
			if r.confirmationTime[i][j] == 1-1 {
				r.mean[j] = -1
				r.variance[j] = -1
				confirmed = false
				break
			}
		}
		if confirmed {
			r.mean[j], r.variance[j] = stat.MeanVariance(col, nil)
		}
	}
	var set []float64
	for j := range r.mean {
		for i := range r.confirmationTime {
			if r.confirmationTime[i][j] > -1 {
				set = append(set, float64(r.confirmationTime[i][j]))
			}
		}
		//fmt.Println("Len col:", len(col))
		r.totalMean, r.totalVariance = stat.MeanVariance(set, nil)
	}

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
	result.mean = a.mean
	result.variance = a.variance
	return result
}

func (a ctResult) ctToString(p Parameters, sample int) string {
	result := "# CT of each tx (descendent order)\n"
	result += "#Index\t\tsample\t\tavg\t\tvar\t\tstd\n"
	for j := range a.confirmationTime[0] {
		result += fmt.Sprintf("%d\t\t%d\t\t%.2f\t\t%.2f\t\t%.4f\n", j, a.confirmationTime[sample][j], a.mean[j], a.variance[j], math.Sqrt(a.variance[j]))
	}
	return result
}

func (a ctResult) Save(p Parameters, sample int) error {
	err := a.SaveCT(p)
	if err != nil {
		fmt.Println("error Saving CW", err)
		return err
	}
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
