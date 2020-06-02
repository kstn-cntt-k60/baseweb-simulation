package main

import (
	"baseweb-simulation/person"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	log.Println("Num CPU:", runtime.NumCPU())

	command := ""
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}

	loopStr := os.Getenv("LOOP")
	loop, err := strconv.Atoi(loopStr)
	if err != nil {
		loop = 10
	}

	if command == "help" {
		fmt.Println("add-person")
	} else if command == "add-person" {
		person.AddPersonBenchmark(loop)
	} else {
		fmt.Println("ERROR: command not supported")
	}
}
