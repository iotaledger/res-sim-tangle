package main

import (
	"math/rand"
	"strings"

	"github.com/schollz/progressbar"
)

// Sim contains the data structure of a Tangle simulation
type Sim struct {
	tangle     []Tx          // A Tangle, i.e., a list of transactions
	tips       []int         // A list of current available/visible tips
	hiddenTips []int         // A list of yet unavailable/hidden tips
	approvers  map[int][]int // A map of direct approvers, e.g., 5 <- 10,13
	cw         [][]uint64    // Matrix of propagated weigth branches (cw[i][] is the column of bit values forthe ith tx, stored as uint64 blocks)
	generator  *rand.Rand    // An unsafe random generator
	param      Parameters    // Set of simulation parameters
	b          Benchmark     // Data structure to save performance of the simulation
}

// RunTangle executes the simulation
func (p *Parameters) RunTangle() (Result, Benchmark) {
	performance := make(Benchmark)
	defer performance.track(runningtime("total"))
	//fmt.Println(p)
	sim := Sim{}
	var nTips int

	var result Result

	p.initSim(&sim)

	// - - - - - - - - - - - - - - - - - - - - -
	// initiate analysis variables
	// - - - - - - - - - - - - - - - - - - - - -
	if p.VelocityEnabled {
		//???is there a way this can be defined in the velocity.go file
		var vr *velocityResult
		if sim.param.TSA != "RW" {
			vr = newVelocityResult([]string{"rw", "all", "first", "last", "second", "third", "fourth", "only-1", "back"})
		} else {
			vr = newVelocityResult([]string{"rw", "all", "first", "last", "back"})
			//vr = newVelocityResult([]string{"rw", "all", "back"})
		}
		result.velocity = *vr
	}
	if p.AnPastCone.Enabled {
		//??? can this be combined into one line?
		r := newPastConeResult([]string{"avg", "1", "2", "3", "4", "5", "rest"})
		result.PastCone = *r
	}
	if p.AnFocusRW.Enabled {
		r := newFocusRWResult([]string{"0.1"})
		result.FocusRW = *r
	}

	//fmt.Println(p.nRun)
	bar := progressbar.New(sim.param.nRun)

	// - - - - - - - - - - - - - - - - - - - - -
	// run nRun tangle sims
	// - - - - - - - - - - - - - - - - - - - - -
	for run := 0; run < sim.param.nRun; run++ {

		clearSim(&sim)
		//fmt.Println(sim)
		sim.generator = rand.New(rand.NewSource(p.Seed + int64(run)))
		//rand.Seed(p.Seed + int64(run))

		sim.tangle[0] = sim.newGenesis()

		if p.Seed == int64(1) {
			bar.Add(1)
		}

		// counter := 0
		for i := 1; i < sim.param.TangleSize; i++ {
			//generate new tx
			t := newTx(&sim, sim.tangle[i-1])
			// fmt.Println("tx", i)

			//update set of tips before running TSA, increase the wb matrix here
			sim.tips = append(sim.tips, sim.tipsUpdate(t)...)

			// fmt.Println("sim.tips", sim.tips)
			// fmt.Println("sim.hiddenTips", sim.hiddenTips)
			// fmt.Println("sim.cw,sim.cw")
			// fmt.Println("len(sim.cw)", len(sim.cw))
			// for ; counter < len(sim.cw); counter++ {
			// 	fmt.Println("Counter now=", counter)
			// 	fmt.Println(sim.cw[counter][0])
			// 	fmt.Println(strconv.FormatInt(int64(sim.cw[counter][0]), 2))
			// }

			//run TSA to select k(2) tips to approve
			t.ref = sim.param.tsa.TipSelect(t, &sim) //sim.tipsSelection(t, sim.vTips)

			//add the new tx to the Tangle and to the hidden tips set
			sim.tangle[i] = t
			sim.hiddenTips = append(sim.hiddenTips, t.id)

			if i > sim.param.minCut && i < sim.param.maxCut {
				nTips += len(sim.tips)
			}

		}

		// - - - - - - - - - - - - - - - - - - - - -
		// data evaluation after each tangle
		// - - - - - - - - - - - - - - - - - - - - -
		result.avgtips.val = float64(nTips) / float64(sim.param.TangleSize-sim.param.minCut-sim.param.maxCutrange) / sim.param.Lambda / float64(sim.param.nRun)
		if p.VelocityEnabled {
			sim.runVelocityStat(&result.velocity)
		}
		if p.AnPastCone.Enabled {
			sim.runAnPastCone(&result.PastCone)
		}
		if p.AnFocusRW.Enabled {
			sim.runAnFocusRW(&result.FocusRW)
		}
	}

	//fmt.Println("E(L):", float64(nTips)/float64(sim.param.TangleSize-sim.param.minCut*2)/sim.param.Lambda/float64(sim.param.nRun))
	return result, performance
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
		sim.param.TangleSize = 0
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
	sim.param.stillrecent = p.stillrecent

	if p.AnPastCone.MaxApp != 0 {
		sim.param.AnPastCone.MaxApp = p.AnPastCone.MaxApp
	} else {
		sim.param.AnPastCone.MaxApp = 2
	}
	if p.AnPastCone.MaxT != 0 {
		sim.param.AnPastCone.MaxT = p.AnPastCone.MaxT
	} else {
		sim.param.AnPastCone.MaxT = 2
	}
	if p.AnPastCone.Resolution != 0 {
		sim.param.AnPastCone.Resolution = p.AnPastCone.Resolution
	} else {
		sim.param.AnPastCone.Resolution = 2
	}

	sim.param.AnFocusRW.murel = p.AnFocusRW.murel
	sim.param.AnFocusRW.nRWs = p.AnFocusRW.nRWs

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

	if p.DataPath != "" {
		sim.param.DataPath = p.DataPath
	}

	sim.param.minCut = p.minCut
	sim.param.maxCutrange = p.maxCutrange
	sim.param.maxCut = p.TangleSize - p.maxCutrange

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
