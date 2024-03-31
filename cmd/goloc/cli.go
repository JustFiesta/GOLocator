package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.String("test", "", "test example")
	fmt.Println("Hello world!")
}
