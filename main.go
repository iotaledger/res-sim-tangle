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
	_ = b
	// Options: RW, URTS
	runSimulation(b, "urts", 100, 0)
	// runSimulation(b, "rw", 10, 0)
	// fmt.Println(runForAlphasLambdas(b))

	//printPerformance(b)
}

func runSimulation(b Benchmark, tsa string, lambda, alpha float64) Result {
	defer b.track(runningtime("TSA=" + strings.ToUpper(tsa) + ", Lambda=" + fmt.Sprintf("%.2f", lambda) + ", Alpha=" + fmt.Sprintf("%.4f", alpha) + "\tTime"))

	//lambda := 100.
	p := Parameters{
		//K:          2,
		//H:          1,
		Lambda:       lambda,
		Alpha:        alpha,
		TangleSize:   200 * int(lambda),
		CWMatrixLen:  200 * int(lambda), // reduce CWMatrix to this len
		minCut:       51 * int(lambda),  // cut data close to the genesis
		maxCutrange:  52 * int(lambda),  // cut data for the most recent txs, not applied for every analysis
		stillrecent:  2 * int(lambda),   // when is a tx considered recent, and when is it a candidate for left behind
		ConstantRate: false,
		// nRun:         int(10000 / lambda),
		nRun: 100,
		TSA:  tsa,

		// - - - Analysis section - - -
		CountTipsEnabled:     false,
		CWAnalysisEnabled:    false,
		SpineEnabled:         false,
		pOrphanEnabled:       false, // calculate orphanage probability
		pOrphanLinFitEnabled: false, // also apply linear fit, numerically expensive
		VelocityEnabled:      false,
		ExitProbEnabled:      false,
		ExitProbNparticle:    10000, // number of sample particles to calculate distribution
		ExitProb2NHisto:      50,    // N of Histogram columns for exitProb2
		DistSlicesEnabled:    true,  // calculate the distances of slices
		// DistSlicesLength:     100 / lambda, //length of Slices
		DistSlicesLength:     0.1, //length of Slices
		DistSlicesResolution: 100, // Number of intervals per distance '1', higher number = higher resolution
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
	f.SaveResults(p)
	return f
}

func run(p Parameters, r *Result, c chan bool) {
	defer func() { c <- true }()
	b := make(Benchmark)
	*r, b = p.RunTangle()
	printPerformance(b)
}

func runForAlphasLambdas(b Benchmark) string {
	// b := make(Benchmark)
	//var ratio string
	var total string
	// lambdas := []float64{50}
	lambdas := []float64{3, 10, 30, 100, 300}
	// lambdas := []float64{600, 300}
	//alphas := []float64{0, 0.01, 0.1, 1}
	alphas := []float64{0.}
	var banner string
	for _, lambda := range lambdas {
		for _, alpha := range alphas {
			//for alpha := 0.001; alpha <= 0.1; alpha += 0.001 {
			//for lambda := 1.; lambda <= 100; lambda++ {
			// if (alpha * lambda) < 10 {
			// r := runSimulation(b, "rw", lambda, alpha)
			r := runSimulation(b, "urts", lambda, alpha)
			if banner == "" {
				banner += fmt.Sprintf("#alpha\t")
				for _, m := range r.velocity.vTime {
					banner += fmt.Sprintf("%v\t", m.desc)
				}
				banner += fmt.Sprintf("OP\tTOP\n")
			}

			output := fmt.Sprintf("%.3f", alpha)
			for _, m := range r.velocity.vTime {
				x, y := r.velocity.getTimeMetric(m.desc)
				output += fmt.Sprintf("\t%.5f", stat.Mean(x, y))
			}
			// output += fmt.Sprintf("\t%.5f", stat.Mean(r.op.op, nil))
			// output += fmt.Sprintf("\t%.5f", stat.Mean(r.op.top, nil))
			output += fmt.Sprintf("\t%.5f", stat.Mean(r.op.op2, nil))
			output += fmt.Sprintf("\t%.5f", stat.Mean(r.op.top2, nil))
			output += fmt.Sprintf("\n")

			total += output
			fmt.Println(output)
		}
	}
	return banner + total
}
