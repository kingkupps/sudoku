use std::collections::HashSet;
use std::result::Result;

use super::cell::Cell;
use super::mask::OptionMask;

#[derive(Debug)]
pub struct Sudoku {
    board: Vec<u32>,
}

impl Sudoku {
    pub fn from(config: &str) -> Self {
        let mut board = Vec::new();
        for ch in config.chars() {
            board.push(if ch == '.' {
                0
            } else {
                ch.to_digit(10).unwrap()
            });
        }
        return Self { board };
    }

    pub fn solve(&self) -> Result<Sudoku, String> {
        let mut copy = self.board.clone();

        let mut empties = vec![];
        for (i, val) in copy.iter().enumerate() {
            if *val == 0 {
                empties.push(Cell::from_idx(i.try_into().unwrap()));
            }
        }

        let mut rv = [OptionMask::all(); 9];
        let mut cv = [OptionMask::all(); 9];
        let mut sv = [OptionMask::all(); 9];
        for (idx, val) in copy.iter().enumerate() {
            let cell = Cell::from_idx(idx.try_into().unwrap());
            if *val == 0 {
                continue;
            }
            rv[cell.row as usize].remove(val);
            cv[cell.col as usize].remove(val);
            sv[cell.sec as usize].remove(val);
        }

        for empty in empties.iter_mut() {
            let rmask = rv[empty.row as usize];
            let cmask = cv[empty.col as usize];
            let smask = sv[empty.sec as usize];
            empty.options.only(&rmask);
            empty.options.only(&cmask);
            empty.options.only(&smask);
        }

        let solved = solve(&mut copy, &mut empties);
        match solved {
            Ok(_) => Ok(Sudoku { board: copy }),
            Err(message) => Err(message),
        }
    }

    #[allow(dead_code)]
    fn to_config(&self) -> String {
        let mut s = String::new();
        for val in self.board.iter() {
            if *val == 0 {
                s.push_str(".");
            } else {
                s.push_str(val.to_string().as_str());
            }
        }
        s
    }

    pub fn pformat(&self) -> String {
        let mut out = String::new();
        for (i, val) in self.board.iter().enumerate() {
            if *val != 0 {
                out.push_str(u32::to_string(val).as_str());
            } else {
                out.push_str("*");
            }

            out.push_str(" ");

            let cell = Cell::from_idx(i.try_into().unwrap());
            if cell.col == 8 {
                out.push_str("\n");
                if cell.row == 2 || cell.row == 5 {
                    out.push_str("\n");
                }
            } else if cell.col == 2 || cell.col == 5 {
                out.push_str(" ");
            }
        }
        return out;
    }
}

fn solve(board: &mut Vec<u32>, empties: &mut Vec<Cell>) -> Result<(), String> {
    if empties.is_empty() {
        return Ok(());
    }

    // find the square with the fewest options remaining
    let mut lowest_options = empties.get(0).unwrap().options.len();
    let mut next_index = 0;
    for i in 1..empties.len() {
        let curr = empties.get(i).unwrap().options.len();
        if curr < lowest_options {
            lowest_options = curr;
            next_index = i;
        }
        if curr == 0 {
            return Err(String::from("no options left"));
        }
    }

    // pop our candidate
    let candidate = empties.swap_remove(next_index);

    for option in candidate.options.iter() {
        // remove this option from squares that share a constraint
        let mut to_update = HashSet::new();
        for empty in empties.iter_mut() {
            if empty.shares_constraint_with(&candidate) && empty.options.has(&option) {
                empty.options.remove(&option);
                to_update.insert(empty.idx);
            }
        }

        // try to fill in the remaining squares assuming we go with this option
        // for the candidate square
        let solved = solve(board, empties);
        if let Ok(_) = solved {
            let index: usize = candidate.idx.try_into().unwrap();
            board[index] = option;
            return Ok(());
        }

        // add back this option to squares that:
        //  1. share a constraint with the candidate
        //  2. had this option as a valid choice prior to the candidate using it
        for empty in empties.iter_mut() {
            if let Some(_) = to_update.get(&empty.idx) {
                empty.options.add(&option);
            }
        }
    }

    return Err(String::from("unable to solve :("));
}

#[cfg(test)]
mod tests {
    use super::Sudoku;

    #[test]
    fn test_solve() {
        let games = vec![
            (
                ".5..83.17...1..4..3.4..56.8....3...9.9.8245....6....7...9....5...729..861.36.72.4",
                "652483917978162435314975628825736149791824563436519872269348751547291386183657294",
            ),
            (
                "2.6.3......1.65.7..471.8.5.5......29..8.194.6...42...1....428..6.93....5.7.....13",
                "256734198891265374347198652514683729728519436963427581135942867689371245472856913",
            ),
            (
                "..45.21781...9..3....8....46..45.....7.9...128.12.35..4.......935..6.8.7.9.3..62.",
                "964532178187694235235817964629451783573986412841273596416728359352169847798345621",
            ),
            (
                "59....147...9....8.72....3.7...4.29..2..3.8.68..17..5...5764..9.36..5...1..8....2",
                "598326147314957628672481935753648291421539876869172453285764319936215784147893562",
            ),
            (
                "9...84.6.6.4..52.7.3..7..8.76...15...53.....1...4.96.31.5.26.9...2.4....8....371.",
                "927384165684915237531672489769231548453768921218459673175826394392147856846593712",
            ),
        ];

        for (to_solve, solution) in games {
            let s = Sudoku::from(to_solve);
            let result = s.solve();
            assert!(result.is_ok());
            assert_eq!(solution, result.unwrap().to_config());
        }
    }
}
