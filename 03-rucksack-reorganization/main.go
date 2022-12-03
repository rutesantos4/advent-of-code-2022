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

	log.Printf("> (1st Puzzle) Find the item type that appears in both compartments of each rucksack. What is the sum of the priorities of those item types?")

	rucksacks, err := parseFile(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	sumOfPriorities := rucksacks.ComputeSumOfFirstSharedItemTypePriorityValues()

	log.Printf("Sum of first shared item type priorities is: %d", sumOfPriorities)

	log.Printf("> (2nd Puzzle) Find the item type that corresponds to the badges of each three-Elf group. What is the sum of the priorities of those item types?")

	sumOfBadgesPriorities := rucksacks.ComputeSumOfGroupsBadgesPriorityValues()

	log.Printf("Sum of badges of each three-Elf group priorities is: %d", sumOfBadgesPriorities)
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

func (rs Rucksacks) ComputeSumOfGroupsBadgesPriorityValues() uint {
	sum := uint(0)

	for i := 0; i < len(rs)-2; i = i + 3 {
		sum += uint(ComputeSumOfBadgesPriorityValues(rs[i], rs[i+1], rs[i+2]))
	}

	return sum
}

func ComputeSumOfBadgesPriorityValues(r1 Rucksack, r2 Rucksack, r3 Rucksack) uint8 {

	r1Items := append(r1.FirstCompartementItems, r1.SecondCompartementItems...)
	r2Items := append(r2.FirstCompartementItems, r2.SecondCompartementItems...)
	r3Items := append(r3.FirstCompartementItems, r3.SecondCompartementItems...)
	for _, r1it := range r1Items {
		for _, r2it := range r2Items {
			for _, r3it := range r3Items {
				if r1it == r2it && r2it == r3it {
					return itemTypeToPriorityValue(r1it)
				}
			}
		}
	}

	return uint8(0)
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
