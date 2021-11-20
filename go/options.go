package sudoku

func numOptions(options int) int {
	uCount := options - ((options >> 1) & 033333333333) - ((options >> 2) & 011111111111)
	return ((uCount + (uCount >> 3)) & 030707070707) % 63
}

func hasOption(curr, op int) bool {
	op = 1 << op
	return curr&op != 0
}

func addOption(curr, op int) int {
	return curr | (1 << op)
}

func removeOption(curr, op int) int {
	return curr & (^(1 << op))
}
