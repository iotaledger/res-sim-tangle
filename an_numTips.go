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
			// fmt.Print(r.nTips[i][j], " ,")
			col = append(col, float64(r.nTips[i][j]))
		}
		r.mean[j], r.variance[j] = MeanVariance(col)
		// fmt.Print("Len col:", len(col), "; ")
		// fmt.Println(r.mean[j], r.variance[j])
	}
	fmt.Println("Len mean:", len(r.mean), "; Len mean cut: ", len(r.mean[p.minCut:]))
	//fmt.Println("Param:", p.minCut, p.TangleSize-p.minCut)
	// r.tAVG = stat.Mean(r.mean[p.minCut:], nil)
	// r.tSTD = math.Sqrt(stat.Mean(r.variance[p.minCut:], nil))
	var variance float64
	r.tAVG, variance = MeanVariance(r.mean[p.minCut:])
	r.tSTD = math.Sqrt(variance)
		for i:=0;i<len(r.nOrphanTips);i++{
		temp[i]=float64(r.nOrphanTips[i])
	}
	r.meanOrphanTips = stat.Mean(temp,nil)
	r.STDOrphanTips = math.Sqrt(stat.Variance(temp,nil))
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

func (a tipsResult) nTipsToString(p Parameters, samples int) string {
	result := "# Number of tips seen by each tx\n"
	result += "#Tx;avg;var;std;samples\n"
	for j := range a.nTips[0][1:] {
		result += fmt.Sprintf("%d;%.2f;%.2f;%.4f", j+1, a.mean[j+1], a.variance[j+1], math.Sqrt(a.variance[j+1]))
		for i := 0; i < samples; i++ {
			result += fmt.Sprintf(";%d", a.nTips[i][j+1])
		}
		result += fmt.Sprintf("\n")
	}
	return result
}

func (a tipsResult) nOrphanTipsToString(p Parameters) string {
	result := "# Number of orphantips seen by each tangle\n"
	result += "#Tangle;nOrphan;OrphanRatio\n"
	for j := range a.nOrphanTips[:] {
		// for orphanratio : adjust Tangle size for D.
		orphanratio := 0.
		if float64(p.TangleSize) > p.D*p.Lambda {
			orphanratio = float64(a.nOrphanTips[j]) / (float64(p.TangleSize) - p.D*p.Lambda)
		}
		result += fmt.Sprintf("%d;%d;%f\n", j+1, a.nOrphanTips[j], orphanratio)
	}
	return result
}

func (a orphanResult) nOrphanTxsToString(p Parameters) string {
	result := "# proportion orphantxs seen by each tangle\n"
	result += "#Tangle;OrphanRatio\n"
	for j := range a.op2[:] {
		result += fmt.Sprintf("%d;%f\n", j+1, a.op2[j])
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
		fmt.Println("error Saving Orphan Tips", err)
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
	str := fmt.Sprintf("%d", int(p.SimStep))
	f, err := os.Create("data/tips_" + str + ".csv")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()

	samples := p.nRun
	if samples > p.recordSamples {
		samples = p.recordSamples
	}
	_, err = f.WriteString(a.nTipsToString(p, samples)) // writing...

	if err != nil {
		fmt.Printf("error writing string: %v", err)
		return err
	}

	return nil

}

func (a tipsResult) SaveOrphanTips(p Parameters) (err error) {
	str := fmt.Sprintf("%d", int(p.SimStep))
	f, err := os.Create("data/orphantips_" + str + ".csv")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(a.nOrphanTipsToString(p)) // writing...

	if err != nil {
		fmt.Printf("error writing string: %v", err)
		return err
	}

	return nil
}

func (a orphanResult) SaveOrphanTxs(p Parameters) (err error) {
	str := fmt.Sprintf("%d", int(p.SimStep))
	f, err := os.Create("data/orphantxs_" + str + ".csv")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(a.nOrphanTxsToString(p)) // writing...

	if err != nil {
		fmt.Printf("error writing string: %v", err)
		return err
	}

	return nil
}

func MeanVariance(x []float64) (mean, variance float64) {
	// Note that this will panic if the slice lengths do not match.
	mean = stat.Mean(x, nil)
	var (
		ss float64
	)
	for _, v := range x {
		d := v - mean
		ss += d * d
	}
	variance = (ss) / float64(len(x))
	return mean, variance
}
