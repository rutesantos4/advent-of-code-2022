package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	path, err := parseFileToMap(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Printf("> (1st Puzzle) Using your scan, simulate the falling sand. How many units of sand come to rest before sand starts flowing into the abyss below?")

	numberSandBeforeAbyss, _ := path.DrawSand()

	log.Printf("The units of sand come to rest before sand starts flowing into the abyss is: %d\n", numberSandBeforeAbyss)

	log.Printf("> (2nd Puzzle) Using your scan, simulate the falling sand until the source of the sand becomes blocked. How many units of sand come to rest?")

	path.DrawFloor()
	numberSandBeforeAbyssFloor, _ := path.DrawSand()

	log.Printf("The units of sand come to rest until the source of the sand becomes blocked is: %d\n", numberSandBeforeAbyssFloor)
}

const (
	Rock      = "#"
	Air       = "."
	Sand      = "o"
	SandStart = "+"
)

const (
	FloorDiff = 2
)

type Element string
type Map [][]Element
type Position struct {
	l int
	c int
}

func (m Map) String() string {
	lines := len(m)
	columns := len(m[0])
	result := "\n"
	for l := 0; l < lines; l++ {
		for c := 0; c < columns; c++ {
			result += string(m[l][c])
		}
		result += "\n"
	}

	return result
}

func (m Map) DrawFloor() {
	height := m.highestPoint() + FloorDiff

	start := Position{
		l: height,
		c: 0,
	}
	end := Position{
		l: height,
		c: len(m[height]) - 1,
	}
	m.fillMapWithRocks(start, end)
}

func (m Map) DrawSand() (numberOfSands int, newMap Map) {

	sizel := len(m)
	newMap = make(Map, sizel)
	for l := 0; l < sizel; l++ {
		sizec := len(m[l])
		newMap[l] = make([]Element, sizec)
		copy(newMap[l], m[l])
	}

	start := Position{
		l: 0,
		c: 500,
	}

	height := m.highestPoint()
	current := start
	numberOfSands = 0
	isEnd := false

	for !isEnd {
		current = newMap.drawNextSand(start, height)
		if current.l >= height {
			isEnd = true
			continue
		}
		numberOfSands++
		if current == start {
			// The point of the sand pouring counts as 1 sand pouring, so this validation is after the incrementation
			isEnd = true
			continue
		}
	}
	return
}

func (m Map) highestPoint() int {
	sizel := len(m)
	height := 0
	for l := 0; l < sizel; l++ {
		sizec := len(m[l])
		for c := 0; c < sizec; c++ {
			if m[l][c] == Rock {
				height = l
				break
			}
		}
	}
	return height
}

func (m Map) drawNextSand(current Position, maxLines int) Position {
	searchPosition := true

	for searchPosition {

		if current.l+1 > maxLines {
			return current
		}

		if m[current.l+1][current.c] == Air {
			current.l = current.l + 1
			continue
		}

		if m[current.l+1][current.c-1] == Air {
			current.l = current.l + 1
			current.c = current.c - 1
			continue
		}

		if m[current.l+1][current.c+1] == Air {
			current.l = current.l + 1
			current.c = current.c + 1
			continue
		}

		m[current.l][current.c] = Sand
		searchPosition = false
	}

	return current
}

func (m Map) fillMapWithRocks(start Position, end Position) {
	if start.l == end.l {
		count := int(math.Abs(float64(end.c - start.c)))
		startfor := start.c
		if start.c > end.c {
			startfor = end.c
		}
		for i := startfor; i <= startfor+count; i++ {
			m[start.l][i] = Rock
		}
	} else if start.c == end.c {
		count := int(math.Abs(float64(end.l - start.l)))
		startfor := start.l
		if start.l > end.l {
			startfor = end.l
		}
		for i := startfor; i <= startfor+count; i++ {
			m[i][start.c] = Rock
		}
	}
}

func (m Map) fillMapWithSand() {
	m[0][500] = SandStart
}

func (m Map) fillMapWithAir() {
	l := len(m)
	c := len(m[0])

	for i := 0; i < l; i++ {
		for j := 0; j < c; j++ {
			if m[i][j] != Rock && m[i][j] != SandStart {
				m[i][j] = Air
			}
		}
	}
}

func fillMapWithLine(line string, m *Map) {
	points := strings.Split(line, " -> ")
	count := len(points)

	for i := 0; i < count-1; i++ {

		pointStart := strings.Split(points[i], ",")
		xS, _ := strconv.Atoi(pointStart[0])
		yS, _ := strconv.Atoi(pointStart[1])
		start := Position{
			l: yS,
			c: xS,
		}

		pointEnd := strings.Split(points[i+1], ",")
		xE, _ := strconv.Atoi(pointEnd[0])
		yE, _ := strconv.Atoi(pointEnd[1])
		end := Position{
			l: yE,
			c: xE,
		}

		m.fillMapWithRocks(start, end)
	}
}

func parseFileToMap(filePath string) (Map, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	size := 1000
	var m = make(Map, size)

	for i := 0; i < size; i++ {
		m[i] = make([]Element, size)
	}

	count := len(lines)
	for i := 0; i < count; i++ {
		fillMapWithLine(lines[i], &m)
	}

	m.fillMapWithSand()
	m.fillMapWithAir()

	return m, nil
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
