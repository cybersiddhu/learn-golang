package main

import (
	 "github.com/cybersiddhu/gobio/seqio"
	 "os"
	 "fmt"
)


func main() {
	r := seqio.NewFastaReader(os.Args[1])
	for {
		f, err := r.NextSeq()
		if err != nil {
			break
		}
		fmt.Fprintf(os.Stdout, "id:%s\tsequence:%s\n", f.Id, f.Sequence)
	}
}
