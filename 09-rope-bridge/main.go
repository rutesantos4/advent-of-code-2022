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
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
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

type Grid [][]Position

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

	cl := len(grid) - 1
	cc := 0
	maxl := cl + 1
	maxc := len(grid[0])

	for _, m := range moves {
		canMarkVisited := false
		isHeadOnTail := true

		for i := 0; i < m.Hops; i++ {
			if canMarkVisited && !isHeadOnTail {
				grid[cl][cc].markVisited()
			}

			// todo: falta validar se H está em cima de T

			switch m.Direction {
			case Up:
				cl--
			case Down:
				cl++
			case Left:
				cc--
			case Right:
				cc++
			}

			if cl == maxl {
				cl--
			} else if cc == maxc {
				cc--
			} else if cl < 0 {
				cl = 0
			} else if cc < 0 {
				cc = 0
			} else if !canMarkVisited {
				canMarkVisited = true
			}

			canMarkVisited = true
		}
	}
}

func (p Position) visited() bool {
	return p != 0
}

func (p *Position) markVisited() {
	*p++
}

func (ms Moves) toGrid() Grid {
	l := 0
	c := 0

	for _, m := range ms {
		if (m.Direction == Up || m.Direction == Down) && m.Hops > l {
			l = m.Hops
		} else if (m.Direction == Left || m.Direction == Right) && m.Hops > c {
			c = m.Hops
		}
	}

	grid := make(Grid, l)

	for i, _ := range grid {
		grid[i] = make([]Position, c)
	}

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
