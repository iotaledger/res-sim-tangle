// organise analysis from here

package main

import (
	"fmt"
)

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
		f.velocity.Save(p)
		f.velocity.SaveStat(p)
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
		// fmt.Println(f.exitProb.Stat(p))
		// f.exitProb.SaveExitProb(p, "ep")
		// //  the following way was the easiest way, otherwise it would have been necessary to copy a huge amount of functions
		// f.exitProb.ep = f.exitProb.ep2
		// f.exitProb.mean = f.exitProb.mean2
		// f.exitProb.median = f.exitProb.median2
		// f.exitProb.std = f.exitProb.std2
		// fmt.Println(f.exitProb.Stat(p))
		// f.exitProb.SaveExitProb(p, "ep2")
		// //f.exitProb.SaveStat(p)
	}
	if p.pOrphanEnabled && p.SpineEnabled {
		fmt.Println(f.op)
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
		sim.runAnFocusRW(&result.FocusRW)
	}
	if p.ExitProbEnabled {
		sim.runExitProbStat(&result.exitProb)
	}
	if p.pOrphanEnabled && p.SpineEnabled {
		sim.runOrphaningP(&result.op)
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
	f.avgtips = f.avgtips.Join(batch.avgtips)
}
