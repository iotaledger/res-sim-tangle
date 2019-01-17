package main

import (
	"fmt"
	"time"
)

//Benchmark contains performance time measures
type Benchmark map[string]time.Duration

func runningtime(s string) (string, time.Time) {
	//fmt.Println("Start:	", s)
	return s, time.Now()
}

func (performance Benchmark) track(s string, startTime time.Time) {
	endTime := time.Now()
	//fmt.Println("End:	", s, "took", endTime.Sub(startTime))
	performance[s] += endTime.Sub(startTime)
}

func printPerformance(performance Benchmark) {
	fmt.Println()
	for metric := range performance {
		fmt.Println(metric, ":", performance[metric])
	}
}
