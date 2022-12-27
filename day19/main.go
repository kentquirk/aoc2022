package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dgryski/go-wyhash"
)

type ResourceStr string

type Blueprint struct {
	Index   int
	Costs   map[ResourceStr]map[ResourceStr]int // map of a resource to a map of resource costs
	Quality int
}

func mcopy[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

type Simulation struct {
	BP        *Blueprint
	Robots    map[ResourceStr]int
	Resources map[ResourceStr]int
	Time      int
}

func (s *Simulation) Clone() *Simulation {
	return &Simulation{
		BP:        s.BP,
		Robots:    mcopy(s.Robots),
		Resources: mcopy(s.Resources),
		Time:      s.Time,
	}
}

func (s *Simulation) String() string {
	return fmt.Sprintf("%v %v %d", s.Robots, s.Resources, s.Time)
}

func (s *Simulation) Hash() uint64 {
	// var b []byte
	return wyhash.Hash([]byte(s.String()), 0x12ad478ef90)
}

func (s *Simulation) Tick() {
	for resource, qty := range s.Robots {
		s.Resources[resource] += qty
	}
	s.Time++
}

func (s *Simulation) BuyRobot(resource ResourceStr) {
	if resource == "" {
		return
	}
	costs := s.BP.Costs[resource]
	for item, cost := range costs {
		s.Resources[item] -= cost
	}
	s.Robots[resource]++
}

type Action string

type MemoCache map[uint64][]Action

var memoCache = make(MemoCache)

// Recursively search for the best result within the given time limit
func (s *Simulation) Run(timelimit int) []Action {
	// first, are we already in the cache?
	// TODO: our current return value includes geodes as part of s, so the hash isn't accurate
	hash := s.Hash()
	if actions, ok := memoCache[hash]; ok {
		return actions
	}

	// see what we can afford to build
	// for all the resources in the blueprint
	bestGeodes := 0
	var bestActions []Action
	for resource, costs := range s.BP.Costs {
		// see if we have enough to make a robot
		canAfford := true
		for item, cost := range costs {
			if s.Resources[item] <= cost {
				canAfford = false
			}
		}
		geodes := 0
		var actions []Action
		if canAfford {
			// ok, let's clone ourselves, run a tick, buy the robot, and then recurse
			t := s.Clone()
			t.Tick()
			t.BuyRobot(resource)
			if t.Time == timelimit {
				actions = []Action{Action(resource)}
			} else {
				actions = t.Run(timelimit)
			}
			geodes = t.Resources["geode"]
		}
		if geodes > bestGeodes {
			bestActions = actions
			bestGeodes = geodes
		}
	}
	s.Resources["geode"] = bestGeodes
	return bestActions
}

func parse(lines []string) []*Blueprint {
	var blueprints []*Blueprint
	robotpat := regexp.MustCompile(`Each ([a-z]+) robot costs (?:([0-9]) ([a-z]+) and )?([0-9]+) ([a-z]+)\.`)
	for i, l := range lines {
		matches := robotpat.FindAllStringSubmatch(l, -1)
		bp := &Blueprint{Index: i + 1, Costs: make(map[ResourceStr]map[ResourceStr]int)}
		for _, m := range matches {
			robot := ResourceStr(m[1])
			resources := make(map[ResourceStr]int)
			if m[3] != "" {
				q, _ := strconv.Atoi(m[2])
				resources[ResourceStr(m[3])] = q
			}
			if m[5] != "" {
				q, _ := strconv.Atoi(m[4])
				resources[ResourceStr(m[5])] = q
			}
			bp.Costs[robot] = resources
		}
		bp.Costs[""] = nil
		blueprints = append(blueprints, bp)
	}
	return blueprints
}

func part1(lines []string) int {
	blueprints := parse(lines)
	bestGeodes := 0
	totalQuality := 0
	var bestActions []Action
	var bestIndex int
	for _, b := range blueprints {
		fmt.Println(b)
		sim := Simulation{
			BP:        b,
			Robots:    map[ResourceStr]int{"ore": 1},
			Resources: make(map[ResourceStr]int),
		}
		actions := sim.Run(18)
		geodes := sim.Resources["geode"]
		quality := geodes * b.Index
		totalQuality += quality
		if geodes > bestGeodes {
			bestActions = actions
			bestGeodes = geodes
			bestIndex = b.Index
		}
	}
	fmt.Println(bestGeodes, bestIndex, bestActions)
	return totalQuality
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
	fmt.Println(part1(lines))
}
