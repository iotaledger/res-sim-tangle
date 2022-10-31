package main

import (
	"fmt"
	"strings"
)

// main routine
func main() {

	b := make(Benchmark)
	_ = b
	//runRealDataEvaluation(10, 0, true)
	runForVariables(b)
	// runSimulation(b, 10)
	// printPerformance(b)
}

func runSimulation(b Benchmark, x float64, simStep int) Result {

	p := newParameters(x, simStep)
	defer b.track(runningtime("TSA=" + strings.ToUpper(p.TSA) + ", X=" + fmt.Sprintf("%.2f", x) + ", " + "\tTime"))
	c := make(chan bool, p.nParallelSims)
	r := make([]Result, p.nParallelSims)
	var f Result
	f.params = p

	//readjusted for number of cores
	p.nRun /= p.nParallelSims

	for i := 0; i < p.nParallelSims; i++ {
		p.Seed = int64(i*p.nRun + 1)
		// run the simulation
		go run(p, &r[i], c)
	}
	for i := 0; i < p.nParallelSims; i++ {
		<-c
	}

	for _, batch := range r {
		f.JoinResults(batch, p)
	}

	fmt.Println("\nTSA=", strings.ToUpper(p.TSA), "\tLambda=", p.Lambda, "\tD=", p.D)
	// save some results in files
	f.FinalEvaluationSaveResults(p)
	fmt.Println("- - - OrphanTips - - -")
	fmt.Println("X\t\tmean\t\tSTD\t\tmean Ratio\t\tSTD Ratio")
	fmt.Println(x, "\t", f.tipsResult.meanOrphanTips, "\t", f.tipsResult.STDOrphanTips, "\t", f.tipsResult.meanOrphanTipsRatio, "\t", f.tipsResult.STDOrphanTipsRatio)
	return f
}

func run(p Parameters, r *Result, c chan bool) {
	defer func() { c <- true }()
	b := make(Benchmark)
	*r, b = p.RunTangle()
	printPerformance(b)
}

func runForVariables(b Benchmark) {
	var total string
	// Xs := []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512}
	// Xs := []float64{0, .1, .2, .3, .4, .5, .6, .7, .8, .9}
	// Xs := []float64{0.5, .55, .6}
	// Xs := []float64{0, .05, .1, .15, .2, .25, .3, .35, .4, .45, .5, .55, .6, .7, .8, .9, .99}
	NXs := 20
	NXs2 := 10
	Xs := make([]float64, 1+NXs+1+NXs2)
	// for i1 := 0; i1 < NXs2; i1++ {
	// 	Xs[i1] = 2./float64(NXs2)*float64(i1+1.) + 1.
	// }
	Xs[0] = .3
	for i1 := 0; i1 < NXs+1; i1++ {
		// Xs[i1] = 1. / float64(NXs) * float64(i1) * .98 * (1. - 1./8.)
		Xs[1+i1] = 1./float64(NXs)*float64(i1)*.2 + .4
		// Xs[i1] = 2 * (float64(i1) + 1)
		// Xs[i1+NXs2] = 4. + float64(i1)
	}
	for i1 := 0; i1 < NXs2; i1++ {
		// Xs[i1] = 1. / float64(NXs) * float64(i1) * .98 * (1. - 1./8.)
		Xs[1+i1+NXs+1] = 1./float64(NXs2)*float64(i1+1)*(1.-.5) + .6
		// Xs[i1] = 2 * (float64(i1) + 1)
	}
	// for i1 := 0; i1 < NXs; i1++ {
	// 	Xs[i1] = .1 * math.Pow(100, float64(i1)/float64(NXs-1))
	// }
	fmt.Println("- - - - - - - - - - - - - - - - ")
	fmt.Println("Variables=", Xs)
	var banner string
	i := 0
	initParamsLog()
	for _, x := range Xs {
		fmt.Println("X=", x)
		r := runSimulation(b, x, i)
		i++
		if banner == "" {
			banner += fmt.Sprintf("#x\tOrphanratio\tSTD\ttipsAVG\ttipsSTD\t#txs\n")
		}

		output := fmt.Sprintf("%.4f", x)
		output += fmt.Sprintf("\t%.8f", r.tipsResult.meanOrphanTipsRatio)
		output += fmt.Sprintf("\t%.8f", r.tipsResult.STDOrphanTipsRatio)
		output += fmt.Sprintf("\t%.8f", r.tipsResult.tAVG)
		output += fmt.Sprintf("\t%.8f", r.tipsResult.tSTD)
		output += fmt.Sprintf("\t%d", r.params.TangleSize-r.params.minCut)
		output += fmt.Sprintf("\n")
		total += output
		fmt.Println(banner + output)
		writetoParamsLog(x)
	}
	fmt.Println(banner + total)

}
