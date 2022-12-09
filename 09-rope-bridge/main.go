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

	grid, err := parseFileToGrid(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Printf("> (1st Puzzle) Consider your map; how many trees are visible from outside the grid?")

	visibleTrees := grid.ComputeNumberOfVisibleTrees()

	log.Printf("Number of trees visible from outside the grid is: %d", visibleTrees)

	log.Printf("> (2nd Puzzle) Consider each tree on your map. What is the highest scenic score possible for any tree?")

	highestScenicScore := grid.ComputeHighestTreeScenicScore()

	log.Printf("Highest scenic score is: %d", highestScenicScore)
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

func (p *Puzzle) SimulatePuzzle() {

	grid := p.Grid
	moves := p.Moves

	cl := len(grid) - 1
	cc := 0

	for _, m := range moves {
		for i := 0; i < m.Hops; i++ {
			if 
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
