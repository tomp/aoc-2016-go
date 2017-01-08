package asmbunny

import (
	"fmt"
	"strconv"
	"strings"
)

type opCode int
type regType string

const REGNAMES string = "abcd"
const NREG int = len(REGNAMES)
const NOTREG regType = ""

type Registers [NREG]int

// Get returns the value of the named register.
func (r Registers) Get(regname string) (value int) {
	i := strings.Index(REGNAMES, regname)
	if i >= 0 {
		value = r[i]
	}
	return
}

// Set sets the value of the named register.
func (r *Registers) Set(regname string, value int) {
	i := strings.Index(REGNAMES, regname)
	if i >= 0 {
		(*r)[i] = value
	}
}

const (
	NOP opCode = iota
	CPY
	INC
	DEC
	JNZ
)

type opType struct {
	code  opCode
	name  string
	nargs int
}

var OPS = [...]opType{
	{NOP, "nop", 0},
	{CPY, "cpy", 2},
	{INC, "inc", 1},
	{DEC, "dec", 1},
	{JNZ, "jnz", 2},
}

type Inst struct {
	op   opType
	x    regType
	xval int
	y    regType
	yval int
}

func (i *Inst) String() string {
	if i.op.nargs == 0 {
		return i.op.name
	}
	var x, y string
	if i.x == NOTREG {
		x = fmt.Sprintf("%d", i.xval)
	} else {
		x = string(i.x)
	}
	if i.op.nargs == 1 {
		return i.op.name + " " + x
	}
	if i.y == NOTREG {
		y = fmt.Sprintf("%d", i.yval)
	} else {
		y = string(i.y)
	}
	return i.op.name + " " + x + " " + y
}

type Program struct {
	inst []Inst
}

func NewProgram() Program {
	return Program{inst: []Inst{}}
}

// parseArg parses a single instruction argument.  If it names a
// register, then the register name and index are returned as 'reg' and
// 'val.  If it's an integer constant, then its value is returned as
// 'val' and 'reg' is the empty string.
func parseArg(arg string) (reg regType, val int, err error) {
	val = strings.Index(REGNAMES, arg)
	if val >= 0 {
		reg = regType(arg)
		return
	}
	val, err = strconv.Atoi(arg)
	return
}

// compileInst compiles a single parsed instruction into an Inst.
// If the instruction isn't recognized, an error is returned.
func compileInst(name string, args []string) (inst Inst, err error) {
	for _, op := range OPS {
		if op.name == name {
			if len(args) != op.nargs {
				err = fmt.Errorf("Op '%s' takes %d arg(s) (found %d)",
					op.name, op.nargs, len(args))
				return
			}
			if op.nargs == 0 {
				inst = Inst{op, NOTREG, 0, NOTREG, 0}
				return
			}
			x, xval, perr := parseArg(args[0])
			if perr != nil {
				err = perr
				return
			}
			if op.nargs == 1 {
				inst = Inst{op, x, xval, NOTREG, 0}
				return
			}
			y, yval, perr := parseArg(args[1])
			if perr != nil {
				err = perr
				return
			}
			inst = Inst{op, x, xval, y, yval}
			return
		}
	}
	err = fmt.Errorf("Unrecognized op '%s'", name)
	return
}

func Compile(source []string) (prog Program, err error) {
	prog = NewProgram()
	for _, line := range source {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix("#", line) {
			continue
		}
		words := strings.Fields(line)
		inst, cerr := compileInst(words[0], words[1:])
		if cerr != nil {
			err = cerr
			return
		}
		prog.inst = append(prog.inst, inst)
	}
	return
}

// Execute the program with the given initial register values.
// The final register values are returned.
func (p *Program) Execute(initreg Registers) (reg Registers, err error) {
	reg, _, err = p.ExecuteFrom(initreg, 0)
	return
}

// Execute the program starting with a particular program counter and the
// given initial register values.
// The final register values and program counter are returned.
func (p *Program) ExecuteFrom(initreg Registers, initpc int) (reg Registers, pc int, err error) {
	for i, val := range initreg {
		reg[i] = val
	}
	pc = initpc
	for pc >= 0 && pc < len(p.inst) {
		inst := p.inst[pc]
		// fmt.Printf("pc:%02d  a:%d\tb:%d\tc:%d\td:%d\t%s\n", pc,
		//   reg[0], reg[1], reg[2], reg[3], inst.String())
		switch inst.op.code {
		case INC:
			reg[inst.xval] += 1
			pc += 1
		case DEC:
			reg[inst.xval] -= 1
			pc += 1
		case CPY:
			if inst.x == NOTREG {
				reg[inst.yval] = inst.xval
			} else {
				reg[inst.yval] = reg[inst.xval]
			}
			pc += 1
		case JNZ:
			jump := false
			if inst.x == NOTREG {
				jump = inst.xval != 0
			} else {
				jump = reg[inst.xval] != 0
			}
			if !jump {
				pc += 1
			} else {
				if inst.y == NOTREG {
					pc += inst.yval
				} else {
					pc += reg[inst.yval]
				}
			}
		case NOP:
			pc += 1
		}
	}
	return
}
