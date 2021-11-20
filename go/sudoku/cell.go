package sudoku

type cell struct {
	row     int
	col     int
	box     int
	options [10]bool
}

func newCell(index int) cell {
	r := index / 9
	c := index % 9

	var b int
	if r < 3 {
		b = c / 3
	} else if r < 6 {
		b = (c / 3) + 3
	} else {
		b = (c / 3) + 6
	}

	return cell{row: r, col: c, box: b}
}

func (c cell) index() int {
	return (c.row * 9) + c.col
}

func (c cell) sharesConstraintWith(o cell) bool {
	return c.row == o.row || c.col == o.col || c.box == o.box
}

func (c cell) optionCount() int {
	count := 0
	for i := 1; i <= 9; i++ {
		if c.options[i] {
			count++
		}
	}
	return count
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
