use super::mask::OptionMask;

#[derive(Debug)]
pub struct Cell {
    pub row: u32,
    pub col: u32,
    pub sec: u32,
    pub idx: u32,
    pub options: OptionMask,
}

impl Cell {
    pub fn from_idx(idx: u32) -> Self {
        let row = idx / 9;
        let col = idx % 9;
        let sec = if row < 3 {
            col / 3
        } else if row < 6 {
            (col / 3) + 3
        } else {
            (col / 3) + 6
        };
        Self {
            row,
            col,
            sec,
            idx,
            options: OptionMask::all(),
        }
    }

    pub fn shares_constraint_with(&self, other: &Cell) -> bool {
        self.row == other.row || self.col == other.col || self.sec == other.sec
    }
}
