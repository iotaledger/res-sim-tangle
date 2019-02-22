package main

import (
	"fmt"
	"math"
	"os"
	"sort"

	"gonum.org/v1/gonum/stat"
)

//Velocity result of simulation
type velocityResult struct {
	vID        []StatInt //???creates a map[int]int with a keyword
	vTime      []StatFloat64
	dApprovers []StatInt
	vCW        []StatInt
	vCWfirst   []StatInt
}

//??? use string to create empty value maps to vID, vTime, dApprovers
func newVelocityResult(veloMetrics []string) *velocityResult {
	// variables initialization for velocity
	var result velocityResult
	for _, metric := range veloMetrics {
		result.vID = append(result.vID, StatInt{metric, make(map[int]int)})
		result.vTime = append(result.vTime, StatFloat64{metric, make(map[float64]int)})
		if metric != "back" {
			result.vCW = append(result.vCW, StatInt{metric, make(map[int]int)})
			if metric != "only-1" {
				result.vCWfirst = append(result.vCWfirst, StatInt{metric, make(map[int]int)})
			}
			if metric == "rw" || metric == "all" {
				result.dApprovers = append(result.dApprovers, StatInt{metric, make(map[int]int)})
			}
		}
	}
	return &result
}

func (sim *Sim) runVelocityStat(result *velocityResult) {
	if sim.param.TSA != "RW" {
		sim.velocityURTS(result.vID[0].v, result.vTime[0].v, result.dApprovers[0].v, result.vCW[0].v, result.vCWfirst[0].v)
		sim.velocityAll(result.vID[1].v, result.vTime[1].v, result.dApprovers[1].v, result.vCW[1].v, result.vCWfirst[1].v)
		sim.velocityOfIndex(result.vID[2].v, result.vTime[2].v, result.vCW[2].v, result.vCWfirst[2].v, 1)
		sim.velocityOfIndex(result.vID[3].v, result.vTime[3].v, result.vCW[3].v, result.vCWfirst[3].v, -1)
		sim.velocityOfIndex(result.vID[4].v, result.vTime[4].v, result.vCW[4].v, result.vCWfirst[4].v, 2)
		sim.velocityOfIndex(result.vID[5].v, result.vTime[5].v, result.vCW[5].v, result.vCWfirst[5].v, 3)
		sim.velocityOfIndex(result.vID[6].v, result.vTime[6].v, result.vCW[6].v, result.vCWfirst[6].v, 4)
		sim.velocityOfOnlyIndex(result.vID[7].v, result.vTime[7].v, result.vCW[7].v, 1)
		sim.velocityBackURTS(result.vID[8].v, result.vTime[8].v)
	} else {
		sim.velocityParticleRW(result.vID[0].v, result.vTime[0].v, result.dApprovers[0].v, result.vCW[0].v, result.vCWfirst[0].v, 100000)
		sim.velocityAll(result.vID[1].v, result.vTime[1].v, result.dApprovers[1].v, result.vCW[1].v, result.vCWfirst[1].v)
		sim.velocityOfIndexRW(result.vID[2].v, result.vTime[2].v, result.vCW[2].v, result.vCWfirst[2].v, 1, 100000)
		sim.velocityOfIndexRW(result.vID[3].v, result.vTime[3].v, result.vCW[3].v, result.vCWfirst[3].v, -1, 100000)
		//sim.velocityOfIndexRW(result.vID[4].v, result.vTime[4].v, result.vCW[4].v, result.vCWfirst[4].v, 2, 100000)
		//sim.velocityOfIndexRW(result.vID[5].v, result.vTime[5].v, result.vCW[5].v, result.vCWfirst[5].v, 3, 100000)
		//sim.velocityOfIndexRW(result.vID[6].v, result.vTime[6].v, result.vCW[6].v, result.vCWfirst[6].v, 4, 100000)
		//sim.velocityOfOnlyIndex(result.vID[7].v, result.vTime[7].v, result.vCW[7].v, 1)
		sim.velocityParticleBackRW(result.vID[4].v, result.vTime[4].v, 100000)
	}

}

func (sim Sim) velocityURTS(v map[int]int, t map[float64]int, d map[int]int, w, wFirst map[int]int) {
	for i := sim.param.minCut; i < sim.param.maxCut; i++ {
		if len(sim.approvers[i]) > 0 {
			l := sim.generator.Intn(len(sim.approvers[i]))
			delta := sim.approvers[i][l] - i
			deltaTime := math.Round((sim.tangle[sim.approvers[i][l]].time-sim.tangle[i].time)*100) / 100
			v[delta]++
			d[l+1]++
			t[deltaTime]++
			deltaCW := sim.tangle[i].cw - sim.tangle[sim.approvers[i][l]].cw
			w[deltaCW]++
			if len(sim.approvers[i]) > 1 {
				wFirst[deltaCW]++
			}
			// if float64(delta)/sim.param.Lambda != deltaTime {
			// 	fmt.Println(sim.approvers[i][l], "-", i, float64(delta)/sim.param.Lambda, "|", deltaTime, sim.tangle[sim.approvers[i][l]].time, "-", sim.tangle[i].time)
			// }
		}
	}
	//fmt.Println(t)
}

func (sim Sim) velocityBackURTS(v map[int]int, t map[float64]int) {
	for i := sim.param.maxCut; i > sim.param.minCut; i-- {

		l := sim.generator.Intn(len(sim.tangle[i].ref))
		delta := sim.tangle[i].ref[l] - i
		deltaTime := math.Round((sim.tangle[sim.tangle[i].ref[l]].time-sim.tangle[i].time)*100) / 100
		v[delta]++
		t[deltaTime]++

	}
}

func (sim *Sim) velocityParticleRW(v map[int]int, t map[float64]int, d map[int]int, w, wFirst map[int]int, nParticles int) {
	for i := 0; i < nParticles; i++ {
		//prev := sim.tangle[0]
		start := sim.generator.Intn(sim.param.minCut)
		prev := sim.tangle[start]
		var tsa RandomWalker
		if sim.param.Alpha != 0 {
			tsa = BRW{}
		} else {
			tsa = URW{}
		}

		for current, currentIdx := tsa.RandomWalk(prev, sim); len(sim.approvers[current.id]) > 0; current, currentIdx = tsa.RandomWalk(current, sim) {
			if current.id > sim.param.minCut && current.id < sim.param.maxCut {
				delta := current.id - prev.id
				v[delta]++
				d[currentIdx+1]++
				deltaTime := math.Round((current.time-prev.time)*100) / 100
				t[deltaTime]++
				deltaCW := prev.cw - current.cw
				w[deltaCW]++
				if len(sim.approvers[prev.id]) > 1 {
					wFirst[deltaCW]++
				}
			}
			prev = current
		}
	}
}

func (sim *Sim) velocityParticleBackRW(v map[int]int, t map[float64]int, nParticles int) {
	for i := 0; i < nParticles; i++ {
		start := sim.generator.Intn(sim.param.minCut) + sim.param.maxCut
		prev := sim.tangle[start]
		var tsa RandomWalker
		if sim.param.Alpha != 0 {
			tsa = BRW{}
		} else {
			tsa = URW{}
		}

		for current := tsa.RandomWalkBack(prev, sim); current.id > sim.param.minCut; current = tsa.RandomWalkBack(current, sim) {
			if current.id > sim.param.minCut && current.id < sim.param.maxCut {
				delta := prev.id - current.id
				v[delta]++
				deltaTime := math.Round((prev.time-current.time)*100) / 100
				t[deltaTime]++
			}
			prev = current
		}
	}
}

func (sim Sim) velocityOfIndex(v map[int]int, t map[float64]int, w, wFirst map[int]int, index int) {
	for i := sim.param.minCut; i < sim.param.maxCut; i++ {
		if index > 0 && len(sim.approvers[i]) > index-1 {
			delta := sim.approvers[i][index-1] - i
			v[delta]++
			deltaTime := math.Round((sim.tangle[sim.approvers[i][index-1]].time-sim.tangle[i].time)*100) / 100
			t[deltaTime]++
			deltaCW := sim.tangle[i].cw - sim.tangle[sim.approvers[i][index-1]].cw
			w[deltaCW]++
			if len(sim.approvers[i]) > 1 {
				wFirst[deltaCW]++
			}
		} else if index < 0 && len(sim.approvers[i]) > 1 {
			delta := sim.approvers[i][len(sim.approvers[i])-1] - i
			v[delta]++
			deltaTime := math.Round((sim.tangle[sim.approvers[i][len(sim.approvers[i])-1]].time-sim.tangle[i].time)*100) / 100
			t[deltaTime]++
			deltaCW := sim.tangle[i].cw - sim.tangle[sim.approvers[i][len(sim.approvers[i])-1]].cw
			w[deltaCW]++
			if len(sim.approvers[i]) > 1 {
				wFirst[deltaCW]++
			}
		}
	}
}

func (sim *Sim) velocityOfIndexRW(v map[int]int, t map[float64]int, w, wFirst map[int]int, index int, nParticles int) {

	for i := 0; i < nParticles; i++ {
		start := sim.generator.Intn(sim.param.minCut)

		for current := sim.tangle[start]; len(sim.approvers[current.id]) > 0 && current.id < sim.param.maxCut; {
			if index > 0 && len(sim.approvers[current.id]) > index-1 {
				delta := sim.approvers[current.id][index-1] - current.id
				deltaTime := math.Round((sim.tangle[sim.approvers[current.id][index-1]].time-sim.tangle[current.id].time)*100) / 100
				deltaCW := sim.tangle[current.id].cw - sim.tangle[sim.approvers[current.id][index-1]].cw
				if current.id > sim.param.minCut {
					v[delta]++
					t[deltaTime]++
					w[deltaCW]++
					if len(sim.approvers[current.id]) > 1 {
						wFirst[deltaCW]++
					}
				}
				current = sim.tangle[sim.approvers[current.id][index-1]]
			} else if index < 0 && len(sim.approvers[current.id]) > 1 {
				delta := sim.approvers[current.id][len(sim.approvers[current.id])-1] - current.id
				deltaTime := math.Round((sim.tangle[sim.approvers[current.id][len(sim.approvers[current.id])-1]].time-sim.tangle[current.id].time)*100) / 100
				deltaCW := sim.tangle[current.id].cw - sim.tangle[sim.approvers[current.id][len(sim.approvers[current.id])-1]].cw
				if current.id > sim.param.minCut {
					v[delta]++
					t[deltaTime]++
					w[deltaCW]++
					if len(sim.approvers[current.id]) > 1 {
						wFirst[deltaCW]++
					}
				}
				current = sim.tangle[sim.approvers[current.id][len(sim.approvers[current.id])-1]]
			} else {
				break
			}
		}
	}
}

func (sim Sim) velocityOfOnlyIndex(v map[int]int, t map[float64]int, w map[int]int, index int) {
	for i := sim.param.minCut; i < sim.param.maxCut; i++ {
		if index > 0 && len(sim.approvers[i]) == index {
			delta := sim.approvers[i][index-1] - i
			v[delta]++
			deltaTime := math.Round((sim.tangle[sim.approvers[i][index-1]].time-sim.tangle[i].time)*100) / 100
			t[deltaTime]++
			deltaCW := sim.tangle[i].cw - sim.tangle[sim.approvers[i][index-1]].cw
			w[deltaCW]++
		}
	}
}

func (sim Sim) velocityAll(v map[int]int, t map[float64]int, d map[int]int, w, wFirst map[int]int) {
	for i := sim.param.minCut; i < sim.param.maxCut; i++ {

		d[len(sim.approvers[i])]++
		for _, a := range sim.approvers[i] {
			delta := a - i
			v[delta]++
			deltaTime := math.Round((sim.tangle[a].time-sim.tangle[i].time)*100) / 100
			t[deltaTime]++
			deltaCW := sim.tangle[i].cw - sim.tangle[a].cw
			w[deltaCW]++
			if len(sim.approvers[i]) > 1 {
				wFirst[deltaCW]++
			}
		}
	}
}

func (p Parameters) printStatVelo(v map[int]int, target string) int {
	var keys []int
	var datapoints int
	for k := range v {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// calculate statistics
	var weigths []float64
	var x []float64
	for k := range keys {
		x = append(x, float64(keys[k])/p.Lambda)
		weigths = append(weigths, float64(v[keys[k]]))
		datapoints = datapoints + v[keys[k]]
	}

	var avg, std = stat.MeanStdDev(x, weigths)
	_, variance := stat.MeanVariance(x, weigths)
	skew := stat.Skew(x, weigths)
	mode, _ := stat.Mode(x, weigths)
	median := median(x, weigths)

	fmt.Println("\n", target)
	fmt.Printf("#Lambda\t\tAlpha\t\tMean\t\tStd\t\tVar\t\tMedian\t\tMode\t\tSkew\t\tMin\t\tMax\t\tN\n")
	if variance > 10000 {
		fmt.Printf("%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t%.2f\t\t%.2f\t\t%.3f\t\t%.2f\t\t%.2f\t\t%d\n", p.Lambda, p.Alpha, avg, std, variance, median, mode, skew, x[0], x[len(x)-1], datapoints)
	} else {
		fmt.Printf("%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.3f\t\t%.2f\t\t%.2f\t\t%d\n", p.Lambda, p.Alpha, avg, std, variance, median, mode, skew, x[0], x[len(x)-1], datapoints)
	}

	// save to file for plot

	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.2f", p.Alpha)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/velocity_" + rateType + "_" + target +
		"_lambda_" + lambdaStr +
		"_alpha_" + alphaStr + "_.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return 0
	}
	defer f.Close()
	for i, k := range x {
		//fmt.Println("Key:", k, "Value:", m[k])
		_, err = f.WriteString(fmt.Sprintf("%f\t%f\n", k, weigths[i]/float64(datapoints)*p.Lambda)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}

	return datapoints
}

func (velo *velocityResult) Join(b velocityResult) (r velocityResult) {
	if velo.vID == nil {
		return b
	}

	for i := range b.vID {
		r.vID = append(r.vID, joinMapStatInt(velo.vID[i], b.vID[i]))
	}

	for i := range b.dApprovers {
		r.dApprovers = append(r.dApprovers, joinMapStatInt(velo.dApprovers[i], b.dApprovers[i]))
	}

	for i := range b.vTime {
		r.vTime = append(r.vTime, joinMapStatFloat64(velo.vTime[i], b.vTime[i]))
	}

	for i := range b.vCW {
		r.vCW = append(r.vCW, joinMapStatInt(velo.vCW[i], b.vCW[i]))
	}

	for i := range b.vCWfirst {
		r.vCWfirst = append(r.vCWfirst, joinMapStatInt(velo.vCWfirst[i], b.vCWfirst[i]))
	}

	return r
}

func (velo velocityResult) Save(p Parameters) (err error) {
	if err = velo.SaveVID(p); err != nil {
		return err
	}
	if err = velo.SaveVTime(p); err != nil {
		return err
	}
	if err = velo.saveApprovers(p); err != nil {
		return err
	}
	if err = velo.saveCW(p); err != nil {
		return err
	}
	if err = velo.saveCWfirst(p); err != nil {
		return err
	}

	return err
}

func (velo velocityResult) SaveStat(p Parameters) (err error) {
	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.2f", p.Alpha)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/velocity_stat_" + p.TSA + "_" + rateType +
		"_lambda_" + lambdaStr +
		"_alpha_" + alphaStr + "_.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(velo.Stat(p)) // writing...

	if err != nil {
		fmt.Printf("error writing string: %v", err)
		return err
	}

	return nil
}

func (velo velocityResult) Stat(p Parameters) (result string) {
	result = velo.StatVID(p)
	result += "\n"
	result += velo.StatVTime(p)
	result += "\n"
	result += velo.StatCW(p)
	result += "\n"
	result += velo.StatCWfirst(p)
	result += "\n"
	result += velo.StatApprovers(p)

	return result
}

// ToString converts a StatInt to a string
func (s StatInt) ToString(p Parameters, normalized bool) (result string) {
	var keys []int
	var datapoints int
	for k := range s.v {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// calculate statistics
	var weigths []float64
	var x []float64
	for k := range keys {
		norm := 1.
		if normalized {
			norm = p.Lambda
		}
		x = append(x, float64(keys[k])/norm)
		weigths = append(weigths, float64(s.v[keys[k]]))
		datapoints = datapoints + s.v[keys[k]]
	}

	var avg, std = stat.MeanStdDev(x, weigths)
	_, variance := stat.MeanVariance(x, weigths)
	skew := stat.Skew(x, weigths)
	mode, _ := stat.Mode(x, weigths)
	median := median(x, weigths)

	//result += fmt.Sprintf("%s\n", s.desc)

	if variance > 10000 {
		result += fmt.Sprintf("%s\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t%.2f\t\t%.2f\t\t%.3f\t\t%.2f\t\t%.2f\t\t%d\n", s.desc, p.Lambda, p.Alpha, avg, std, variance, median, mode, skew, x[0], x[len(x)-1], datapoints)
	} else {
		result += fmt.Sprintf("%s\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.3f\t\t%.2f\t\t%.2f\t\t%d\n", s.desc, p.Lambda, p.Alpha, avg, std, variance, median, mode, skew, x[0], x[len(x)-1], datapoints)
	}
	return result
}

// ToString converts a StatFloat64 to a string
func (s StatFloat64) ToString(p Parameters, normalized bool) (result string) {
	var keys []float64
	var datapoints int
	for k := range s.v {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	// calculate statistics
	var weigths []float64
	var x []float64
	for k := range keys {
		norm := 1.
		if normalized {
			norm = p.Lambda
		}
		x = append(x, float64(keys[k])/norm)
		weigths = append(weigths, float64(s.v[keys[k]]))
		datapoints = datapoints + s.v[keys[k]]
	}

	var avg, std = stat.MeanStdDev(x, weigths)
	_, variance := stat.MeanVariance(x, weigths)
	skew := stat.Skew(x, weigths)
	mode, _ := stat.Mode(x, weigths)
	median := median(x, weigths)

	//result += fmt.Sprintf("%s\n", s.desc)

	if variance > 10000 {
		result += fmt.Sprintf("%s\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t%.2f\t\t%.2f\t\t%.3f\t\t%.2f\t\t%.2f\t\t%d\n", s.desc, p.Lambda, p.Alpha, avg, std, variance, median, mode, skew, x[0], x[len(x)-1], datapoints)
	} else {
		result += fmt.Sprintf("%s\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.3f\t\t%.2f\t\t%.2f\t\t%d\n", s.desc, p.Lambda, p.Alpha, avg, std, variance, median, mode, skew, x[0], x[len(x)-1], datapoints)
	}
	return result
}

// Save saves a StatInt on a file
func (s StatInt) Save(p Parameters, target string, normalized bool) error {
	var keys []int
	var datapoints int
	for k := range s.v {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// calculate statistics
	var weigths []float64
	var x []float64
	norm := 1.
	for k := range keys {
		if normalized {
			norm = p.Lambda
		}
		x = append(x, float64(keys[k])/norm)
		weigths = append(weigths, float64(s.v[keys[k]]))
		datapoints = datapoints + s.v[keys[k]]
	}
	// save to file for plot

	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.2f", p.Alpha)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/velocity_" + target + "_" + p.TSA + "_" + rateType + "_" + s.desc +
		"_lambda_" + lambdaStr +
		"_alpha_" + alphaStr + "_.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()
	for i, k := range x {
		//fmt.Println("Key:", k, "Value:", m[k])
		if target == "approvers" {
			_, err = f.WriteString(fmt.Sprintf("%d\t%f\n", int(k), weigths[i]/float64(datapoints)*norm)) // writing...
		} else {
			_, err = f.WriteString(fmt.Sprintf("%f\t%f\n", k, weigths[i]/float64(datapoints)*norm)) // writing...
		}
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
	return nil
}

// Save saves a StatFloat64 as a file
func (s StatFloat64) Save(p Parameters, target string, normalized bool) error {
	var keys []float64
	var datapoints int
	for k := range s.v {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	var weigths []float64
	var x []float64
	norm := 1.
	for k := range keys {
		if normalized {
			norm = p.Lambda
		}
		x = append(x, float64(keys[k])/norm)
		weigths = append(weigths, float64(s.v[keys[k]]))
		datapoints = datapoints + s.v[keys[k]]
	}
	// save to file for plot

	lambdaStr := fmt.Sprintf("%.2f", p.Lambda)
	alphaStr := fmt.Sprintf("%.2f", p.Alpha)
	var rateType string
	if p.ConstantRate {
		rateType = "constant"
	} else {
		rateType = "poisson"
	}
	f, err := os.Create("data/velocity_" + target + "_" + p.TSA + "_" + rateType + "_" + s.desc +
		"_lambda_" + lambdaStr +
		"_alpha_" + alphaStr + "_.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()
	for i, k := range x {
		_, err = f.WriteString(fmt.Sprintf("%f\t%f\n", k, weigths[i]/float64(datapoints)*norm)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
	return nil
}

func (velo velocityResult) SaveVID(p Parameters) error {
	for _, velocity := range velo.vID {
		velocity.Save(p, "ID", true)
	}
	return nil
}

func (velo velocityResult) SaveVTime(p Parameters) error {
	for _, velocity := range velo.vTime {
		velocity.Save(p, "time", false)
	}
	return nil
}

func (velo velocityResult) saveApprovers(p Parameters) error {
	for _, velocity := range velo.dApprovers {
		velocity.Save(p, "approvers", false)
	}
	return nil
}

func (velo velocityResult) saveCW(p Parameters) error {
	for _, velocity := range velo.vCW {
		velocity.Save(p, "cw", true)
	}
	return nil
}

func (velo velocityResult) saveCWfirst(p Parameters) error {
	for _, velocity := range velo.vCWfirst {
		velocity.Save(p, "cw-first", true)
	}
	return nil
}

func (velo velocityResult) StatVID(p Parameters) (result string) {
	result += "#Velocity Stats [ID]\n"
	result += "#Stat Type\tLambda\t\tAlpha\t\tMean\t\tStdDev\t\tVariance\tMedian\t\tMode\t\tSkew\t\tMinVal\t\tMaxVal\t\tN\n"
	for _, velocity := range velo.vID {
		result += velocity.ToString(p, true)
	}
	return result
}

func (velo velocityResult) StatApprovers(p Parameters) (result string) {
	result += "#Direct Approvers Stats\n"
	result += "#Stat Type\tLambda\t\tAlpha\t\tMean\t\tStdDev\t\tVariance\tMedian\t\tMode\t\tSkew\t\tMinVal\t\tMaxVal\t\tN\n"
	for _, velocity := range velo.dApprovers {
		result += velocity.ToString(p, false)
	}
	return result

}

func (velo velocityResult) StatVTime(p Parameters) (result string) {
	result += "#Velocity Stats [Time]\n"
	result += "#Stat Type\tLambda\t\tAlpha\t\tMean\t\tStdDev\t\tVariance\tMedian\t\tMode\t\tSkew\t\tMinVal\t\tMaxVal\t\tN\n"
	for _, velocity := range velo.vTime {
		result += velocity.ToString(p, false)
	}
	return result

}

func (velo velocityResult) StatCW(p Parameters) (result string) {
	result += "#Velocity Stats CW\n"
	result += "#Stat Type\tLambda\t\tAlpha\t\tMean\t\tStdDev\t\tVariance\tMedian\t\tMode\t\tSkew\t\tMinVal\t\tMaxVal\t\tN\n"
	for _, velocity := range velo.vCW {
		result += velocity.ToString(p, true)
	}
	return result
}

func (velo velocityResult) StatCWfirst(p Parameters) (result string) {
	result += "#Velocity Stats CW-first\n"
	result += "#Stat Type\tLambda\t\tAlpha\t\tMean\t\tStdDev\t\tVariance\tMedian\t\tMode\t\tSkew\t\tMinVal\t\tMaxVal\t\tN\n"
	for _, velocity := range velo.vCWfirst {
		result += velocity.ToString(p, true)
	}
	return result
}
