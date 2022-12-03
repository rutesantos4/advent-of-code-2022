// go run main.go
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	log.Printf("> Find the item type that appears in both compartments of each rucksack. What is the sum of the priorities of those item types?")

	rucksacks, err := parseFile(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	sumOfPriorities := rucksacks.ComputeSumOfFirstSharedItemTypePriorityValues()

	log.Printf("Sum of first shared item type priorities is: %d", sumOfPriorities)

}

const (
	CapitalAlphabetStartByteValue                 int = 65
	LowerAlphabetStartByteValue                   int = 97
	SpaceBetweenCapitalAndLowerAsciiAlphabetStart int = 32
	SpaceBetweenAsciiAlphabetEndToStart           int = 5
)

type Rucksack struct {
	FirstCompartementItems  []byte
	SecondCompartementItems []byte
	SharedItems             []byte
}

type Rucksacks []Rucksack

func (rs Rucksacks) ComputeSumOfFirstSharedItemTypePriorityValues() uint {
	sum := uint(0)

	for _, v := range rs {
		sum += uint(v.FirstSharedItemTypePriority())
	}

	return sum
}

func (r Rucksack) FirstSharedItemTypePriority() uint8 {
	return itemTypeToPriorityValue(r.SharedItems[0])
}

func lineToRucksack(l string) Rucksack {
	items := []byte(l)
	itemsCount := len(items)
	firstCompartementItems := items[0 : itemsCount/2]
	secondCompartementItems := items[itemsCount/2 : itemsCount]
	sharedItems := computeSharedItems(firstCompartementItems, secondCompartementItems)

	return Rucksack{
		FirstCompartementItems:  firstCompartementItems,
		SecondCompartementItems: secondCompartementItems,
		SharedItems:             sharedItems,
	}
}

func computeSharedItems(fci []byte, sci []byte) []byte {
	var sharedItemTypes []byte

	for _, fit := range fci {
		for _, sit := range sci {
			if fit == sit {
				sharedItemTypes = append(sharedItemTypes, fit)
			}
		}
	}

	return sharedItemTypes
}

func itemTypeToPriorityValue(it byte) uint8 {
	// byte(a) = 97
	// byte(A) = 65
	// space between ascii alphabet = 32
	// priority(a) = byte(a) - SpaceBetweenCapitalAndLowerAsciiAlphabetStart - CapitalAlphabetStartByteValue + 1 = 97 - 32 - 65 + 1 = 1
	// priority(b) = byte(b) - SpaceBetweenCapitalAndLowerAsciiAlphabetStart - CapitalAlphabetStartByteValue + 1 = 98 - 32 - 65 + 1 = 2
	// ----------
	// priority(A) = byte(A) - SpaceBetweenCapitalAndLowerAsciiAlphabetStart - SpaceBetweenAsciiAlphabetEndToStart - 1 = 65 - 32 - 5 - 1 = 27
	// priority(B) = byte(B) - SpaceBetweenCapitalAndLowerAsciiAlphabetStart - SpaceBetweenAsciiAlphabetEndToStart - 1 = 66 - 32 - 5 - 1 = 28
	bv := int(it)

	var pv int

	if bv < LowerAlphabetStartByteValue {
		pv = bv - SpaceBetweenCapitalAndLowerAsciiAlphabetStart - SpaceBetweenAsciiAlphabetEndToStart - 1
	} else {
		pv = bv - SpaceBetweenCapitalAndLowerAsciiAlphabetStart - CapitalAlphabetStartByteValue + 1
	}

	return uint8(pv)
}

func parseFile(filePath string) (Rucksacks, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	result := make([]Rucksack, len(lines))

	for i, line := range lines {
		result[i] = lineToRucksack(line)
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
