package main

import (
	"fmt"
	"os"
	"sort"

	"gonum.org/v1/gonum/stat"
)

// Result is the data structure containing all the results of a simulation
type Result struct {
	tips       tipsResult
	velocity   velocityResult
	PastCone   PastConeResult
	FocusRW    FocusRWResult
	exitProb   exitProbResult
	op         pOrphanResult
	cw         cwResult
	avgtips    avgTips
	DistSlices DistSlicesResult
	AppStatsRW AppStatsRWResult
}

type avgTips struct {
	val []float64
}

func (a avgTips) Join(b avgTips) avgTips {
	if a.val == nil {
		return b
	}
	var result avgTips
	result.val = append(a.val, b.val...)
	return result
}

func (a avgTips) String() string {
	return fmt.Sprintln("E(L):", stat.Mean(a.val, nil))
}

// MetricIntInt defines a metric of ints
type MetricIntInt struct {
	desc string
	v    map[int]int
}

// MetricIntFloat64 defines a map from int to floats
type MetricIntFloat64 struct {
	desc string
	v    map[int]float64
}

// MetricFloat64Int defines a metric of float64s
type MetricFloat64Int struct {
	desc string
	v    map[float64]int
}

// MetricFloat64Float64 defines a metric of float64s
type MetricFloat64Float64 struct {
	desc string
	v    map[float64]float64
}

func joinMapIntInt(a, b map[int]int) map[int]int {
	if a == nil {
		return b
	}
	for k, v := range b {
		a[k] += v
	}
	return a
}

// func avgMapIntInt(a, b map[int]int) map[int]float64 {
// 	r := make(map[int]float64)

// 	//copy b to r
// 	for k, v := range b {
// 		r[k] = float64(v)
// 	}

// 	if a == nil {
// 		return r
// 	}

// 	//fill b with same keys as a
// 	for k := range a {
// 		if _, ok := b[k]; !ok {
// 			b[k] = 0
// 		}
// 	}
// 	//compute avg between a and b
// 	for k, v := range b {
// 		a[k] += v
// 		a[k] /= 2
// 	}
// 	return a
// }

func joinMapFloat64Int(a, b map[float64]int) map[float64]int {
	if a == nil {
		return b
	}
	for k, v := range b {
		a[k] += v
	}
	return a
}

func joinMapIntFloat64(a, b map[int]float64) map[int]float64 {
	if a == nil {
		return b
	}
	for k, v := range b {
		a[k] += v
	}
	return a
}

func joinMapFloat64Float64(a, b map[float64]float64) map[float64]float64 {
	if a == nil {
		return b
	}
	for k, v := range b {
		a[k] += v
	}
	return a
}

func (s MetricFloat64Float64) getkeys() []float64 {
	// var datapoints int
	var keys []float64
	for k := range s.v {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	return keys
}

func getkeysMapFloat64Float64(s map[float64]float64) []float64 {
	// var datapoints int
	var keys []float64
	for k, _ := range s {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	return keys
}

// this can be simplified also in the application of it by (a *MetricIntInt) joinMapMetricIntInt (b MetricIntInt) {}
func joinMapMetricIntInt(a, b MetricIntInt) MetricIntInt {
	a.desc = b.desc
	a.v = joinMapIntInt(a.v, b.v)
	return a
}

func joinMapMetricIntFloat64(a, b MetricIntFloat64) MetricIntFloat64 {
	a.desc = b.desc
	a.v = joinMapIntFloat64(a.v, b.v)
	return a
}

func joinMapMetricFloat64Int(a, b MetricFloat64Int) MetricFloat64Int {
	a.desc = b.desc
	a.v = joinMapFloat64Int(a.v, b.v)
	return a
}

func joinMapMetricFloat64Float64(a, b MetricFloat64Float64) MetricFloat64Float64 {
	a.desc = b.desc
	a.v = joinMapFloat64Float64(a.v, b.v)
	return a
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//  saves an array of MetricFloat64Float64 to a file
func SaveArrayMetricFloat64Float64(p Parameters, filename string, r []MetricFloat64Float64) error {
	for _, v := range r {
		v.SaveMetricFloat64Float64(p, filename, true)
	}
	return nil
}
func (m MetricFloat64Float64) SaveMetricFloat64Float64(p Parameters, filename string, normalized bool) error {
	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.4f", p.Alpha)

	if len(m.getkeys()) > 0 {
		f, err := os.Create("data/" + filename + "__" + m.desc + "__TSA=" + p.TSA + "_lambda=" + lambdaStr + "_alpha=" + alphaStr + ".txt")
		if err != nil {
			fmt.Printf("error creating file: %v", err)
			return err
		}
		defer f.Close()
		// for i, k := range x {
		for _, k := range m.getkeys() {
			_, err = f.WriteString(fmt.Sprintf("%f\t%f\n", k, m.v[k])) // writing...
			// _, err = f.WriteString(fmt.Sprintf("%f\t%f\n", k, weigths[i]/float64(datapoints)*norm)) // writing...
			if err != nil {
				fmt.Printf("error writing string: %v", err)
			}
		}
	}

	return nil
}
