package main

import (
	"fmt"
	gobam "github.com/cybersiddhu/biogo.boom"
	"os"
)

binSize := 1

func main() {
	bam, err := gobam.OpenBAM(os.Args[1])
	if err != nil {
		panic(err)
	}

	for i, name := range bam.RefNames() {
		lengths := bam.RefLengths()
		if id, ok := bam.RefID(name); ok {
			fmt.Printf("refname:%s\ttid:%d\tlength:%d\n", name, id, lengths[i])

		}
	}

}
