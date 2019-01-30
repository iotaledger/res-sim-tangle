package main

import (
	"time"
)

func (sim *Sim) directOrder(tip Tx) (int64, int) {
	set := make(map[int]bool)
	start := time.Now()
	dfs(tip, set, sim)
	end := time.Since(start)
	return end.Nanoseconds(), len(set)
}

func (sim *Sim) reverseOrder(tip Tx) (int64, int) {
	set := make(map[int]bool)
	start := time.Now()
	dfsPastToPresent(tip, set, sim)
	end := time.Since(start)
	return end.Nanoseconds(), len(set)
}

func dfs(t Tx, visited map[int]bool, sim *Sim) {
	if t.id > 0 {
		for _, id := range t.ref {
			if !visited[id] {
				visited[id] = true
				dfs(sim.tangle[id], visited, sim)
			}
		}
	}
}

func dfsPastToPresent(t Tx, visited map[int]bool, sim *Sim) {
	if len(sim.approvers[t.id]) > 0 {
		for _, id := range sim.approvers[t.id] {
			if !visited[id] {
				visited[id] = true
				dfsPastToPresent(sim.tangle[id], visited, sim)
			}
		}
	}
}
