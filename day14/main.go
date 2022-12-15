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

type State int

const (
	Empty State = iota
	Wall
	Sand
	Origin
)

type Point struct {
	X int
	Y int
}

type Cave struct {
	Min   Point
	Max   Point
	Cells map[Point]State
}

func NewCave() *Cave {
	return &Cave{
		Min:   Point{10000, 10000},
		Max:   Point{0, 0},
		Cells: make(map[Point]State),
	}
}

func (c *Cave) Print() {
	fmt.Println("---")
	for y := c.Min.Y; y <= c.Max.Y; y++ {
		for x := c.Min.X; x <= c.Max.X; x++ {
			switch c.Cells[Point{x, y}] {
			case Empty:
				fmt.Print(".")
			case Wall:
				fmt.Print("#")
			case Sand:
				fmt.Print("o")
			case Origin:
				fmt.Print("+")
			default:
				fmt.Print("!")
			}
		}
		fmt.Println()
	}
}

func sign(v1, v2 int) int {
	switch {
	case v1 < v2:
		return 1
	case v1 > v2:
		return -1
	default:
		return 0
	}
}

func (c *Cave) CheckLimits(pt Point) {
	if pt.X < c.Min.X {
		c.Min.X = pt.X
	}
	if pt.X > c.Max.X {
		c.Max.X = pt.X
	}
	if pt.Y < c.Min.Y {
		c.Min.Y = pt.Y
	}
	if pt.Y > c.Max.Y {
		c.Max.Y = pt.Y
	}
}

func (c *Cave) DrawWall(from Point, to Point) {
	c.CheckLimits(from)
	c.CheckLimits(to)
	if from.Y == to.Y {
		sgn := sign(from.X, to.X)
		for x := from.X; x != to.X+sgn; x += sgn {
			c.Cells[Point{x, from.Y}] = Wall
		}
	} else {
		sgn := sign(from.Y, to.Y)
		for y := from.Y; y != to.Y+sgn; y += sgn {
			c.Cells[Point{from.X, y}] = Wall
		}
	}
}

func (c *Cave) Drop(pt Point, until func(Point) bool) Point {
	for {
		newPt := Point{pt.X, pt.Y + 1}
		if until(newPt) {
			c.Cells[pt] = Sand
			return newPt
		}
		if c.Cells[newPt] == Empty {
			pt = newPt
			continue
		}
		newPt.X--
		if c.Cells[newPt] == Empty {
			pt = newPt
			continue
		}
		newPt.X += 2
		if c.Cells[newPt] == Empty {
			pt = newPt
			continue
		}
		if s, ok := c.Cells[pt]; ok {
			fmt.Println("reuse at ", pt, s)
		}
		c.Cells[pt] = Sand
		return pt
	}
}

func (c *Cave) Parse(lines []string) {
	linepat := regexp.MustCompile(`[0-9]+,[0-9]+`)
	for _, l := range lines {
		pts := linepat.FindAllString(l, -1)
		var points []Point
		for _, p := range pts {
			values := strings.Split(p, ",")
			x, _ := strconv.Atoi(values[0])
			y, _ := strconv.Atoi(values[1])
			points = append(points, Point{x, y})
		}
		for i := 1; i < len(points); i++ {
			c.DrawWall(points[i-1], points[i])
		}
	}
}

func part1(lines []string) {
	c := NewCave()
	origin := Point{500, 0}
	c.Parse(lines)
	c.CheckLimits(origin)
	c.Cells[origin] = Origin

	grains := 0
	for {
		pt := c.Drop(origin, func(p Point) bool { return p.Y > c.Max.Y })
		if pt.Y > c.Max.Y {
			break
		}
		grains++
	}
	c.Print()
	fmt.Println(grains)
}

func part2(lines []string) {
	c := NewCave()
	origin := Point{500, 0}
	c.Parse(lines)
	c.DrawWall(Point{c.Min.X - 1, c.Max.Y + 2}, Point{c.Max.X + 1, c.Max.Y + 2})
	c.CheckLimits(origin)

	grains := 0
	for {
		grains++
		pt := c.Drop(origin, func(p Point) bool { return p.Y == c.Max.Y })
		c.CheckLimits(pt)
		if pt == origin {
			break
		}
	}
	c.DrawWall(Point{c.Min.X, c.Max.Y}, c.Max)
	// c.Print()
	fmt.Println(grains)
	n := 0
	for _, s := range c.Cells {
		if s == Sand {
			n++
		}
	}
	fmt.Println(n)
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
	part1(lines)
	part2(lines)
}
