package main

import (
	"math"
	"strings"
)

// variable initialization
func newParameters(tsa string, lambda, alpha float64) Parameters {
	lambdaForSize := int(math.Max(1, lambda)) // make sure this value is at least 1
	p := Parameters{
		//K:          2,
		//H:          1,
		Lambda:       lambda,
		Alpha:        alpha,
		K:            2, // Num of tips to select
		H:            1, // ???
		Seed:         1, // ???
		TangleSize:   200 * lambdaForSize,
		CWMatrixLen:  40 * lambdaForSize, // reduce CWMatrix to this len
		minCut:       51 * lambdaForSize, // cut data close to the genesis
		maxCutrange:  52 * lambdaForSize, // cut data for the most recent txs, not applied for every analysis
		stillrecent:  2 * lambdaForSize,  // when is a tx considered recent, and when is it a candidate for left behind
		ConstantRate: false,
		nRun:         int(math.Min(100., 1000/lambda)),
		// nRun:              1,
		TSA:               tsa,
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
		DistSlicesEnabled:    false, // calculate the distances of slices
		DistSlicesByTime:     false, // true = tx time slices, false= tx ID slices
		// DistSlicesLength:     100 / lambda, //length of Slices
		DistSlicesLength:     1,    //length of Slices
		DistSlicesResolution: 100,  // Number of intervals per distance '1', higher number = higher resolution
		AppStatsRWEnabled:    true, // Approver Stats along the RW
		AppStatsRW_NumRWs:    100,
		AppStatsAllEnabled:   true, // Approver stats for all txs
		//{Enabled, Resolution, MaxT, MaxApp}
		AnPastCone: AnPastCone{false, 5, 40, 5},
		//{Enabled, maxiMT, murel, nRW}
		AnFocusRW: AnFocusRW{false, 0.2, 30},
	}

	// - - - - setup some of the parameter values - - -
	switch strings.ToUpper(p.TSA) {
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
	K                 int
	H                 int
	Lambda            float64
	Alpha             float64
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
	AnFocusRW            AnFocusRW
	DistSlicesEnabled    bool
	DistSlicesByTime     bool
	DistSlicesLength     float64
	DistSlicesResolution int
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

// AnFocusRW Analysis Focus RW
type AnFocusRW struct {
	Enabled bool
	murel   float64 // tx by adversary = murel * lambda
	nRWs    int     // number of RWs per data point
}
