package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	vm := NewVM()
	vm.Load(lines)
	vm.Reset()
	sum := 0
	for vm.Tick() {
		fmt.Printf("%s\n", vm)
		if vm.Ticks%40 == 20 {
			signalStrength := vm.Ticks * vm.LastX
			sum += signalStrength
			fmt.Printf("-- sig: %d, sum: %d\n", signalStrength, sum)
		}
	}
	fmt.Println(sum)
}
