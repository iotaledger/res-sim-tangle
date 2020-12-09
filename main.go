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
	runForLambdas(b)
	// runSimulation(b, 20)
	// printPerformance(b)
}

func runSimulation(b Benchmark, x float64) Result {

	p := newParameters(x)
	defer b.track(runningtime("TSA=" + strings.ToUpper(p.TSA) + ", X=" + fmt.Sprintf("%.2f", x) + ", " + "\tTime"))
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
	fmt.Println("X\t\tD\t\tmean\t\tSTD\t\tmean Ratio\t\tSTD Ratio")
	fmt.Println(x, "\t", p.D, "\t", f.tips.meanOrphanTips, "\t", f.tips.STDOrphanTips, "\t", f.tips.meanOrphanTipsRatio, "\t", f.tips.STDOrphanTipsRatio)
	return f
}

func run(p Parameters, r *Result, c chan bool) {
	defer func() { c <- true }()
	b := make(Benchmark)
	*r, b = p.RunTangle()
	printPerformance(b)
}

func runForLambdas(b Benchmark) {
	var total string
	Xs := []float64{2, 3, 4, 5, 6, 7, 8}
	// Nlambdas := 10
	// lambdas := make([]float64, Nlambdas)
	// for i1 := 0; i1 < Nlambdas; i1++ {
	// 	lambdas[i1] = .1 * math.Pow(100, float64(i1)/float64(Nlambdas-1))
	// }

	var banner string
	for _, x := range Xs {
		if x > 0 {
			r := runSimulation(b, x)
			if banner == "" {
				banner += fmt.Sprintf("#lambda\tOrphanratio\tSTD\n")
			}

			output := fmt.Sprintf("%.0f", x)
			output += fmt.Sprintf("\t%.8f", r.tips.meanOrphanTipsRatio)
			output += fmt.Sprintf("\t%.8f", r.tips.STDOrphanTipsRatio)
			output += fmt.Sprintf("\n")

			total += output
			fmt.Println(banner + output)
		}
	}
	fmt.Println(banner + total)
}
