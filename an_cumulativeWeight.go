package main

import (
	"fmt"
	"math"
	"os"
	"sort"

	"gonum.org/v1/gonum/stat"
)

type cwResult struct {
	cw       [][]int   // cw of each tx
	mean     []float64 // avg of cw in descending order over different Tangles
	variance []float64
}

func newCWResult(p Parameters) cwResult {
	// variables initialization for entropy
	var result cwResult
	result.cw = make([][]int, p.nRun)
	for i := range result.cw {
		result.cw[i] = make([]int, p.TangleSize)
	}
	result.mean = make([]float64, p.TangleSize)
	result.variance = make([]float64, p.TangleSize)
	return result
}

func (sim *Sim) fillCW(run int, r *cwResult) {
	for _, tx := range sim.tangle {
		r.cw[run][tx.id] = tx.cw
	}
	//order row
	sort.Sort(sort.Reverse(sort.IntSlice(r.cw[run])))
}

func (r *cwResult) Statistics(p Parameters) {
	for j := range r.mean {
		var col []float64
		for i := range r.cw {
			col = append(col, float64(r.cw[i][j]))
		}
		//fmt.Println("Len col:", len(col))
		r.mean[j], r.variance[j] = stat.MeanVariance(col, nil)
	}
	//fmt.Println("Len mean:", len(r.mean))
	//fmt.Println("Param:", p.minCut, p.TangleSize-p.minCut)
	//r.tAVG = stat.Mean(r.mean[p.minCut:], nil)
	//r.tSTD = math.Sqrt(stat.Mean(r.variance[p.minCut:], nil))

	// total pdf
	// r.tPDF = MetricIntInt{"pdf", make(map[int]int)}
	// for _, row := range r.pdf {
	// 	r.tPDF = joinMapMetricIntInt(r.tPDF, row)
	// 	//fmt.Println(r.tPDF)
	// }
}

func (a cwResult) Join(b cwResult) cwResult {
	if a.mean == nil {
		return b
	}
	var result cwResult
	result.cw = append(a.cw, b.cw...)
	result.mean = a.mean
	result.variance = a.variance
	return result
}

func (a cwResult) cwToString(p Parameters, sample int) string {
	result := "# cw of each tx (descendent order)\n"
	result += "#Index\t\tsample\t\tavg\t\tvar\t\tstd\n"
	for j := range a.cw[0] {
		result += fmt.Sprintf("%d\t\t%d\t\t%.2f\t\t%.2f\t\t%.4f\n", j, a.cw[sample][j], a.mean[j], a.variance[j], math.Sqrt(a.variance[j]))
	}
	return result
}

func (a cwResult) Save(p Parameters, sample int) error {
	err := a.SaveCW(p)
	if err != nil {
		fmt.Println("error Saving CW", err)
		return err
	}
	return err
}

func (a cwResult) SaveCW(p Parameters) (err error) {
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

	_, err = f.WriteString(a.cwToString(p, 0)) // writing...

	if err != nil {
		fmt.Printf("error writing string: %v", err)
		return err
	}

	return nil

}
