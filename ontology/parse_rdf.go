package main

import (
	"bitbucket.org/ww/goraptor"
	"fmt"
	"os"
)

func main() {
	parser := goraptor.NewParser("rdfxml")
	defer parser.Free()

	ch := parser.ParseFile(os.Args[1], "")
	for {
		statement, ok := <-ch
		if !ok {
			break
		}
		fmt.Printf("Subject:%s Predicate:%s Object:%s\n", statement.Subject.String(), statement.Predicate.String(), statement.Object.String())
	}
}
