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
	if p.VelocityEnabled {
		//???is there a way this can be defined in the velocity.go file
		if p.TSA != "RW" {
			result.velocity = newVelocityResult([]string{"rw", "all", "first", "last", "second", "third", "fourth", "only-1", "CW-Max", "CW-Min", "CWMaxRW", "CWMinRW", "backU"}, *p)
		} else {
			result.velocity = newVelocityResult([]string{"rw", "all", "first", "last", "CW-Max", "CW-Min", "backU", "backB", "URW", "backG"}, *p)
			//vr = newVelocityResult([]string{"rw", "all", "first"}, sim.param)
			//fmt.Println(*vr)
			//vr = newVelocityResult([]string{"rw", "all", "back"})
		}
	}
	if p.AnPastCone.Enabled {
		result.PastCone = newPastConeResult([]string{"avg", "1", "2", "3", "4", "5", "rest"})
	}
	if p.AnFocusRW.Enabled {
		result.FocusRW = newFocusRWResult([]string{"0.1"})
	}
	if p.ExitProbEnabled {
		result.exitProb = newExitProbResult()
	}
	if p.pOrphanEnabled {
		result.op = newPOrphanResult(p)
	}
	if p.DistSlicesEnabled {
		result.DistSlices = newDistSlicesResult()
	}
	if p.DistRWsEnabled {
		result.DistRWs = newDistRWsResult()
	}
	if p.AppStatsRWEnabled {
		result.AppStatsRW = newAppStatsRWResult()
	}
	if p.AppStatsAllEnabled {
		result.AppStatsAll = newAppStatsAllResult()
	}

}

//save results at end of simulation
func (f *Result) SaveResults(p Parameters) {
	if p.CountTipsEnabled {
		f.tips.Statistics(p)
		fmt.Println(f.tips.ToString(p))
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
	if p.VelocityEnabled {
		fmt.Println(f.velocity.Stat(p))
		//f.velocity.Save(p)
		//f.velocity.SaveStat(p)
	}
	if p.AnPastCone.Enabled {
		f.PastCone.finalprocess(p)
		f.PastCone.Save(p)
	}
	if p.AnFocusRW.Enabled {
		f.FocusRW.finalprocess(p)
		f.FocusRW.Save(p)
	}
	if p.ExitProbEnabled {
		f.exitProb.Save(p)
	}
	if p.pOrphanEnabled && p.SpineEnabled {
		fmt.Println(f.op)
	}
	if p.DistSlicesEnabled {
		f.DistSlices.finalprocess()
		f.DistSlices.Save(p)
	}
	if p.DistRWsEnabled {
		f.DistRWs.finalprocess()
		f.DistRWs.Save(p)
	}
	if p.AppStatsRWEnabled {
		f.AppStatsRW.finalprocess()
		f.AppStatsRW.Save(p)
	}
	if p.AppStatsAllEnabled {
		f.AppStatsAll.finalprocess()
		f.AppStatsAll.Save(p)
	}
	return
}

//Evaluate after each tx
func (result *Result) EvaluateAfterTx(sim *Sim, p *Parameters, run, i int) {
	// ??? the following lines seems to make no sense. can we remove it?
	// if i > sim.param.minCut && i < sim.param.maxCut {
	// 	nTips += len(sim.tips)
	// }
	if p.CountTipsEnabled {
		sim.countTips(i, run, &result.tips)
	}
	if p.pOrphanEnabled && p.pOrphanLinFitEnabled && p.SpineEnabled {
		sim.runAnOPLinfit(i, &result.op, run)
	}
}

//Evaluate after each Tangle
func (result *Result) EvaluateTangle(sim *Sim, p *Parameters, run int) {
	if p.SpineEnabled {
		sim.computeSpine()
		//printApprovers(sim.spineApprovers)
	}

	if p.CountTipsEnabled {
		//sim.runTipsStat(&result.tips)
	}
	if p.CWAnalysisEnabled {
		sim.fillCW(run, &result.cw)
	}
	if p.VelocityEnabled {
		sim.runVelocityStat(&result.velocity)
	}
	if p.AnPastCone.Enabled {
		sim.runAnPastCone(&result.PastCone)
	}
	if p.AnFocusRW.Enabled {
		sim.evalTangle_AnFocusRW(&result.FocusRW)
	}
	if p.ExitProbEnabled {
		sim.runExitProbStat(&result.exitProb)
	}
	if p.pOrphanEnabled && p.SpineEnabled {
		sim.runOrphaningP(&result.op)
	}
	if p.DistSlicesEnabled {
		sim.evalTangle_DistSlices(&result.DistSlices)
	}
	if p.DistRWsEnabled {
		sim.evalTangle_DistRWs(&result.DistRWs)
	}
	if p.AppStatsRWEnabled {
		sim.evalTangle_AppStatsRW(&result.AppStatsRW)
	}
	if p.AppStatsAllEnabled {
		sim.evalTangle_AppStatsAll(&result.AppStatsAll)
	}
}

//JoinResults joins result
func (f *Result) JoinResults(batch Result, p Parameters) {
	if p.VelocityEnabled {
		f.velocity = f.velocity.Join(batch.velocity)
	}
	if p.AnPastCone.Enabled {
		f.PastCone = f.PastCone.Join(batch.PastCone)
	}
	if p.AnFocusRW.Enabled {
		f.FocusRW = f.FocusRW.Join(batch.FocusRW)
	}
	if p.ExitProbEnabled {
		f.exitProb = f.exitProb.Join(batch.exitProb)
	}
	if p.pOrphanEnabled && p.SpineEnabled {
		f.op = f.op.Join(batch.op)
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
	if p.DistRWsEnabled {
		f.DistRWs.Join(batch.DistRWs)
	}
	if p.AppStatsRWEnabled {
		f.AppStatsRW.Join(batch.AppStatsRW)
	}
	if p.AppStatsAllEnabled {
		f.AppStatsAll.Join(batch.AppStatsAll)
	}
	//f.avgtips = f.avgtips.Join(batch.avgtips)
}
