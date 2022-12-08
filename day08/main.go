package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Tree struct {
	Height  int
	Visible bool
}

type Forest struct {
	Trees [][]Tree
}

func (f *Forest) markVisiblesByRows(nrows, c1, c2, ci int) {
	for r := 0; r < nrows; r++ {
		max := -1
		for c := c1; c != c2; c += ci {
			if f.Trees[r][c].Height > max {
				max = f.Trees[r][c].Height
				f.Trees[r][c].Visible = true
			}
		}
	}
}

func (f *Forest) markVisiblesByCols(ncols, r1, r2, ri int) {
	for c := 0; c < ncols; c++ {
		max := -1
		for r := r1; r != r2; r += ri {
			if f.Trees[r][c].Height > max {
				max = f.Trees[r][c].Height
				f.Trees[r][c].Visible = true
			}
		}
	}
}

func (f *Forest) MarkVisibles() {
	nRows := len(f.Trees)
	nCols := len(f.Trees[0])
	f.markVisiblesByRows(nRows, 0, nCols, 1)     // from left
	f.markVisiblesByRows(nRows, nCols-1, -1, -1) // from right
	f.markVisiblesByCols(nCols, 0, nRows, 1)     // from top
	f.markVisiblesByCols(nCols, nRows-1, -1, -1) // from bottom
}

func (f *Forest) CountVisibles() int {
	nRows := len(f.Trees)
	nCols := len(f.Trees[0])
	total := 0
	for r := 0; r < nRows; r++ {
		for c := 0; c < nCols; c++ {
			if f.Trees[r][c].Visible {
				total++
			}
		}
	}
	return total
}

func (f *Forest) Print(visible bool) {
	nRows := len(f.Trees)
	nCols := len(f.Trees[0])
	for r := 0; r < nRows; r++ {
		for c := 0; c < nCols; c++ {
			t := f.Trees[r][c]
			if visible {
				if t.Visible {
					fmt.Print("*")
				} else {
					fmt.Print(".")
				}
			} else {
				fmt.Printf("%d", t.Height)
			}
		}
		fmt.Println()
	}
}

func (f *Forest) rowViewingDistances(r, c int) (int, int) {
	nCols := len(f.Trees[r])
	h := f.Trees[r][c].Height

	right := 0
	for cx := c + 1; cx < nCols; cx++ {
		if f.Trees[r][cx].Height >= h {
			right++
			break
		}
		right++
	}

	left := 0
	for cx := c - 1; cx >= 0; cx-- {
		if f.Trees[r][cx].Height >= h {
			left++
			break
		}
		left++
	}
	return left, right
}

func (f *Forest) colViewingDistances(r, c int) (int, int) {
	nRows := len(f.Trees)
	h := f.Trees[r][c].Height

	down := 0
	for rx := r + 1; rx < nRows; rx++ {
		if f.Trees[rx][c].Height >= h {
			down++
			break
		}
		down++
	}

	up := 0
	for rx := r - 1; rx >= 0; rx-- {
		if f.Trees[rx][c].Height >= h {
			up++
			break
		}
		up++
	}
	return up, down
}

func (f *Forest) ViewingDistanceFor(row, col int) int {
	l, r := f.rowViewingDistances(row, col)
	u, d := f.colViewingDistances(row, col)
	// fmt.Printf("(R%d, C%d) l:%d r:%d u:%d d:%d\n", row, col, l, r, u, d)
	return l * r * u * d
}

func (f *Forest) BestViewingDistance() int {
	best := 0
	nRows := len(f.Trees)
	nCols := len(f.Trees[0])
	for r := 0; r < nRows; r++ {
		for c := 0; c < nCols; c++ {
			v := f.ViewingDistanceFor(r, c)
			if v > best {
				best = v
				// fmt.Printf("better view %d at (R%d, C%d)\n", best, r, c)
			}
		}
	}
	return best
}

func Parse(lines []string) *Forest {
	forest := &Forest{Trees: make([][]Tree, 0)}
	for _, l := range lines {
		row := []Tree{}
		for _, ch := range l {
			row = append(row, Tree{Height: int(ch - '0')})
		}
		forest.Trees = append(forest.Trees, row)
	}
	return forest
}

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	forest := Parse(lines)
	forest.MarkVisibles()
	// forest.Print(true)
	// forest.Print(false)
	fmt.Println(forest.CountVisibles())
	fmt.Println(forest.BestViewingDistance())
}
