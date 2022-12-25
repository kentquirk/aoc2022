package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// We're going to make a slice of items that retains the original
// value ordering, and turn it into a linked list we can adjust
// as we move things.
type Item struct {
	Value int
	Prev  *Item
	Next  *Item
}

func (i *Item) String() string {
	return fmt.Sprintf("[%d N->%d P->%d]", i.Value, i.Next.Value, i.Prev.Value)
}

type Sequence struct {
	Len   int
	Items []*Item
	Head  *Item
	Zero  *Item
}

func (s *Sequence) Print() {
	for _, it := range s.Items {
		fmt.Printf("%d ", it.Value)
	}
	fmt.Println()
}

func (s *Sequence) PrintFromZero() {
	head := s.Zero
	for i := 0; i < s.Len; i++ {
		fmt.Printf("%d ", head.Value)
		head = head.Next
	}
	fmt.Println()
}

func (s *Sequence) Reorder() {
	head := s.Zero
	for i := 0; i < s.Len; i++ {
		s.Items[i] = head
		head = head.Next
	}
}

func (s *Sequence) NthItem(ix int) int {
	return s.Items[ix%s.Len].Value
}

func (s *Sequence) Coords() int {
	return s.NthItem(1000) + s.NthItem(2000) + s.NthItem(3000)
}

// There are 2 "correct" ways to do mod, and go chose the one where
// negative numbers stay negative. We need both kinds in this
// particular puzzle. This is the other one.
func mod(n int, z int) int {
	n %= z
	if n < 0 {
		n += z
	}
	return n
}

func BuildSequence(lines []string, key int) *Sequence {
	nitems := len(lines)
	seq := &Sequence{
		Len:   nitems,
		Items: make([]*Item, nitems),
	}
	for i, l := range lines {
		n, _ := strconv.Atoi(l)
		seq.Items[i] = &Item{Value: n * key}
	}
	for i, it := range seq.Items {
		it.Prev = seq.Items[mod(i-1, nitems)]
		it.Next = seq.Items[mod(i+1, nitems)]
		if it.Value == 0 {
			seq.Zero = it
		}
	}
	seq.Head = seq.Items[0]
	return seq
}

func (s *Sequence) Mix() {
	for _, it := range s.Items {
		// fmt.Printf("moving %d\n", it.Value)
		// if we move something by s.Len-1 it stays in the same place
		steps := it.Value % (s.Len - 1)
		// where we came from
		hp := s.Head.Prev
		prev := it.Prev
		// remove it from the old location
		it.Prev.Next = it.Next
		it.Next.Prev = it.Prev

		moveafter := prev
		if steps < 0 {
			for ; steps != 0; steps++ {
				moveafter = moveafter.Prev
			}
		} else {
			for ; steps != 0; steps-- {
				moveafter = moveafter.Next
			}
		}
		// now put ourselves back in the list
		it.Prev = moveafter
		it.Next = moveafter.Next
		it.Prev.Next = it
		it.Next.Prev = it
		s.Head = hp.Next
		// s.PrintFromZero()
	}
}

var dbg bool = false

func part1(lines []string) int {
	seq := BuildSequence(lines, 1)
	if dbg {
		seq.Print()
		seq.PrintFromZero()
	}
	seq.Mix()
	if dbg {
		seq.PrintFromZero()
	}
	seq.Reorder()
	if dbg {
		seq.PrintFromZero()
		seq.Print()
	}
	return seq.Coords()
}

func part2(lines []string) int {
	seq := BuildSequence(lines, 811589153)
	if dbg {
		seq.Print()
		seq.PrintFromZero()
	}
	for i := 0; i < 10; i++ {
		seq.Mix()
	}
	if dbg {
		seq.PrintFromZero()
	}
	seq.Reorder()
	if dbg {
		seq.PrintFromZero()
		seq.Print()
	}
	return seq.Coords()
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
