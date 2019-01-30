package main

//Parameters of the simulation
type Parameters struct {
	K                      int
	H                      int
	Lambda                 float64
	Alpha                  float64
	TangleSize             int
	Seed                   int64
	TSA                    string
	tsa                    TipSelector
	ConstantRate           bool
	DataPath               string
	minCut                 int
	maxCut                 int
	nRun                   int
	VelocityEnabled        bool
	ReusableAddressEnabled bool
}
