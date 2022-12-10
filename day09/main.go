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
	return d.X <= 1 && d.X >= -1 && d.Y <= 1 && d.Y >= -1
}

type Move struct {
	Delta Point
	Count int
}

type Rope struct {
	Head          Point
	Tail          Point
	TailPositions map[Point]int
}

func NewRope() *Rope {
	return &Rope{
		// put first tail position into the map
		TailPositions: map[Point]int{{}: 1},
	}
}

func (r *Rope) Move(move Move) {
	r.Head = r.Head.Add(move.Delta)
	if !r.Head.Adjacent(r.Tail) {
		// figure out where to move the tail
		diff := r.Head.Sub(r.Tail)
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
		r.Tail = r.Tail.Add(diff)
		r.TailPositions[r.Tail]++
	}
}

func (r *Rope) ExecuteMoves(moves []Move) {
	for _, m := range moves {
		for n := 0; n < m.Count; n++ {
			r.Move(m)
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
	rope := NewRope()
	rope.ExecuteMoves(moves)
	fmt.Println(len(rope.TailPositions))
}
