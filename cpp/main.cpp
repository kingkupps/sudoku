#include <iostream>
#include "sudoku/sudoku.h"

int main()
{
    sudoku::Sudoku board(".5..83.17...1..4..3.4..56.8....3...9.9.8245....6....7...9....5...729..861.36.72.4");

    bool solved = board.solve();
    if (!solved)
    {
        std::cout << "unsolved :(" << std::endl;
    }
    else
    {
        std::cout << board.pformat() << std::endl;
    }
}
