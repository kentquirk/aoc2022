package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func MulOp(operand int, worryLess func(int) int) func(int) int {
	return func(v int) int {
		return worryLess(v * operand)
	}
}

func SquareOp(worryLess func(int) int) func(int) int {
	return func(v int) int {
		return worryLess(v * v)
	}
}

func AddOp(operand int, worryLess func(int) int) func(int) int {
	return func(v int) int {
		return worryLess(v + operand)
	}
}

func DivTest(operand int) func(int) bool {
	return func(v int) bool {
		return v%operand == 0
	}
}

type Item struct {
	Worry int
}

func (i *Item) String() string {
	return strconv.Itoa(i.Worry)
}

type Monkey struct {
	Items     []*Item
	Op        func(int) int
	Test      func(int) bool
	Throw     func(*Item, int)
	TrueDest  int
	FalseDest int
	NInspect  int
}

func (m *Monkey) Evaluate() {
	for _, item := range m.Items {
		item.Worry = m.Op(item.Worry)
		if m.Test(item.Worry) {
			m.Throw(item, m.TrueDest)
		} else {
			m.Throw(item, m.FalseDest)
		}
		m.NInspect++
	}
	m.Items = make([]*Item, 0)
}

// Pandemonium is the collective name for a group of flying monkeys.
type Pandemonium struct {
	Monkeys []*Monkey
	lcm     int
}

func (p *Pandemonium) Throw(item *Item, newMonkey int) {
	p.Monkeys[newMonkey].Items = append(p.Monkeys[newMonkey].Items, item)
}

func getNumbers(s string) []int {
	var nums []int
	numpat := regexp.MustCompile("[0-9]+")
	a := numpat.FindAllString(s, -1)
	for _, v := range a {
		n, _ := strconv.Atoi(v)
		nums = append(nums, n)
	}
	return nums
}

func buildOp(s string, worryLess func(int) int) func(int) int {
	oppat := regexp.MustCompile(`new = old (\*|\+) (old|[0-9]+)`)
	parts := oppat.FindStringSubmatch(s)
	if parts[2] == "old" {
		return SquareOp(worryLess)
	}

	n, _ := strconv.Atoi(parts[2])
	if parts[1] == "*" {
		return MulOp(n, worryLess)
	}
	return AddOp(n, worryLess)
}

func (p *Pandemonium) AddMonkey(setup []string, worryLess func(int) int) {
	monkeyID := getNumbers(setup[0])[0]
	if len(p.Monkeys) != monkeyID {
		panic("wrong monkey!")
	}
	m := &Monkey{}
	for _, w := range getNumbers(setup[1]) {
		m.Items = append(m.Items, &Item{Worry: w})
	}
	m.Op = buildOp(setup[2], worryLess)
	modulus := getNumbers(setup[3])[0]
	p.lcm *= modulus
	m.Test = DivTest(modulus)
	m.TrueDest = getNumbers(setup[4])[0]
	m.FalseDest = getNumbers(setup[5])[0]
	m.Throw = p.Throw
	p.Monkeys = append(p.Monkeys, m)
}

func (p *Pandemonium) Round() {
	for _, m := range p.Monkeys {
		m.Evaluate()
	}
}
func (p *Pandemonium) Print() {
	fmt.Println("---")
	for i, m := range p.Monkeys {
		fmt.Printf("Monkey %d has %v, NInspect=%d\n", i, m.Items, m.NInspect)
	}
}

func (p *Pandemonium) MonkeyBusiness() int {
	var biz []int
	for _, m := range p.Monkeys {
		biz = append(biz, m.NInspect)
	}
	sort.Ints(biz)
	return biz[len(biz)-2] * biz[len(biz)-1]
}

func (p *Pandemonium) WorryMod(x int) int {
	return x % p.lcm
}

func Parse(lines []string, worryLess func(int) int) *Pandemonium {
	p := &Pandemonium{lcm: 1}
	for len(lines) > 6 {
		p.AddMonkey(lines[:7], worryLess)
		lines = lines[7:]
	}
	return p
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
	p := Parse(lines, func(x int) int { return x / 3 })
	for i := 0; i < 20; i++ {
		// p.Print()
		p.Round()
	}
	p.Print()
	fmt.Println(p.MonkeyBusiness())

	p = Parse(lines, p.WorryMod)
	for i := 0; i < 10000; i++ {
		// p.Print()
		p.Round()
	}
	p.Print()
	fmt.Println(p.MonkeyBusiness())
}
