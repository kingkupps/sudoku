package sudoku_test

import (
	"testing"

	"github.com/kingkupps/sudoku"
	"github.com/stretchr/testify/assert"
)

var games = []struct {
	game   string
	solved string
}{
	{
		game:   ".5..83.17...1..4..3.4..56.8....3...9.9.8245....6....7...9....5...729..861.36.72.4",
		solved: "652483917978162435314975628825736149791824563436519872269348751547291386183657294",
	},
	{
		game:   "2.6.3......1.65.7..471.8.5.5......29..8.194.6...42...1....428..6.93....5.7.....13",
		solved: "256734198891265374347198652514683729728519436963427581135942867689371245472856913",
	},
	{
		game:   "..45.21781...9..3....8....46..45.....7.9...128.12.35..4.......935..6.8.7.9.3..62.",
		solved: "964532178187694235235817964629451783573986412841273596416728359352169847798345621",
	},
	{
		game:   "59....147...9....8.72....3.7...4.29..2..3.8.68..17..5...5764..9.36..5...1..8....2",
		solved: "598326147314957628672481935753648291421539876869172453285764319936215784147893562",
	},
	{
		game:   "9...84.6.6.4..52.7.3..7..8.76...15...53.....1...4.96.31.5.26.9...2.4....8....371.",
		solved: "927384165684915237531672489769231548453768921218459673175826394392147856846593712",
	},
}

func TestSolve(t *testing.T) {
	for _, testCase := range games {
		t.Run(testCase.game, func(t *testing.T) {
			g, err := sudoku.LoadGame(testCase.game)
			assert.Nil(t, err)
			err = g.Solve()
			assert.Nil(t, err)
			assert.Equal(t, testCase.solved, g.String())
		})
	}
}

func TestUnsolvable(t *testing.T) {
	game, err := sudoku.LoadGame("9...8446.6.4..5227.3..7..8.76...15...53.....1...4.96.31.5.26.9...2.4....8....371.")
	assert.Nil(t, err)
	err = game.Solve()
	assert.ErrorIs(t, err, sudoku.ErrUnsolvable)
}

func TestIsSolved(t *testing.T) {
	for _, testCase := range games {
		t.Run(testCase.game, func(t *testing.T) {
			unsolved, err := sudoku.LoadGame(testCase.game)
			assert.Nil(t, err)
			assert.False(t, unsolved.IsSolved())

			solved, err := sudoku.LoadGame(testCase.solved)
			assert.Nil(t, err)
			assert.True(t, solved.IsSolved())
		})
	}
}

func TestPformat(t *testing.T) {
	game, err := sudoku.LoadGame("9...8446.6.4..5227.3..7..8.76...15...53.....1...4.96.31.5.26.9...2.4....8....371.")
	assert.Nil(t, err)

	want := "9 * *  * 8 4  4 6 *\n6 * 4  * * 5  2 2 7\n* 3 *  * 7 *  * 8 *\n\n7 6 *  * * 1  5 * *\n* 5 3  * * *  * * 1\n* * *  4 * 9  6 * 3\n\n1 * 5  * 2 6  * 9 *\n* * 2  * 4 *  * * *\n8 * *  * * 3  7 1 *\n"
	got := game.Pformat()
	assert.Equal(t, want, got)
}

func BenchmarkSolve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		game := sudoku.Game([81]int{0, 5, 0, 0, 8, 3, 0, 1, 7, 0, 0, 0, 1, 0, 0, 4, 0, 0, 3, 0, 4, 0, 0, 5, 6, 0, 8, 0, 0, 0, 0, 3, 0, 0, 0, 9, 0, 9, 0, 8, 2, 4, 5, 0, 0, 0, 0, 6, 0, 0, 0, 0, 7, 0, 0, 0, 9, 0, 0, 0, 0, 5, 0, 0, 0, 7, 2, 9, 0, 0, 8, 6, 1, 0, 3, 6, 0, 7, 2, 0, 4})
		game.Solve()
	}
}
