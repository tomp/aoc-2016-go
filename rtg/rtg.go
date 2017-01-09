package rtg

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const NFLOORS int = 4

type State struct {
	elevator  int
	generator []int
	chip      []int
	isotopes  []string
	nisotopes int
}

type objType int

const (
	NOOBJ objType = iota
	GENERATOR
	CHIP
)

type Object struct {
	Typ objType
	Iso int
}

func InitialState(elevator int, generator []int, chip []int,
	isotopes []string) (s State, err error) {

	nisotopes := len(isotopes)
	if len(generator) != nisotopes {
		err = fmt.Errorf("%d generator values supplied for %d isotopes",
			len(generator), nisotopes)
	}
	if len(chip) != nisotopes {
		err = fmt.Errorf("%d chip values supplied for %d isotopes",
			len(chip), nisotopes)
	}
	if elevator < 1 || elevator > NFLOORS {
		err = fmt.Errorf("Elevator cannot be on floor %d", elevator)
		return
	}
	for i := 0; i < nisotopes; i++ {
		if generator[i] < 1 || generator[i] > NFLOORS {
			err = fmt.Errorf("Generator cannot be on floor %d", generator[i])
			return
		}
		if chip[i] < 1 || chip[i] > NFLOORS {
			err = fmt.Errorf("Generator cannot be on floor %d", chip[i])
			return
		}
	}
	s = State{elevator: elevator, generator: generator, chip: chip,
		isotopes: isotopes, nisotopes: nisotopes}
	return
}

func abs(val int) int {
	if val < 0 {
		return -val
	} else {
		return val
	}
}

func (s *State) validMove(elevator int, obj1 Object, obj2 Object) (err error) {
	if abs(elevator-s.elevator) > 1 || elevator < 1 || elevator > NFLOORS {
		err = fmt.Errorf("Elevator cannot move from floor %d to floor %d",
			s.elevator, elevator)
		return
	}
	if obj1.Typ == GENERATOR && s.generator[obj1.Iso] != s.elevator {
		err = fmt.Errorf("%sG cannot be moved to floor %d", s.isotopes[obj1.Iso], elevator)
		return
	}
	if obj1.Typ == CHIP && s.chip[obj1.Iso] != s.elevator {
		err = fmt.Errorf("%sM cannot be moved to floor %d", s.isotopes[obj1.Iso], elevator)
		return
	}
	if obj2.Typ == GENERATOR && s.generator[obj2.Iso] != s.elevator {
		err = fmt.Errorf("%sG cannot be moved to floor %d", s.isotopes[obj2.Iso], elevator)
		return
	}
	if obj2.Typ == CHIP && s.chip[obj2.Iso] != s.elevator {
		err = fmt.Errorf("%sM cannot be moved to floor %d", s.isotopes[obj2.Iso], elevator)
		return
	}
	return
}

func (s *State) MoveObjects(elevator int, obj1 Object, obj2 Object) (snew *State) {
	if err := s.validMove(elevator, obj1, obj2); err != nil {
		panic(err)
	}
	generator := make([]int, s.nisotopes)
	chip := make([]int, s.nisotopes)
	snew = &State{elevator, generator, chip, s.isotopes, s.nisotopes}
	for iso := 0; iso < s.nisotopes; iso++ {
		snew.generator[iso] = s.generator[iso]
		snew.chip[iso] = s.chip[iso]
	}
	if obj1.Typ == GENERATOR {
		snew.generator[obj1.Iso] = elevator
	} else if obj1.Typ == CHIP {
		snew.chip[obj1.Iso] = elevator
	}
	if obj2.Typ == GENERATOR {
		snew.generator[obj2.Iso] = elevator
	} else if obj2.Typ == CHIP {
		snew.chip[obj2.Iso] = elevator
	}
	return
}

func (s *State) String() string {
	lines := []string{}
	for f := NFLOORS; f > 0; f-- {
		parts := []string{fmt.Sprintf("F%d ", f)}
		if s.elevator == f {
			parts = append(parts, "E  ")
		} else {
			parts = append(parts, ".  ")
		}
		for iso, name := range s.isotopes {
			if s.generator[iso] == f {
				parts = append(parts, fmt.Sprintf("%sG ", name))
			} else {
				parts = append(parts, ".  ")
			}
			if s.chip[iso] == f {
				parts = append(parts, fmt.Sprintf("%sM ", name))
			} else {
				parts = append(parts, ".  ")
			}
		}
		lines = append(lines, strings.Join(parts, ""))
	}
	return strings.Join(lines, "\n")
}

func (s *State) Heuristic() int {
	dist := 0
	for iso := 0; iso < s.nisotopes; iso++ {
		dist += NFLOORS - s.generator[iso]
		dist += NFLOORS - s.chip[iso]
	}
	return dist
}

func (s *State) Hash() string {
	parts := []string{}
	for iso := range s.isotopes {
		parts = append(parts, strconv.Itoa(s.generator[iso])+":"+
			strconv.Itoa(s.chip[iso]))
	}
	sort.StringSlice(parts).Sort()
	parts = append(parts, strconv.Itoa(s.elevator))
	return strings.Join(parts, ",")
}

// Done returns true if all objects are on the top floor
func (s *State) Done() bool {
	for iso := 0; iso < s.nisotopes; iso++ {
		if s.chip[iso] < NFLOORS || s.generator[iso] < NFLOORS {
			return false
		}
	}
	return true
}

// Fried returns true if any chips would be fried in the current state.
func (s *State) Fried() bool {
	for iso := 0; iso < s.nisotopes; iso++ {
		if s.chip[iso] != s.generator[iso] {
			// unshielded chip found
			for iso2 := 0; iso2 < s.nisotopes; iso2++ {
				if s.generator[iso2] == s.chip[iso] {
					return true
				}
			}
		}
	}
	return false
}

type floorState struct {
	generator []int
	chip      []int
	shielded  []int
}

func (fs *floorState) String() string {
	parts := []string{"gens:"}
	for _, iso := range fs.generator {
		parts = append(parts, fmt.Sprintf(" %d", iso))
	}
	parts = append(parts, ", chips:")
	for _, iso := range fs.chip {
		parts = append(parts, fmt.Sprintf(" %d", iso))
	}
	parts = append(parts, ", shielded:")
	for _, iso := range fs.shielded {
		parts = append(parts, fmt.Sprintf(" %d", iso))
	}
	return strings.Join(parts, "")
}

func (s *State) getFloorState(floor int) (fs *floorState) {
	fs = &floorState{[]int{}, []int{}, []int{}}
	for iso := 0; iso < s.nisotopes; iso++ {
		if s.generator[iso] == floor {
			fs.generator = append(fs.generator, iso)
			if s.chip[iso] == floor {
				fs.shielded = append(fs.shielded, iso)
			}
		}
		if s.chip[iso] == floor {
			fs.chip = append(fs.chip, iso)
		}
	}
	return
}

// NextStates returns a list of all states we can reach from the current
// state.
func (s *State) NextStates(visited map[string]bool) (states []*State) {
	floor := s.elevator
	states = []*State{}
	if floor < NFLOORS {
		for _, snew := range s.NextStatesOnFloor(floor + 1) {
			newHash := snew.Hash()
			if !visited[newHash] && !snew.Fried() {
				visited[newHash] = true
				states = append(states, snew)
			}
		}
	}
	if floor > 1 {
		for _, snew := range s.NextStatesOnFloor(floor - 1) {
			newHash := snew.Hash()
			if !visited[newHash] && !snew.Fried() {
				visited[newHash] = true
				states = append(states, snew)
			}
		}
	}
	return
}

// NextStates returns a list of all states we can reach from the current
// state where the elevator is on the given new floor..
func (s *State) NextStatesOnFloor(newFloor int) (states []*State) {
	floor := s.elevator
	fs := s.getFloorState(floor)
	nchips := len(fs.chip)
	ngens := len(fs.generator)

	states = []*State{}
	for _, iso := range fs.shielded {
		states = append(states, s.MoveObjects(newFloor,
			Object{GENERATOR, iso}, Object{CHIP, iso}))
	}
	if nchips > 1 {
		for i := 0; i < nchips-1; i++ {
			for j := i + 1; j < nchips; j++ {
				states = append(states, s.MoveObjects(newFloor,
					Object{CHIP, fs.chip[i]}, Object{CHIP, fs.chip[j]}))
			}
		}
	}
	if ngens > 1 {
		for i := 0; i < ngens-1; i++ {
			for j := i + 1; j < ngens; j++ {
				states = append(states, s.MoveObjects(newFloor,
					Object{GENERATOR, fs.generator[i]},
					Object{GENERATOR, fs.generator[j]}))
			}
		}
	}
	if nchips > 0 {
		for i := 0; i < nchips; i++ {
			states = append(states, s.MoveObjects(newFloor,
				Object{CHIP, fs.chip[i]}, Object{}))
		}
	}
	if ngens > 0 {
		for i := 0; i < ngens; i++ {
			states = append(states, s.MoveObjects(newFloor,
				Object{GENERATOR, fs.generator[i]}, Object{}))
		}
	}
	return
}
