package main

// Result is the data structure containing all the results of a simulation
type Result struct {
	tips     tipsResult
	velocity velocityResult
	PastCone PastConeResult
	FocusRW  FocusRWResult
	entropy  entropyResult
	op       pOrphanResult
	cw       cwResult
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
