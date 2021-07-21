package main

import (
	"github.com/iotaledger/hive.go/configuration"
	"math"
	"runtime"
	"strconv"
	"strings"
)

// variable initialization
func newParameters() Parameters {
	p := Parameters{}
	config := configuration.New()
	err := config.LoadFile("./parameters.yml")
	if err != nil {
		panic(err)
	}
	configuration.BindParameters(&p.ConfParameters, "simulation")
	configuration.UpdateBoundParameters(config)

	// - - - - setup some of the parameter values - - -

	if p.NParallelSims == -1 {
		// factor 2 is to use the physical cores, whereas NumCPU returns double the number due to hyper-threadingif
		p.NParallelSims = runtime.NumCPU()/2 - 1
	}

	lambdaForSize := int(math.Max(1, p.Lambda)) // make sure this value is at least 1
	p.TangleSize = p.TangleSizeNormalized * lambdaForSize
	p.MinCut = p.MinCutNormalized * lambdaForSize
	p.MaxCutRange = p.MaxCutRangeNormalized * lambdaForSize
	p.MaxCut = p.TangleSize - p.MaxCutRange
	p.StillRecent = p.StillRecent * lambdaForSize
	p.AcceptableNumberTips = p.AcceptableNumberTips * lambdaForSize

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

	if p.MaxK < p.K {
		p.MaxK = p.K
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
	return AnCone{
		enabled,
		resolution,
		maxT,
		int(maxApp),
	}
}

// Parameters defined in Yaml
type ConfParameters struct {
	NParallelSims         int
	K                     int
	H                     int
	D                     int
	Lambda                float64
	TangleSizeNormalized  int
	MinCutNormalized      int
	MaxCutRangeNormalized int
	Seed                  int64
	TSA                   string
	TSAAdversary          string
	SingleEdgeEnabled     bool
	ConstantRate          bool
	DataPath              string
	NRun                  int
	StillRecent           int
	// CWMatrixLen       int

	Q                       float64
	ResponseSpamTipsEnabled bool
	AcceptableNumberTips    int
	ResponseKIncrease       float64
	MaxK                    int
	// - - - Analysis - - -
	CountTipsEnabled bool
	// CWAnalysisEnabled bool

	POrphanEnabled       bool
	POrphanLinFitEnabled bool
	PastCone             []string
	FutureCone           []string
	DistSlicesEnabled    bool
	DistSlicesByTime     bool
	DistSlicesLength     float64
	DistSlicesResolution int
	AppStatsAllEnabled   bool
	// - - - Drawing - - -
	//DrawTangleMode = 0: drawing disabled
	//DrawTangleMode = 1: simple Tangle with/without highlighed path
	//DrawTangleMode = 2: Ghost path, Ghost cone, Orphans + tips (TODO: clustering needs to be done manually)
	//DrawTangleMode = 3: Tangle with tx visiting probability in red gradients
	//DrawTangleMode = 4: Tangle with highlighted path of random walker transitioning to first approver
	//DrawTangleMode = 5: Tangle with highlighted path of random walker transitioning to last approver
	//DrawTangleMode = -1: 10 random walk and draws the Tangle at each step (for GIF or video only)
	DrawTangleMode        int
	HorizontalOrientation bool
}

//Parameters to be used by simulation
type Parameters struct {
	ConfParameters
	tsa          TipSelector
	tsaAdversary TipSelectorAdversary
	AnPastCone   AnCone
	AnFutureCone AnCone
	TangleSize   int
	MinCut       int
	MaxCutRange  int
	MaxCut       int
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

func (p *Parameters) setVariable(variable float64) {
	if variable != -1 {
		/* comment out chosen variable */
		p.D = int(variable)
		// p.K = int(variable)
		// p.lambda = variable
	}
}
