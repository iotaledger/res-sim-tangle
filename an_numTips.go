package main

import (
	"fmt"
	"math"
	"os"

	"gonum.org/v1/gonum/stat"
)

type tipsResult struct {
	nTips               [][]int        // # of tips seen by each tx ([sim-tangle number][tx])
	mean                []float64      // avg of tips seen by each tx over different Tangles
	variance            []float64      // var of tips seen by each tx over different Tangles
	pdf                 []MetricIntInt // probability density function for each run
	tAVG                float64        // total # of tips avg
	tSTD                float64        // total # of tips std
	tPDF                MetricIntInt   // total probability density function
	nOrphanTips         []int          // # of orphanTips seen at the end of simulation for each Tangle
	meanOrphanTips      float64        // average of orphaned tips over all Tangles
	meanOrphanTipsRatio float64        // ratio of orphaned tips over all Tangles
	STDOrphanTips       float64        // variance of orphaned tips over all Tangles
	STDOrphanTipsRatio  float64        // variance of ratio of orphaned tips over all Tangles
}

func newTipsResult(p Parameters) tipsResult {
	var result tipsResult
	result.nTips = make([][]int, p.nRun)
	result.pdf = make([]MetricIntInt, p.nRun)
	for i := range result.nTips {
		result.nTips[i] = make([]int, p.TangleSize)
		result.pdf[i] = MetricIntInt{"pdf", make(map[int]int)}
	}
	result.mean = make([]float64, p.TangleSize)
	result.variance = make([]float64, p.TangleSize)
	result.nOrphanTips = make([]int, p.nRun)
	//result.tPDF = MetricIntInt{"total_tips_pdf", make(map[int]int)}
	return result
}

func (sim *Sim) countTips(tx int, run int, r *tipsResult) {
	r.nTips[run][tx] = len(sim.tips)
	if tx > sim.param.minCut {
		r.pdf[run].v[len(sim.tips)]++
	}
}

func (sim *Sim) countOrphanTips(run int, r *tipsResult) {
	r.nOrphanTips[run] = len(sim.orphanTips)
}

func (r *tipsResult) Statistics(p Parameters) {
	for j := range r.mean {
		var col []float64
		for i := range r.nTips {
			col = append(col, float64(r.nTips[i][j]))
		}
		//fmt.Println("Len col:", len(col))
		r.mean[j], r.variance[j] = stat.MeanVariance(col, nil)
	}
	//fmt.Println("Len mean:", len(r.mean))
	//fmt.Println("Param:", p.minCut, p.TangleSize-p.minCut)
	r.tAVG = stat.Mean(r.mean[p.minCut:], nil)
	r.tSTD = math.Sqrt(stat.Mean(r.variance[p.minCut:], nil))
	r.meanOrphanTips = MeanInt(r.nOrphanTips)
	r.STDOrphanTips = math.Sqrt(VarInt(r.nOrphanTips))
	r.meanOrphanTipsRatio = r.meanOrphanTips / float64(p.TangleSize)
	r.STDOrphanTipsRatio = r.STDOrphanTips / float64(p.TangleSize)
	// total pdf
	r.tPDF = MetricIntInt{"pdf", make(map[int]int)}
	for _, row := range r.pdf {
		r.tPDF = joinMapMetricIntInt(r.tPDF, row)
		//fmt.Println(r.tPDF)
	}
}

func (a tipsResult) Join(b tipsResult) tipsResult {
	if a.mean == nil {
		return b
	}
	var result tipsResult
	result.nTips = append(a.nTips, b.nTips...)
	result.pdf = append(a.pdf, b.pdf...)
	// result.mean = a.mean     // this is just a gap filler and will be replaced later
	// result.variance = a.variance // this is just a gap filler and will be replaced later
	result.nOrphanTips = append(a.nOrphanTips, b.nOrphanTips...)
	return result
}

// func (a tipsResult) ToString(p Parameters) string {
// 	//result := fmt.Sprintln("E(L):", a.tAVG, a.tSTD)
// 	result := "#Tips Statistics\n"
// 	result += "#Stat Type\tLambda\t\tAlpha\t\tMean\t\tStdDev\t\tVariance\tMedian\t\tMode\t\tSkew\t\tMinVal\t\tMaxVal\t\tN\n"
// 	result += a.tPDF.ToString(p, false)
// 	return result
// }

func (a tipsResult) nTipsToString(p Parameters, sample int) string {
	result := "# Number of tips seen by each tx\n"
	result += "#Tx\t\tsample\t\tavg\t\tvar\t\tstd\n"
	for j := range a.nTips[0][1:] {
		result += fmt.Sprintf("%d\t\t%d\t\t%.2f\t\t%.2f\t\t%.4f\n", j+1, a.nTips[sample][j+1], a.mean[j+1], a.variance[j+1], math.Sqrt(a.variance[j+1]))
	}
	return result
}

func (a tipsResult) nOrphanTipsToString(p Parameters, sample int) string {
	result := "# Number of orphantips seen by each tangle\n"
	result += "#Tangle\t\tnOrphan\t\tmean\t\tratio\n"
	for j := range a.nOrphanTips[:] {
		result += fmt.Sprintf("%d\t\t%d\t\t%.2f\t\t%.2f\n", j+1, a.nOrphanTips[j])
	}
	return result
}

func (a tipsResult) Save(p Parameters, sample int) error {
	err := a.SaveTips(p)
	if err != nil {
		fmt.Println("error Saving Tips", err)
		return err
	}
	err = a.SaveOrphanTips(p)
	if err != nil {
		fmt.Println("error Saving Tips", err)
		return err
	}
	// err = a.tPDF.Save(p, "tips_pdf", "avg", false)
	// if err != nil {
	// 	fmt.Println("error Saving Tips PDF avg", err)
	// 	return err
	// }
	// err = a.pdf[sample].Save(p, "tips_pdf", "sample", false)
	// if err != nil {
	// 	fmt.Println("error Saving Tips PDF sample", err)
	// 	return err
	// }
	return err
}

func (a tipsResult) SaveTips(p Parameters) (err error) {
	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/tips_" + p.TSA + "_" + rateType +
		"_lambda_" + lambdaStr + "_.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(a.nTipsToString(p, 0)) // writing...

	if err != nil {
		fmt.Printf("error writing string: %v", err)
		return err
	}

	return nil

}

func (a tipsResult) SaveOrphanTips(p Parameters) (err error) {
	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/orphantips_" + p.TSA + "_" + rateType +
		"_lambda_" + lambdaStr + "_.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(a.nOrphanTipsToString(p, 0)) // writing...

	if err != nil {
		fmt.Printf("error writing string: %v", err)
		return err
	}

	return nil

}

func MeanInt(x []int) float64 {
	sum := 0
	for _, xs := range x {
		sum += xs
	}
	return float64(sum) / float64(len(x))
}

func VarInt(x []int) float64 {
	sum := 0.
	sum2 := 0.
	for _, xs := range x {
		sum += float64(xs)
	}
	for _, xs := range x {
		a := float64(xs) - sum/float64(len(x))
		sum2 += a * a
	}
	return float64(sum2) / float64(len(x))
}
