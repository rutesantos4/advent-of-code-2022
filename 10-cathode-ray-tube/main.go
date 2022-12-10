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

	program, err := parseFileToProgram(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Printf("> (1st Puzzle) Find the signal strength during the 20th, 60th, 100th, 140th, 180th, and 220th cycles. What is the sum of these six signal strengths?")

	cycles := []int{20, 60, 100, 140, 180, 220}
	sumOfCycleSignalStrengths := program.ComputeCyclesSignalStrengthSum(cycles)

	log.Printf("The sum of these six signal strengths is : %d", sumOfCycleSignalStrengths)

}

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

func (p Program) ComputeCyclesSignalStrengthSum(cycles []int) int {
	xValues := p.computeXValue()

	sum := 0
	for _, cycle := range cycles {
		sum += xValues[cycle-1].ComputeCycleSignalStrength()
	}

	return sum
}

func (p Program) computeXValue() []XValue {

	var xValues []XValue

	cycle := 0
	xDuring := 1
	xAfter := 1

	for _, instruction := range p.Instructions {
		instructionDuration := instruction.Duration()

		for i := 0; i < instructionDuration; i++ {
			cycle++

			xValues = append(xValues, XValue{
				Cycle:  cycle,
				During: xDuring,
				After:  xAfter,
			})

			if i == instructionDuration-2 { //Increment before the last one
				xAfter += instruction.IncreaseValue
			}
		}

		xDuring = xValues[cycle-1].After
		xAfter = xDuring
	}

	return xValues
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
