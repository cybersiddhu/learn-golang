package main

import (
	"fmt"
	gobam "github.com/cybersiddhu/biogo.boom"
	"io"
	"os"
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
	bam, err := gobam.OpenBAM(os.Args[1])
	dieIfError(err)

	idx, err := gobam.LoadIndex(os.Args[1])
	dieIfError(err)

	for i, name := range bam.RefNames() {
		lengths := bam.RefLengths()
		if id, ok := bam.RefID(name); ok {

			w, err := os.Create(name + ".wig")
			dieIfError(err)
			defer w.Close()
			fmt.Fprintf(w, "fixedStep chrom=%s start=1 step=%d span=%d\n", name, binSize, binSize)

			cov := InitCoverage(0, int(lengths[i]-1))

			for start := 0; start < int(lengths[i]); start += dumpSize {
				// calculate end
				//0 based indexing
				end := start + dumpSize - 1
				if end > int(lengths[i]) {
					end = int(lengths[i])
				}
				_, err := bam.Fetch(idx, id, start, end, cov.calculate)
				dieIfError(err)
			}
			wg.Add(1)
			go WriteCoverage(w, cov)
		}
	}
	wg.Wait()
}

func WriteCoverage(w io.Writer, cov *Coverage) {
	defer wg.Done()
	var sortedLoc []int
	for k := range cov.counter {
		sortedLoc = append(sortedLoc, k)
	}
	sort.Ints(sortedLoc)
	for _, k := range sortedLoc {
		fmt.Fprintln(w, cov.counter[k])
	}
}

func dieIfError(e error) {
	if e != nil {
		panic(e)
	}
}
