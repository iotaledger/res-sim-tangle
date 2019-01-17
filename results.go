package main

import "fmt"

// Result is the data structure containing all the results of a simulation
type Result struct {
	tips     avgTips
	velocity velocityResult
}

type avgTips struct {
	tips float64
}

func (a avgTips) Join(b avgTips) avgTips {
	if a.tips == 0 {
		return b
	}
	var result avgTips
	result.tips = (a.tips + b.tips) / 2.
	return result
}

func (a avgTips) String() string {
	return fmt.Sprintln("E(L):", a.tips)
}

// StatInt defines a metric of ints
type StatInt struct {
	desc string
	v    map[int]int
}

// StatFloat64 defines a metric of float64s
type StatFloat64 struct {
	desc string
	v    map[float64]int
}

func joinMapInt(a, b map[int]int) map[int]int {
	if a == nil {
		return b
	}
	for k, v := range b {
		a[k] += v
	}
	return a
}

func joinMapFloat64(a, b map[float64]int) map[float64]int {
	if a == nil {
		return b
	}
	for k, v := range b {
		a[k] += v
	}
	return a
}

func joinMapStatInt(a, b StatInt) StatInt {
	a.desc = b.desc
	a.v = joinMapInt(a.v, b.v)
	return a
}

func joinMapStatFloat64(a, b StatFloat64) StatFloat64 {
	a.desc = b.desc
	a.v = joinMapFloat64(a.v, b.v)
	return a
}
