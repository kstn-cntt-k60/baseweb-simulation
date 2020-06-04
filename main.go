package main

import (
	importProduct "baseweb-simulation/import"
	"baseweb-simulation/person"
	"baseweb-simulation/product"
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

	if command == "add-person" {
		person.AddPersonBenchmark(loop)
	} else if command == "add-product" {
		product.AddProductBenchmark(loop)
	} else if command == "add-inventory-item" {
		importProduct.AddInventoryItemBenchmark(loop)
	} else if command == "add-many-items" {
		for i := 0; i < loop; i++ {
			importProduct.AddInventoryItemBenchmark(200)
		}
	} else {
		fmt.Println("ERROR: command not supported")
		fmt.Println("Using one of the following commands:")
		fmt.Println("    - add-person")
		fmt.Println("    - add-product")
		fmt.Println("    - add-inventory-item")
	}
}
