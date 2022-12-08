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
}

type Grid map[int]map[int]int

func (g Grid) ComputeNumberOfVisibleTrees() int {
	nLines := len(g)
	nColumns := len(g[0])
	visibleInteriorTrees := 0

	for l := 1; l < nLines-1; l++ {
		for c := 1; c < nColumns-1; c++ {
			element := g[l][c]
			log.Printf("LINHA %v COLUNA %v -> %v", l, c, element)

			if g.isVisibleFromLeft(l, c) {
				log.Printf("Tree is visible from left-> %v\n", element)
				visibleInteriorTrees++
				continue
			}
			if g.isVisibleFromRigth(l, c) {
				log.Printf("Tree is visible from rigth-> %v\n", element)
				visibleInteriorTrees++
				continue
			}
			if g.isVisibleFromTop(l, c) {
				log.Printf("Tree is visible from top-> %v\n", element)
				visibleInteriorTrees++
				continue
			}
			if g.isVisibleFromBottom(l, c) {
				log.Printf("Tree is visible from bottom-> %v\n", element)
				visibleInteriorTrees++
				continue
			}
		}
		log.Println("\n*****")
	}

	edges := nLines*2 + (nColumns-2)*2
	log.Printf("edges -> %v", edges)
	return edges + visibleInteriorTrees
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

func (g Grid) isVisibleFromRigth(line int, column int) bool {
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

	for i := line - 1; i < end; i++ {
		if g[i][column] >= treeToBeCompared {
			return false
		}
	}

	return true
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
