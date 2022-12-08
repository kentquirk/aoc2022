package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type FileTree struct {
	Name   string
	Dirs   map[string]*FileTree
	Files  map[string]int
	Parent *FileTree
}

func NewFileTree(name string, parent *FileTree) *FileTree {
	return &FileTree{
		Name:   name,
		Dirs:   make(map[string]*FileTree),
		Files:  make(map[string]int),
		Parent: parent,
	}
}

func (f *FileTree) Print(indent int) {
	fmt.Printf("%s- %s (dir, size = %d)\n", strings.Repeat(" ", indent*2), f.Name, f.Size())
	for _, d := range f.Dirs {
		d.Print(indent + 2)
	}
	for name, size := range f.Files {
		fmt.Printf("  %s- %s (file, size=%d)\n", strings.Repeat(" ", indent*2), name, size)
	}
}

func (f *FileTree) Size() int {
	t := 0
	for _, d := range f.Dirs {
		t += d.Size()
	}
	for _, size := range f.Files {
		t += size
	}
	return t
}

type Visitor interface {
	Visit(f *FileTree)
}

func (f *FileTree) Walk(visitor Visitor) {
	visitor.Visit(f)
	for _, d := range f.Dirs {
		d.Walk(visitor)
	}
}

func (f *FileTree) Path() string {
	if f.Parent == nil {
		return f.Name
	}
	return f.Parent.Path() + f.Name + "/"
}

type Parser struct {
	ix       int
	cmds     []string
	pushback []string
}

func NewParser(lines []string) *Parser {
	return &Parser{cmds: lines}
}

func (p *Parser) Pushback(sa []string) {
	p.pushback = sa
}

func (p *Parser) Get() []string {
	if p.pushback != nil {
		r := p.pushback
		p.pushback = nil
		return r
	}
	if p.ix >= len(p.cmds) {
		return nil
	}

	splits := strings.Split(p.cmds[p.ix], " ")
	p.ix++
	if len(splits) < 2 {
		return nil
	}
	return splits
}

func parse(lines []string) *FileTree {
	p := NewParser(lines)
	root := NewFileTree("/", nil)
	current := root

outer:
	for cmd := p.Get(); cmd != nil; cmd = p.Get() {
		if cmd[0] == "$" {
			switch cmd[1] {
			case "cd":
				switch cmd[2] {
				case "/":
					current = root
				case "..":
					if current.Parent != nil {
						current = current.Parent
					}
				default:
					if sub, ok := current.Dirs[cmd[2]]; ok {
						current = sub
					}
				}
			case "ls":
				for {
					item := p.Get()
					if item == nil || item[0] == "$" {
						p.Pushback(item)
						continue outer
					}
					if item[0] == "dir" {
						current.Dirs[item[1]] = NewFileTree(item[1], current)
					} else {
						sz, _ := strconv.Atoi(item[0])
						current.Files[item[1]] = sz
					}
				}
			}
		}
	}
	return root
}

type totaller struct {
	total int
}

func (t *totaller) Visit(f *FileTree) {
	if f.Size() <= 100000 {
		t.total += f.Size()
	}
}

type freeSpace struct {
	removeAtLeast int
	best          *FileTree
}

func (f *freeSpace) Visit(t *FileTree) {
	sz := t.Size()
	if sz > f.removeAtLeast && sz < f.best.Size() {
		f.best = t
	}
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
	root := parse(lines)
	// root.Print(0)
	tt := &totaller{}
	root.Walk(tt)
	fmt.Println(tt.total)

	disksize := 70_000_000
	needed := 30_000_000
	used := root.Size()
	unused := disksize - used
	fs := &freeSpace{removeAtLeast: needed - unused, best: root}
	root.Walk(fs)
	fmt.Println("remove at least ", fs.removeAtLeast)
	fmt.Println(fs.best.Path(), fs.best.Size())
}
