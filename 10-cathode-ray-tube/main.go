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

	program, err := parseFileToProgram(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Printf("> (1st Puzzle) Find the signal strength during the 20th, 60th, 100th, 140th, 180th, and 220th cycles. What is the sum of these six signal strengths?")

	program.executeCycles()

	cycles := []int{20, 60, 100, 140, 180, 220}
	sumOfCycleSignalStrengths := program.ComputeCyclesSignalStrengthSum(cycles)

	log.Printf("The sum of these six signal strengths is : %d", sumOfCycleSignalStrengths)

	log.Printf("> (2nd Puzzle) Render the image given by your program. What eight capital letters appear on your CRT?")

	log.Printf("CRT Image rendered \n%v\n", program.CRTImage)
}

const (
	Point string = "."
	Hash  string = "#"
)

const (
	CRTWide int = 40
	CRTHigh int = 6

	SpriteWide int = 3
)

const (
	AddxCycles int = 2
	NoopCycles int = 1
)

const (
	Addx byte = 0
	Noop byte = 1
)

type Program struct {
	Instructions []Instruction
	XValues      []XValue
	CRTImage     Grid
}

type Instruction struct {
	Type          byte
	IncreaseValue int
}

type XValue struct {
	Cycle  int
	During int
	After  int
}
type Grid [][]string
type Sprite []string

func (p Program) ComputeCyclesSignalStrengthSum(cycles []int) int {
	xValues := p.XValues

	sum := 0
	for _, cycle := range cycles {
		sum += xValues[cycle-1].ComputeCycleSignalStrength()
	}

	return sum
}

func (p *Program) executeCycles() {

	var xValues []XValue
	grid := buildGrid()
	sprite := buildSprite()

	cycle := 0
	xDuring := 1
	xAfter := 1

	for _, instruction := range p.Instructions {
		instructionDuration := instruction.Duration()

		for i := 0; i < instructionDuration; i++ {
			cycle++

			grid.changeValue(cycle, sprite)

			xValues = append(xValues, XValue{
				Cycle:  cycle,
				During: xDuring,
				After:  xAfter,
			})

			if i == instructionDuration-2 { //Increment before the last one
				xAfter += instruction.IncreaseValue
			}
		}
		sprite.shiftPositions(instruction.IncreaseValue)

		xDuring = xValues[cycle-1].After
		xAfter = xDuring
	}

	p.XValues = xValues
	p.CRTImage = grid
}

func (g Grid) changeValue(cycle int, sprite Sprite) {
	l := int(cycle / CRTWide)
	remaining := cycle % CRTWide
	c := remaining - 1

	if remaining == 0 {
		c = CRTWide - 1
		l--
	}

	v := sprite[c]
	g[l][c] = v
}

func (g Grid) String() string {
	lines := len(g)
	columns := len(g[0])
	result := "\n"
	for l := 0; l < lines; l++ {
		for c := 0; c < columns; c++ {
			result += g[l][c]
		}
		result += "\n"
	}

	return result
}

func (s *Sprite) shiftPositions(moves int) {

	//right rotation
	i := len(*s) - moves

	if moves < 0 {
		//left rotation
		i = int(math.Abs(float64(moves)))
	}

	x, b := (*s)[:i], (*s)[i:]
	*s = append(b, x...)
}

func (s Sprite) String() string {

	result := "\n"
	for _, v := range s {
		result += v
	}

	return result
}

func buildGrid() Grid {
	l := CRTHigh
	c := CRTWide
	grid := make(Grid, l)

	for i := range grid {
		grid[i] = make([]string, c)
	}

	return grid
}

func buildSprite() Sprite {
	sprite := make(Sprite, CRTWide)

	for i := 0; i < SpriteWide; i++ {
		sprite[i] = Hash
	}

	for i := SpriteWide; i < CRTWide; i++ {
		sprite[i] = Point
	}

	return sprite
}

func (i Instruction) Duration() int {

	instructionDuration := 0

	switch i.Type {
	case Noop:
		instructionDuration = NoopCycles
	case Addx:
		instructionDuration = AddxCycles
	}

	return instructionDuration
}

func (x XValue) ComputeCycleSignalStrength() int {
	return x.Cycle * x.During
}

func lineToInstruction(l string) Instruction {
	lineSplit := strings.Split(l, " ")
	insType := lineSplit[0]
	constType := Noop
	increaseValue := 0

	switch insType {
	case "noop":
		constType = Noop
	case "addx":
		constType = Addx
		increaseValue, _ = strconv.Atoi(lineSplit[1])
	}

	return Instruction{
		Type:          constType,
		IncreaseValue: increaseValue,
	}
}

func parseFileToProgram(filePath string) (*Program, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	instructions := make([]Instruction, len(lines))
	nLines := len(lines)

	for i := 0; i < nLines; i++ {
		instructions[i] = lineToInstruction(lines[i])
	}

	program := Program{
		Instructions: instructions,
	}

	return &program, nil
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
