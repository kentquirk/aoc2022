package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
)

type Loc struct {
	R int
	C int
}

type Elf struct {
	Proposal *Loc
}

type Field struct {
	Elves     map[Loc]*Elf
	PrevElves map[Loc]*Elf
	Moveorder string
}

func (f *Field) Bounds() (Loc, Loc) {
	min := Loc{math.MaxInt, math.MaxInt}
	max := Loc{math.MinInt, math.MinInt}

	for loc := range f.Elves {
		if loc.R < min.R {
			min.R = loc.R
		}
		if loc.R > max.R {
			max.R = loc.R
		}
		if loc.C < min.C {
			min.C = loc.C
		}
		if loc.C > max.C {
			max.C = loc.C
		}
	}
	return min, max
}

func (f *Field) NEmpty() int {
	min, max := f.Bounds()
	area := (max.R - min.R + 1) * (max.C - min.C + 1)
	return area - len(f.Elves)
}

func (f *Field) HasElf(r, c int) bool {
	_, ok := f.Elves[Loc{R: r, C: c}]
	return ok
}

func (f *Field) Neighbors(loc Loc) map[string]struct{} {
	m := make(map[string]struct{})
	if f.HasElf(loc.R-1, loc.C-1) {
		m["NW"] = struct{}{}
	}
	if f.HasElf(loc.R-1, loc.C) {
		m["N"] = struct{}{}
	}
	if f.HasElf(loc.R-1, loc.C+1) {
		m["NE"] = struct{}{}
	}
	if f.HasElf(loc.R, loc.C-1) {
		m["W"] = struct{}{}
	}
	if f.HasElf(loc.R, loc.C+1) {
		m["E"] = struct{}{}
	}
	if f.HasElf(loc.R+1, loc.C-1) {
		m["SW"] = struct{}{}
	}
	if f.HasElf(loc.R+1, loc.C) {
		m["S"] = struct{}{}
	}
	if f.HasElf(loc.R+1, loc.C+1) {
		m["SE"] = struct{}{}
	}
	return m
}

func isClear(neighbors map[string]struct{}, direction rune) bool {
	for n := range neighbors {
		if strings.ContainsRune(n, direction) {
			return false
		}
	}
	return true
}

func (f *Field) Generation() int {
	offsets := map[rune]Loc{
		'N': {R: -1, C: 0},
		'S': {R: 1, C: 0},
		'W': {R: 0, C: -1},
		'E': {R: 0, C: 1},
	}
	proposals := make(map[Loc]int)

nextelf:
	// first half, all the elves make proposals
	for loc, elf := range f.Elves {
		neighbors := f.Neighbors(loc)
		elf.Proposal = nil
		if len(neighbors) == 0 {
			continue
		}
		for _, dir := range f.Moveorder {
			if isClear(neighbors, dir) {
				prop := &Loc{R: loc.R + offsets[dir].R, C: loc.C + offsets[dir].C}
				elf.Proposal = prop
				proposals[*prop]++
				continue nextelf
			}
		}
	}

	if len(proposals) == 0 {
		return 0
	}
	// second half, everyone moves if they won't collide
	moveCount := 0
	newElves := make(map[Loc]*Elf)
	for loc, elf := range f.Elves {
		if elf.Proposal != nil && proposals[*elf.Proposal] == 1 {
			if _, ok := newElves[*elf.Proposal]; ok {
				fmt.Println("proposal was on another elf!", loc, *elf.Proposal)
			}
			newElves[*elf.Proposal] = elf
			moveCount++
		} else {
			if _, ok := newElves[loc]; ok {
				fmt.Println("unmoved elf overwrote!", loc)
			}
			newElves[loc] = elf
		}
	}
	f.PrevElves = f.Elves
	f.Elves = newElves
	// now rotate the ordering
	f.Moveorder = f.Moveorder[1:] + f.Moveorder[:1]
	return moveCount
}

func (f *Field) Print(g int, min Loc, max Loc, animate bool) {
	bmin, bmax := f.Bounds()
	if bmin.C < min.C || bmin.R < min.R {
		min = bmin
	}
	if bmax.C > max.C || bmax.R > max.R {
		max = bmax
	}

	if animate {
		fmt.Print("\x1b[0;0H")
	}
	fmt.Printf("---- Generation %d ---- %v %v (%d elves)\n", g, bmin, bmax, len(f.Elves))
	for r := min.R; r <= max.R; r++ {
		if animate {
			fmt.Print("\x1b[K")
		}
		for c := min.C; c <= max.C; c++ {
			if r >= bmin.R && r <= bmax.R && c >= bmin.C && c <= bmax.C {
				fmt.Print("\x1b[41m")
			}
			loc := Loc{R: r, C: c}
			if _, ok := f.Elves[loc]; ok {
				fmt.Print("#")
			} else {
				if _, ok := f.PrevElves[loc]; ok {
					fmt.Print("\x1b[44m.\x1b[41m")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Print("\x1b[40m")
		}
		fmt.Println()
	}
	fmt.Println()
}

func Parse(lines []string) *Field {
	f := &Field{
		Elves:     make(map[Loc]*Elf),
		Moveorder: "NSWE",
	}
	elfpat := regexp.MustCompile("#")
	for row, l := range lines {
		elfCols := elfpat.FindAllStringIndex(l, -1)
		for _, e := range elfCols {
			f.Elves[Loc{C: e[0], R: row}] = &Elf{}
		}
	}
	return f
}

func part1(lines []string, maxgenerations int, animate bool) int {
	field := Parse(lines)
	if animate {
		fmt.Print("\x1b[2J")
	}
	min, max := field.Bounds()
	min.C -= 3
	min.R -= 2
	max.C += 4
	max.R += 3
	for g := 0; g < maxgenerations; g++ {
		field.Print(g, min, max, animate)
		if field.Generation() == 0 {
			break
		}
		if animate {
			time.Sleep(1000 * time.Millisecond)
		}
	}
	field.Print(maxgenerations, min, max, animate)
	return field.NEmpty()
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
	fmt.Println(part1(lines, 10, true))
}
