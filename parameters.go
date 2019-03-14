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
	stillrecent  int
	CWMatrixLen  int
	// - - - Analysis - - -
	CountTipsEnabled     bool
	CWAnalysisEnabled    bool
	VelocityEnabled      bool
	ExitProbEnabled      bool
	ExitProbNparticle    int
	SpineEnabled         bool
	pOrphanEnabled       bool
	pOrphanLinFitEnabled bool
	AnPastCone           AnPastCone
	AnFocusRW            AnFocusRW
	// - - - Drawing - - -
	//drawTangleMode = 0: drawing disabled
	//drawTangleMode = 1: simple Tangle with/without highlighed path
	//drawTangleMode = 2: Ghost path, Ghost cone, Orphans + tips (TODO: clustering needs to be done manually)
	//drawTangleMode = 3: Tangle with tx visiting probability in red gradients
	//drawTangleMode = 4: Tangle with highlighted path of random walker transitioning to first approver
	//drawTangleMode = 5: Tangle with highlighted path of random walker transitioning to last approver
	//drawTangleMode = -1: 10 random walk and draws the Tangle at each step (for GIF or video only)
	drawTangleMode int
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
