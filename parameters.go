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
	DistSlicesLength     float64
	DistSlicesResolution int
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
