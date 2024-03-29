package main

import (
	"math/rand"

	"github.com/schollz/progressbar"
)

// ??? is there a reason why approvers is not part of the tx variable, i.e. Tx has the field app []int? This would seem much more intuitive...
// Sim contains the data structure of a Tangle simulation
type Sim struct {
	tangle     []Tx  // A Tangle, i.e., a list of transactions
	tips       []int // A list of current available/visible tips
	orphanTips []int // A list of old tips for RURTS
	hiddenTips []int // A list of yet unavailable/hidden tips
	// approvers      map[int][]int // A map of direct approvers, e.g., 5 <- 10,13
	cwMatrix      [][]uint64 // Matrix of propagated weigth branches (cw[i][] is the column of bit values for the ith tx, stored as uint64 blocks)
	generator     *rand.Rand // An unsafe random generator
	param         Parameters // Set of simulation parameters
	b             Benchmark  // Data structure to save performance of the simulation
	spinePastCone map[int]Tx
	// spineApprovers map[int][]int
}

// RunTangle executes the simulation
func (p *Parameters) RunTangle() (Result, Benchmark) {
	performance := make(Benchmark)
	defer performance.track(runningtime("total time"))
	sim := Sim{}

	var result Result
	sim.param = *p
	result.initResults(p)
	sim.clearSim()
	//fmt.Println(p.nRun)
	bar := progressbar.New(sim.param.nRun)

	// - - - - - - - - - - - - - - - - - - - - -
	// run nRun tangle sims
	// - - - - - - - - - - - - - - - - - - - - -
	for run := 0; run < sim.param.nRun; run++ {

		sim.clearSim()
		//fmt.Println(sim)
		sim.generator = rand.New(rand.NewSource(p.Seed + int64(run)))
		//rand.Seed(p.Seed + int64(run))

		sim.tangle[0] = sim.newGenesis()
		//nTips := 0

		if p.Seed == int64(1) {
			bar.Add(1)
		}

		// counter := 0
		for i := 1; i < sim.param.TangleSize; i++ {

			//generate new tx
			t := newTx(&sim, sim.tangle[i-1])
			// fmt.Println("tx", i)

			//update set of tips before running TSA, increase the wb matrix here
			sim.removeOldTips(t)
			sim.tips = append(sim.tips, sim.revealTips(t)...)

			//run TSA to select tips to approve
			if sim.isAdverse(i) {
				t.ref = sim.param.tsaAdversary.TipSelectAdversary(t, &sim) // adversary tip selection
			} else {
				t.ref = sim.param.tsa.TipSelect(t, &sim) //sim.tipsSelection(t, sim.vTips)
			}

			//add the new tx to the Tangle and to the hidden tips set
			sim.tangle[i] = t
			sim.hiddenTips = append(sim.hiddenTips, t.id)

			result.EvaluateAfterTx(&sim, p, run, i)

		}
		//saveTangle(sim.tangle)
		//fmt.Println("\n\n")
		//fmt.Println("Tangle size: ", sim.param.TangleSize)

		//	fmt.Println(getCWgrowth(sim.tangle[sim.param.TangleSize-10*int(sim.param.Lambda)], &sim))
		//fmt.Println(sim.tangle[sim.param.TangleSize-10*int(sim.param.Lambda)].cw)

		//Compare CWs
		//fmt.Println("CW comparison:", sim.compareCW())
		// data evaluation after each tangle
		//result.avgtips.val = append(result.avgtips.val, float64(nTips)/float64(sim.param.TangleSize-sim.param.minCut-sim.param.maxCutrange)/sim.param.Lambda)
		result.EvaluateTangle(&sim, p, run)

		//Visualize the Tangle
		if p.drawTangleMode > 0 {
			sim.visualizeTangle(nil, p.drawTangleMode)
		}

	}

	//fmt.Println("E(L):", float64(nTips)/float64(sim.param.TangleSize-sim.param.minCut*2)/sim.param.Lambda/float64(sim.param.nRun))
	return result, performance
}

func (sim *Sim) clearSim() {
	// sim.approvers = make(map[int][]int)
	sim.b = make(Benchmark)

	//sim.cw = [][]uint64{}
	// sim.cw = make([][]uint64, sim.param.CWMatrixLen)

	sim.tangle = make([]Tx, sim.param.TangleSize)
	sim.tips = []int{}
	sim.orphanTips = []int{}
	sim.hiddenTips = []int{}

	// sim.spinePastCone = make(map[int]Tx)
	// sim.spineApprovers = make(map[int][]int)
}

func (sim Sim) isAdverse(i int) bool {
	isAdverse := false
	attackActive := true
	if sim.param.qPartiallyActive {
		if i < sim.param.TangleSize/3 || i > sim.param.TangleSize*2/3 {
			attackActive = false
		}
	}
	if (sim.generator.Float64() < sim.param.q) && attackActive {
		isAdverse = true
	}
	return isAdverse
}
