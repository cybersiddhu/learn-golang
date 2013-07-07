package main

import (
	"fmt"
	gobam "github.com/cybersiddhu/biogo.boom"
	"log"
	"os"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

//default bin size in 1
var binSize = 1
var dumpSize = 10000
var covStack map[int]map[int]int = make(map[int]map[int]int)

func RecordAlignment(r *gobam.Record) bool {
	if stack, ok := covStack[r.RefID()]; ok {
		for start := r.Start(); start <= r.End(); start += 1 {
			if value, ok := stack[start]; ok {
				value += 1
				stack[start] = value
			} else {
				stack[start] = value
			}
		}
		covStack[r.RefID()] = stack
	} else {
		stack := make(map[int]int)
		for start := r.Start(); start <= r.End(); start += 1 {
			stack[start] = 1 
		}
		covStack[r.RefID()] = stack
	}
	return false
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU() + 2)
	//log.Println("starting up")
	bam, err := gobam.OpenBAM(os.Args[1])
	dieIfError(err)
	//log.Println("got bam reader")

	//log.Println("before loading index")
	idx, err := gobam.LoadIndex(os.Args[1])
	dieIfError(err)
	//log.Println("loaded index")

	//log.Println("before maxprocs")
	//log.Println("after maxprocs")

	lengths := bam.RefLengths()
	//covStack = make(map[string]*Coverage, len(lengths))
	cm := &CoverageHandler{
		bam:     bam,
		index:   idx,
		chunk:   dumpSize,
		binSize: binSize,
	}

	var id2len map[int]int = make(map[int]int)
	var id2name map[int]string = make(map[int]string)

	for i, name := range bam.RefNames() {
		if id, ok := bam.RefID(name); ok {
			//	log.Printf("before sending %s\n", name)
			length := int(lengths[i])
			id2len[id] = length
			id2name[id] = name
			cm.Generate(id, name, length)
			//	log.Printf("before writing %s\n", name)
		}
	}

	for id, cov := range covStack {
		if name, ok := id2name[id]; ok {
			log.Printf("%s with length %d will have %d entries with read\n",name, id2len[id], len(cov))
			wg.Add(1)
			go cm.Write(name, id2len[id], cov)
		}
	}
	wg.Wait()
}

type CoverageHandler struct {
	bam     *gobam.BAMFile
	index   *gobam.Index
	chunk   int
	binSize int
}

func (cm *CoverageHandler) Generate(id int, name string, length int) {
	//log.Printf("going to generate coverage for %s\n", name)

	for start := 0; start < length; start += cm.chunk {
		// calculate end
		//0 based indexing
		end := start + cm.chunk - 1
		if end > length {
			end = length
		}
		_, err := cm.bam.Fetch(cm.index, id, start, end, RecordAlignment)
		dieIfError(err)
	}
	log.Printf("finished coverage for %s with %d bases with read\n", name, len(covStack[id]))
}

func (cm *CoverageHandler) Write(name string, length int, cov map[int]int) {
	log.Printf("receiving coverage for %s\n", name)
	defer wg.Done()

	w, err := os.Create(name + ".wig")
	dieIfError(err)
	defer w.Close()

	fmt.Fprintf(w, "fixedStep chrom=%s start=1 step=%d span=%d\n", name, cm.binSize, cm.binSize)
	for pos := 0; pos < length; pos += 1 {
		if value, ok := cov[pos]; ok {
			fmt.Fprintln(w, value)
		} else {
			fmt.Fprintln(w, 0)
		}
	}
	log.Printf("done writing coverage for %s\n", name)
}

func dieIfError(e error) {
	if e != nil {
		panic(e)
	}
}
