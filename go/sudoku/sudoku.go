package sudoku

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrInvalidConfig = errors.New("config must specify 81 cells and contain only 1-9 and \".\"")
var ErrUnsolvable = errors.New("puzzle is not solveable")

// Game represents a standard 9x9 Sudoku game.
type Game [81]int

func LoadGame(config string) (Game, error) {
	var g [81]int
	if len(config) != 81 {
		return g, ErrInvalidConfig
	}

	for i, rn := range config {
		if !unicode.IsDigit(rn) && rn != '.' {
			return g, ErrInvalidConfig
		}
		if rn != '.' {
			g[i] = int(rn - '0')
		} else {
			g[i] = 0
		}
	}
	return g, nil
}

// Solve fills all empty spaces of g or returns a non-nil error if g is unsolvable.
func (g *Game) Solve() error {
	rows, cols, boxes := valuesByGroup(), valuesByGroup(), valuesByGroup()
	for i, val := range g {
		r, c := rowColOf(i)
		b := box(r, c)
		rows[r][val] = true
		cols[c][val] = true
		boxes[b][val] = true
	}

	options := make(map[int]options)
	for i, val := range g {
		if val != 0 {
			continue
		}

		r, c := rowColOf(i)
		b := box(r, c)

		optionsForCell := make(map[int]bool)
		for op := 1; op <= 9; op++ {
			if !rows[r][op] && !cols[c][op] && !boxes[b][op] {
				optionsForCell[op] = true
			}
		}
		options[i] = optionsForCell
	}

	return backtrackSolve(g, options)
}

func backtrackSolve(g *Game, cellOptions map[int]options) error {
	if len(cellOptions) == 0 {
		return nil
	}

	var nextCell int
	var nextOps options
	for cell, cellOps := range cellOptions {
		if nextOps == nil || len(cellOps) < len(nextOps) {
			nextCell = cell
			nextOps = cellOps
		}
	}

	if len(nextOps) == 0 {
		return ErrUnsolvable
	}

	delete(cellOptions, nextCell)
	for op := range nextOps {
		var cellsToRestore []int
		for cell := range cellOptions {
			if sharesConstraintWith(nextCell, cell) {
				cellsToRestore = append(cellsToRestore, cell)
				delete(cellOptions[cell], op)
			}
		}

		if backtrackSolve(g, cellOptions) == nil {
			g[nextCell] = op
			return nil
		}

		for _, cell := range cellsToRestore {
			cellOptions[cell][op] = true
		}
	}
	cellOptions[nextCell] = nextOps

	return ErrUnsolvable
}

func (g Game) IsSolved() bool {
	rows, cols, boxes := valuesByGroup(), valuesByGroup(), valuesByGroup()
	for i, val := range g {
		r, c := rowColOf(i)
		b := box(r, c)

		if rows[r][val] {
			return false
		}
		rows[r][val] = true

		if cols[c][val] {
			return false
		}
		cols[c][val] = true

		if boxes[b][val] {
			return false
		}
		boxes[b][val] = true
	}
	return true
}

func valuesByGroup() []options {
	groups := make([]options, 10)
	for i := range groups {
		groups[i] = make(map[int]bool)
	}
	return groups
}

func box(r, c int) int {
	if r < 3 {
		return c / 3
	}
	if r < 6 {
		return (c / 3) + 3
	}
	return (c / 3) + 6
}

func rowColOf(index int) (int, int) {
	return (index / 9), (index % 9)
}

func sharesConstraintWith(first, second int) bool {
	r1, c1 := rowColOf(first)
	r2, c2 := rowColOf(second)
	if r1 == r2 || c1 == c2 {
		return true
	}
	b1 := box(r1, c1)
	b2 := box(r2, c2)
	return b1 == b2
}

// String returns a minimal string representation of the board, using "." for empty spaces.
func (g Game) String() string {
	var s strings.Builder
	for _, val := range g {
		if val == 0 {
			s.WriteRune('.')
		} else {
			s.WriteRune(rune('0' + val))
		}
	}
	return s.String()
}

func (g Game) Pformat() string {
	var s strings.Builder
	for i, val := range g {
		if val == 0 {
			s.WriteRune('*')
		} else {
			s.WriteRune(rune('0' + val))
		}

		r, c := rowColOf(i)
		if c != 8 {
			s.WriteRune(' ')
		} else {
			s.WriteRune('\n')
			if r == 2 || r == 5 {
				s.WriteRune('\n')
			}
		}

		if c == 2 || c == 5 {
			s.WriteRune(' ')
		}
	}
	return s.String()
}

func (g Game) Display() {
	fmt.Println(g.Pformat())
}

type options map[int]bool
