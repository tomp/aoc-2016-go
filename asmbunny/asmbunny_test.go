package asmbunny

import (
	"testing"
)

func TestParseArg(t *testing.T) {
	cases := [...]struct {
		name string
		args []string
		inst Inst
	}{
		{"cpy", []string{"41", "a"},
			Inst{OPS[CPY], NONE, 41, "a", 0}},
		{"cpy", []string{"a", "b"},
			Inst{OPS[CPY], "a", 0, "b", 1}},
		{"inc", []string{"c"},
			Inst{OPS[INC], "c", 2, NONE, 0}},
		{"dec", []string{"d"},
			Inst{OPS[DEC], "d", 3, NONE, 0}},
		{"jnz", []string{"a", "2"},
			Inst{OPS[JNZ], "a", 0, NONE, 2}},
		{"jnz", []string{"1", "2"},
			Inst{OPS[JNZ], NONE, 1, NONE, 2}},
		{"jnz", []string{"1", "b"},
			Inst{OPS[JNZ], NONE, 1, "b", 1}},
	}

	for ncase, item := range cases {
		inst, err := compileInst(item.name, item.args)
		if err != nil {
			t.Errorf("[%d] error: %s", ncase, err)
		}
		if inst.op.code != item.inst.op.code {
			t.Errorf("[%d] Got op code %d  (expected %d)", ncase,
				inst.op.code, item.inst.op.code)
		}
		if inst.x != item.inst.x || inst.xval != item.inst.xval {
			t.Errorf("[%d] x=%s xval=%d  (expected %s, %d)", ncase,
				inst.x, inst.xval, item.inst.x, item.inst.xval)
		}
		if inst.y != item.inst.y || inst.yval != item.inst.yval {
			t.Errorf("[%d] y=%s yval=%d  (expected %s, %d)", ncase,
				inst.y, inst.yval, item.inst.y, item.inst.yval)
		}
	}
}

func TestExecute(t *testing.T) {
	source := []string{
		"cpy 41 a",
		"inc a",
		"inc a",
		"dec a",
		"jnz a 2",
		"dec a",
	}

	prog, err := Compile(source)
	if err != nil {
		t.Errorf("Program didnot compile")
	}
	expected := 42

	init := Registers{}
	reg, err := prog.Execute(init)
	if reg[0] != expected {
		t.Errorf("Register a had value %d  (expected %d)", reg[0], expected)
	}
}
