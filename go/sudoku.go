// Package sudoku implements a type for visualizing and solving 9x9 Sudoku puzzles.
package sudoku

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// ErrInvalidConfig is returned if sudoku.LoadGame is passed a string with
// invalid symbols or a length not equal to 81.
var ErrInvalidConfig = errors.New("config must specify 81 cells and contain only 1-9 and \".\"")

// ErrUnsolvable is returned by Game.Solve
var ErrUnsolvable = errors.New("puzzle is not solveable")

// Game represents a standard 9x9 Sudoku game.
type Game [81]int

// LoadGame returns a Game from a flattened string representation of a Sudoku
// puzzle. See sudoku_test.go for examples of the format.
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

// Solve fills all empty spaces of g or returns a non-nil error if g is
// unsolvable.
func (g *Game) Solve() error {
	var rows [9]int
	var cols [9]int
	var boxes [9]int
	for i, val := range g {
		r, c := rowColOf(i)
		b := box(r, c)
		rows[r] = addOption(rows[r], val)
		cols[c] = addOption(cols[c], val)
		boxes[b] = addOption(boxes[b], val)
	}

	var empties []cell
	for i, val := range g {
		if val != 0 {
			continue
		}

		empty := newCell(i)
		for op := 1; op <= 9; op++ {
			if !hasOption(rows[empty.row], op) && !hasOption(cols[empty.col], op) && !hasOption(boxes[empty.box], op) {
				empty.options = addOption(empty.options, op)
			}
		}
		empties = append(empties, empty)
	}

	return backtrackSolve(g, empties)
}

func backtrackSolve(g *Game, empties []cell) error {
	numEmpties := len(empties)
	if numEmpties == 0 {
		return nil
	}

	// Pick the empty cell with the fewest remaining possibilities
	nextEmptyIndex := -1
	nextOptionsCount := -1
	for i, cell := range empties {
		if opCount := numOptions(cell.options); nextEmptyIndex == -1 || opCount < nextOptionsCount {
			nextEmptyIndex = i
			nextOptionsCount = opCount
		}
	}

	// We took a wrong turn and need to backtrack
	if nextOptionsCount == 0 {
		return ErrUnsolvable
	}

	// Moves the next candidate cell to the end of empties so we can
	// efficiently pass the remaining empty cells to further calls
	empties[nextEmptyIndex], empties[numEmpties-1] = empties[numEmpties-1], empties[nextEmptyIndex]

	// Try to fill nextCell with each avaialble option
	nextCell := empties[numEmpties-1]
	for value := 1; value <= 9; value++ {
		if !hasOption(nextCell.options, value) {
			continue
		}

		cellsToRestore := make([]int, 0, len(empties))
		for i := 0; i < len(empties)-1; i++ {
			if nextCell.sharesConstraintWith(empties[i]) {
				cellsToRestore = append(cellsToRestore, i)
				empties[i].options = removeOption(empties[i].options, value)
			}
		}

		if backtrackSolve(g, empties[:numEmpties-1]) == nil {
			g[nextCell.index()] = value
			return nil
		}

		for _, index := range cellsToRestore {
			empties[index].options = addOption(empties[index].options, value)
		}
	}

	return ErrUnsolvable
}

// IsSolved returns whether g is correctly filled out according to the rules
// of Sudoku.
func (g Game) IsSolved() bool {
	var rows [9]int
	var cols [9]int
	var boxes [9]int
	for i, val := range g {
		r, c := rowColOf(i)
		b := box(r, c)

		if hasOption(rows[r], val) {
			return false
		}
		rows[r] = addOption(rows[r], val)

		if hasOption(cols[c], val) {
			return false
		}
		cols[c] = addOption(cols[c], val)

		if hasOption(boxes[b], val) {
			return false
		}
		boxes[b] = addOption(boxes[b], val)
	}
	return true
}

// String returns a minimal string representation of the board, using "." for
// empty spaces.
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

// Pformat returns a string representation of g that's closer to the typical
// human representation of a Sudoku puzzle (i.e. 9x9 grid).
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

// Display is equivalent to fmt.Println(g.Pformat()).
func (g Game) Display() {
	fmt.Println(g.Pformat())
}
