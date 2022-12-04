// go run main.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	log.Printf("> (1st Puzzle) In how many assignment pairs does one range fully contain the other?")

	elvesPair, err := parseFile(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	fullyOverlappingSectionsCount := elvesPair.ComputeNumberOfFullyOverlappingSections()

	log.Printf("The number of assignment pairs that one range fully contain the other is: %d", fullyOverlappingSectionsCount)

	log.Printf("> (2nd Puzzle) In how many assignment pairs does one range fully contain the other?")

	overlappingSectionsCount := elvesPair.ComputeNumberOfOverlappingSections()

	log.Printf("The number of assignment pairs that overlap the other is: %d", overlappingSectionsCount)
}

type ElfPair struct {
	FirstElfSections  []int
	SecondElfSections []int
	Overlaps          bool
	FullyOverlaps     bool
}

type ElvesPair []ElfPair

func (this ElvesPair) ComputeNumberOfFullyOverlappingSections() uint {
	sum := uint(0)

	for _, pair := range this {
		if pair.FullyOverlaps {
			sum++
		}
	}

	return sum

}

func (this ElvesPair) ComputeNumberOfOverlappingSections() uint {
	sum := uint(0)

	for _, pair := range this {
		if pair.Overlaps {
			sum++
		}
	}

	return sum

}

func lineToElfPair(l string) ElfPair {
	elves := strings.Split(l, ",")
	firstElfSections := sectionRangeToArray(elves[0])
	secondElfSections := sectionRangeToArray(elves[1])
	overlapingSections := computeOverlapingSections(firstElfSections, secondElfSections)
	fullyOverlapingSections := computeFullyOverlapingSections(firstElfSections, secondElfSections)

	return ElfPair{
		FirstElfSections:  firstElfSections,
		SecondElfSections: secondElfSections,
		Overlaps:          len(overlapingSections) > 0,
		FullyOverlaps:     len(fullyOverlapingSections) > 0,
	}
}

func computeFullyOverlapingSections(firstElfSections []int, secondElfSections []int) []int {
	fs := firstElfSections[0]
	fe := firstElfSections[len(firstElfSections)-1]
	ss := secondElfSections[0]
	se := secondElfSections[len(secondElfSections)-1]

	var sections []int

	if fs <= ss && fe >= se {
		sections = secondElfSections
	} else if fs >= ss && fe <= se {
		sections = firstElfSections
	} else {
		sections = []int{}
	}

	return sections
}

func computeOverlapingSections(firstElfSections []int, secondElfSections []int) []int {
	fs := firstElfSections[0]
	fe := firstElfSections[len(firstElfSections)-1]
	ss := secondElfSections[0]
	se := secondElfSections[len(secondElfSections)-1]

	var sections []int

	if fs <= ss && fe >= se {
		//fully overlapping
		sections = secondElfSections
	} else if fs >= ss && fe <= se {
		//fully overlapping
		sections = firstElfSections
	} else if fs < ss && fe >= ss && fe < se {
		//overlapping
		sections = sectionRangeToArray(fmt.Sprintf("%v-%v", ss, fe))
	} else if fs > ss && fs <= se && fe > se {
		//overlapping
		sections = sectionRangeToArray(fmt.Sprintf("%v-%v", fs, se))
	} else {
		sections = []int{}
	}

	return sections
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
