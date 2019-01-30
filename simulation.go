package main

import (
	"fmt"
	"math/rand"
	"strings"

	"gonum.org/v1/gonum/stat"
)

// Sim contains the data structure of a Tangle simulation
type Sim struct {
	tangle     []Tx          // A Tangle, i.e., a list of transactions
	tips       []int         // A list of current available/visible tips
	hiddenTips []int         // A list of yet unavailable/hidden tips
	approvers  map[int][]int // A map of direct approvers, e.g., 5 <- 10,13
	cw         [][]uint64    // Matrix of cumulative weigth branches
	generator  *rand.Rand    // An unsafe random generator
	param      Parameters    // Set of simulation parameters
	b          Benchmark     // Data structure to save performance of the simulation
}

func (p Parameters) initSim(sim *Sim) {

	clearSim(sim)

	if p.K != 0 {
		sim.param.K = p.K
	} else {
		sim.param.K = 2
	}

	if p.H != 0 {
		sim.param.H = p.H
	} else {
		sim.param.H = 1
	}

	if p.Lambda != 0 {
		sim.param.Lambda = p.Lambda
	} else {
		sim.param.Lambda = 1
	}

	if p.Alpha != 0 {
		sim.param.Alpha = p.Alpha
	} else {
		sim.param.Alpha = 0
	}

	if p.TangleSize != 0 {
		sim.param.TangleSize = p.TangleSize
	} else {
		sim.param.Alpha = 0
	}

	if p.Seed != 0 {
		sim.param.Seed = p.Seed
	} else {
		sim.param.Seed = 1
	}

	if p.nRun != 0 {
		sim.param.nRun = p.nRun
	} else {
		sim.param.nRun = 1
	}

	switch strings.ToUpper(p.TSA) {
	case "URTS":
		sim.param.TSA = p.TSA
		sim.param.tsa = URTS{}
	case "RW":
		if p.Alpha == 0 {
			sim.param.TSA = "RW"
			sim.param.tsa = URW{}
		} else {
			sim.param.TSA = "RW"
			sim.param.tsa = BRW{}
		}
	default:
		sim.param.TSA = "URTS"
		sim.param.tsa = URTS{}
	}

	sim.param.ConstantRate = p.ConstantRate
	sim.param.VelocityEnabled = p.VelocityEnabled
	sim.param.ReusableAddressEnabled = p.ReusableAddressEnabled

	if p.DataPath != "" {
		sim.param.DataPath = p.DataPath
	}

	sim.param.minCut = 30 * int(sim.param.Lambda)
	sim.param.maxCut = sim.param.TangleSize - sim.param.minCut

	createDirIfNotExist("data")

}

func clearSim(sim *Sim) {
	sim.approvers = make(map[int][]int)
	sim.b = make(Benchmark)

	sim.cw = [][]uint64{}
	sim.tangle = make([]Tx, sim.param.TangleSize)
	sim.tips = []int{}
	sim.hiddenTips = []int{}
}

// RunTangle executes the simulation
func (p *Parameters) RunTangle() (Result, Benchmark) {
	performance := make(Benchmark)
	defer performance.track(runningtime("total"))
	//fmt.Println(p)
	sim := Sim{}
	var nTips int

	var result Result

	if p.VelocityEnabled {
		vr := newVelocityResult([]string{"rw", "all", "first", "last", "second", "third", "fourth", "only-1"})
		result.velocity = *vr
	}

	p.initSim(&sim)
	//fmt.Println(p.nRun)
	//bar := progressbar.New(sim.param.nRun)
	//bar := progressbar.New(sim.param.TangleSize)

	timesDirect := make(map[int64]int)
	timesReverse := make(map[int64]int)
	visitedDirect := make(map[int]int)
	visitedReverse := make(map[int]int)

	for run := 0; run < sim.param.nRun; run++ {

		clearSim(&sim)
		//fmt.Println(sim)
		sim.generator = rand.New(rand.NewSource(p.Seed + int64(run)))
		//rand.Seed(p.Seed + int64(run))

		sim.tangle[0] = sim.newGenesis()

		// if p.Seed == int64(1) {
		// 	bar.Add(1)
		// }

		for i := 1; i < sim.param.TangleSize; i++ {

			// if p.Seed == int64(1) {
			// 	bar.Add(1)
			// }

			//generate new tx
			t := newTx(&sim, sim.tangle[i-1])

			//update set of tips before running TSA
			sim.tips = append(sim.tips, sim.tipsUpdate(t)...)

			//run TSA to select k(2) tips to approve
			t.ref = sim.param.tsa.TipSelect(t, &sim) //sim.tipsSelection(t, sim.vTips)

			//add the new tx to the Tangle and to the hidden tips set
			sim.tangle[i] = t
			sim.hiddenTips = append(sim.hiddenTips, t.id)

			if i > sim.param.minCut && i < sim.param.maxCut {
				nTips += len(sim.tips)
			}

		}
		//printApprovers(sim.approvers)
		//fmt.Println("CW comparison:", sim.compareCW())
		//printPerformance(sim.b)
		//printCWRef(sim.cw)
		//fmt.Println(sim.tangle)

		result.tips.tips = float64(nTips) / float64(sim.param.TangleSize-sim.param.minCut*2) / sim.param.Lambda / float64(sim.param.nRun)
		if p.VelocityEnabled {
			sim.runVelocityStat(&result.velocity)
		}
		if sim.param.ReusableAddressEnabled {
			for i := 0; i < 1; i++ {
				t, v := sim.directOrder(sim.tangle[sim.tips[len(sim.tips)-1]]) //start from last visible tip
				timesDirect[t]++
				visitedDirect[v]++
				t, v = sim.reverseOrder(sim.tangle[0]) //start from genesis
				timesReverse[t]++
				visitedReverse[v]++
			}
		}
	}

	if sim.param.ReusableAddressEnabled {
		// calculate statistics for direct order
		fmt.Printf("\nVersion \tavg[ms] \tstd \tvar\t#visited\n")
		var weigths []float64
		var x []float64
		for k, v := range timesDirect {
			x = append(x, float64(k)/1000000)
			weigths = append(weigths, float64(v))
		}
		avgT, std := stat.MeanStdDev(x, weigths)
		_, variance := stat.MeanVariance(x, weigths)

		weigths = []float64{}
		x = []float64{}
		for k, v := range visitedDirect {
			x = append(x, float64(k))
			weigths = append(weigths, float64(v))
		}
		avg := stat.Mean(x, weigths)

		fmt.Printf("Direct order: \t%0.4f \t%0.4f \t%0.4f \t%d\n", avgT, std, variance, int(avg))

		// calculate statistics for reverse order
		weigths = []float64{}
		x = []float64{}
		for k, v := range timesReverse {
			x = append(x, float64(k)/1000000)
			weigths = append(weigths, float64(v))
		}
		avgT, std = stat.MeanStdDev(x, weigths)
		_, variance = stat.MeanVariance(x, weigths)

		weigths = []float64{}
		x = []float64{}
		for k, v := range visitedReverse {
			x = append(x, float64(k))
			weigths = append(weigths, float64(v))
		}
		avg = stat.Mean(x, weigths)

		fmt.Printf("Reverse order: \t%0.4f \t%0.4f \t%0.4f \t%d\n", avgT, std, variance, int(avg))
	}
	//saveTangle(sim.tangle)

	//fmt.Println("E(L):", float64(nTips)/float64(sim.param.TangleSize-sim.param.minCut*2)/sim.param.Lambda/float64(sim.param.nRun))
	return result, performance
}
