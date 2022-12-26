package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var debug bool = false

type Monkey interface {
	Yell(depth int) int
}

type MathMonkey struct {
	Name     string
	M1       string
	M2       string
	Operator string
	barrel   map[string]Monkey
}

func (mm MathMonkey) Yell(depth int) int {
	if debug {
		fmt.Printf("%s - M(%s) = %s %s %s\n", strings.Repeat("  ", depth), mm.Name, mm.M1, mm.Operator, mm.M2)
	}
	m1 := mm.barrel[mm.M1].Yell(depth + 1)
	m2 := mm.barrel[mm.M2].Yell(depth + 1)
	result := 0
	switch mm.Operator {
	case "+":
		result = m1 + m2
	case "*":
		result = m1 * m2
	case "-":
		result = m1 - m2
	case "/":
		result = m1 / m2
	default:
		panic("bad operator!")
	}
	if debug {
		fmt.Printf("%s - M(%s) result = %d\n", strings.Repeat("  ", depth), mm.Name, result)
	}
	return result
}

func (mm MathMonkey) Test() bool {
	fmt.Printf("   - M(%s) Testing %s == %s\n", mm.Name, mm.M1, mm.M2)
	m1 := mm.barrel[mm.M1].Yell(1)
	m2 := mm.barrel[mm.M2].Yell(1)
	fmt.Printf("   - M(%s) result = %t\n", mm.Name, m1 == m2)
	return m1 == m2
}

type NumberMonkey struct {
	Name  string
	Value int
}

func (nm NumberMonkey) Yell(depth int) int {
	if debug {
		fmt.Printf("%s - N(%s) = %d\n", strings.Repeat("  ", depth), nm.Name, nm.Value)
	}
	return nm.Value
}

type Human struct {
	Name  string
	Root  MathMonkey
	Guess int
}

// Humans still fulfill the Monkey interface
func (h Human) Yell(depth int) int {
	if debug {
		fmt.Printf("%s - N(%s) = %d\n", strings.Repeat("  ", depth), h.Name, h.Guess)
	}
	return h.Guess
}

func (h Human) Try(v int) int {
	h.Guess = v
	h.Root.barrel[h.Name] = h
	m1 := h.Root.barrel[h.Root.M1].Yell(1)
	m2 := h.Root.barrel[h.Root.M2].Yell(1)
	diff := m2 - m1
	return diff
}

// using Newton's method to iterate on a solution
func (h Human) Iterate() int {
	// calculate the slope between f(x) and the origin
	// as the first slope
	// picking x=100 as our first guess
	x := 10

	limit := 2000
	for i := 0; i < limit; i++ {
		if i%100 == 0 {
			fmt.Println("trying ", x)
		}
		// slope is dy/dx, so we take two readings at our sample point exactly 1.0 apart
		// so that dx is always exactly 1
		y1 := h.Try(x)
		if y1 == 0 {
			return x
		}
		y2 := h.Try(x + 1)
		if y2 == 0 {
			return x + 1
		}
		// if we get the same value for two adjacent tests that isn't useful
		if y1 == y2 {
			x++
			continue
		}
		slope := y2 - y1
		// now calculate our next trial
		x = x - int(float64(y1)/float64(slope))
	}
	fmt.Println("hit the limit")
	return -1
}

func BuildBarrel(lines []string) map[string]Monkey {
	barrel := make(map[string]Monkey)
	mathpat := regexp.MustCompile("([a-z]{4}): ([a-z]{4}) (.) ([a-z]{4})")
	numpat := regexp.MustCompile("([a-z]{4}): ([0-9]+)")
	for _, l := range lines {
		if numpat.MatchString(l) {
			parts := numpat.FindStringSubmatch(l)
			n, _ := strconv.Atoi(parts[2])
			barrel[parts[1]] = NumberMonkey{Name: parts[1], Value: n}
		} else {
			parts := mathpat.FindStringSubmatch(l)
			m := MathMonkey{
				Name:     parts[1],
				M1:       parts[2],
				M2:       parts[4],
				Operator: parts[3],
				barrel:   barrel,
			}
			barrel[parts[1]] = m
		}
	}
	return barrel
}

func part1(lines []string) int {
	barrel := BuildBarrel(lines)
	return barrel["root"].Yell(1)
}

func part2(lines []string) int {
	barrel := BuildBarrel(lines)
	root := barrel["root"].(MathMonkey)
	h := Human{
		Name: "humn",
		Root: root,
	}
	return h.Iterate()
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
	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}
