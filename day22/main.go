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

type Orientation byte

const (
	Right Orientation = iota
	Down
	Left
	Up
)

func (o Orientation) LeftOf() Orientation {
	if o == Right {
		return Up
	}
	return o - 1
}

func (o Orientation) RightOf() Orientation {
	if o == Up {
		return Right
	}
	return o + 1
}

func (o Orientation) Value() int {
	return int(o)
}

func (o Orientation) String() string {
	return ">v<^"[o : o+1]
}

type TileState byte

func (t TileState) String() string {
	return ".#"[t : t+1]
}

type Tile struct {
	State      TileState
	LastFacing Orientation
	Visited    bool
}

const (
	Empty TileState = iota
	Wall
)

type Row struct {
	Tiles  []*Tile
	Offset int
}

func (r *Row) IsOnRow(c int) bool {
	return c >= r.Offset && c < r.Offset+len(r.Tiles)
}

func (r *Row) First() int {
	return r.Offset
}

func (r *Row) Last() int {
	return r.Offset + len(r.Tiles) - 1
}

func (r *Row) Get(c int) *Tile {
	return r.Tiles[c-r.Offset]
}

type Player struct {
	Facing Orientation
	Row    int
	Col    int
}

func (p *Player) Password() int {
	r := p.Row + 1
	c := p.Col + 1
	return 1000*r + 4*c + p.Facing.Value()
}

type Board struct {
	Rows []*Row
	Path string
}

func NewBoard(lines []string) *Board {
	b := &Board{}
	rowpat := regexp.MustCompile("( *)([.#]+)")
	for _, l := range lines {
		data := rowpat.FindSubmatch([]byte(l))
		if len(data) < 3 {
			break
		}
		var tiles []*Tile
		for _, ch := range data[2] {
			switch ch {
			case '#':
				tiles = append(tiles, &Tile{State: Wall})
			case '.':
				tiles = append(tiles, &Tile{State: Empty})
			default:
				panic("oops")
			}
		}
		b.Rows = append(b.Rows, &Row{
			Tiles:  tiles,
			Offset: len(data[1]),
		})
	}
	b.Path = lines[len(lines)-1]
	return b
}

// Modifies player to be one character further right in the row taking toroidal
// movement into account. Returns true if the move was successful, and false if
// the player hit a wall and can't go farther. Assumes that starting player
// position is valid.
func (b *Board) MoveRight(p *Player) bool {
	row := b.Rows[p.Row]
	cix := p.Col + 1
	if !row.IsOnRow(cix) {
		cix = row.First()
	}
	if row.Get(cix).State == Empty {
		p.Col = cix
		return true
	}
	return false
}

func (b *Board) MoveLeft(p *Player) bool {
	row := b.Rows[p.Row]
	cix := p.Col - 1
	if !row.IsOnRow(cix) {
		cix = row.Last()
	}
	if row.Get(cix).State == Empty {
		p.Col = cix
		return true
	}
	return false
}

func (b *Board) MoveDown(p *Player) bool {
	rix := p.Row
	for {
		rix++
		if rix >= len(b.Rows) {
			rix = 0
		}
		row := b.Rows[rix]
		if !row.IsOnRow(p.Col) {
			continue
		}
		if row.Get(p.Col).State == Empty {
			p.Row = rix
			return true
		}
		return false
	}
}

func (b *Board) MoveUp(p *Player) bool {
	rix := p.Row
	for {
		rix--
		if rix < 0 {
			rix = len(b.Rows) - 1
		}
		row := b.Rows[rix]
		if !row.IsOnRow(p.Col) {
			continue
		}
		if row.Get(p.Col).State == Empty {
			p.Row = rix
			return true
		}
		return false
	}
}

func (b *Board) StampTile(p *Player) {
	tile := b.Rows[p.Row].Get(p.Col)
	tile.LastFacing = p.Facing
	tile.Visited = true
}

// returns true if full move was successful, false if it hit a wall
func (b *Board) MovePlayer(p *Player, n int) bool {
	for i := 0; i < n; i++ {
		b.StampTile(p)
		switch p.Facing {
		case Right:
			if !b.MoveRight(p) {
				return false
			}
		case Left:
			if !b.MoveLeft(p) {
				return false
			}
		case Up:
			if !b.MoveUp(p) {
				return false
			}
		case Down:
			if !b.MoveDown(p) {
				return false
			}
		}
	}
	return true
}

func (b *Board) Follow() int {
	p := &Player{Col: b.Rows[0].Offset}
	movepat := regexp.MustCompile("[LR]|[0-9]+")
	for _, move := range movepat.FindAllString(b.Path, -1) {
		switch move {
		case "L":
			p.Facing = p.Facing.LeftOf()
		case "R":
			p.Facing = p.Facing.RightOf()
		default:
			n, _ := strconv.Atoi(move)
			b.MovePlayer(p, n)
		}
		b.StampTile(p)
	}
	return p.Password()
}

func (b *Board) Print() {
	for _, row := range b.Rows {
		fmt.Print(strings.Repeat(" ", row.Offset))
		for _, tile := range row.Tiles {
			if tile.Visited {
				fmt.Print(tile.LastFacing)
			} else {
				fmt.Print(tile.State)
			}
		}
		fmt.Println()
	}
}

func part1(lines []string) int {
	b := NewBoard(lines)
	pw := b.Follow()
	// b.Print()
	return pw
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
}
