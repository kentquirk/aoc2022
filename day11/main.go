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

func MulOp(operand int) func(int) int {
	return func(v int) int {
		return v * operand
	}
}

func SquareOp() func(int) int {
	return func(v int) int {
		return v * v
	}
}

func AddOp(operand int) func(int) int {
	return func(v int) int {
		return v + operand
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
		item.Worry = m.Op(item.Worry) / 3
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

func buildOp(s string) func(int) int {
	oppat := regexp.MustCompile(`new = old (\*|\+) (old|[0-9]+)`)
	parts := oppat.FindStringSubmatch(s)
	if parts[2] == "old" {
		return SquareOp()
	}

	n, _ := strconv.Atoi(parts[2])
	if parts[1] == "*" {
		return MulOp(n)
	}
	return AddOp(n)
}

func (p *Pandemonium) AddMonkey(setup []string) {
	monkeyID := getNumbers(setup[0])[0]
	if len(p.Monkeys) != monkeyID {
		panic("wrong monkey!")
	}
	m := &Monkey{}
	for _, w := range getNumbers(setup[1]) {
		m.Items = append(m.Items, &Item{Worry: w})
	}
	m.Op = buildOp(setup[2])
	m.Test = DivTest(getNumbers(setup[3])[0])
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

func Parse(lines []string) *Pandemonium {
	p := &Pandemonium{}
	for len(lines) > 6 {
		p.AddMonkey(lines[:7])
		lines = lines[7:]
	}
	return p
}

func main() {
	f, err := os.Open("./inputsample.txt")
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	p := Parse(lines)
	for i := 0; i < 20; i++ {
		// p.Print()
		p.Round()
	}
	p.Print()
	fmt.Println(p.MonkeyBusiness())
}
