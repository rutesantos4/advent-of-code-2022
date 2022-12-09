// go run main.go
package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
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
	Up    = 0
	Left  = 1
	Right = 2
	Down  = 3
)

type Grid map[int]map[int]int

func (g Grid) ComputeNumberOfVisibleTrees() int {
	nLines := len(g)
	nColumns := len(g[0])
	visibleInteriorTrees := 0

	for l := 1; l < nLines-1; l++ {
		for c := 1; c < nColumns-1; c++ {

			if g.isVisibleFromLeft(l, c) {
				visibleInteriorTrees++
				continue
			}

			if g.isVisibleFromRight(l, c) {
				visibleInteriorTrees++
				continue
			}

			if g.isVisibleFromTop(l, c) {
				visibleInteriorTrees++
				continue
			}

			if g.isVisibleFromBottom(l, c) {
				visibleInteriorTrees++
				continue
			}

		}
	}

	edges := nLines*2 + (nColumns-2)*2
	return edges + visibleInteriorTrees
}

func (g Grid) ComputeHighestTreeScenicScore() int {
	nLines := len(g)
	nColumns := len(g[0])
	highestScore := 0

	for l := 0; l < nLines; l++ {
		for c := 0; c < nColumns; c++ {
			scenicScore := g.scenicScore(l, c)

			if scenicScore > highestScore {
				highestScore = scenicScore
			}
		}
	}

	return highestScore
}

func (g Grid) scenicScore(line int, column int) int {
	return g.countNonViewBlockingTrees(Up, line, column) * g.countNonViewBlockingTrees(Down, line, column) * g.countNonViewBlockingTrees(Left, line, column) * g.countNonViewBlockingTrees(Right, line, column)
}

func (g Grid) countNonViewBlockingTrees(side int, line int, column int) int {
	nbt := 0

	tree := g[line][column]
	stop := false
	l := line
	c := column
	maxl := len(g) - 1
	maxc := len(g[0]) - 1

	for !stop {
		if side == Up {
			l--
			stop = isSideTreeBlockingView(g[l][c], tree)
		} else if side == Left {
			c--
			stop = isSideTreeBlockingView(g[l][c], tree)
		} else if side == Right {
			c++
			stop = isSideTreeBlockingView(g[l][c], tree)
		} else {
			l++
			stop = isSideTreeBlockingView(g[l][c], tree)
		}

		if l < 0 || c < 0 || l > maxl || c > maxc {
			stop = true
		} else {
			nbt++
		}
	}

	return nbt
}

func (g Grid) isVisibleFromLeft(line int, column int) bool {
	treeToBeCompared := g[line][column]

	for i := column - 1; i >= 0; i-- {
		if g[line][i] >= treeToBeCompared {
			return false
		}
	}

	return true
}

func (g Grid) isVisibleFromRight(line int, column int) bool {
	treeToBeCompared := g[line][column]
	end := len(g[line])

	for i := column + 1; i < end; i++ {
		if g[line][i] >= treeToBeCompared {
			return false
		}
	}

	return true
}

func (g Grid) isVisibleFromTop(line int, column int) bool {
	treeToBeCompared := g[line][column]

	for i := line - 1; i >= 0; i-- {
		if g[i][column] >= treeToBeCompared {
			return false
		}
	}

	return true
}

func (g Grid) isVisibleFromBottom(line int, column int) bool {
	treeToBeCompared := g[line][column]
	end := len(g[line])

	for i := line + 1; i < end; i++ {
		if g[i][column] >= treeToBeCompared {
			return false
		}
	}

	return true
}

func isSideTreeBlockingView(sideTree int, tree int) bool {
	return sideTree >= tree
}

func parseFileToGrid(filePath string) (Grid, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	grid := make(map[int]map[int]int)
	nLines := len(lines)

	for i := 0; i < nLines; i++ {

		chars := []rune(lines[i])
		nChars := len(chars)
		grid[i] = make(map[int]int)

		for j := 0; j < nChars; j++ {
			intVar, _ := strconv.Atoi(string(chars[j]))
			grid[i][j] = intVar
		}
	}

	return grid, nil
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
