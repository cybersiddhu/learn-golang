package main

import (
	"fmt"
	gobam "github.com/cybersiddhu/biogo.boom"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"
)

//default bin size in 1
var binSize = 1
var dumpSize = 10000
var wg sync.WaitGroup

type Coverage struct {
	counter   map[int]int
	calculate func(*gobam.Record) bool
}

func (c *Coverage) AddCoverage(pos int) {
	if value, ok := c.counter[pos]; ok {
		value = value + 1
		c.counter[pos] = value
	} else {
		c.counter[pos] = 1
	}
}

func InitCoverage(start int, end int) *Coverage {
	c := new(Coverage)
	c.counter = make(map[int]int, end+1)
	for i := start; i <= end; i += 1 {
		c.counter[i] = 0
	}
	c.calculate = func(r *gobam.Record) bool {
		for start := r.Start(); start <= r.End(); start += 1 {
			c.AddCoverage(start)
		}
		return false
	}
	return c
}

func main() {

	log.Println("starting up")
	bam, err := gobam.OpenBAM(os.Args[1])
	dieIfError(err)
	log.Println("got bam reader")


	log.Println("before loading index")
	idx, err := gobam.LoadIndex(os.Args[1])
	dieIfError(err)
	log.Println("loaded index")


	log.Println("before maxprocs")
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("after maxprocs")

	lengths := bam.RefLengths()
	cm := &CoverageHandler{
		bam:     bam,
		index:   idx,
		chunk:   dumpSize,
		binSize: binSize,
		channel: make(chan *Coverage, len(lengths)),
	}

	for i, name := range bam.RefNames() {
		if id, ok := bam.RefID(name); ok {
			log.Printf("before sending %s\n", name)
			wg.Add(1)
			go cm.Generate(id, int(lengths[i]))
			log.Printf("before writing %s\n", name)
			go cm.Write(name)
		}
	}
	wg.Wait()
}

type CoverageHandler struct {
	bam     *gobam.BAMFile
	index   *gobam.Index
	chunk   int
	binSize int
	channel chan *Coverage
}

func (cm *CoverageHandler) Generate(id, length int) {
	log.Printf("going to generate coverage for %d\n", id)
	cov := InitCoverage(0, (length - 1))

	for start := 0; start < length; start += cm.chunk {
		// calculate end
		//0 based indexing
		end := start + cm.chunk - 1
		if end > length {
			end = length
		}
		_, err := cm.bam.Fetch(cm.index, id, start, end, cov.calculate)
		dieIfError(err)
	}
	log.Printf("sending coverage for %d\n", id)
	cm.channel <- cov
}

func (cm *CoverageHandler) Write(name string) {
	cov := <-cm.channel
	log.Printf("receiving coverage for %s\n", name)
	defer wg.Done()
	var sortedLoc []int
	for k := range cov.counter {
		sortedLoc = append(sortedLoc, k)
	}
	sort.Ints(sortedLoc)

	w, err := os.Create(name + ".wig")
	dieIfError(err)
	defer w.Close()

	fmt.Fprintf(w, "fixedStep chrom=%s start=1 step=%d span=%d\n", name, binSize, binSize)
	for _, k := range sortedLoc {
		fmt.Fprintln(w, cov.counter[k])
	}
	log.Printf("done writing coverage for %s\n", name)
}

func dieIfError(e error) {
	if e != nil {
		panic(e)
	}
}
