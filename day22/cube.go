package main

import (
	"regexp"
	"strconv"
)

type Neighbor struct {
	Index  int
	Facing Orientation
}

type Face struct {
	Rows      []*Row
	Neighbors []*Neighbor // index is Orientation
}

type Cube struct {
	Faces []*Face
}

func NewCube(lines []string) *Board {
	cube := &Cube{
		Faces: make([]*Face, 6),
	}
	rowpat := regexp.MustCompile("[.#]{50}")
	faceIndex := 0
	for lineIx, l := range lines {
		if lineIx % 50 == 0 {
			faceIndex+=
		data := rowpat.FindAll([]byte(l), -1)
		for _, d := range data {
		var tiles []*Tile
		for tileIx, ch := range data[2] {
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
