package main

import (
	"fmt"
	"log"

	"github.com/kingkupps/sudoku/sudoku"
)

func main() {
	s, err := sudoku.LoadGame(".5..83.17...1..4..3.4..56.8....3...9.9.8245....6....7...9....5...729..861.36.72.4")
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	s.Display()

	fmt.Printf("%#v\n", s)

	if err := s.Solve(); err != nil {
		log.Fatalf("%s\n", err)
	}
	s.Display()
}
