package main

import (
	"fmt"
	gobam "github.com/cybersiddhu/biogo.boom"
	"log"
	"os"
	"runtime"
	//"sync"
)

//default bin size in 1
var binSize = 1
var dumpSize = 10000
//var wg sync.WaitGroup
var covStack map[string]*Coverage = make(map[string]*Coverage)
//var covMap map[int]int = make(map[int]int)

type Coverage struct {
	counter map[int]int
	end     int
	name    string
}

func (c *Coverage) AddCoverage(pos int) {
	if value, ok := c.counter[pos]; ok {
		value += 1
		c.counter[pos] = value
	} else {
		c.counter[pos] = 1
	}
}

func InitCoverage(name string, end int) *Coverage {
	c := new(Coverage)
	c.name = name
	c.counter = make(map[int]int, end+1)
	c.end = end
	//for i := start; i <= end; i += 1 {
	//c.counter[i] = 0
	//}
	//c.calculate = func(r *gobam.Record) bool {
	//for start := r.Start(); start <= r.End(); start += 1 {
	//c.AddCoverage(start)
	//}
	//return false
	//}
	return c
}

func RecordAlignment(r *gobam.Record) bool {
	if cov, ok := covStack[r.Name()]; ok {
		for start := r.Start(); start <= r.End(); start += 1 {
			cov.AddCoverage(start)
//			covMap[r.RefID()] = start
		}
	}
	return false
}

func main() {

	//log.Println("starting up")
	bam, err := gobam.OpenBAM(os.Args[1])
	dieIfError(err)
	//log.Println("got bam reader")

	//log.Println("before loading index")
	idx, err := gobam.LoadIndex(os.Args[1])
	dieIfError(err)
	//log.Println("loaded index")

	//log.Println("before maxprocs")
	runtime.GOMAXPROCS(runtime.NumCPU())
	//log.Println("after maxprocs")

	lengths := bam.RefLengths()
	//covStack = make(map[string]*Coverage, len(lengths))
	cm := &CoverageHandler{
		bam:     bam,
		index:   idx,
		chunk:   dumpSize,
		binSize: binSize,
	}

	for i, name := range bam.RefNames() {
		if id, ok := bam.RefID(name); ok {
			//	log.Printf("before sending %s\n", name)
			cm.Generate(id, name, int(lengths[i]))
			//	log.Printf("before writing %s\n", name)
		}
	}
	//for _, cov := range covStack {
		//wg.Add(1)
		//go cm.Write(cov)
	//}
	//wg.Wait()
}

type CoverageHandler struct {
	bam     *gobam.BAMFile
	index   *gobam.Index
	chunk   int
	binSize int
}

func (cm *CoverageHandler) Generate(id int, name string, length int) {
	log.Printf("going to generate coverage for %d\n", id)
	cov := InitCoverage(name, (length - 1))
	covStack[name] = cov

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
	log.Printf("finished coverage for %s\n", name)
}

func (cm *CoverageHandler) Write(cov *Coverage) {
	log.Printf("receiving coverage for %s\n", cov.name)
	//defer wg.Done()

	w, err := os.Create(cov.name + ".wig")
	dieIfError(err)
	defer w.Close()

	fmt.Fprintf(w, "fixedStep chrom=%s start=1 step=%d span=%d\n", cov.name, cm.binSize, cm.binSize)
	for pos := 0; pos <= cov.end; pos += 1 {
		if value, ok := cov.counter[pos]; ok {
			fmt.Fprintln(w, value)
		} else {
			fmt.Fprintln(w, 0)
		}
	}
	log.Printf("done writing coverage for %s\n", cov.name)
}

func dieIfError(e error) {
	if e != nil {
		panic(e)
	}
}
