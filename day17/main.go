package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type RockSpace struct {
	Contents [][]byte
	Width    int
}

func RS(a ...string) RockSpace {
	rs := RockSpace{Width: len(a[0])}
	for _, s := range a {
		rs.Contents = append(rs.Contents, []byte(s))
	}
	return rs
}

// Returns the height of the highest row that contains a nonempty cell
func (r *RockSpace) Height() int {
	for i := len(r.Contents) - 1; i >= 0; i-- {
		for _, b := range r.Contents[i] {
			if b != byte('.') {
				return i + 1
			}
		}
	}
	return 0
}

// returns true if a rock in the given position intersects the sides, floor, or
// an occupied square in the space. Position of rock is measured with its lower
// left corner at (0,0) in the rock's space, and the lower left of the container
// is also (0,0).
func (r *RockSpace) Collides(rock RockSpace, x int, y int) bool {
	if x < 0 || y < 0 {
		return true
	}
	if x+rock.Width > r.Width {
		return true
	}
	if y > len(r.Contents) {
		return false
	}
	for rockY := 0; rockY < len(rock.Contents); rockY++ {
		spaceY := y + rockY
		if spaceY >= len(r.Contents) {
			// if it hasn't collided yet, it can't possibly do it now
			return false
		}
		for rockX := 0; rockX < len(rock.Contents[rockY]); rockX++ {
			if rock.Contents[rockY][rockX] == byte('#') && r.Contents[spaceY][x+rockX] == byte('#') {
				return true
			}
		}
	}
	return false
}

// Places a rock that's passed the collide function already
func (r *RockSpace) Place(rock RockSpace, x int, y int) {
	for y+len(rock.Contents) >= len(r.Contents) {
		r.Contents = append(r.Contents, bytes.Repeat([]byte{'.'}, r.Width))
	}
	for rockY := 0; rockY < len(rock.Contents); rockY++ {
		for rockX := 0; rockX < len(rock.Contents[rockY]); rockX++ {
			if rock.Contents[rockY][rockX] == byte('#') {
				r.Contents[y+rockY][x+rockX] = byte('#')
			}
		}
	}
}

func (r *RockSpace) Print() {
	for i := len(r.Contents) - 1; i >= 0; i-- {
		fmt.Printf("%s\n", r.Contents[i])
	}
	fmt.Printf("%s\n", strings.Repeat("-", r.Width))
}

var rocks = []RockSpace{
	RS("####"),
	RS(".#.", "###", ".#."),
	RS("###", "..#", "..#"), // upside down since we measure contents from bottom up
	RS("#", "#", "#", "#"),
	RS("##", "##"),
}

func part1(breeze string, iterations int) int {
	chamber := RockSpace{Width: 7}
	windIx := 0
	for rockIx := 0; rockIx < iterations; rockIx++ {
		rock := rocks[rockIx%len(rocks)]
		dropY := chamber.Height() + 3
		dropX := 2
		for {
			wind := breeze[windIx%len(breeze)]
			windIx++
			switch wind {
			case '<':
				if !chamber.Collides(rock, dropX-1, dropY) {
					dropX--
				}
			case '>':
				if !chamber.Collides(rock, dropX+1, dropY) {
					dropX++
				}
			default:
				panic("GAAAAH")
			}
			if !chamber.Collides(rock, dropX, dropY-1) {
				dropY--
			} else {
				chamber.Place(rock, dropX, dropY)
				break
			}
		}
		// chamber.Print()
	}
	return chamber.Height()
}

func part2(breeze string) int {
	const iterations = 1_000_000_000_000
	// we need the least common multiple between the length of breeze
	// and 5, the number of rocks we're repeatedly dropping.
	// We need to drop that many rocks, measure the height, drop
	// that many again, measure the change in height, and then figure out where
	// the appropriate number of iterations will land us based on the modulus of
	// iterations and the lcm.
	lcm := len(breeze)
	if lcm%len(rocks) != 0 {
		lcm *= len(rocks)
	}
	if lcm%7 != 0 { // width of the chamber
		lcm *= 7
	}
	for i := 0; i < 10; i++ {
		h1 := part1(breeze, lcm*(i+1))
		h2 := part1(breeze, lcm*(i+2))
		delta := h2 - h1
		fmt.Println(h1, h2, delta)
	}
	return lcm
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
	fmt.Println(part1(lines[0], 2022))
	fmt.Println(part2(lines[0]))
}
