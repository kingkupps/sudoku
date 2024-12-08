#include <algorithm>
#include <iostream>
#include <iterator>
#include <map>
#include <set>
#include <string>

#include "sudoku/sudoku.h"

template <typename T>
bool set_has(const std::set<T> &s, T val)
{
    return s.find(val) != s.end();
}

class Cell
{
public:
    unsigned short index;
    unsigned short row;
    unsigned short col;
    unsigned short box;
    std::set<unsigned short> options;

    Cell(unsigned short index)
    {
        this->index = index;
        this->row = index / 9;
        this->col = index % 9;
        if (this->row < 3)
        {
            this->box = col / 3;
        }
        else if (this->row < 6)
        {
            this->box = (col / 3) + 3;
        }
        else
        {
            this->box = (col / 3) + 6;
        }
    }

    bool shares_constraint_with(Cell *other)
    {
        return this->row == other->row || this->col == other->col || this->box == other->box;
    }

    void init_options(const std::set<unsigned short> &rows,
                      const std::set<unsigned short> &cols,
                      const std::set<unsigned short> &boxes)
    {
        for (unsigned short i = 1; i <= 9; ++i)
        {
            if (!set_has(rows, i) && !set_has(cols, i) && !set_has(boxes, i))
            {
                this->options.insert(i);
            }
        }
    }

    /** For debugging.. */
    std::string to_string() const
    {
        std::string optionsPart = "";
        if (this->options.empty())
        {
            optionsPart += "EMPTY";
        }
        else
        {
            optionsPart += "[ ";
            for (auto op : this->options)
            {
                optionsPart += std::to_string(op);
                optionsPart += " ";
            }
            optionsPart += "]";
        }

        return "row = " + std::to_string(this->row) + ", col = " + std::to_string(this->col) + ", box = " + std::to_string(this->box) + ", options = " + optionsPart;
    }
};

bool backtrackSolve(sudoku::Sudoku *sudoku,
                    std::vector<Cell *> &empties)
{
    // No more empties, success!
    if (empties.empty())
    {
        return true;
    }

    // Find the cell with the fewest open options
    size_t nextIndex = 0;
    Cell *next = nullptr;
    for (size_t i = 0; i < empties.size(); ++i)
    {
        auto cell = empties[i];
        if (!next || cell->options.size() < next->options.size())
        {
            next = cell;
            nextIndex = i;
        }
    }

    // We've hit a cell with no option, we goofed, time to backtrack
    if (next->options.size() == 0)
    {
        return false;
    }

    // Remove chosen cell
    size_t last = empties.size() - 1;
    std::swap(empties[last], empties[nextIndex]);
    empties.pop_back();

    // Try to fill in the cell with an option and continue
    for (auto option : next->options)
    {
        std::vector<Cell *> affected;
        for (auto cell : empties)
        {
            if (next->shares_constraint_with(cell))
            {
                affected.push_back(cell);
                cell->options.erase(option);
            }
        }

        auto solved = backtrackSolve(sudoku, empties);
        if (solved)
        {
            sudoku->set(next->index, option);
            return true;
        }

        for (auto cell : affected)
        {
            cell->options.insert(option);
        }
    }

    // We we're able to solve past this point, backtrack
    empties.push_back(next);
    return false;
};

sudoku::Sudoku::Sudoku(const std::string &config)
{
    for (auto ch : config)
    {
        if (ch == '.')
        {
            this->cells.push_back(0);
        }
        else
        {
            this->cells.push_back(ch - '0');
        }
    }
}

std::string sudoku::Sudoku::pformat() const
{
    std::string display;
    for (unsigned short i = 0; i < 81; i++)
    {
        unsigned short row = i / 9;
        unsigned short col = i % 9;
        unsigned short val = this->cells[i];

        if (val == 0)
        {
            display.push_back('*');
        }
        else
        {
            display.push_back(val + '0');
        }

        display.push_back(' ');

        if (col == 2 || col == 5)
        {
            display.push_back(' ');
        }

        if (col == 8)
        {
            display.push_back('\n');
            if (row == 2 || row == 5)
            {
                display.push_back('\n');
            }
        }
    }
    return display;
}

std::string sudoku::Sudoku::config() const
{
    std::string config;
    for (auto val : this->cells)
    {
        if (val == 0)
        {
            config.push_back('.');
        }
        else
        {
            config.push_back('0' + val);
        }
    }
    return config;
}

bool sudoku::Sudoku::solve()
{
    std::map<unsigned short, std::set<unsigned short>> rows;
    std::map<unsigned short, std::set<unsigned short>> cols;
    std::map<unsigned short, std::set<unsigned short>> boxes;

    std::vector<Cell *> empties;
    for (size_t index = 0; index < this->cells.size(); ++index)
    {
        auto value = this->cells[index];
        auto cell = new Cell(index);
        if (value != 0)
        {
            rows[cell->row].insert(value);
            cols[cell->col].insert(value);
            boxes[cell->box].insert(value);
            delete cell;
        }
        else
        {
            empties.push_back(cell);
        }
    }

    for (auto cell : empties)
    {
        cell->init_options(rows[cell->row], cols[cell->col], boxes[cell->box]);
    }

    auto solved = backtrackSolve(this, empties);
    for (auto e : empties)
    {
        delete e;
    }
    return solved;
}

void sudoku::Sudoku::set(size_t index, unsigned short option)
{
    this->cells[index] = option;
}