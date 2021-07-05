package main

import (
	"fmt"
	"math/rand"

	"github.com/schollz/progressbar"
)

// ??? is there a reason why approvers is not part of the tx variable, i.e. Tx has the field app []int? This would seem much more intuitive...
// Sim contains the data structure of a Tangle simulation
type Sim struct {
	tangle     []Tx      // A Tangle, i.e., a list of transactions
	tangleAge  int       // number of issued transactions
	mana       []float64 // A list of the mana of the issuing nodes, total mana normalized to 1
	tips       []int     // A list of current available/visible tips
	orphanTips []int     // A list of old tips for RURTS
	hiddenTips []int     // A list of yet unavailable/hidden tips
	// approvers      map[int][]int // A map of direct approvers, e.g., 5 <- 10,13
	cw            [][]uint64 // Matrix of propagated weigth branches (cw[i][] is the column of bit values forthe ith tx, stored as uint64 blocks)
	nodeApprover  [][]bool   // slice of approver for each tx; nodeApprover[i][x] meaning that node i (in-directly) approves transactions x
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
	var nodeID int
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

		sim.generator = rand.New(rand.NewSource(p.Seed + int64(run)))
		//rand.Seed(p.Seed + int64(run))

		sim.tangle[0] = sim.newGenesis()

		for i := 0; i < p.numberNodes; i++ {
			sim.mana[i] = float64(1) / float64(p.numberNodes)
		}
		//nTips := 0

		if p.Seed == int64(1) {
			bar.Add(1)
		}
		fmt.Println(sim.mana)
		// counter := 0
		for i := 1; i < sim.param.TangleSize; i++ {
			// choose node to issue next transaction
			nodeID = sim.generator.Int() % p.numberNodes
			//generate new tx
			t := newTx(&sim, sim.tangle[i-1], nodeID)
			//fmt.Println("tx", i)

			//update set of tips before running TSA, increase the wb matrix here
			sim.removeOldTips(t)
			sim.tips = append(sim.tips, sim.revealTips(t)...)
			fmt.Println("Tip sets", sim.tips)
			//run TSA to select tips to approve
			if sim.isAdverse(i) {
				t.ref = sim.param.tsaAdversary.TipSelectAdversary(t, &sim) // adversary tip selection
			} else {
				t.ref = sim.param.tsa.TipSelect(t, &sim) //sim.tipsSelection(t, sim.vTips)
			}
			fmt.Println(t.ref)
			//add the new tx to the Tangle and to the hidden tips set
			sim.tangle[i] = t
			//fmt.Println("Increase tangle age")
			sim.tangleAge += 1
			sim.hiddenTips = append(sim.hiddenTips, t.id)
			//update approving nodes ( atm this is done immedidately, and not after tips becomes visible)
			//fmt.Println("Nex tx", t)
			sim.updateAW(t, nodeID)
			//fmt.Println(sim.tangle[1].aw)
			// gather results
			result.EvaluateAfterTx(&sim, p, run, i)

		}
		for i := 0; i < p.TangleSize; i++ {
			fmt.Println(sim.tangle[i].confirmationTime)
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
	sim.tangleAge = 0
	//sim.nodeApprover = [][]bool{}
	sim.nodeApprover = make([][]bool, sim.param.numberNodes)
	for i := 0; i < sim.param.numberNodes; i++ {
		sim.nodeApprover[i] = make([]bool, sim.param.TangleSize)
		for j := 0; j < sim.param.TangleSize; j++ {
			sim.nodeApprover[i][j] = false
		}
	}
	//sim.mana = []float64{}
	sim.mana = make([]float64, sim.param.numberNodes)
	sim.tangle = make([]Tx, sim.param.TangleSize)
	sim.tips = []int{}
	sim.orphanTips = []int{}
	sim.hiddenTips = []int{}

	// sim.spinePastCone = make(map[int]Tx)
	// sim.spineApprovers = make(map[int][]int)
}

func (sim Sim) isAdverse(i int) bool {
	isAdverse := false
	if i > sim.param.TangleSize/3 {
		if i < sim.param.TangleSize*2/3 {
			if sim.generator.Float64() < sim.param.q {
				isAdverse = true
			}
		}
	}
	return isAdverse
}
