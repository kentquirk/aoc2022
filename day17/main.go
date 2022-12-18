package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dgryski/go-wyhash"
)

type RockSpace struct {
	Contents []byte
	Width    int
}

// converts a string like "..#.#" to a left-justified byte based on the original
// string length (in this case, the result is 00101000)
func toBinary(s string) byte {
	r := strings.NewReplacer(".", "0", "#", "1")
	s = r.Replace(s)
	b, _ := strconv.ParseUint(s, 2, 8)
	return byte(b << (8 - len(s)))
}

func RS(a ...string) RockSpace {
	rs := RockSpace{Width: len(a[0])}
	for _, s := range a {
		rs.Contents = append(rs.Contents, toBinary(s))
	}
	return rs
}

// Returns the height of the highest row that contains a nonempty cell
func (r *RockSpace) Height() int {
	for i := len(r.Contents) - 1; i >= 0; i-- {
		if r.Contents[i] != 0 {
			return i + 1
		}
	}
	return 0
}

// Calculates a hash for the top of the tower by accumulating it for
// the values of the top slices until there has been a rock in every
// location. The hash also includes the current wind and rock indices; when
// the hash repeats, we've got a repeated situation.
func (r *RockSpace) Hash(windIx int, rockIx int) uint64 {
	key := []byte{byte((windIx >> 8) & 0xFF), byte(windIx & 0xFF), byte(rockIx)}
	var allbits byte = 0
	for i := len(r.Contents) - 1; i >= 0; i-- {
		b := r.Contents[i]
		if b == 0 {
			continue
		}
		key = append(key, b)
		allbits |= b
		if allbits&0xFE == 0xFE {
			break
		}
	}
	return wyhash.Hash(key, 0x14534fe78bc)
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
		rb := rock.Contents[rockY] >> x
		if rb&r.Contents[spaceY] != 0 {
			return true
		}
	}
	return false
}

// Places a rock that's passed the collide function already
func (r *RockSpace) Place(rock RockSpace, x int, y int) {
	for y+len(rock.Contents) >= len(r.Contents) {
		r.Contents = append(r.Contents, 0)
	}
	for rockY := 0; rockY < len(rock.Contents); rockY++ {
		r.Contents[y+rockY] |= rock.Contents[rockY] >> x
	}
}

func (r *RockSpace) Print() {
	for i := len(r.Contents) - 1; i >= 0; i-- {
		b := r.Contents[i]
		rep := strings.NewReplacer("0", ".", "1", "#")
		s := fmt.Sprintf("%08b", b)
		fmt.Printf("%s\n", rep.Replace(s)[:r.Width])
	}
	fmt.Printf("%s\n\n", strings.Repeat("-", r.Width))
}

var rocks = []RockSpace{
	RS("####"),
	RS(".#.", "###", ".#."),
	RS("###", "..#", "..#"), // upside down since we measure contents from bottom up
	RS("#", "#", "#", "#"),
	RS("##", "##"),
}

type hashState struct {
	RockCount int
	Height    int
}

func part1(breeze string, iterations int) int {
	hashes := make(map[uint64]hashState)
	var heightOffset = 0

	chamber := RockSpace{Width: 7}
	windIx := 0
	doHashes := true
	for rockCount := 0; rockCount < iterations; rockCount++ {
		rockIx := rockCount % len(rocks)
		if doHashes {
			h := chamber.Hash(windIx, rockIx)
			state := hashState{RockCount: rockCount, Height: chamber.Height()}
			if prevState, ok := hashes[h]; ok {
				nrocks := state.RockCount - prevState.RockCount
				deltaHeight := state.Height - prevState.Height
				fmt.Printf("Dup at %v, prev %v, nrocks = %d, deltaH = %d\n", state, prevState, nrocks, deltaHeight)
				iterationsLeft := iterations - rockCount
				rockCount += (iterationsLeft / nrocks) * nrocks
				heightOffset = deltaHeight * (iterationsLeft / nrocks)
				doHashes = false
			}
			hashes[h] = state
		}

		rock := rocks[rockIx]
		dropY := chamber.Height() + 3
		dropX := 2
		for {
			wind := breeze[windIx]
			windIx++
			windIx = windIx % len(breeze)
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
	return chamber.Height() + heightOffset
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
	fmt.Println(part1(lines[0], 1_000_000_000_000))
}
