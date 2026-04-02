const ALL_OPTIONS: u16 = 0b111111111;

pub struct OptionMaskIterator {
    mask: OptionMask,
    curr: u32,
}

impl Iterator for OptionMaskIterator {
    type Item = u32;

    fn next(&mut self) -> Option<Self::Item> {
        let mut curr = self.curr;
        while !self.mask.has(&curr) && curr <= 9 {
            curr += 1;
        }
        if curr > 9 {
            return None;
        }
        self.curr = curr + 1;
        return Some(curr);
    }
}

#[derive(Debug, Copy, Clone)]
pub struct OptionMask {
    mask: u16,
    size: u16,
}

impl OptionMask {
    pub fn all() -> Self {
        OptionMask {
            mask: ALL_OPTIONS,
            size: 9,
        }
    }

    #[allow(dead_code)]
    fn empty() -> Self {
        return OptionMask { mask: 0, size: 0 };
    }

    pub fn add(&mut self, option: &u32) {
        if self.has(option) {
            return;
        }
        self.size += 1;
        self.mask = self.mask | self.option_to_mask(option);
    }

    pub fn remove(&mut self, option: &u32) {
        if !self.has(option) {
            return;
        }
        self.size -= 1;
        self.mask = (self.mask ^ self.option_to_mask(option)) & ALL_OPTIONS;
    }

    pub fn only(&mut self, other: &OptionMask) {
        self.mask = (self.mask & other.mask) & ALL_OPTIONS;

        let mut copy = self.mask;
        let mut size = 0;
        while copy > 0 {
            if copy & 1 != 0 {
                size += 1;
            }
            copy = copy >> 1;
        }

        self.size = size;
    }

    pub fn has(&self, option: &u32) -> bool {
        self.option_to_mask(option) & self.mask != 0
    }

    pub fn len(&self) -> u16 {
        self.size
    }

    pub fn iter(&self) -> OptionMaskIterator {
        OptionMaskIterator {
            mask: OptionMask {
                mask: self.mask,
                size: self.size,
            },
            curr: 1,
        }
    }

    fn option_to_mask(&self, option: &u32) -> u16 {
        1 << (option - 1)
    }
}

#[cfg(test)]
mod tests {
    use crate::sudoku::mask::OptionMask;

    #[test]
    fn test_iteration() {
        let mut mask = OptionMask::empty();
        mask.add(&4);
        mask.add(&7);
        mask.add(&8);

        let mut vals = Vec::new();
        for val in mask.iter() {
            vals.push(val);
        }

        assert_eq!(3, mask.len());
        assert_eq!(vec![4, 7, 8], vals);
    }

    #[test]
    fn test_add_remove() {
        let mut mask = OptionMask::empty();
        assert_eq!(0, mask.len());

        mask.add(&1);
        mask.add(&8);
        mask.add(&9);
        mask.add(&2);
        assert_eq!(4, mask.len());
        assert!(mask.has(&1));
        assert!(mask.has(&9));
        assert!(mask.has(&2));
        assert!(mask.has(&8));

        mask.remove(&8);
        mask.remove(&1);
        assert_eq!(2, mask.len());
        assert!(!mask.has(&1));
        assert!(!mask.has(&8));
    }

    #[test]
    fn test_all_vals() {
        let mask = OptionMask::all();
        assert_eq!(9, mask.len());
        assert_eq!(
            vec![1, 2, 3, 4, 5, 6, 7, 8, 9],
            mask.iter().collect::<Vec<u32>>()
        )
    }

    #[test]
    fn test_only() {
        let mut first = OptionMask::empty();
        first.add(&1);
        first.add(&3);
        first.add(&5);
        first.add(&8);
        first.add(&9);
        let mut second = OptionMask::empty();
        second.add(&2);
        second.add(&3);
        second.add(&5);

        first.only(&second);

        assert_eq!(2, first.len());
        assert_eq!(vec![3, 5], first.iter().collect::<Vec<u32>>())
    }
}
