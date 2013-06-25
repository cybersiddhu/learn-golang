package main;

import "fmt"

func main() {
	 b := []byte("hello")
	 fmt.Printf("first %d\n",len(b))
	 
	 b = []byte{}
	 fmt.Printf("second %d\n",len(b))
}
