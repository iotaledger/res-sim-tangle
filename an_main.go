// organise analysis from here

package main

import (
	"fmt"
)

// initiate analysis variables
func (result *Result) initResults(p *Parameters) {
	if p.CountTipsEnabled {
		result.tips = newTipsResult(*p)
	}
	if p.CWAnalysisEnabled {
		result.cw = newCWResult(*p)
	}
	if p.AnPastCone.Enabled {
		result.PastCone = newPastConeResult([]string{"avg", "1", "2", "3", "4", "5", "rest"})
	}
	if p.pOrphanEnabled {
		result.op = newPOrphanResult(p)
	}
	if p.DistSlicesEnabled {
		result.DistSlices = newDistSlicesResult()
	}
	if p.AppStatsAllEnabled {
		result.AppStatsAll = newAppStatsAllResult()
	}

}

//save results at end of simulation
func (f *Result) FinalEvaluationSaveResults(p Parameters) {
	if p.CountTipsEnabled {
		f.tips.Statistics(p)
		// fmt.Println(f.tips.ToString(p))
		//fmt.Println(f.tips.nTipsToString(p, 0))
		f.tips.Save(p, 0)
		// //debug
		// var keys []int
		// for k := range f.tips.tPDF.v {
		// 	keys = append(keys, k)
		// }
		// sort.Ints(keys)

		// for _, v := range keys {
		// 	fmt.Println(v, f.tips.tPDF.v[v], f.tips.pdf[0].v[v])
		// }
		// //fmt.Println(f.tips.pdf[0])
		// //fmt.Println(f.tips.tPDF)
	}
	if p.CWAnalysisEnabled {
		f.cw.Statistics(p)
		//fmt.Println(f.cw.ToString(p))
		fmt.Println(f.cw.cwToString(p, 0))
		f.cw.Save(p, 0)
	}
	if p.AnPastCone.Enabled {
		f.PastCone.finalprocess(p)
		f.PastCone.Save(p)
	}
	if p.DistSlicesEnabled {
		f.DistSlices.finalprocess()
		f.DistSlices.Save(p)
	}
	if p.AppStatsAllEnabled {
		f.AppStatsAll.finalprocess()
		f.AppStatsAll.Save(p)
	}
	return
}

//Evaluate after each tx
func (result *Result) EvaluateAfterTx(sim *Sim, p *Parameters, run, i int) {
	if p.CountTipsEnabled {
		sim.countTips(i, run, &result.tips)
	}
}

//Evaluate after each Tangle
func (result *Result) EvaluateTangle(sim *Sim, p *Parameters, run int) {
	if p.CountTipsEnabled {
		sim.countOrphanTips(run, &result.tips)
		//sim.runTipsStat(&result.tips)
	}
	if p.CWAnalysisEnabled {
		sim.fillCW(run, &result.cw)
	}
	if p.AnPastCone.Enabled {
		sim.runAnPastCone(&result.PastCone)
	}
	if p.DistSlicesEnabled {
		sim.evalTangle_DistSlices(&result.DistSlices)
	}
	if p.AppStatsAllEnabled {
		sim.evalTangle_AppStatsAll(&result.AppStatsAll)
	}
}

//JoinResults joins result
func (f *Result) JoinResults(batch Result, p Parameters) {
	if p.AnPastCone.Enabled {
		f.PastCone = f.PastCone.Join(batch.PastCone)
	}
	if p.CountTipsEnabled {
		f.tips = f.tips.Join(batch.tips)
	}
	if p.CWAnalysisEnabled {
		f.cw = f.cw.Join(batch.cw)
	}
	if p.DistSlicesEnabled {
		f.DistSlices.Join(batch.DistSlices)
	}
	if p.AppStatsAllEnabled {
		f.AppStatsAll.Join(batch.AppStatsAll)
	}
	//f.avgtips = f.avgtips.Join(batch.avgtips)
}
