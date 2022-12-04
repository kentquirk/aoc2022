package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Range struct {
	Min int
	Max int
}

func (r *Range) Contains(other *Range) bool {
	return r.Min <= other.Min && r.Max >= other.Max
}

// if both of one are outside the bounds of the other one;
// we don't have to do a symmetric test.
func (r *Range) Disjoint(other *Range) bool {
	return r.Max < other.Min || r.Min > other.Max
}

func parse(line string) (*Range, *Range) {
	pat := regexp.MustCompile(`([0-9]+)-([0-9]+),([0-9]+)-([0-9]+)`)
	all := pat.FindStringSubmatch(line)
	toint := func(s string) int {
		n, _ := strconv.Atoi(s)
		return n
	}
	return &Range{Min: toint(all[1]), Max: toint(all[2])}, &Range{Min: toint(all[3]), Max: toint(all[4])}
}

func part1(lines []string) {
	containsCount := 0
	for _, l := range lines {
		r1, r2 := parse(l)
		if r1.Contains(r2) || r2.Contains(r1) {
			containsCount++
		}
	}
	fmt.Println(containsCount)
}

func part2(lines []string) {
	overlapCount := 0
	for _, l := range lines {
		r1, r2 := parse(l)
		if !r1.Disjoint(r2) {
			overlapCount++
		}
	}
	fmt.Println(overlapCount)
}

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
	part1(lines)
	part2(lines)
}
