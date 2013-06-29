package main

import (
	"fmt"
	"github.com/cybersiddhu/gobio/seqio"
	"os"
)

func main() {
	r := seqio.NewFastaReader(os.Args[1])
	for r.HasEntry() {
		f := r.NextEntry()
		fmt.Fprintf(os.Stdout, "Entry>>>>\n Id:%s\tsequence:%s\n", f.Id, f.Sequence)
	}
}
