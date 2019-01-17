package main

import "fmt"

var limit = 5

func main() {
	b := make(Benchmark)
	// lambdas := []float64{1, 3, 5, 10, 30, 50}
	// alphas := []float64{0, 0.01, 0.1, 1}
	// for _, lambda := range lambdas {
	// 	for _, alpha := range alphas {
	// 		runSimulation(b, lambda, alpha)
	// 	}
	// }
	runSimulation(b, 1, 0.01)
	// runSimulation(b, 0.01)
	// runSimulation(b, 0.1)
	// runSimulation(b, 1.0)
	printPerformance(b)
}

func runSimulation(b Benchmark, lambda, alpha float64) {
	defer b.track(runningtime("Alpha: " + fmt.Sprintf("%f", alpha)))

	//lambda := 100.
	p := Parameters{
		//K:          2,
		//H:          1,
		Lambda:       lambda,
		Alpha:        alpha,
		TangleSize:   1000 * int(lambda),
		ConstantRate: false,
		nRun:         2,
		TSA:          "rw",
	}
	c := make(chan bool, limit)
	r := make([]velocityResult, limit)
	var f velocityResult

	for i := 0; i < limit; i++ {
		p.Seed = int64(i*p.nRun + 1)
		go run(p, &r[i], c)
	}
	for i := 0; i < limit; i++ {
		<-c
	}

	for _, batch := range r {
		f = f.Join(batch)
	}

	fmt.Println("\nTSA:", p.TSA, "Lambda:", p.Lambda, "Alpha", p.Alpha)
	fmt.Println(f.Stat(p))
	f.Save(p)
	f.SaveStat(p)
}

func run(p Parameters, r *velocityResult, c chan bool) {
	defer func() { c <- true }()
	b := make(Benchmark)
	*r, b = p.RunTangle()
	printPerformance(b)
}
