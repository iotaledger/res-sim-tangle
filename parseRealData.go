package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/iotaledger/iota.go/transaction"
	"github.com/iotaledger/iota.go/trinary"

	"github.com/iotaledger/iota.go/api"
)

// TxsByTimestamp are sorted txs by timestamp.
type TxsByTimestamp []iriTx

func (a TxsByTimestamp) Len() int               { return len(a) }
func (a TxsByTimestamp) Swap(i int, j int)      { a[i], a[j] = a[j], a[i] }
func (a TxsByTimestamp) Less(i int, j int) bool { return a[i].time < a[j].time }

// // TxsByTimestamp are sorted txs by timestamp.
// type TxsByTrunkID []iriTx

// func (a TxsByTrunkID) Len() int          { return len(a) }
// func (a TxsByTrunkID) Swap(i int, j int) { a[i], a[j] = a[j], a[i] }
// func (a TxsByTrunkID) Less(i int, j int) bool {
// 	if a[i].bundle == a[j].bundle {
// 		//return (a[i].ref[0]) <= (a[j].ref[0])
// 		return (a[i].bundleCurrentIndex) > (a[j].bundleCurrentIndex)
// 	}
// 	return a[i].time < a[j].time
// }

type iriTx struct {
	Hash               trinary.Hash
	id                 int
	time               int64
	cw                 int
	cw2                int // TODO: to remove, used only to compare different CW update mechanisms
	ref                []int
	refHash            []trinary.Hash
	app                []int
	appHash            []trinary.Hash
	bundle             trinary.Hash
	firstApproval      float64
	bundleCurrentIndex uint64
	bundleLastIndex    uint64
	isTip              bool
	//trunkHash          trinary.Hash
	//branchHash         trinary.Hash
}

func pullData(filename, iriURI string, numberOfTxs int) error {

	// create a new API instance
	api, err := api.ComposeAPI(api.HTTPClientSettings{URI: iriURI})
	if err != nil {
		log.Fatal(err)
		return err
	}

	//overwrite old file
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error creating file: %v", err)
	}
	defer f.Close()

	// create BFS list for txs to visit
	set := make(map[trinary.Hash]bool)
	tips := make(map[trinary.Hash]bool)
	toVisit, err := api.GetTips()
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, tx := range toVisit {
		set[tx] = true
		tips[tx] = true
	}

	// perform BFS
	for len(toVisit) > 0 && len(set) < numberOfTxs {
		toVisit = bfs(toVisit[0], toVisit, set, api)
	}

	// save trytes on a file
	for k := range set {
		err = saveTrytes(k, api, filename)
	}
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

func (sim *Sim) buildTangleFromFile(filename string) error {
	// read trytes from a file
	trytes, err := readTrytes(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//fmt.Println(trytes)

	// create Tangle array with hashes and times
	iriTangle := []iriTx{}
	for _, k := range trytes {
		tx := fillTx(k)
		if tx.time != 0 {
			iriTangle = append(iriTangle, tx)
		}
	}
	//building appHash
	for _, tx := range iriTangle {
		for _, k := range tx.ref {
			iriTangle[k].appHash = append(iriTangle[k].appHash, tx.Hash)
		}
	}

	sort.Sort(TxsByTimestamp(iriTangle)) // sort txs by time
	//sort.Sort(TxsByTrunkID(tangle))

	hashToID := make(map[trinary.Hash]int)
	// create ID for []Tx
	for i, t := range iriTangle {
		iriTangle[i].id = i
		hashToID[t.Hash] = i
	}
	// create refs for []Tx
	for i, t := range iriTangle {
		for _, v := range t.refHash {
			if _, ok := hashToID[v]; ok {
				//fmt.Println(api.get)
				iriTangle[i].ref = append(iriTangle[i].ref, hashToID[v])
			}
		}
	}

	buildApprovers(iriTangle)
	// print output
	// for i, t := range iriTangle {
	// 	fmt.Println(i, t.id, "\t", t.ref, "\t", t.app, "\t", t.time, "\t", t.bundleCurrentIndex, t.bundleLastIndex, "\t", t.bundle[:5], t.Hash[:5], "\t", t.refHash[0][:5], t.refHash[1][:5])
	// }
	// fmt.Println("tps=", float64(len(tangle))/2/(tangle[len(tangle)*3/4].time-tangle[len(tangle)/4].time))

	// converting iri tangle to sim tangle
	sim.tangle = make([]Tx, len(iriTangle))
	for i, iriTx := range iriTangle {
		sim.tangle[i] = iriTx.ToTx()
	}

	//TODO init function
	sim.param.CWMatrixLen = len(sim.tangle)
	sim.param.TangleSize = len(sim.tangle)
	sim.cw = make([][]uint64, sim.param.CWMatrixLen)

	sim.initializeCW(sim.tangle)
	sim.computeCW()
	fmt.Println(len(sim.cw))
	sim.computeCWDFS(sim.tangle)

	return err

}

func (sim *Sim) initializeCW(tangle []Tx) {
	base := 64
	for i, tx := range tangle {
		if len(tx.ref) == 0 {
			sim.cw[i] = make([]uint64, (i/base)+1)
			setCW(sim.cw[i], tx.id)
		}
	}
}

func (sim *Sim) computeCW() {
	//fmt.Println(len(sim.tangle))
	//fmt.Println(len(sim.cw))
	//pauseit()
	for _, tx := range sim.tangle {
		//fmt.Println(tx.id)
		//printCWRef(sim.cw[tx.id])
		if len(tx.ref) > 0 {
			sim.updateCWOpt(tx)
			//	printCWRef(sim.cw[tx.id])
			//pauseit()
		}
	}
}

func (sim *Sim) computeCWDFS(tangle []Tx) {
	for _, tx := range tangle {
		sim.updateCWDFS(tx)
	}
}

// func addGenesis(tangle []iriTx) {
// 	newTangle := make([]iriTx, len(tangle)+1)
// 	//add Genesis
// 	newTangle[0] = iriTx{
// 		id:            0,
// 		time:          0,
// 		cw:            1,
// 		firstApproval: -1,
// 		cw2:           1,
// 	}
// 	//fill genesis approvers

// }

func (a *iriTx) ToTx() Tx {
	return Tx{
		id:   a.id,
		time: float64(a.time),
		ref:  a.ref,
		app:  a.app,
		cw:   1,
		cw2:  1,
	}
}

func compareTangle(A, B []Tx) bool {
	for i := 0; i < len(A); i++ {
		if len(A[i].ref) != len(B[i].ref) {
			return false
		}
		if len(A[i].ref) > 0 {
			if A[i].ref[0] != B[i].ref[0] {
				return false
			}
		}
	}
	return true
}

// dfs type tx search
func bfs(tx trinary.Hash, toVisit []trinary.Hash, seenTxs map[trinary.Hash]bool, api *api.API) []trinary.Hash {
	// get trytes of tx
	// call bfs for both branch and trunk
	txTrytes, _ := api.GetTrytes(tx)
	toVisit = toVisit[1:]
	txObject, _ := transaction.AsTransactionObject(txTrytes[0])
	if !seenTxs[txObject.TrunkTransaction] {
		seenTxs[txObject.TrunkTransaction] = true
		toVisit = append(toVisit, txObject.TrunkTransaction)
	}
	if !seenTxs[txObject.BranchTransaction] {
		seenTxs[txObject.BranchTransaction] = true
		toVisit = append(toVisit, txObject.BranchTransaction)
	}
	return toVisit
}

// dfs type tx search
func iriDfs(tx trinary.Hash, visited map[trinary.Hash]bool, size int, api *api.API) {
	if len(visited) < size {
		// get trytes of tx
		// call bfs for both branch and trunk
		txTrytes, _ := api.GetTrytes(tx)
		txObject, _ := transaction.AsTransactionObject(txTrytes[0])
		if !visited[txObject.TrunkTransaction] {
			visited[txObject.TrunkTransaction] = true
			iriDfs(txObject.TrunkTransaction, visited, size, api)
		}
		if !visited[txObject.BranchTransaction] {
			visited[txObject.BranchTransaction] = true
			iriDfs(txObject.BranchTransaction, visited, size, api)
		}
	}
}

// collect data and store it into []iriTx
func fillTx(txTrytes trinary.Trytes) iriTx {
	var goTx iriTx
	//txTrytes, _ := api.GetTrytes(tx)
	txObject, _ := transaction.AsTransactionObject(txTrytes)
	goTx.Hash = txObject.Hash
	goTx.time = txObject.AttachmentTimestamp
	goTx.refHash = append(goTx.refHash, txObject.TrunkTransaction)
	goTx.refHash = append(goTx.refHash, txObject.BranchTransaction)
	goTx.bundle = txObject.Bundle
	goTx.bundleCurrentIndex = txObject.CurrentIndex
	goTx.bundleLastIndex = txObject.LastIndex
	return goTx
}

func saveTrytes(tx trinary.Hash, api *api.API, filename string) error {
	txTrytes, _ := api.GetTrytes(tx)

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if _, err := f.Write([]byte(txTrytes[0] + "\n")); err != nil {
		log.Fatal(err)
		return err
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
		return err
	}

	return err
}

func readTrytes(filename string) ([]trinary.Trytes, error) {
	var trytes []trinary.Trytes
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	// Start reading from the file using a scanner.
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		txTrytes := scanner.Text()
		trytes = append(trytes, txTrytes)
	}

	if scanner.Err() != nil {
		fmt.Printf(" > Failed!: %v\n", scanner.Err())
	}

	return trytes, err
}

func buildApprovers(tangle []iriTx) {
	for _, v := range tangle {
		//v is the tx I have to use for checking
		if len(v.ref) > 0 {
			for _, approvee := range v.ref {
				tangle[approvee].app = appendUnique(tangle[approvee].app, v.id)
			}
		}
	}
}

func isRefConsistent(tangle []Tx) bool {
	for _, tx := range tangle {
		// check ref
		for _, ref := range tx.ref {
			if tx.id < ref {
				return false
			}
		}
	}
	return true
}

//func fillFirstTimeApproval
