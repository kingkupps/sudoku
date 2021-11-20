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
	var rows [9][10]bool
	var cols [9][10]bool
	var boxes [9][10]bool
	for i, val := range g {
		r, c := rowColOf(i)
		b := box(r, c)
		rows[r][val] = true
		cols[c][val] = true
		boxes[b][val] = true
	}

	var empties []cell
	for i, val := range g {
		if val != 0 {
			continue
		}

		empty := newCell(i)
		for op := 1; op <= 9; op++ {
			if !rows[empty.row][op] && !cols[empty.col][op] && !boxes[empty.box][op] {
				empty.options[op] = true
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
		if opCount := cell.optionCount(); nextEmptyIndex == -1 || opCount < nextOptionsCount {
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
		if !nextCell.options[value] {
			continue
		}

		cellsToRestore := make([]int, 0, len(empties))
		for i := 0; i < len(empties)-1; i++ {
			if nextCell.sharesConstraintWith(empties[i]) {
				cellsToRestore = append(cellsToRestore, i)
				empties[i].options[value] = false
			}
		}

		if backtrackSolve(g, empties[:numEmpties-1]) == nil {
			g[nextCell.index()] = value
			return nil
		}

		for _, index := range cellsToRestore {
			empties[index].options[value] = true
		}
	}

	return ErrUnsolvable
}

func (g Game) IsSolved() bool {
	var rows [9][10]bool
	var cols [9][10]bool
	var boxes [9][10]bool
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
