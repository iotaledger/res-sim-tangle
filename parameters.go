package main

import (
	"github.com/iotaledger/hive.go/configuration"
	"math"
	"runtime"
	"strconv"
	"strings"
)

// variable initialization
func newParameters(age float64) Parameters {
	p := Parameters{}
	config := configuration.New()
	err := config.LoadFile("./parameters.yml")
	if err != nil {
		panic(err)
	}
	configuration.BindParameters(p)
	configuration.UpdateBoundParameters(config)

	// - - - - setup some of the parameter values - - -

	if p.nParallelSims == -1 {
		// factor 2 is to use the physical cores, whereas NumCPU returns double the number due to hyper-threadingif
		p.nParallelSims = runtime.NumCPU()/2 - 1
	}

	if age != -1 {
		p.D = int(age)
	}
	lambdaForSize := int(math.Max(1, p.Lambda)) // make sure this value is at least 1
	p.TangleSize = p.TangleSize * lambdaForSize
	p.minCut = p.minCut * lambdaForSize
	p.maxCutrange = p.maxCutrange * lambdaForSize
	p.maxCut = p.TangleSize - p.maxCutrange
	p.stillrecent = p.stillrecent * lambdaForSize
	p.acceptableNumberTips = p.acceptableNumberTips * lambdaForSize

	p.AnPastCone = coneFromParameters(p.PastCone)
	p.AnFutureCone = coneFromParameters(p.FutureCone)

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

	if p.maxK < p.K {
		p.maxK = p.K
	}

	createDirIfNotExist("data")
	createDirIfNotExist("graph")

	return p
}

func coneFromParameters(cone []string) AnCone {
	enabled, err := strconv.ParseBool(cone[0])
	if err != nil {
		panic(err)
	}
	resolution, err := strconv.ParseFloat(cone[1], 64)
	if err != nil {
		panic(err)
	}
	maxT, err := strconv.ParseFloat(cone[2], 64)
	if err != nil {
		panic(err)
	}
	maxApp, err := strconv.ParseInt(cone[3], 10, 32)
	if err != nil {
		panic(err)
	}
	return AnCone{enabled,
		resolution,
		maxT,
		int(maxApp),
	}
}

// Parameters define Parameters types
type Parameters struct {
	nParallelSims     int
	K                 int
	H                 int
	D                 int
	Lambda            float64
	TangleSize        int
	minCut            int
	maxCutrange       int
	maxCut            int
	Seed              int64
	TSA               string
	tsa               TipSelector
	TSAAdversary      string
	tsaAdversary      TipSelectorAdversary
	SingleEdgeEnabled bool
	ConstantRate      bool
	DataPath          string
	nRun              int
	stillrecent       int
	// CWMatrixLen       int

	q                       float64
	attackType              string
	responseSpamTipsEnabled bool
	acceptableNumberTips    int
	responseKIncrease       float64
	maxK                    int
	// - - - Analysis - - -
	CountTipsEnabled bool
	// CWAnalysisEnabled bool

	pOrphanEnabled       bool
	pOrphanLinFitEnabled bool
	PastCone             []string
	AnPastCone           AnCone
	FutureCone           []string
	AnFutureCone         AnCone
	DistSlicesEnabled    bool
	DistSlicesByTime     bool
	DistSlicesLength     float64
	DistSlicesResolution int
	AppStatsAllEnabled   bool
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

// AnCone Analysis results
type AnCone struct {
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
