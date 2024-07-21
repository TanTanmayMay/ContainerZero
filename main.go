package main

import (
	"fmt"
	"os"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("Wrong Command")
	}
}

func run() {
	fmt.Printf("Running %v\n", os.Args[2:])
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}