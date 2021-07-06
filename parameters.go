package main

import (
	"math"
	"strings"
)

// variable initialization
func newParameters(variable float64) Parameters {
	lambda := 100. // transaction per seconds
	// lambda := variable
	lambdaForSize := int(math.Max(1, lambda)) // make sure this value is at least 1
	hlarge := 1
	numberNodes := 2000
	p := Parameters{
		numberNodes: numberNodes,
		zipf:        variable,
		// factor 2 is to use the physical cores, whereas NumCPU returns double the number due to hyper-threading
		//nParallelSims: runtime.NumCPU()/2 - 1,
		nParallelSims: 1,
		// nRun:          int(math.Min(10000., 10000/lambda)),
		nRun:   1,
		Lambda: lambda,
		TSA:    "RURTS",
		// TSA:               "URTS",
		K:                 2,      // Num of tips to select
		Hsmall:            1,      // Delay for first type of tx, should be set to 1. Delay in seconds
		Hlarge:            hlarge, // Delay for second type of tx
		p:                 0,      //proportion of second type of tx
		D:                 10000,  // max age for RURTS
		Seed:              1,      //
		TangleSize:        1000 * lambdaForSize,
		minCut:            10 * hlarge * lambdaForSize, // cut data close to the genesis
		maxCutrange:       10 * hlarge * lambdaForSize, // cut data for the most recent txs, not applied for every analysis
		stillrecent:       20 * lambdaForSize,          // when is a tx considered recent, and when is it a candidate for left behind
		ConstantRate:      false,
		SingleEdgeEnabled: false, // true = SingleEdge model, false = MultiEdge model

		// - - - Attacks - - -
		q:            .5,              // proportion of adversary txs
		TSAAdversary: "SpamGenesis",   // spam tips linked to the genesis,
		adversaryID:  numberNodes - 1, // nodeID of adversary
		// - - - Response - - -
		responseSpamTipsEnabled: false,           // response dynamically to the tip spam attack
		acceptableNumberTips:    int(2 * lambda), // when we should start to increase K
		responseKIncrease:       3.,              // at which rate do we increase K
		maxK:                    20,              // maximum K used for protection, value will get replaced when K is larger
		// - - - Analysis section - - -
		CountTipsEnabled:  true,
		CTAnalysisEnabled: true,

		// - - - Drawing - - -
		//
		//drawTangleMode = 0: drawing disabled
		//drawTangleMode = 1: simple Tangle with/without highlighed path
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
	nParallelSims int
	numberNodes   int
	zipf          float64
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

	q                       float64
	attackType              string
	adversaryID             int
	responseSpamTipsEnabled bool
	acceptableNumberTips    int
	responseKIncrease       float64
	maxK                    int
	// - - - Analysis - - -
	CountTipsEnabled  bool
	CTAnalysisEnabled bool

	// - - - Drawing - - -
	//drawTangleMode = 0: drawing disabled
	//drawTangleMode = 1: simple Tangle with/without highlighed path
	drawTangleMode        int
	horizontalOrientation bool
}
