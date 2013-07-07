package main

import (
	"fmt"
	gobam "github.com/cybersiddhu/biogo.boom"
	"os"
)

//default bin size in 1
var binSize = 1
var dumpSize = 10000

func main() {

	//log.Println("starting up")
	bam, err := gobam.OpenBAM(os.Args[1])
	dieIfError(err)
	//log.Println("got bam reader")

	//log.Println("before loading index")
	idx, err := gobam.LoadIndex(os.Args[1])
	dieIfError(err)

	lengths := bam.RefLengths()

	for i, name := range bam.RefNames() {
		if id, ok := bam.RefID(name); ok {
			length := int(lengths[i])
			w, err := os.Create(name + ".wig")
			dieIfError(err)
			fmt.Fprintf(w, "fixedStep chrom=%s start=1 step=%d span=%d\n", name, binSize, binSize)

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
		}
	}
}

func dieIfError(e error) {
	if e != nil {
		panic(e)
	}
}
