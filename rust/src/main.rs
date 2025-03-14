mod sudoku;

use crate::sudoku::sudoku::Sudoku;
use std::result::Result;

fn main() {
    let game = Sudoku::from(
        ".5..83.17...1..4..3.4..56.8....3...9.9.8245....6....7...9....5...729..861.36.72.4",
    );
    println!("Original\n{}", game.pformat());

    match game.solve() {
        Result::Ok(solved) => println!("Solved\n{}", solved.pformat()),
        Result::Err(message) => println!("Unsolvable game: {}\n{}", message, game.pformat()),
    }
}
