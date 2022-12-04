// go run main.go
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputr.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	log.Printf("> (1st Puzzle) In how many assignment pairs does one range fully contain the other?")

	elvesPair, err := parseFile(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	sumOfPriorities := elvesPair.ComputeNumberOfFullyOverlapingSections()

	log.Printf("The number of assignment pairs that one range fully contain the other is: %d", sumOfPriorities)
}

type ElfPair struct {
	FirstElfSections        []int
	SecondElfSections       []int
	FullyOverlapingSections []int
}

type ElvesPair []ElfPair

func (this ElvesPair) ComputeNumberOfFullyOverlapingSections() uint {
	sum := uint(0)

	for _, pair := range this {
		if len(pair.FullyOverlapingSections) > 0 {
			sum++
		}
	}

	return sum

}

func lineToElfPair(l string) ElfPair {
	elves := strings.Split(l, ",")
	firstElfSections := sectionRangeToArray(elves[0])
	secondElfSections := sectionRangeToArray(elves[1])
	fullyOverlapingSections := computeFullyOverlapingSections(firstElfSections, secondElfSections)

	return ElfPair{
		FirstElfSections:        firstElfSections,
		SecondElfSections:       secondElfSections,
		FullyOverlapingSections: fullyOverlapingSections,
	}
}

func computeFullyOverlapingSections(firstElfSections []int, secondElfSections []int) []int {
	fs := firstElfSections[0]
	fe := firstElfSections[len(firstElfSections)-1]
	ss := secondElfSections[0]
	se := secondElfSections[len(secondElfSections)-1]

	if fs <= ss && fe >= se {
		return secondElfSections
	}
	if fs >= ss && fe >= se {
		return firstElfSections
	}
	return []int{}
}

func sectionRangeToArray(sectionRange string) []int {
	values := strings.Split(sectionRange, "-")
	start, _ := strconv.Atoi(strings.TrimSpace(values[0]))
	end, _ := strconv.Atoi(strings.TrimSpace(values[1]))
	result := make([]int, end-start+1)
	for i := range result {
		result[i] = i + start
	}
	return result
}

func parseFile(filePath string) (ElvesPair, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	result := make([]ElfPair, len(lines))

	for i, line := range lines {
		result[i] = lineToElfPair(line)
	}

	return result, nil
}

func getFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	rawBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error while reading file: %v", err)
		return nil, err
	}
	return strings.Split(string(rawBytes), "\n"), nil
}
