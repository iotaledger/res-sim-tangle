package main

import (
	"fmt"
	"runtime"
	"strings"
)

// var nParallelSims = 1

// factor 2 is to use the physical cores, whereas NumCPU returns double the number due to hyper-threading
var nParallelSims = runtime.NumCPU()/2 - 1

func main() {

	b := make(Benchmark)
	// lambdas := []float64{1, 2, 3, 6, 10, 20, 30, 60, 100, 200, 300, 600}
	// lambdas := []float64{600, 300}
	// // alphas := []float64{0.001, 0.01, 0.1}
	// alphas := []float64{0.1}
	// for _, lambda := range lambdas {
	// 	for _, alpha := range alphas {
	// 		// if (alpha * lambda) < 10 {
	// 		runSimulation(b, "rw", lambda, alpha)
	// 		// }
	// 	}
	// }

	// Options: RW, URTS
	// runSimulation(b, "urts", 10, 0)
	runSimulation(b, "rw", 100, 0)

	printPerformance(b)
}

func runSimulation(b Benchmark, tsa string, lambda, alpha float64) {
	defer b.track(runningtime("TSA=" + strings.ToUpper(tsa) + ", Lambda=" + fmt.Sprintf("%.2f", lambda) + ", Alpha=" + fmt.Sprintf("%.4f", alpha) + "\tTime"))

	//this is only a temporary code to provide the correct tangle size
	tmpHelp := 30.0 / alpha / lambda
	_ = tmpHelp
	if alpha*lambda > 1 {
		tmpHelp = 100
	}

	//lambda := 100.
	p := Parameters{
		//K:          2,
		//H:          1,
		Lambda:     lambda,
		Alpha:      alpha,
		TangleSize: 500 * int(lambda),
		// TangleSize:   int(math.Min(3000, (100+math.Max(100, 30.0/alpha/lambda)))) * int(lambda),
		minCut:       30 * int(lambda),
		maxCutrange:  30 * int(lambda),
		ConstantRate: false,
		// nRun:         int(math.Max(10000/lambda, 100)),
		nRun:        1,
		TSA:         tsa,
		stillrecent: 2 * int(lambda), // when is a tx considered recent, and when is it a candidate for left behind

		// - - - Analysis section - - -
		VelocityEnabled: false,
		EntropyEnabled:  false,
		//{Enabled, Resolution, MaxT, MaxApp}
		AnPastCone: AnPastCone{false, 40, 10, 5},
		//{Enabled, maxiMT, murel, nRW}
		AnFocusRW: AnFocusRW{false, 0.3, 30},
	}

	c := make(chan bool, nParallelSims)
	r := make([]Result, nParallelSims)
	var f Result

	for i := 0; i < nParallelSims; i++ {
		p.Seed = int64(i*p.nRun + 1)
		go run(p, &r[i], c)
	}
	for i := 0; i < nParallelSims; i++ {
		<-c
	}

	for _, batch := range r {
		f.JoinResults(batch, p)
	}

	fmt.Println("\nTSA=", strings.ToUpper(p.TSA), "\tLambda=", p.Lambda, "\tAlpha=", p.Alpha)
	fmt.Println(f.tips)
	// if p.VelocityEnabled {
	// 	fmt.Println(f.velocity.Stat(p))
	// 	f.velocity.Save(p)
	// 	f.velocity.SaveStat(p)
	// }
	// if p.AnPastEdgeEnabled {
	// 	f.PastEdge.Save(p)
	// }
	f.SaveResults(p)
}

func run(p Parameters, r *Result, c chan bool) {
	defer func() { c <- true }()
	b := make(Benchmark)
	*r, b = p.RunTangle()
	printPerformance(b)
}
