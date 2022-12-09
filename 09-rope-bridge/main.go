// go run main.go
package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputr.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	puzzle, err := parseFileToPuzzle(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Printf("> (1st Puzzle) Simulate your complete hypothetical series of motions. How many positions does the tail of the rope visit at least once?")

	puzzle.SimulatePuzzle()

	positionsVisitCount := puzzle.CountPositionsVisited()

	log.Printf("Number of positions visited by tail: %d", positionsVisitCount)

	puzzle, err = parseFileToPuzzle(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Printf("> (2nd Puzzle) Simulate your complete series of motions on a larger rope with ten knots. How many positions does the tail of the rope visit at least once?")

	puzzle.SimulatePuzzle2()

	positionsVisitCount = puzzle.CountPositionsVisited()

	log.Printf("Number of positions visited by the last tail: %d", positionsVisitCount)
}

const (
	Up    byte = 0
	Left  byte = 1
	Right byte = 2
	Down  byte = 3
)

type Move struct {
	Direction byte
	Hops      int
}

type Position byte

type Moves []Move

// type Grid [][]Position
type Grid map[int]map[int]Position

type Puzzle struct {
	Grid  Grid
	Moves Moves
}

func (p Puzzle) CountPositionsVisited() int {
	count := 0

	for _, pl := range p.Grid {
		for _, pc := range pl {
			if pc.visited() {
				count++
			}
		}
	}

	return count
}

func (p *Puzzle) SimulatePuzzle() {

	grid := p.Grid
	moves := p.Moves

	cl := 0
	cc := 0

	tailL := cl
	tailC := cc
	headL := cl
	headC := cc

	if grid[tailL] == nil {
		grid[tailL] = make(map[int]Position)
		grid[tailL][tailC] = 1
	}

	for _, m := range moves {

		for i := 0; i < m.Hops; i++ {
			cl = headL
			cc = headC

			switch m.Direction {
			case Up:
				headL--
			case Down:
				headL++
			case Left:
				headC--
			case Right:
				headC++
			}

			if tailCanBeMoved(tailL, tailC, headL, headC) {
				tailL = cl
				tailC = cc
				if grid[tailL] == nil {
					grid[tailL] = make(map[int]Position)
					grid[tailL][tailC] = 0
				}
				grid[tailL][tailC]++
			}
		}
	}
}

func (p *Puzzle) SimulatePuzzle2() {

	numberOfTails := 9
	tailsL := make([]int, numberOfTails)
	tailsC := make([]int, numberOfTails)

	grid := p.Grid
	moves := p.Moves

	cl := 0
	cc := 0

	tailL := cl
	tailC := cc
	headL := cl
	headC := cc

	if grid[tailL] == nil {
		grid[tailL] = make(map[int]Position)
		grid[tailL][tailC] = 1
	}

	for _, m := range moves {

		for i := 0; i < m.Hops; i++ {

			cl = headL
			cc = headC

			switch m.Direction {
			case Up:
				headL--
			case Down:
				headL++
			case Left:
				headC--
			case Right:
				headC++
			}

			pl := headL
			pc := headC

			for tail := 0; tail < numberOfTails; tail++ {
				tailL = tailsL[tail]
				tailC = tailsC[tail]

				if tailCanBeMoved(tailL, tailC, pl, pc) {

					tailsL[tail] = cl
					tailsC[tail] = cc
					tailL = cl
					tailC = cc
					// Only store the last one
					if tail == numberOfTails-1 {
						log.Println("** Storing **")
						if grid[tailL] == nil {
							grid[tailL] = make(map[int]Position)
							grid[tailL][tailC] = 0
						}
						grid[tailL][tailC]++
					}
				}
				pl = cl
				pc = cc
			}
		}
	}
}

func tailCanBeMoved(tailL int, tailC int, headL int, headC int) bool {
	// on top
	if tailC == headC && tailL == headL {
		return false
	}

	// up or down
	if tailC == headC && (tailL+1 == headL || tailL-1 == headL) {
		return false
	}

	// left or right
	if tailL == headL && (tailC+1 == headC || tailC-1 == headC) {
		return false
	}

	// diagonal
	if (tailL-1 == headL && tailC-1 == headC) || (tailL-1 == headL && tailC+1 == headC) {
		return false
	}

	// diagonal
	if (tailL+1 == headL && tailC-1 == headC) || (tailL+1 == headL && tailC+1 == headC) {
		return false
	}

	return true
}

func (p Position) visited() bool {
	return p != 0
}

func (ms Moves) toGrid() Grid {
	grid := make(Grid)

	return grid
}

func parseDirection(direction string) byte {
	dir := Down

	switch direction {
	case "U":
		dir = Up
	case "L":
		dir = Left
	case "R":
		dir = Right
	case "D":
		dir = Down
	}

	return dir
}

func lineToMove(l string) Move {
	lineSplit := strings.Split(l, " ")
	direction := lineSplit[0]
	fileSizeParse, _ := strconv.Atoi(lineSplit[1])

	return Move{
		Hops:      fileSizeParse,
		Direction: parseDirection(direction),
	}
}

func parseFileToPuzzle(filePath string) (*Puzzle, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	moves := make(Moves, len(lines))
	nLines := len(lines)

	for i := 0; i < nLines; i++ {
		moves[i] = lineToMove(lines[i])
	}

	grid := moves.toGrid()

	puzzle := Puzzle{
		Grid:  grid,
		Moves: moves,
	}

	return &puzzle, nil
}

func getFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	return fileLines, nil
}
