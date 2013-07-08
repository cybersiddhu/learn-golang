package main

import (
	"fmt"
	gobam "github.com/cybersiddhu/biogo.boom"
	"log"
	"os"
	"runtime"
	"sync"
)

//default bin size in 1
var wg sync.WaitGroup

func main() {
	binSize := 1
	dumpSize := 100000
	runtime.GOMAXPROCS(4)

	bam, err := gobam.OpenBAM(os.Args[1])
	dieIfError(err)

	idx, err := gobam.LoadIndex(os.Args[1])
	dieIfError(err)

	lengths := bam.RefLengths()

	for i, name := range bam.RefNames() {
		if id, ok := bam.RefID(name); ok {
			length := int(lengths[i])
			wg.Add(1)
			go ProcessTarget(bam, idx, dumpSize, binSize, length, name, id)
		}
	}
	wg.Wait()
}

func dieIfError(e error) {
	if e != nil {
		panic(e)
	}
}

func ProcessTarget(bam *gobam.BAMFile, idx *gobam.Index, dumpSize int, binSize int, length int, name string, id int) {
	log.Printf("going to write %s\n", name)
	defer wg.Done()
	w, err := os.Create(name + ".wig")
	dieIfError(err)
	fmt.Fprintf(w, "fixedStep chrom=%s start=1 step=%d span=%d\n", name, binSize, binSize)
	defer w.Close()

	for start := 0; start < length; start += dumpSize {
		// calculate end 0 based indexing
		end := start + dumpSize - 1
		if end > length {
			end = length
		}
		coverage, err := idx.Coverage(bam, id, start, end)
		dieIfError(err)
		for _, reads := range coverage {
			fmt.Fprintf(w, "%d\n", reads)
		}
	}
	log.Printf("Finished writing %s\n", name)
}
