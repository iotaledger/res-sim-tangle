package main

import (
	"fmt"
	"math"
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

		// Populating  mana
		//f	or i := 0; i < p.numberNodes; i++ {
		//		sim.mana[i] = float64(1) / float64(p.numberNodes)
		//}
		sim.setMana(p.zipf)
		//fmt.Println(sim.mana)
		if p.Seed == int64(1) {
			bar.Add(1)
		}

		// counter := 0
		for i := 1; i < sim.param.TangleSize; i++ {
			// choose node to issue next transaction
			//nodeID = sim.generator.Int() % p.numberNodes // every node has same probability

			//choose nodes proportional to their mana
			var t Tx
			random := sim.generator.Float64()
			if random > sim.param.q { // tx is issued by honest node
				nodeID = sim.chooseNodeID()
				t = newTx(&sim, sim.tangle[i-1], nodeID)
				sim.removeOldTips(t)
				sim.tips = append(sim.tips, sim.revealTips(t)...)
				t.ref = sim.param.tsa.TipSelect(t, &sim) // honest TSA
			} else {
				nodeID = sim.param.adversaryID
				t = newTx(&sim, sim.tangle[i-1], nodeID)
				sim.removeOldTips(t)
				sim.tips = append(sim.tips, sim.revealTips(t)...)
				t.ref = sim.param.tsaAdversary.TipSelectAdversary(t, &sim) // adversary TSA
			}

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
		//for i := 0; i < p.TangleSize; i++ {
		//	fmt.Println(sim.tangle[i].confirmationTime)
		//}
		fmt.Println("\n ")
		fmt.Println("Simulations ended. Started to calculate statistics")
		// calculating the mean confirmation time
		var meanConfirmation float64
		counter := 0
		meanConfirmation = 0
		for _, tx := range sim.tangle {
			if tx.confirmationTime > -1 {
				counter += 1
				meanConfirmation += float64(tx.confirmationTime)
			}
		}
		meanConfirmation = meanConfirmation / float64(counter)
		//  calculating the sd confirmation time

		var sdConfirmation float64
		counter = 0
		sdConfirmation = 0
		for _, tx := range sim.tangle {
			if tx.confirmationTime > -1 {
				counter += 1
				sdConfirmation += math.Pow(float64(tx.confirmationTime)-meanConfirmation, 2)
			}
		}
		sdConfirmation = math.Sqrt(sdConfirmation / float64(counter))

		fmt.Println("Mean confirmation time is:", meanConfirmation)
		fmt.Println("SD confirmation time is:", sdConfirmation)
		//saveTangle(sim.tangle)
		//fmt.Println("\n\n")
		//fmt.Println("Tangle size: ", sim.param.TangleSize)

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
	if i > sim.param.TangleSize/10 {
		if i < sim.param.TangleSize*9/10 {
			if sim.generator.Float64() < sim.param.q {
				isAdverse = true
			}
		}
	}
	return isAdverse
}

func (sim *Sim) setMana(s float64) {
	var sum float64
	sum = 0
	for i := 0; i < int(sim.param.numberNodes); i++ {
		sum += (math.Pow(float64(i+1), -s))
	}
	for i := 0; i < int(sim.param.numberNodes); i++ {
		sim.mana[i] = (math.Pow(float64(i+1), -s)) / sum
	}
	return
}

func (sim *Sim) chooseNodeID() int {
	random := sim.generator.Float64()
	nodeID := 0
	cumsum := sim.mana[nodeID]
	for {
		if cumsum > random {
			break
		}
		nodeID++
		cumsum += sim.mana[nodeID]
	}
	return nodeID
}
