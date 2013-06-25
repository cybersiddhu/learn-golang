package main

import (
	 "fmt"
	 "regexp"
)

func main() {
	 h := []byte(">tora bora")
	 c := regexp.MustCompile(`^>(\S+)`)
	 if m := c.FindSubmatch(h); m != nil {
	 		fmt.Printf("match %s\n",m[1])
	 }
}
