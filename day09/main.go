package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

func (p Point) Add(other Point) Point {
	return Point{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

func (p Point) Sub(other Point) Point {
	return Point{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}

func (p Point) Adjacent(other Point) bool {
	d := p.Sub(other)
	// fmt.Printf("%v %v %v %t\n", p, other, d, d.X <= 1 && d.X >= -1 && d.Y <= 1 && d.Y >= -1)
	return d.X <= 1 && d.X >= -1 && d.Y <= 1 && d.Y >= -1
}

type Move struct {
	Delta Point
	Count int
}

// Head is Nodes[0]
type Rope struct {
	Nodes         []Point
	TailPositions map[Point]struct{}
}

func NewRope(length int) *Rope {
	return &Rope{
		Nodes: make([]Point, length),
		// put first tail position into the map
		TailPositions: map[Point]struct{}{{}: struct{}{}},
	}
}

func (r *Rope) MoveOne(delta Point) {
	r.Nodes[0] = r.Nodes[0].Add(delta)
	for i := 1; i < len(r.Nodes); i++ {
		if r.Nodes[i-1].Adjacent(r.Nodes[i]) {
			break // nothing else moves
		}
		// figure out where to move the next node
		diff := r.Nodes[i-1].Sub(r.Nodes[i])
		if diff.X > 1 {
			diff.X = 1
		}
		if diff.X < -1 {
			diff.X = -1
		}
		if diff.Y > 1 {
			diff.Y = 1
		}
		if diff.Y < -1 {
			diff.Y = -1
		}
		r.Nodes[i] = r.Nodes[i].Add(diff)
	}
	// fmt.Println(r)
	r.TailPositions[r.Nodes[len(r.Nodes)-1]] = struct{}{}
}

func (r *Rope) ExecuteMoves(moves []Move) {
	for _, m := range moves {
		for n := 0; n < m.Count; n++ {
			r.MoveOne(m.Delta)
		}
		// fmt.Println(r)
	}
}

func Parse(lines []string) []Move {
	moves := make([]Move, 0)
	for _, l := range lines {
		splits := strings.Split(l, " ")
		d := splits[0]
		n, _ := strconv.Atoi(splits[1])
		switch d {
		case "R":
			moves = append(moves, Move{Delta: Point{X: 1, Y: 0}, Count: n})
		case "L":
			moves = append(moves, Move{Delta: Point{X: -1, Y: 0}, Count: n})
		case "U":
			moves = append(moves, Move{Delta: Point{X: 0, Y: 1}, Count: n})
		case "D":
			moves = append(moves, Move{Delta: Point{X: 0, Y: -1}, Count: n})
		default:
			panic("oops")
		}
	}
	return moves
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
	moves := Parse(lines)
	rope := NewRope(2)
	rope.ExecuteMoves(moves)
	fmt.Println(len(rope.TailPositions))

	longrope := NewRope(10)
	longrope.ExecuteMoves(moves)
	fmt.Println(len(longrope.TailPositions))
}
