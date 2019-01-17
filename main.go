package main

import (
	"fmt"
	"runtime"
	"strings"
)

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
	runSimulation(b, "urts", 100, 0)

	printPerformance(b)
}

func runSimulation(b Benchmark, tsa string, lambda, alpha float64) {
	defer b.track(runningtime("TSA=" + strings.ToUpper(tsa) + ", Lambda=" + fmt.Sprintf("%.2f", lambda) + ", Alpha=" + fmt.Sprintf("%.4f", alpha) + "\tTime"))

	//lambda := 100.
	p := Parameters{
		//K:          2,
		//H:          1,
		Lambda:          lambda,
		Alpha:           alpha,
		TangleSize:      1000 * int(lambda),
		ConstantRate:    false,
		nRun:            2,
		TSA:             tsa,
		VelocityEnabled: false,
	}
	c := make(chan bool, nParallelSims)
	r := make([]result, nParallelSims)
	var f result

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
		f.tips = f.tips.Join(batch.tips)
	}

	fmt.Println("\nTSA=", strings.ToUpper(p.TSA), "\tLambda=", p.Lambda, "\tAlpha=", p.Alpha)
	fmt.Println(f.tips)
	if p.VelocityEnabled {
		fmt.Println(f.velocity.Stat(p))
		f.velocity.Save(p)
		f.velocity.SaveStat(p)
	}
}

func run(p Parameters, r *result, c chan bool) {
	defer func() { c <- true }()
	b := make(Benchmark)
	*r, b = p.RunTangle()
	printPerformance(b)
}
