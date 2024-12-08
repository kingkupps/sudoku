#ifndef SUDOKU_LIB
#define SUDOKU_LIB

#include <string>
#include <vector>

namespace sudoku
{
    class Sudoku
    {
    public:
        Sudoku(const std::string &config);
        std::string pformat() const;
        std::string config() const;
        bool solve();
        void set(size_t index, unsigned short option);

    private:
        std::vector<unsigned short> cells;
    };
}

#endif