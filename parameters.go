package main

//Parameters of the simulation
type Parameters struct {
	K            int
	H            int
	Lambda       float64
	Alpha        float64
	TangleSize   int
	Seed         int64
	TSA          string
	tsa          TipSelector
	ConstantRate bool
	DataPath     string
	minCut       int
	maxCutrange  int
	maxCut       int
	nRun         int
	// - - - Analysis - - -
	VelocityEnabled bool
	AnPastCone      AnPastCone
	AnFocusRW       AnFocusRW
}

// Analysis Past Cone
type AnPastCone struct {
	Enabled    bool
	Resolution float64
	MaxT       float64
	MaxApp     int
}

// Analysis Focus RW
type AnFocusRW struct {
	Enabled bool
	maxiMT  int     // maximum PC length (measured in MT txs)
	murel   float64 // tx by adversary = murel * lambda
	nRWs    int     // number of RWs per data point
}
