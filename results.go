package main

type Result interface {
	Join(Result) Result
	Save(Parameters) error
	SaveStat(Parameters) error
	Stat(Parameters) string
}

type StatInt struct {
	desc string
	v    map[int]int //description: -1=RW, 0=all links, 1=first link, 2=second link, etc.
}

type StatFloat64 struct {
	desc string
	v    map[float64]int //description: -1=RW, 0=all links, 1=first link, 2=second link, etc.
}

func joinMapInt(a, b map[int]int) map[int]int {
	if a == nil {
		return b
	}
	for k, v := range b {
		a[k] += v
	}
	return a
}

func joinMapFloat64(a, b map[float64]int) map[float64]int {
	if a == nil {
		return b
	}
	for k, v := range b {
		a[k] += v
	}
	return a
}

func joinMapStatInt(a, b StatInt) StatInt {
	a.desc = b.desc
	a.v = joinMapInt(a.v, b.v)
	return a
}

func joinMapStatFloat64(a, b StatFloat64) StatFloat64 {
	a.desc = b.desc
	a.v = joinMapFloat64(a.v, b.v)
	return a
}
