// organise analysis from here

package main

import "fmt"

//SaveResults saves result
func (f *Result) SaveResults(p Parameters) {
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
	if p.EntropyEnabled {
		fmt.Println(f.entropy.Stat(p))
		f.entropy.Save(p)
		//f.entropy.SaveStat(p)
	}
	return
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
	if p.EntropyEnabled {
		f.entropy = f.entropy.Join(batch.entropy)
	}
	f.tips = f.tips.Join(batch.tips)
}
