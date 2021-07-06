// organise analysis from here

package main

// initiate analysis variables
func (result *Result) initResults(p *Parameters) {
	if p.CountTipsEnabled {
		result.tips = newTipsResult(*p)
	}
	if p.CTAnalysisEnabled {
		result.confirmationTime = newCTResult(*p)
	}
}

//save results at end of simulation
func (f *Result) FinalEvaluationSaveResults(p Parameters) {
	if p.CountTipsEnabled {
		f.tips.Statistics(p)
		// fmt.Println(f.tips.ToString(p))
		//fmt.Println(f.tips.nTipsToString(p, 0))
		//f.tips.Save(p, 0)
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
	if p.CTAnalysisEnabled {
		f.confirmationTime.Statistics(p)
		//fmt.Println(f.cw.ToString(p))
		//fmt.Println(f.confirmationTime.ctToString(p, 0))
		f.confirmationTime.Save(p, 0)
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
	if p.CTAnalysisEnabled {
		sim.fillCT(run, &result.confirmationTime)
	}
}

//JoinResults joins result
func (f *Result) JoinResults(batch Result, p Parameters) {
	if p.CountTipsEnabled {
		f.tips = f.tips.Join(batch.tips)
	}
	if p.CTAnalysisEnabled {
		f.confirmationTime = f.confirmationTime.Join(batch.confirmationTime)
	}
}
