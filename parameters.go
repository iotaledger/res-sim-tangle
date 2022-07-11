package main

import (
	"math"
	"strings"
)

// variable initialization
func newParameters(variable float64, simStep int) Parameters {
	lambda := 20.
	// lambda := variable
	lambdaForSize := int(math.Max(1, lambda)) // make sure this value is at least 1
	hlarge := 1
	p := Parameters{
		Variable: variable,
		SimStep:  simStep, // at which # of variable of the simulation we are

		// factor 2 is to use the physical cores, whereas NumCPU returns double the number due to hyper-threading
		// nParallelSims: runtime.NumCPU()/2 - 1,
		nParallelSims: 1,
		// nRun:          int(math.Min(10000., 10000/lambda)),
		nRun:   100,
		Lambda: lambda,
		TSA:    "RURTS",
		// TSA:               "URTS",
		K:          2,             // Num of tips to select
		Hsmall:     1,             // Delay for first type of tx,
		Hlarge:     hlarge,        // Delay for second type of tx
		p:          0.,            //proportion of second type of tx
		D:          int(variable), // max age for RURTS
		Seed:       1,             //
		TangleSize: (10*hlarge + 500) * lambdaForSize,
		// CWMatrixLen:       300 * lambdaForSize, // reduce CWMatrix to this len
		minCut:            20 * hlarge * lambdaForSize, // cut data close to the genesis
		maxCutrange:       20 * hlarge * lambdaForSize, // cut data for the most recent txs, not applied for every analysis
		stillrecent:       2 * lambdaForSize,           // when is a tx considered recent, and when is it a candidate for left behind
		ConstantRate:      false,
		SingleEdgeEnabled: false, // true = SingleEdge model, false = MultiEdge model

		// - - - Attacks - - -
		q:                0.25,          // proportion of adversary txs
		qPartiallyActive: false,         // attack only active between [1/3,2/3] of the Tangle
		TSAAdversary:     "SpamGenesis", // spam tips linked to the genesis,
		// - - - Response - - -
		responseSpamTipsEnabled: false,           // response dynamically to the tip spam attack
		acceptableNumberTips:    int(2 * lambda), // when we should start to increase K
		responseKIncrease:       3.,              // at which rate do we increase K
		maxK:                    20,              // maximum K used for protection, value will get replaced when K is larger
		// - - - Analysis section - - -
		CountTipsEnabled: true, // including orphan tips
		// CWAnalysisEnabled:    false,
		recordOrphansForEachSim: true,  // save for each Tangle the orphan number
		pOrphanLinFitEnabled:    false, // also apply linear fit, numerically expensive
		// measure distance of slices compared to the expected distribution
		DistSlicesEnabled:    false,
		DistSlicesByTime:     false, // true = tx time slices, false= tx ID slices
		DistSlicesLength:     1,     //length of Slices
		DistSlicesResolution: 100,   // Number of intervals per distance '1', higher number = higher resolution
		// measure app stats for all txs
		AppStatsAllEnabled: false, // Approver stats for all txs
		// AnPastCone Analysis
		AnPastCone: AnPastCone{false, 5, 40, 5}, //{Enabled, Resolution, MaxT, MaxApp}
		// AnFutureCone Analysis
		AnFutureCone: AnFutureCone{false, 5, 40, 5}, //{Enabled, Resolution, MaxT, MaxApp}

		// - - - Drawing - - -
		//
		//drawTangleMode = 0: drawing disabled
		//drawTangleMode = 1: simple Tangle with/without highlighed path
		//drawTangleMode = 2: Ghost path, Ghost cone, Orphans + tips (TODO: clustering needs to be done manually)
		//drawTangleMode = 3: Tangle with tx visiting probability in red gradients
		//drawTangleMode = 4: Tangle with highlighted path of random walker transitioning to first approver
		//drawTangleMode = 5: Tangle with highlighted path of random walker transitioning to last approver
		//drawTangleMode = -1: 10 random walk and draws the Tangle at each step (for GIF or video only)
		drawTangleMode:        0,
		horizontalOrientation: true,
	}

	// - - - - setup some of the parameter values - - -
	p.TSA = strings.ToUpper(p.TSA) // make sure string is upper case
	switch p.TSA {
	case "HPS":
		p.tsa = HPS{}
	case "RURTS":
		p.tsa = RURTS{}
	case "URTS":
		p.tsa = URTS{}
	default:
		p.TSA = "URTS"
		p.tsa = URTS{}
	}

	switch p.TSAAdversary {
	case "spamGenesis":
		p.tsaAdversary = SpamGenesis{}
	default:
		p.TSAAdversary = "spamGenesis"
		p.tsaAdversary = SpamGenesis{}
	}

	p.maxCut = p.TangleSize - p.maxCutrange

	if p.maxK < p.K {
		p.maxK = p.K
	}

	createDirIfNotExist("data")
	createDirIfNotExist("graph")

	return p
}

//define Parameters types
type Parameters struct {
	Variable      float64
	SimStep       int
	nParallelSims int
	K             int
	Hsmall        int
	Hlarge        int
	p             float64
	D             int
	Lambda        float64
	// tsaType           string
	TangleSize        int
	Seed              int64
	TSA               string
	tsa               TipSelector
	TSAAdversary      string
	tsaAdversary      TipSelectorAdversary
	SingleEdgeEnabled bool
	ConstantRate      bool
	DataPath          string
	minCut            int
	maxCutrange       int
	maxCut            int
	nRun              int
	stillrecent       int
	// CWMatrixLen       int

	q                       float64
	qPartiallyActive        bool
	attackType              string
	responseSpamTipsEnabled bool
	acceptableNumberTips    int
	responseKIncrease       float64
	maxK                    int
	// - - - Analysis - - -
	CountTipsEnabled bool
	// CWAnalysisEnabled bool

	recordOrphansForEachSim bool
	pOrphanLinFitEnabled    bool
	AnPastCone              AnPastCone
	AnFutureCone            AnFutureCone
	DistSlicesEnabled       bool
	DistSlicesByTime        bool
	DistSlicesLength        float64
	DistSlicesResolution    int
	AppStatsAllEnabled      bool
	// - - - Drawing - - -
	//drawTangleMode = 0: drawing disabled
	//drawTangleMode = 1: simple Tangle with/without highlighed path
	//drawTangleMode = 2: Ghost path, Ghost cone, Orphans + tips (TODO: clustering needs to be done manually)
	//drawTangleMode = 3: Tangle with tx visiting probability in red gradients
	//drawTangleMode = 4: Tangle with highlighted path of random walker transitioning to first approver
	//drawTangleMode = 5: Tangle with highlighted path of random walker transitioning to last approver
	//drawTangleMode = -1: 10 random walk and draws the Tangle at each step (for GIF or video only)
	drawTangleMode        int
	horizontalOrientation bool
}

// AnPastCone Analysis Past Cone
type AnPastCone struct {
	Enabled    bool
	Resolution float64
	MaxT       float64
	MaxApp     int
}

// AnFutureCone Analysis
type AnFutureCone struct {
	Enabled    bool
	Resolution float64
	MaxT       float64
	MaxApp     int
}

// AnFocusRW Analysis Focus RW
type AnFocusRW struct {
	Enabled bool
	murel   float64 // tx by adversary = murel * lambda
	nRWs    int     // number of RWs per data point
}
