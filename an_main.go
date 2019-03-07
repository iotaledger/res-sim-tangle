// organise analysis from here

package main

import (
	"fmt"
)

//SaveResults saves result
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
	if p.EntropyEnabled {
		fmt.Println(f.entropy.Stat(p))
		f.entropy.Save(p)
		//f.entropy.SaveStat(p)
	}
	if p.pOrphanEnabled && p.SpineEnabled {
		fmt.Println(f.op)
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
