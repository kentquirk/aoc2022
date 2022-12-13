package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"github.com/beefsack/go-astar"
)

type Square struct {
	Row    int
	Col    int
	Height int
	Links  []*Square
	PathIx int
}

func (s *Square) PathNeighbors() []astar.Pather {
	var neighbors []astar.Pather
	for _, l := range s.Links {
		neighbors = append(neighbors, l)
	}
	return neighbors
}

func (s *Square) PathNeighborCost(to astar.Pather) float64 {
	return 1
}

func (s *Square) PathEstimatedCost(to astar.Pather) float64 {
	dx := s.Col - to.(*Square).Col
	dy := s.Row - to.(*Square).Row
	return math.Sqrt(float64(dx*dx + dy*dy))
}

type Grid struct {
	Squares  [][]*Square
	StartRow int
	StartCol int
	EndRow   int
	EndCol   int
}

// generates possible edges from this square to neighbors
func (g *Grid) Connect(r, c int) {
	nRows := len(g.Squares)
	nCols := len(g.Squares[0])
	self := g.Squares[r][c]
	if r > 0 {
		if g.Squares[r-1][c].Height <= self.Height+1 {
			self.Links = append(self.Links, g.Squares[r-1][c])
		}
	}
	if c > 0 {
		if g.Squares[r][c-1].Height <= self.Height+1 {
			self.Links = append(self.Links, g.Squares[r][c-1])
		}
	}
	if r < nRows-1 {
		if g.Squares[r+1][c].Height <= self.Height+1 {
			self.Links = append(self.Links, g.Squares[r+1][c])
		}
	}
	if c < nCols-1 {
		if g.Squares[r][c+1].Height <= self.Height+1 {
			self.Links = append(self.Links, g.Squares[r][c+1])
		}
	}
}

func (g *Grid) GetLowestPoints() []*Square {
	var candidates []*Square

	nRows := len(g.Squares)
	nCols := len(g.Squares[0])
	for r := 0; r < nRows; r++ {
		for c := 0; c < nCols; c++ {
			if g.Squares[r][c].Height == 0 {
				candidates = append(candidates, g.Squares[r][c])
			}
		}
	}
	return candidates
}

func (g *Grid) GenerateGraph() {
	nRows := len(g.Squares)
	nCols := len(g.Squares[0])
	for r := 0; r < nRows; r++ {
		for c := 0; c < nCols; c++ {
			g.Connect(r, c)
		}
	}
}

func (g *Grid) Print(sr int, sc int, distance int) {
	nRows := len(g.Squares)
	nCols := len(g.Squares[0])
	for r := 0; r < nRows; r++ {
		for c := 0; c < nCols; c++ {
			sq := g.Squares[r][c]
			fg := 34
			bg := 0
			if r == sr && c == sc {
				bg = 42
			}
			if sq.PathIx >= 0 {
				fg = 37
				bg = 44
				if sq.PathIx == distance-1 {
					bg = 41
				}
			}
			fmt.Printf("\x1b[%d;%dm%c", bg, fg, 'a'+sq.Height)
		}
		fmt.Println("\x1b[0m")
	}
}

func Parse(lines []string) *Grid {
	grid := &Grid{Squares: make([][]*Square, 0)}
	for r, l := range lines {
		row := []*Square{}
		for c, ch := range l {
			hgt := int(ch - 'a')
			switch ch {
			case 'S':
				hgt = 0
				grid.StartRow = r
				grid.StartCol = c
			case 'E':
				hgt = 25
				grid.EndRow = r
				grid.EndCol = c
			}
			row = append(row, &Square{Row: r, Col: c, Height: hgt, PathIx: -1})
		}
		grid.Squares = append(grid.Squares, row)
	}
	grid.GenerateGraph()
	return grid
}

func part1(lines []string) int {
	grid := Parse(lines)

	path, distance, _ := astar.Path(
		grid.Squares[grid.StartRow][grid.StartCol],
		grid.Squares[grid.EndRow][grid.EndCol],
	)
	for i, p := range path {
		// pather returns the path in reverse order, so compensate
		p.(*Square).PathIx = int(distance) - i - 1
	}

	grid.Print(grid.StartRow, grid.StartCol, int(distance))
	return int(distance)
}

func part2(lines []string) int {
	grid := Parse(lines)

	candidates := grid.GetLowestPoints()
	var bestsq *Square
	var bestdist float64 = 1000000
	var bestpath []astar.Pather
	for _, c := range candidates {
		path, distance, found := astar.Path(
			c,
			grid.Squares[grid.EndRow][grid.EndCol],
		)
		if found && distance < bestdist {
			bestdist = distance
			bestpath = path
			bestsq = c
		}
	}
	fmt.Println()
	for i, p := range bestpath {
		// pather returns the path in reverse order, so compensate
		p.(*Square).PathIx = int(bestdist) - i - 1
	}

	grid.Print(bestsq.Row, bestsq.Col, int(bestdist))
	return int(bestdist)
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
