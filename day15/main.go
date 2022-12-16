package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	Min int
	Max int
}

func (r Range) Contains(other Range) bool {
	return r.Min <= other.Min && r.Max >= other.Max
}

func (r Range) Disjoint(other Range) bool {
	return r.Min > other.Max+1 || r.Max+1 < other.Min
}

func (r Range) Len() int {
	return r.Max - r.Min + 1
}

func AddRanges(ranges []Range) int {
	sum := 0
	for _, r := range ranges {
		sum += r.Len()
	}
	return sum
}

func Combine(ranges []Range) []Range {
	if len(ranges) < 2 {
		return ranges
	}
	sort.Slice(ranges, func(i int, j int) bool {
		return ranges[i].Min < ranges[j].Min
	})
	var output []Range
	newrange := ranges[0]
	for i := 1; i < len(ranges); {
		switch {
		case newrange.Disjoint(ranges[i]):
			output = append(output, newrange)
			newrange = ranges[i]
		case newrange.Contains(ranges[i]):
			i++
		default:
			newrange.Max = ranges[i].Max
			i++
		}
	}
	output = append(output, newrange)
	return output
}

type Point struct {
	X int
	Y int
}

func (p Point) TuningFreq() int {
	return 4_000_000*p.X + p.Y
}

type Sensor struct {
	Location      Point
	ClosestBeacon Point
	Distance      int
}

func NewSensor(loc Point, beacon Point) *Sensor {
	return &Sensor{
		Location:      loc,
		ClosestBeacon: beacon,
		Distance:      loc.Manhattan(beacon),
	}
}

func (s *Sensor) String() string {
	return fmt.Sprintf("L:%v B:%v D:%d\n", s.Location, s.ClosestBeacon, s.Distance)
}

// Returns true if pt is closer to Location than the beacon.
func (s *Sensor) Inside(pt Point) bool {
	return pt.Manhattan(s.Location) <= s.Distance
}

func (s *Sensor) RangesFor(line int, beacons []int) []Range {
	var ranges []Range
	if line < s.Location.Y-s.Distance || line > s.Location.Y+s.Distance {
		return ranges
	}

	xmin := s.Location.X - s.Distance + iabs(line-s.Location.Y)
	xmax := s.Location.X + s.Distance - iabs(line-s.Location.Y)

	for b := range beacons {
		switch {
		case b == xmin:
			xmin++
		case b == xmax:
			xmax--
		case b > xmin && b < xmax:
			ranges = append(ranges, Range{xmin, b - 1})
			xmin = b + 1
		}
	}
	ranges = append(ranges, Range{xmin, xmax})
	return ranges
}

func iabs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (p Point) Manhattan(other Point) int {
	dx := iabs(p.X - other.X)
	dy := iabs(p.Y - other.Y)
	return dx + dy
}

type Cave struct {
	Sensors []*Sensor
	Beacons map[Point]struct{}
	Min     Point
	Max     Point
}

func (c *Cave) CheckLimits(pt Point, dist int) {
	if pt.X-dist < c.Min.X {
		c.Min.X = pt.X - dist
	}
	if pt.X+dist > c.Max.X {
		c.Max.X = pt.X + dist
	}
	if pt.Y-dist < c.Min.Y {
		c.Min.Y = pt.Y - dist
	}
	if pt.Y+dist > c.Max.Y {
		c.Max.Y = pt.Y + dist
	}
}

func NewCave(lines []string) *Cave {
	c := &Cave{
		Beacons: make(map[Point]struct{}),
		Min:     Point{math.MaxInt, math.MaxInt},
		Max:     Point{math.MinInt, math.MinInt},
	}
	pat := regexp.MustCompile("[0-9-]+")
	for _, l := range lines {
		var numbers []int
		for _, v := range pat.FindAllString(l, -1) {
			n, _ := strconv.Atoi(v)
			numbers = append(numbers, n)
		}
		loc := Point{numbers[0], numbers[1]}
		beacon := Point{numbers[2], numbers[3]}
		s := NewSensor(loc, beacon)
		c.Sensors = append(c.Sensors, s)
		c.Beacons[beacon] = struct{}{}
		c.CheckLimits(loc, s.Distance)
		c.CheckLimits(beacon, s.Distance)
	}
	return c
}

func (c *Cave) BeaconsFor(line int) []int {
	var beacons []int
	for b := range c.Beacons {
		if b.Y == line {
			beacons = append(beacons, b.X)
		}
	}
	sort.Ints(beacons)
	return beacons
}

func (c *Cave) RangesForRow(n int) []Range {
	beacons := c.BeaconsFor(n)
	var ranges []Range
	for _, s := range c.Sensors {
		ranges = append(ranges, s.RangesFor(n, beacons)...)
	}
	// fmt.Println("R:", ranges)
	ranges = Combine(ranges)
	// fmt.Println(ranges)
	return ranges
}

func (c *Cave) CheckRowWithRanges(n int) int {
	ranges := c.RangesForRow(n)
	return AddRanges(ranges)
}

func part1(lines []string, row int) int {
	c := NewCave(lines)
	fmt.Println(c.Min, c.Max)
	return c.CheckRowWithRanges(row)
}

func part2(lines []string, lastrow int) int {
	c := NewCave(lines)
	for r := 0; r <= lastrow; r++ {
		ranges := c.RangesForRow(r)
		switch len(ranges) {
		case 1:
			continue
		case 2:
			pt := Point{ranges[0].Max + 1, r}
			fmt.Println("Found?", pt, pt.TuningFreq(), ranges)
		default:
			fmt.Println("too big!", ranges)
		}
		fmt.Println(r, ranges)
	}
	return 0
}

func main() {
	fn := "./input.txt"
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	testrow := 10
	lastrow := 20
	if fn == "./input.txt" {
		testrow = 2_000_000
		lastrow = 4_000_000
	}
	testrow = 397470
	fmt.Println(part1(lines, testrow))
	fmt.Println(part2(lines, lastrow))
}
