package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	heightmap, err := parseFileToHeightPositionMap(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Printf("> (1st Puzzle) What is the fewest steps required to move from your current position to the location that should get the best signal?")

	shp := heightmap.StartPosition()
	ehp := heightmap.EndPosition()

	shortestPath := heightmap.FindShortestPath(shp, ehp)

	log.Printf("The number of fewest steps required to move from your current position to the location is: %d", len(shortestPath)-1)

	log.Printf("> (2nd Puzzle) What is the fewest steps required to move starting from any square with elevation a to the location that should get the best signal?")

	lhps := heightmap.PositionsByHeight(LowestPositionHeight)

	shortestPath = heightmap.FindShortestPath2(lhps, ehp)

	log.Printf("The number of fewest steps required to move from your current position to the location is: %d", len(shortestPath)-1)

}

const (
	StartPositionHeight   = Height('S')
	EndPositionHeight     = Height('E')
	LowestPositionHeight  = Height('a')
	HighestPositionHeight = Height('z')
)

type Height byte

type HeightPosition struct {
	X      int
	Y      int
	Height Height
}

type Path []HeightPosition

type HeightPositionMap [][]HeightPosition

func (hm HeightPositionMap) StartPosition() HeightPosition {
	return hm.positionByHeight(StartPositionHeight)
}

func (hm HeightPositionMap) EndPosition() HeightPosition {
	return hm.positionByHeight(EndPositionHeight)
}

func (hm HeightPositionMap) positionByHeight(h Height) HeightPosition {
	var hp *HeightPosition

	for _, l := range hm {
		if hp == nil {
			for _, php := range l {
				if php.Height == h {
					hp = &php

					break
				}
			}
		}
	}

	return *hp
}

func (hm HeightPositionMap) PositionsByHeight(h Height) []HeightPosition {

	hp := []HeightPosition{}

	for _, l := range hm {
		for _, php := range l {
			if php.Height == h || (php.Height == StartPositionHeight && h == LowestPositionHeight) {
				hp = append(hp, php)
			}
		}

	}

	return hp
}

func (hm HeightPositionMap) FindShortestPath(shp, ehp HeightPosition) Path {
	possiblePaths := hm.findAllPaths(shp, ehp)

	shortestPath := possiblePaths[0]

	lsp := len(shortestPath)

	for _, p := range possiblePaths {
		if len(p) < lsp {
			shortestPath = p
		}
	}

	return shortestPath
}

func (hm HeightPositionMap) FindShortestPath2(hps []HeightPosition, ehp HeightPosition) Path {
	possiblePaths := []Path{}

	for _, hp := range hps {
		possiblePaths = append(possiblePaths, hm.FindShortestPath(hp, ehp))
	}

	shortestPath := possiblePaths[0]

	lsp := len(shortestPath)

	for _, p := range possiblePaths {
		if len(p) < lsp {
			shortestPath = p
		}
	}

	return shortestPath
}

func (hm HeightPositionMap) findAllPaths(shp, ehp HeightPosition) []Path {
	possiblePaths := []Path{}
	possiblePathsQueue := []Path{}

	possiblePathsQueue = append(possiblePathsQueue, Path{shp})

	for {
		ppqLen := len(possiblePathsQueue)

		foundAnyMovablePosition := false

		for i := 0; i < ppqLen; i++ {
			np := possiblePathsQueue[i]

			chp := np[len(np)-1]
			mp := hm.MovablePositions(chp, np)

			if len(mp) > 0 {
				foundAnyMovablePosition = true

				for _, hmp := range mp[1:] {
					npCpy := make(Path, len(np))
					copy(npCpy, np)
					npCpy = append(npCpy, hmp)

					possiblePathsQueue = append(possiblePathsQueue, npCpy)
				}

				np = append(np, mp[0])
				possiblePathsQueue[i] = np

			} else if chp != ehp {
				possiblePathsQueue = append(possiblePathsQueue[:i], possiblePathsQueue[i+1:]...)

				ppqLen--
			}
		}

		if !foundAnyMovablePosition {
			break
		}
	}

	for _, p := range possiblePathsQueue {
		if p[len(p)-1] == ehp {
			possiblePaths = append(possiblePaths, p)
		}
	}

	return possiblePaths
}

func (hm HeightPositionMap) MovablePositions(hp HeightPosition, exhp []HeightPosition) []HeightPosition {
	positions := []HeightPosition{}

	minX := 0
	minY := 0

	maxX := len(hm) - 1
	maxY := len(hm[0]) - 1

	x := hp.X
	y := hp.Y

	if x > minX {
		positions = append(positions, hm[x-1][y])
	}

	if x < maxX {
		positions = append(positions, hm[x+1][y])
	}

	if y > minY {
		positions = append(positions, hm[x][y-1])
	}

	if y < maxY {
		positions = append(positions, hm[x][y+1])
	}

	movablePositions := []HeightPosition{}

	for _, hpp := range positions {
		if doesNotRequireGearChange(hp, hpp) && doesNotExist(hpp, exhp) {
			movablePositions = append(movablePositions, hpp)
		}
	}

	return movablePositions
}

func (hm HeightPositionMap) String() string {
	lines := len(hm)
	columns := len(hm[0])
	result := "\n"
	for l := 0; l < lines; l++ {
		for c := 0; c < columns; c++ {
			result += string(hm[l][c].Height)
		}
		result += "\n"
	}

	return result
}

func doesNotRequireGearChange(hpSrc, hpDest HeightPosition) bool {

	hpSrcHeight := hpSrc.Height
	hpDestHeight := hpDest.Height

	if hpDestHeight == EndPositionHeight {
		hpDestHeight = HighestPositionHeight
	} else if hpSrcHeight == StartPositionHeight {
		hpSrcHeight = LowestPositionHeight
	}

	diff := int(hpDestHeight) - int(hpSrcHeight)

	return diff <= 1
}

func doesNotExist(hp HeightPosition, exhp []HeightPosition) bool {
	for _, hpp := range exhp {
		if hp.X == hpp.X && hp.Y == hpp.Y {
			return false
		}
	}

	return true
}

func parseLineToHeightPositions(line string, l int) []HeightPosition {
	count := len(line)
	heights := make([]HeightPosition, count)

	for i := 0; i < count; i++ {
		heights[i] = HeightPosition{
			X:      l,
			Y:      i,
			Height: Height(line[i]),
		}
	}

	return heights
}

func parseFileToHeightPositionMap(filePath string) (HeightPositionMap, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	l := len(lines)
	heightmap := make(HeightPositionMap, l)

	for i, line := range lines {
		heightmap[i] = parseLineToHeightPositions(line, i)
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
