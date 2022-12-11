package main

import (
	"bytes"
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

type CRT struct {
	Pixels [240]byte
}

func (c *CRT) Display() {
	for r := 0; r < 6; r++ {
		b := c.Pixels[r*40 : (r+1)*40]
		s := string(bytes.Replace(b, []byte{0}, []byte(" "), -1))
		fmt.Println(s)
	}
}

func (c *CRT) Draw(cycle int, x int) {
	offset := (cycle / 40) * 40
	fmt.Printf("cycle %d, offset %d, x %d", cycle, offset, x)
	if offset+x >= cycle-2 && offset+x <= cycle+0 {
		c.Pixels[cycle-1] = byte('#')
		fmt.Printf("-- draw at %d", cycle-1)
	}
	fmt.Println()
}

func part2(lines []string) {
	vm := NewVM()
	vm.Load(lines)
	vm.Reset()
	crt := CRT{}
	for vm.Tick() {
		// fmt.Printf("%s\n", vm)
		crt.Draw(vm.Ticks, vm.LastX)
	}
	crt.Display()
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
	part2(lines)
}
