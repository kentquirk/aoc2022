package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func part1(lines []string) {
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

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	part1(lines)
}
