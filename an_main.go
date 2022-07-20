// organise analysis from here

package main

// initiate analysis variables
func (result *Result) initResults(p *Parameters) {
	if p.CountTipsEnabled {
		result.tipsResult = newTipsResult(*p)
	}
	// if p.CWAnalysisEnabled {
	// 	result.cw = newCWResult(*p)
	// }
	if p.AnPastCone.Enabled {
		result.PastConeResult = newPastConeResult([]string{"avg", "1", "2", "3", "4", "5", "rest"})
	}
	result.opResult = newOrphanResult(p)
	if p.DistSlicesEnabled {
		result.DistSlicesResult = newDistSlicesResult()
	}
	if p.AppStatsAllEnabled {
		result.AppStatsAllResult = newAppStatsAllResult()
	}

}

//save results at end of simulation
func (f *Result) FinalEvaluationSaveResults(p Parameters) {
	if p.CountTipsEnabled {
		f.tipsResult.Statistics(p)
		// fmt.Println(f.tips.ToString(p))
		//fmt.Println(f.tips.nTipsToString(p, 0))
		f.tipsResult.Save(p, 0)
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
	// if p.CWAnalysisEnabled {
	// 	f.cw.Statistics(p)
	// 	//fmt.Println(f.cw.ToString(p))
	// 	fmt.Println(f.cw.cwToString(p, 0))
	// 	f.cw.Save(p, 0)
	// }
	if p.AnPastCone.Enabled {
		f.PastConeResult.finalprocess(p)
		f.PastConeResult.Save(p)
	}
	if p.DistSlicesEnabled {
		f.DistSlicesResult.finalprocess()
		f.DistSlicesResult.Save(p)
	}
	if p.AppStatsAllEnabled {
		f.AppStatsAllResult.finalprocess()
		f.AppStatsAllResult.Save(p)
	}
	if p.AnOrphanageEnabled {
		f.opResult.Save(p)
	}

	return
}

//Evaluate after each tx
func (result *Result) EvaluateAfterTx(sim *Sim, p *Parameters, run, i int) {
	if p.CountTipsEnabled {
		sim.countTips(i, run, &result.tipsResult)
	}
}

//Evaluate after each Tangle
func (result *Result) EvaluateTangle(sim *Sim, p *Parameters, run int) {
	if p.CountTipsEnabled {
		sim.countOrphanTips(run, &result.tipsResult)
		//sim.runTipsStat(&result.tips)
	}
	// if p.CWAnalysisEnabled {
	// 	sim.fillCW(run, &result.cw)
	// }
	if p.AnPastCone.Enabled {
		sim.runAnPastCone(&result.PastConeResult)
	}
	if p.DistSlicesEnabled {
		sim.evalTangle_DistSlices(&result.DistSlicesResult)
	}
	if p.AppStatsAllEnabled {
		sim.evalTangle_AppStatsAll(&result.AppStatsAllResult)
	}
	if p.AnOrphanageEnabled {
		sim.runOrphanageRecent(&result.opResult) // calculate op2
	}
}

//JoinResults joins result
func (f *Result) JoinResults(batch Result, p Parameters) {
	if p.AnPastCone.Enabled {
		f.PastConeResult = f.PastConeResult.Join(batch.PastConeResult)
	}
	if p.CountTipsEnabled {
		f.tipsResult = f.tipsResult.Join(batch.tipsResult)
	}
	// if p.CWAnalysisEnabled {
	// 	f.cw = f.cw.Join(batch.cw)
	// }
	if p.DistSlicesEnabled {
		f.DistSlicesResult.Join(batch.DistSlicesResult)
	}
	if p.AppStatsAllEnabled {
		f.AppStatsAllResult.Join(batch.AppStatsAllResult)
	}
	if p.AnOrphanageEnabled {
		f.opResult = f.opResult.Join(batch.opResult)
	}
	//f.avgtips = f.avgtips.Join(batch.avgtips)
}
