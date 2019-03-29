package main

import (
	"math"
	"runtime"
	"strings"
)

// variable initialization
func newParameters(lambda, alpha float64) Parameters {
	lambdaForSize := int(math.Max(1, lambda)) // make sure this value is at least 1
	p := Parameters{

		// factor 2 is to use the physical cores, whereas NumCPU returns double the number due to hyper-threading
		nParallelSims: runtime.NumCPU()/2 - 1,
		// nParallelSims: runtime.NumCPU(),
		// nParallelSims: 1,
		// nRun: int(math.Min(1000., 1000/lambda)),
		nRun:   1,
		Lambda: lambda,
		Alpha:  alpha,
		TSA:    "RW",
		// TSA:               "URTS",
		K:                 2, // Num of tips to select
		H:                 1, // ???
		Seed:              1, // ???
		TangleSize:        300 * lambdaForSize,
		CWMatrixLen:       300 * lambdaForSize, // reduce CWMatrix to this len
		minCut:            51 * lambdaForSize,  // cut data close to the genesis
		maxCutrange:       52 * lambdaForSize,  // cut data for the most recent txs, not applied for every analysis
		stillrecent:       2 * lambdaForSize,   // when is a tx considered recent, and when is it a candidate for left behind
		ConstantRate:      false,
		SingleEdgeEnabled: true, // true = SingleEdge model, false = MultiEdge model

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
		// measure distance of slices compared to the expected distribution
		DistSlicesEnabled:    false,
		DistSlicesByTime:     false, // true = tx time slices, false= tx ID slices
		DistSlicesLength:     1,     //length of Slices
		DistSlicesResolution: 100,   // Number of intervals per distance '1', higher number = higher resolution
		// measure distance of RWs compared to the expected distribution
		DistRWsEnabled:      false,
		DistRWsSampleLength: 20,                 // Length of considered RWs
		DistRWsSampleRWNum:  lambdaForSize * 10, // Number of sample RWs per Tangle
		DistRWsResolution:   100,                // Number of intervals per distance '1', higher number = higher resolution
		// measure Approver stats during RW
		AppStatsRWEnabled: false, // Approver Stats along the RW
		AppStatsRW_NumRWs: max2Int(100*lambdaForSize, 100),
		// measure app stats for all txs
		AppStatsAllEnabled: false, // Approver stats for all txs
		// AnPastCone Analysis
		AnPastCone: AnPastCone{false, 5, 40, 5}, //{Enabled, Resolution, MaxT, MaxApp}
		// AnFutureCone Analysis
		AnFutureCone: AnFutureCone{false, 5, 40, 5}, //{Enabled, Resolution, MaxT, MaxApp}
		// AnFocusRW Analysis Focus RW
		AnFocusRW: AnFocusRW{false, 0.2, 30}, //{Enabled, maxiMT, murel, nRW}

	}

	// - - - - setup some of the parameter values - - -
	p.TSA = strings.ToUpper(p.TSA) // make sure string is upper case
	switch p.TSA {
	case "URTS":
		p.tsa = URTS{}
	case "RW":
		if p.Alpha == 0 {
			p.tsa = URW{}
		} else {
			p.tsa = BRW{}
		}
	default:
		p.TSA = "URTS"
		p.tsa = URTS{}
	}

	if p.TSA == "URTS" || p.Alpha == 0 {
		p.SpineEnabled = false
	}

	p.maxCut = p.TangleSize - p.maxCutrange

	createDirIfNotExist("data")

	return p
}

//define Parameters types
type Parameters struct {
	nParallelSims int
	K             int
	H             int
	Lambda        float64
	Alpha         float64
	// tsaType           string
	TangleSize        int
	Seed              int64
	TSA               string
	tsa               TipSelector
	SingleEdgeEnabled bool
	ConstantRate      bool
	DataPath          string
	minCut            int
	maxCutrange       int
	maxCut            int
	nRun              int
	stillrecent       int
	CWMatrixLen       int
	// - - - Analysis - - -
	CountTipsEnabled  bool
	CWAnalysisEnabled bool
	VelocityEnabled   bool
	ExitProbEnabled   bool
	ExitProbNparticle int
	ExitProb2NHisto   int

	SpineEnabled         bool
	pOrphanEnabled       bool
	pOrphanLinFitEnabled bool
	AnPastCone           AnPastCone
	AnFutureCone         AnFutureCone
	AnFocusRW            AnFocusRW
	DistSlicesEnabled    bool
	DistSlicesByTime     bool
	DistSlicesLength     float64
	DistSlicesResolution int
	DistRWsEnabled       bool
	DistRWsSampleLength  int
	DistRWsSampleRWNum   int
	DistRWsResolution    int
	AppStatsRWEnabled    bool
	AppStatsRW_NumRWs    int
	AppStatsAllEnabled   bool
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
