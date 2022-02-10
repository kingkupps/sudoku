import unittest

from parameterized import param, parameterized
from sudoku import Sudoku, to_box_num, to_row_col


class TestSudoku(unittest.TestCase):

    @parameterized.expand([
        param(
            game=".5..83.17...1..4..3.4..56.8....3...9.9.8245....6....7...9....5...729..861.36.72.4",
            want="652483917978162435314975628825736149791824563436519872269348751547291386183657294",
        ),
        param(
            game="2.6.3......1.65.7..471.8.5.5......29..8.194.6...42...1....428..6.93....5.7.....13",
            want="256734198891265374347198652514683729728519436963427581135942867689371245472856913",
        ),
        param(
            game="..45.21781...9..3....8....46..45.....7.9...128.12.35..4.......935..6.8.7.9.3..62.",
            want="964532178187694235235817964629451783573986412841273596416728359352169847798345621",
        ),
        param(
            game="59....147...9....8.72....3.7...4.29..2..3.8.68..17..5...5764..9.36..5...1..8....2",
            want="598326147314957628672481935753648291421539876869172453285764319936215784147893562",
        ),
        param(
            game="9...84.6.6.4..52.7.3..7..8.76...15...53.....1...4.96.31.5.26.9...2.4....8....371.",
            want="927384165684915237531672489769231548453768921218459673175826394392147856846593712",
        )
    ])
    def test_solve_positive(self, game, want):
        puzzle = Sudoku.parse(game)
        self.assertTrue(puzzle.solve())
        self.assertEqual(want, str(puzzle))

    def test_solve_negative(self):
        puzzle = Sudoku.parse(
            '9...8446.6.4..5227.3..7..8.76...15...53.....1...4.96.31.5.26.9...2.4....8....371.')
        puzzle.solve()
        self.assertFalse(puzzle.is_solved)

    def test_pformat(self):
        puzzle = Sudoku.parse(
            '.5..83.17...1..4..3.4..56.8....3...9.9.8245....6....7...9....5...729..861.36.72.4')
        self.assertEqual(
            '* 5 *  * 8 3  * 1 7\n'
            '* * *  1 * *  4 * *\n'
            '3 * 4  * * 5  6 * 8\n'
            '\n'
            '* * *  * 3 *  * * 9\n'
            '* 9 *  8 2 4  5 * *\n'
            '* * 6  * * *  * 7 *\n'
            '\n'
            '* * 9  * * *  * 5 *\n'
            '* * 7  2 9 *  * 8 6\n'
            '1 * 3  6 * 7  2 * 4\n',
            puzzle.pformat()
        )

    @parameterized.expand([
        param(
            index=0,
            want=(0, 0)
        ),
        param(
            index=14,
            want=(1, 5)
        ),
        param(
            index=80,
            want=(8, 8)
        )
    ])
    def test_index_to_row_col(self, index=None, want=None):
        self.assertEqual(want, to_row_col(index))

    @parameterized.expand([
        param(
            row=2,
            col=1,
            want=0
        ),
        param(
            row=8,
            col=8,
            want=8
        ),
        param(
            row=5,
            col=6,
            want=5
        )
    ])
    def test_to_box_num(self, row=None, col=None, want=None):
        self.assertEqual(want, to_box_num(row, col))


if __name__ == '__main__':
    unittest.main()
