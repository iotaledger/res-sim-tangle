package main

import (
	"fmt"
	"runtime"
	"strings"
)

// var nParallelSims = 1

var nParallelSims = runtime.NumCPU()/2 - 1

func main() {
	b := make(Benchmark)
	// lambdas := []float64{1, 3, 5, 10, 30, 50}
	// alphas := []float64{0, 0.01, 0.1, 1}
	// for _, lambda := range lambdas {
	// 	for _, alpha := range alphas {
	// 		runSimulation(b, "rw", lambda, alpha)
	// 	}
	// }

	// Options: RW, URTS
	// runSimulation(b, "urts", 100, 0)
	runSimulation(b, "rw", 100, 0.001)

	printPerformance(b)
}

func runSimulation(b Benchmark, tsa string, lambda, alpha float64) {
	defer b.track(runningtime("TSA=" + strings.ToUpper(tsa) + ", Lambda=" + fmt.Sprintf("%.2f", lambda) + ", Alpha=" + fmt.Sprintf("%.4f", alpha) + "\tTime"))

	//lambda := 100.
	p := Parameters{
		//K:          2,
		//H:          1,
		Lambda:       lambda,
		Alpha:        alpha,
		TangleSize:   202 * int(lambda),
		minCut:       100 * int(lambda),
		maxCutrange:  100 * int(lambda),
		ConstantRate: false,
		nRun:         1,
		TSA:          tsa,
		// - - - Analysis section - - -
		VelocityEnabled: false,
		//{Enabled, Resolution, MaxT, MaxApp}
		AnPastCone: AnPastCone{false, 40, 10, 5},
		//{Enabled, maxiMT, murel, nRW}
		AnFocusRW: AnFocusRW{true, 300 * int(lambda), 0.3, 10},
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
		if p.VelocityEnabled {
			f.velocity = f.velocity.Join(batch.velocity)
		}
		if p.AnPastCone.Enabled {
			f.PastCone = f.PastCone.Join(batch.PastCone)
		}
		if p.AnFocusRW.Enabled {
			f.FocusRW = f.FocusRW.Join(batch.FocusRW)
		}
		f.tips = f.tips.Join(batch.tips)
	}

	fmt.Println("\nTSA=", strings.ToUpper(p.TSA), "\tLambda=", p.Lambda, "\tAlpha=", p.Alpha)
	fmt.Println(f.tips)
	if p.VelocityEnabled {
		// fmt.Println(f.velocity.Stat(p))
		f.velocity.Save(p)
		f.velocity.SaveStat(p)
	}
	if p.AnPastCone.Enabled {
		f.PastCone.finalprocess(p)
		f.PastCone.Save(p)
	}
	if p.AnFocusRW.Enabled {
		f.FocusRW.finalprocess(p)
		f.FocusRW.Save(p)
	}
}

func run(p Parameters, r *Result, c chan bool) {
	defer func() { c <- true }()
	b := make(Benchmark)
	*r, b = p.RunTangle()
	printPerformance(b)
}
