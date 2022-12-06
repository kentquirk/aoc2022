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

type Command struct {
	Qty  int
	FrIx int
	ToIx int
}

type Cargo struct {
	stacks   []Stack
	commands []Command
}

func Parse(lines []string) *Cargo {
	// we're gonna use a Stack to parse the cargo
	// because we need to iterate it in reverse
	items := NewStack()
	var nstacks int
	var ix int
	for _, l := range lines {
		ix++
		if !strings.HasPrefix(strings.TrimSpace(l), "[") {
			// we have a line of numbers here, grab the last one
			nstacks = int(l[len(l)-1] - '0')
			break
		}
		items.Push(l)
	}

	cargo := &Cargo{
		stacks:   make([]Stack, nstacks),
		commands: make([]Command, 0),
	}

	for l, ok := items.Pop(); ok; l, ok = items.Pop() {
		for i := 0; i < nstacks; i++ {
			if len(l) > i*4+1 {
				crate := l[i*4+1 : i*4+2]
				if crate != " " {
					cargo.stacks[i].Push(crate)
				}
			}
		}
	}

	pat := regexp.MustCompile("move ([0-9]+) from ([0-9]) to ([0-9])")
	for ; ix < len(lines); ix++ {
		m := pat.FindStringSubmatch(lines[ix])
		if len(m) == 4 {
			qty, _ := strconv.Atoi(m[1])
			fr, _ := strconv.Atoi(m[2])
			to, _ := strconv.Atoi(m[3])
			cargo.commands = append(cargo.commands, Command{Qty: qty, FrIx: fr - 1, ToIx: to - 1})
		}
	}
	return cargo
}

func (c *Cargo) MoveSingle(cmd Command) {
	for i := 0; i < cmd.Qty; i++ {
		if crate, ok := c.stacks[cmd.FrIx].Pop(); ok {
			c.stacks[cmd.ToIx].Push(crate)
		}
	}
}

func (c *Cargo) ExecSingle() {
	for _, cmd := range c.commands {
		c.MoveSingle(cmd)
	}
}

func (c *Cargo) Tops() string {
	tops := ""
	for _, s := range c.stacks {
		if t, ok := s.Top(); ok {
			tops += t
		}
	}
	return tops
}

type Stack struct {
	s []string
}

func NewStack() *Stack {
	return &Stack{
		s: make([]string, 0),
	}
}

func (s *Stack) Push(str string) {
	s.s = append(s.s, str)
}

func (s *Stack) Pop() (string, bool) {
	if len(s.s) > 0 {
		r := s.s[len(s.s)-1]
		s.s = s.s[:len(s.s)-1]
		return r, true
	}
	return "", false
}

func (s *Stack) Top() (string, bool) {
	if len(s.s) > 0 {
		r := s.s[len(s.s)-1]
		return r, true
	}
	return "", false
}

func (s *Stack) String() string {
	return strings.Join(s.s, "|")
}

func part1(lines []string) {
	cargo := Parse(lines)
	cargo.ExecSingle()
	fmt.Println(cargo.Tops())
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
}
