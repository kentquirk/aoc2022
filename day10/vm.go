package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type Opcode struct {
	Name string
	Code []func(vm *VM, ins Instruction)
}

type Instruction struct {
	Opcode *Opcode
	Arg    int
}

type VM struct {
	Memory []Instruction
	Ticks  int
	X      int
	LastX  int
	IP     int
	OpIx   int
}

var Microcode = map[string]*Opcode{
	"addx": {
		Name: "addx",
		Code: []func(vm *VM, ins Instruction){
			func(vm *VM, ins Instruction) {
				vm.Ticks++
				vm.OpIx++
				vm.LastX = vm.X
			},
			func(vm *VM, ins Instruction) {
				vm.Ticks++
				vm.OpIx = 0
				vm.IP++
				vm.LastX = vm.X
				vm.X += ins.Arg
			},
		},
	},
	"noop": {
		Name: "noop",
		Code: []func(vm *VM, ins Instruction){
			func(vm *VM, ins Instruction) {
				vm.Ticks++
				vm.IP++
				vm.LastX = vm.X
			},
		},
	},
}

func NewVM() *VM {
	return &VM{X: 1}
}

func (vm *VM) Load(lines []string) error {
	vm.Memory = make([]Instruction, 0)
	vm.X = 1
	linepat := regexp.MustCompile(`([a-z]+) ?(-?[0-9]+)?`)
	for _, l := range lines {
		tokens := linepat.FindStringSubmatch(l)
		if len(tokens) < 2 || len(tokens) > 3 {
			return fmt.Errorf("parse error: %s", tokens[0])
		}
		op, ok := Microcode[tokens[1]]
		if !ok {
			return fmt.Errorf("invalid opcode %s", tokens[1])
		}
		ins := Instruction{Opcode: op}
		if len(tokens) == 3 {
			a, _ := strconv.Atoi(tokens[2])
			ins.Arg = a
		}
		vm.Memory = append(vm.Memory, ins)
	}
	return nil
}

func (vm *VM) Reset() {
	vm.X = 1
	vm.Ticks = 0
	vm.IP = 0
	vm.OpIx = 0
}

func (vm *VM) Tick() bool {
	if vm.IP >= len(vm.Memory) {
		return false
	}
	vm.Memory[vm.IP].Opcode.Code[vm.OpIx](vm, vm.Memory[vm.IP])
	return true
}

func (vm *VM) String() string {
	return fmt.Sprintf("VM: %d instr, IP: %3d, OpIx %1d, Ticks %3d, LastX: %4d, X: %4d", len(vm.Memory), vm.IP, vm.OpIx, vm.Ticks, vm.LastX, vm.X)
}
