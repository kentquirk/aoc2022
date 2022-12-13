package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

type Pair [2]any

func ParsePairs(lines []string) []Pair {
	var pairs []Pair
	for i := 0; i < len(lines); i += 3 {
		p := Pair{}
		json.Unmarshal([]byte(lines[i]), &p[0])
		json.Unmarshal([]byte(lines[i+1]), &p[1])
		pairs = append(pairs, p)
	}
	return pairs
}

func ParseLines(lines []string) []any {
	var result []any
	for _, l := range lines {
		if l == "" {
			continue
		}

		var item any
		json.Unmarshal([]byte(l), &item)
		result = append(result, item)
	}
	return result
}

func toString(a any) string {
	switch va := a.(type) {
	case float64:
		return fmt.Sprint(va)
	case []any:
		all := []string{}
		for _, i := range va {
			all = append(all, toString(i))
		}
		return "[" + strings.Join(all, ",") + "]"
	default:
		return "OOPS"
	}
}

func compare(a, b any) int {
	switch va := a.(type) {
	case float64:
		switch vb := b.(type) {
		case float64:
			return int(va - vb)
		case []any:
			return compare([]any{a}, b)
		default:
			panic("b is bad")
		}
	case []any:
		switch vb := b.(type) {
		case float64:
			return compare(a, []any{b})
		case []any:
			la := len(va)
			lb := len(vb)
			for i := 0; ; i++ {
				// if we ran off the end of both lists, they're equal
				if i == la && i == lb {
					return 0
				}
				// first list shorter, list a is less
				if i == la {
					return -1
				}
				// second list shorter, list b is less
				if i == lb {
					return 1
				}
				// we're still in both lists, compare the elements of the lists
				c := compare(va[i], vb[i])
				if c != 0 {
					return c
				}
			}
		default:
			fmt.Printf("%#v (%T)\n", b, b)
			panic("b is bad")
		}
	default:
		fmt.Printf("%#v (%T)\n", a, a)
		panic("a is bad")
	}
}

func part1(lines []string) int {
	pairs := ParsePairs(lines)
	sum := 0
	for i, p := range pairs {
		if compare(p[0], p[1]) < 0 {
			sum += i + 1
		}
	}
	return sum
}

func part2(lines []string) int {
	divider1 := "[[2]]"
	divider2 := "[[6]]"
	lines = append(lines, divider1, divider2)
	all := ParseLines(lines)

	sort.Slice(all, func(i, j int) bool {
		return compare(all[i], all[j]) < 0
	})

	var first, second int
	for i, p := range all {
		s := toString(p)
		if s == divider1 {
			first = i + 1
		}
		if s == divider2 {
			second = i + 1
		}
		// fmt.Println(p)
	}
	return first * second
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
