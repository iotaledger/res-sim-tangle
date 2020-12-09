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
	runSimulation(b, 20)
	// printPerformance(b)
}

func runSimulation(b Benchmark, lambda float64) Result {

	p := newParameters(lambda)
	defer b.track(runningtime("TSA=" + strings.ToUpper(p.TSA) + ", Lambda=" + fmt.Sprintf("%.2f", lambda) + ", " + "\tTime"))
	c := make(chan bool, p.nParallelSims)
	r := make([]Result, p.nParallelSims)
	var f Result

	for i := 0; i < p.nParallelSims; i++ {
		p.Seed = int64(i*p.nRun + 1)
		go run(p, &r[i], c)
	}
	for i := 0; i < p.nParallelSims; i++ {
		<-c
	}

	for _, batch := range r {
		f.JoinResults(batch, p)
	}

	fmt.Println("\nTSA=", strings.ToUpper(p.TSA), "\tLambda=", p.Lambda, "\tD=", p.D)
	f.FinalEvaluationSaveResults(p)
	fmt.Println("- - - OrphanTips - - -")
	fmt.Println("Lambda\t\tD\t\tmean\t\tSTD\t\tmean Ratio\t\tSTD Ratio")
	fmt.Println(p.Lambda, p.D, f.tips.meanOrphanTips, f.tips.STDOrphanTips, f.tips.meanOrphanTipsRatio, f.tips.STDOrphanTipsRatio)
	return f
}

func run(p Parameters, r *Result, c chan bool) {
	defer func() { c <- true }()
	b := make(Benchmark)
	*r, b = p.RunTangle()
	printPerformance(b)
}
