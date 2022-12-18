package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Action interface {
	Do(sys *System)
	String() string
}

// valve is a node, tunnel is an edge
type Valve struct {
	Name        string
	FlowRate    int
	OpenCost    int
	Tunnels     Set[string]
	ValveAction Action
	MoveActions []Action
}

type MoveAction struct {
	From string
	To   string
	Via  string
}

func (a MoveAction) Do(sys *System) {
	sys.Current = a.To
	sys.Ticks++
}

func (a MoveAction) String() string {
	return fmt.Sprintf("Move from %s to %s via %s", a.From, a.To, a.Via)
}

type OpenAction struct {
	At string
}

func (a OpenAction) Do(sys *System) {
	sys.OpenValves.Add(a.At)
	sys.Ticks++
}

func (a OpenAction) String() string {
	return fmt.Sprintf("Open valve in %s", a.At)
}

type Tunnel struct {
	Name         string
	TraverseCost int
	From         string
	To           string
}

type System struct {
	Valves     map[string]*Valve
	Tunnels    map[string]*Tunnel
	OpenValves Set[string]
	Start      string
	Current    string
	Ticks      int
}

func NewSystem(lines []string) *System {
	s := &System{
		Valves:  map[string]*Valve{},
		Tunnels: map[string]*Tunnel{},
		Start:   "AA",
	}

	namepat := regexp.MustCompile("[A-Z]{2}")
	numpat := regexp.MustCompile("[0-9]+")
	connections := map[string][]string{}
	// make all the valves and record their connections
	for _, l := range lines {
		names := namepat.FindAllString(l, -1)
		rate, _ := strconv.Atoi(numpat.FindString(l))
		var va Action // default to nil
		if rate > 0 {
			va = OpenAction{At: names[0]}
		}
		s.Valves[names[0]] = &Valve{
			Name:        names[0],
			FlowRate:    rate,
			OpenCost:    1,
			Tunnels:     NewSet[string](),
			ValveAction: va,
			MoveActions: make([]Action, 0),
		}
		sort.Strings(names[1:])
		connections[names[0]] = names[1:]
	}
	// now that the valves exist, we can make tunnels to link them up
	for v, names := range connections {
		for _, name := range names {
			t := &Tunnel{
				Name:         v + " -- " + name,
				TraverseCost: 1,
				From:         v,
				To:           name,
			}
			s.Tunnels[t.Name] = t
			s.Valves[v].Tunnels.Add(t.Name)
			s.Valves[v].MoveActions = append(s.Valves[v].MoveActions, MoveAction{From: v, To: name, Via: t.Name})
		}
	}
	return s
}

// func (s *System) Consolidate() {
// 	var removes []string
// 	for n, v := range s.Valves {
// 		if v.FlowRate == 0 && v.Name != "AA" {
// 			removes = append(removes, n)
// 		}
// 	}

// 	for _, r := range removes {
// 		for n, v := range s.Valves {
// 			if v.Tunnels.Contains(r) {
// 				v.Tunnels.Add(s.Tunnels)
// 		}
// 	}
// 	// delete(s.Valves, r)
// }

func (s *System) GenerateGraphviz() string {
	b := strings.Builder{}
	b.WriteString("graph G {\n")
	for n, v := range s.Valves {
		col := "cyan"
		if n == "AA" {
			col = "yellow"
		}
		if v.FlowRate != 0 {
			col = "orange"
		}
		b.WriteString(fmt.Sprintf("  %s [color=%s style=filled label=\"%s (%d)\"];\n", n, col, n, v.FlowRate))
	}
	for _, t := range s.Tunnels {
		b.WriteString(fmt.Sprintf("  %s -- %s [color=blue label=%d];\n", t.From, t.To, t.TraverseCost))
	}
	b.WriteString("}")
	return b.String()
}

type cacheKey struct {
	Start   string
	MaxTime int
}

type bestPath struct {
	Valves   []string
	Pressure int
}

var pathCache = make(map[cacheKey]bestPath)

// Inputs:
//   Memoize key:
//     Current valve
//     time remaining
//   existing state:
//     valves already open
//     tunnels already traversed
// Returns:
//     best path found given the above -- a list of actions
//       Go (tunnel)
//       Open current valve

// Return the path taken and total pressure released for the best path starting
// from start that a total time of maxtime or less.
// Store the result in pathcache to shortcut recursion.
func (s *System) bestPathFrom(start string, time int, openValves Set[string], visitedTunnels Set[string]) bestPath {
	key := cacheKey{Start: start, MaxTime: time}
	if bp, ok := pathCache[key]; ok {
		return bp
	}
	me := s.Valves[start]
	best := bestPath{}
	me.Tunnels.Each(func(name string) bool {
		visited := visitedTunnels.Clone()
		visited.Add(name)
		candidate := s.bestPathFrom(name, time-1, openValves.Clone(), visited)
		if best.Pressure < candidate.Pressure {
			best = candidate
		}
		if me.FlowRate > 0 && !openValves.Contains(me.Name) {
			open := openValves.Clone()
			open.Add(me.Name)
			candidate := s.bestPathFrom(name, time-2, open, visited)
			if best.Pressure < candidate.Pressure {
				best = candidate
			}
		}
		return false
	})
	return best
}

func (s *System) TraverseAll(start string) {
	// for
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
	s := NewSystem(lines)
	fmt.Println(s.GenerateGraphviz())
}
