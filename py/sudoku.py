from __future__ import annotations

from collections import defaultdict
from mask import Mask
from typing import Dict, List, Tuple


class Sudoku:

    BOARD_SIZE = 81

    def __init__(self, game: List[int]) -> None:
        self._game = game

    @property
    def is_solved(self) -> bool:
        return not any(i == 0 for i in self._game)

    def solve(self) -> bool:
        rows = defaultdict(Mask)
        cols = defaultdict(Mask)
        boxes = defaultdict(Mask)
        for i, val in enumerate(self._game):
            if val == 0:
                continue
            r, c = to_row_col(i)
            rows[r].add(val)
            cols[c].add(val)
            boxes[to_box_num(r, c)].add(val)

        possibilities = {}
        for i, val in enumerate(self._game):
            if val != 0:
                continue

            r, c = to_row_col(i)
            b = to_box_num(r, c)

            row_mask, col_mask, box_mask = rows[r], cols[c], boxes[b]
            m = Mask()
            possibilities[i] = m
            for val in range(1, 10):
                if val not in row_mask and val not in col_mask and val not in box_mask:
                    m.add(val)

        return self._backtrack(possibilities)

    def _backtrack(self, to_fill: Dict[int, Mask]) -> bool:
        if not to_fill:
            return True

        next_to_fill, options = min(
            to_fill.items(),
            key=lambda entry: len(entry[1])
        )
        if len(options) == 0:
            return False

        to_fill.pop(next_to_fill)
        for op in options:
            to_update = []
            for cell, mask in to_fill.items():
                if cells_collide(cell, next_to_fill) and op in mask:
                    mask.pop(op)
                    to_update.append(mask)

            if self._backtrack(to_fill):
                self._game[next_to_fill] = op
                return True

            for mask in to_update:
                mask.add(op)

        to_fill[next_to_fill] = options
        return False

    def pformat(self) -> str:
        out = []
        for i, val in enumerate(self._game):
            r, c = to_row_col(i)
            if val != 0:
                out.append(str(val))
            else:
                out.append('*')

            if c < 8:
                out.append(' ')
            else:
                out.append('\n')

            if c in [2, 5]:
                out.append(' ')

            if c == 8 and r in [2, 5]:
                out.append('\n')

        return ''.join(out)

    def pprint(self) -> None:
        print(self.pformat())

    def __str__(self) -> str:
        return ''.join([str(val) if val != 0 else '.' for val in self._game])

    @classmethod
    def parse(cls, config: str) -> Sudoku:
        if len(config) != cls.BOARD_SIZE:
            raise ValueError('Sudoku config must be 81 characters long.')

        game = [int(i) if i != '.' else 0 for i in config]
        return Sudoku(game)


def to_row_col(index: int) -> Tuple[int, int]:
    return index // 9, index % 9


def to_box_num(row: int, col: int) -> int:
    if row < 3:
        return col // 3
    if row < 6:
        return (col // 3) + 3
    return (col // 3) + 6


def cells_collide(first: int, second: int) -> bool:
    r1, c1 = to_row_col(first)
    r2, c2 = to_row_col(second)
    b1, b2 = to_box_num(r1, c1), to_box_num(r2, c2)
    return r1 == r2 or c1 == c2 or b1 == b2
