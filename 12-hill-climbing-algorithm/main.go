package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	heightmap, err := parseFileToHeightmap(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Println(heightmap)

	log.Printf("> (1st Puzzle) What is the fewest steps required to move from your current position to the location that should get the best signal?")

	possiblePaths := heightmap.FindPossiblePaths()

	println(possiblePaths)
	// fewestSteps := heightmap.PlayMonkeyInTheMiddleFor(1)

	// log.Printf("The number of fewest steps required to move from your current position to the location is: %d")

}

const (
	HeightDiffToGo = 1
)
const (
	StartPosition rune = 'S'
	EndPosition   rune = 'E'
)

type Path []Position

type Position struct {
	l int
	c int
}
type Height byte

type Heightmap [][]Height

func (h Heightmap) FindPossiblePaths() []Path {
	paths := []Path{}

	sp := h.findStartPosition()
	h[sp.l][sp.l] = Height('a')

	path := []Position{}
	findPath(h, sp, path)

	return append(paths, path)
}

func (h Heightmap) findStartPosition() Position {

	for l, heights := range h {
		for c, height := range heights {
			if height == Height(byte(StartPosition)) {
				return Position{
					l: l,
					c: c,
				}
			}
		}
	}

	return Position{} //Should never get here
}

func findPath(h Heightmap, currentPosition Position, path Path) {

	if currentPosition.l == 2 && currentPosition.c == 5 {
		log.Println("END!")
	}
	
	if h[currentPosition.l][currentPosition.c] == Height(EndPosition) {
		return
	}

	if positionVisited(path, currentPosition) {
		return
	}

	path = append(path, currentPosition)

	if possible, next := canGoToLeft(h, currentPosition); possible {
		findPath(h, next, path)
	}

	if possible, next := canGoToRight(h, currentPosition); possible {
		findPath(h, next, path)
	}

	if possible, next := canGoToUp(h, currentPosition); possible {
		findPath(h, next, path)
	}

	if possible, next := canGoToDown(h, currentPosition); possible {
		findPath(h, next, path)
	}
}

func canGoToLeft(h Heightmap, currentPosition Position) (bool, Position) {
	c := currentPosition.c - 1
	if c < 0 {
		return false, Position{}
	}

	next := Position{
		l: currentPosition.l,
		c: c,
	}

	if !h.isValidNextPosition(currentPosition, next) {
		return false, Position{}
	}

	return true, next
}

func canGoToRight(h Heightmap, currentPosition Position) (bool, Position) {
	c := currentPosition.c + 1
	if c >= len(h[0]) {
		return false, Position{}
	}

	next := Position{
		l: currentPosition.l,
		c: c,
	}

	if !h.isValidNextPosition(currentPosition, next) {
		return false, Position{}
	}

	return true, next
}

func canGoToUp(h Heightmap, currentPosition Position) (bool, Position) {
	l := currentPosition.l - 1
	if l < 0 {
		return false, Position{}
	}

	next := Position{
		l: l,
		c: currentPosition.c,
	}

	if !h.isValidNextPosition(currentPosition, next) {
		return false, Position{}
	}

	return true, next
}

func canGoToDown(h Heightmap, currentPosition Position) (bool, Position) {
	l := currentPosition.l + 1
	if l >= len(h) {
		return false, Position{}
	}

	next := Position{
		l: l,
		c: currentPosition.c,
	}

	if !h.isValidNextPosition(currentPosition, next) {
		return false, Position{}
	}

	return true, next
}

func (h Heightmap) isValidNextPosition(currentPosition Position, nextPosition Position) bool {

	currentHeight := h[currentPosition.l][currentPosition.c]
	nextHeight := h[nextPosition.l][nextPosition.c]

	return Height(math.Abs(float64(nextHeight)-float64(currentHeight))) <= Height(HeightDiffToGo)
}

func positionVisited(p Path, position Position) bool {

	for _, v := range p {
		if v == position {
			return true
		}
	}

	return false
}

func (h Heightmap) String() string {
	lines := len(h)
	columns := len(h[0])
	result := "\n"
	for l := 0; l < lines; l++ {
		for c := 0; c < columns; c++ {
			result += string(h[l][c])
		}
		result += "\n"
	}

	return result
}

func parseLineToHeighs(line string) []Height {
	count := len(line)
	heights := make([]Height, count)

	for i := 0; i < count; i++ {
		heights[i] = Height(line[i])
	}

	return heights
}

func parseFileToHeightmap(filePath string) (Heightmap, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	l := len(lines)
	heightmap := make(Heightmap, l)

	for i, line := range lines {
		heightmap[i] = parseLineToHeighs(line)
	}

	return heightmap, nil
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
