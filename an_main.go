// organise analysis from here

package main

//Save Results
func (f *Result) SaveResults(p Parameters) {
	if p.VelocityEnabled {
		// fmt.Println(f.velocity.Stat(p))
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
	return
}

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
	f.tips = f.tips.Join(batch.tips)
}
