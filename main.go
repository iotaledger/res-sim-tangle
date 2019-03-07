package main

import (
	"fmt"
	"runtime"
	"strings"

	"gonum.org/v1/gonum/stat"
)

// var nParallelSims = 1

// factor 2 is to use the physical cores, whereas NumCPU returns double the number due to hyper-threading
var nParallelSims = runtime.NumCPU()/2 - 1

func main() {

	b := make(Benchmark)
	//var ratio string
	var total string
	lambdas := []float64{100}
	// lambdas := []float64{1, 2, 3, 6, 10, 20, 30, 60, 100, 200, 300, 600}
	// lambdas := []float64{600, 300}
	//alphas := []float64{0}
	for _, lambda := range lambdas {
		//for _, alpha := range alphas {
		for alpha := 0.001; alpha <= 0.1; alpha += 0.001 {
			//for lambda := 1.; lambda <= 100; lambda++ {
			// if (alpha * lambda) < 10 {
			r := runSimulation(b, "rw", lambda, alpha)
			ratio := fmt.Sprintf("%.3f", alpha)

			for _, m := range []string{"rw", "backU", "backB", "CW-Max"} {
				x, y := r.velocity.getTimeMetric(m)
				ratio += fmt.Sprintf("\t%.5f", stat.Mean(x, y))
			}
			ratio += fmt.Sprintf("\t%.5f", stat.Mean(r.op.op, nil))
			ratio += fmt.Sprintf("\t%.5f", stat.Mean(r.op.top, nil))
			ratio += fmt.Sprintf("\t%.5f", stat.Mean(r.op.op2, nil))
			ratio += fmt.Sprintf("\t%.5f", stat.Mean(r.op.top2, nil))
			ratio += fmt.Sprintf("\n")

			total += ratio
			fmt.Println(ratio)
		}
	}

	// Options: RW, URTS
	// runSimulation(b, "urts", 10, 0)
	//r := runSimulation(b, "rw", 100, 0.005)
	fmt.Println(total)

	//printPerformance(b)
}

func runSimulation(b Benchmark, tsa string, lambda, alpha float64) Result {
	defer b.track(runningtime("TSA=" + strings.ToUpper(tsa) + ", Lambda=" + fmt.Sprintf("%.2f", lambda) + ", Alpha=" + fmt.Sprintf("%.4f", alpha) + "\tTime"))

	//lambda := 100.
	p := Parameters{
		//K:          2,
		//H:          1,
		Lambda:      lambda,
		Alpha:       alpha,
		TangleSize:  200 * int(lambda),
		CWMatrixLen: 200 * int(lambda), // reduce CWMatrix to this len
		// TangleSize:   int(math.Min(3000, (100+math.Max(100, 30.0/alpha/lambda)))) * int(lambda),
		minCut:       51 * int(lambda), // cut data close to the genesis
		maxCutrange:  50 * int(lambda), // cut data for the most recent txs, not applied for every analysis
		stillrecent:  2 * int(lambda),  // when is a tx considered recent, and when is it a candidate for left behind
		ConstantRate: false,
		// nRun:         int(math.Max(10000/lambda, 100)),
		nRun: 20,
		TSA:  tsa,

		// - - - Analysis section - - -
		CountTipsEnabled:  false,
		CWAnalysisEnabled: false,
		SpineEnabled:      true,
		pOrphanEnabled:    true,
		VelocityEnabled:   true,
		EntropyEnabled:    false,
		// VelocityEnabled: false,
		//{Enabled, Resolution, MaxT, MaxApp}
		AnPastCone: AnPastCone{false, 5, 40, 5},
		//{Enabled, maxiMT, murel, nRW}
		AnFocusRW: AnFocusRW{false, 0.2, 30},
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
	fmt.Println(f.avgtips)
	// if p.VelocityEnabled {
	// 	fmt.Println(f.velocity.Stat(p))
	// 	f.velocity.Save(p)
	// 	f.velocity.SaveStat(p)
	// }
	// if p.AnPastEdgeEnabled {
	// 	f.PastEdge.Save(p)
	// }
	f.SaveResults(p)
	return f
}

func run(p Parameters, r *Result, c chan bool) {
	defer func() { c <- true }()
	b := make(Benchmark)
	*r, b = p.RunTangle()
	printPerformance(b)
}
