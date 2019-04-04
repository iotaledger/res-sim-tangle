package main

import (
	"fmt"
	"os"
	"sort"
)

//func printCWRef(a map[int][]uint64) {
// func printCWRef(a [][]uint64) {
// 	for i := 0; i < len(a); i++ {
// 		fmt.Printf("%d: ", i)
// 		for j := len(a[i]) - 1; j >= 0; j-- {
// 			fmt.Printf("\t%064b\n", a[i][j])
// 		}
// 		fmt.Println()
// 	}
// }

func printCWRef(a []uint64) {
	for j := len(a) - 1; j >= 0; j-- {
		fmt.Printf("\t%064b\n", a[j])
	}
	fmt.Println()

}

func printApprovers(a []Tx) {
	for _, approvee := range a {
		fmt.Println(approvee.id, "<-", unique(approvee.app))
	}
}

func saveApprovers(a map[int][]int) {
	var keys []int
	for key := range a {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	f, err := os.Create("data/approvers.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
	}
	defer f.Close()

	for _, t := range keys {
		_, err = f.WriteString(fmt.Sprintln(t, "<-", unique(a[t]))) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
}

func saveTangle(tangle []Tx) {
	f, err := os.Create("data/tangle.txt")
	if err != nil {
		fmt.Printf("error creating file: %v", err)
	}
	defer f.Close()

	for _, t := range tangle {
		_, err = f.WriteString(fmt.Sprintln(t)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
}

func printCW(sim Sim) {
	for _, t := range sim.tangle {
		fmt.Println(t.id, t.cw)
	}
}

func printTips(a map[int]bool) {
	var keys []int
	for key := range a {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	fmt.Println(keys)
}
