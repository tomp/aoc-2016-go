package main

import (
	"testing"
)

func TestOpenDoors(t *testing.T) {
	cases := [...]struct {
		passcode string
		steps  History
		expected string
	}{
		{"hijkl", History{""}, "D"},
		{"hijkl", History{"D"}, "UR"},
		{"hijkl", History{"DR"}, ""},
		{"hijkl", History{"DU"}, "R"},
		{"hijkl", History{"DUR"}, ""},
	}

	for _, item := range cases {
		doors := openDoors(item.passcode, item.steps)
		if doors != item.expected {
			t.Errorf("passcode '%s' + history '%s' -> '%s'  (expected '%s')",
				item.passcode, item.steps, doors[:4], item.expected)
		}
	}
}

func TestPosition(t *testing.T) {
	cases := [...]struct {
		steps    History
		expected_x int
		expected_y int
	}{
		{History{""}, 0, 0},
		{History{"DRUL"}, 0, 0},
		{History{"DRRDLU"}, 1, 1},
	}

	for _, item := range cases {
		x, y := position(item.steps)
		if x != item.expected_x || y != item.expected_y {
			t.Errorf("(%d, %d) -- %s -> (%d, %d)  (expected %d, %d)",
				0, 0, item.steps, x, y,
				item.expected_x, item.expected_y)
		}
	}
}
