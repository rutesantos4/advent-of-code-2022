// go run main.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	log.Printf("> (1st Puzzle) After the rearrangement procedure completes, what crate ends up on top of each stack?")

	rearrangement, err := parseFile(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	rearrangement.ProcessRearrangement()
	topCratesOfStacks := rearrangement.GetTopCratesStacks()

	log.Printf("The crates that end up on top of each stack are: %v", topCratesOfStacks)
}

const (
	RegularExpressinOfMoves string = `move (.*?) from (.*?) to (.?)`
)

type Rearrangement struct {
	CratesStack map[int]Crates //int - stack number | Crates - list if crates
	Moves       []Move
}
type Crates []byte //byte - crate id

type Move struct {
	// move {{Count}} from {{From}} to {{To}}
	Count int
	From  int
	To    int
}

func (r Rearrangement) GetTopCratesStacks() string {
	var result string
	for i := 1; i < len(r.CratesStack)+1; i++ {
		crate := r.CratesStack[i][0]
		result += fmt.Sprintf("%c", crate)
	}
	return result
}

func (r Rearrangement) ProcessRearrangement() {
	for _, move := range r.Moves {
		for i := 0; i < move.Count; i++ {
			crates := r.CratesStack[move.From]
			value := crates[0]
			r.CratesStack[move.From] = crates[1:]
			r.CratesStack[move.To] = append([]byte{value}, r.CratesStack[move.To]...)
		}
	}
}

func linesToCratesStack(lines []string) map[int]Crates {
	indexOfLineOfStacksIds := len(lines) - 1
	result := make(map[int]Crates)
	lineStacksNumbers := lines[indexOfLineOfStacksIds]
	stacksNumbers := strings.Split(lineStacksNumbers, " ")
	for _, stackNumberString := range stacksNumbers {
		if stackNumberString == "" || stackNumberString == "\r" {
			continue
		}
		index := strings.Index(lineStacksNumbers, stackNumberString)
		stackNumber, _ := strconv.Atoi(stackNumberString)
		var crates []byte
		for i := 0; i < indexOfLineOfStacksIds; i++ {
			stacksIds := lines[i][index]
			if string(stacksIds) != " " {
				crates = append(crates, stacksIds)
			}
		}
		result[stackNumber] = crates
	}
	return result
}

func lineToMove(l string) Move {
	re := regexp.MustCompile(RegularExpressinOfMoves)
	match := re.FindStringSubmatch(l)

	count, _ := strconv.Atoi(match[1])
	from, _ := strconv.Atoi(match[2])
	to, _ := strconv.Atoi(match[3])

	return Move{
		Count: count,
		From:  from,
		To:    to,
	}
}

func parseFile(filePath string) (Rearrangement, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return Rearrangement{}, err
	}

	var cratesStackLines []string
	lineNumer := 0
	for lineNumer = 0; lineNumer < len(lines); lineNumer++ {
		line := lines[lineNumer]
		if line == "\r" {
			lineNumer++
			break
		}
		cratesStackLines = append(cratesStackLines, line)
	}

	cratesStack := linesToCratesStack(cratesStackLines)

	numberMoves := len(lines) - lineNumer
	moves := make([]Move, len(lines)-lineNumer)
	for i := 0; i < numberMoves; i++ {
		moves[i] = lineToMove(lines[i+lineNumer])
	}

	return Rearrangement{
		CratesStack: cratesStack,
		Moves:       moves,
	}, nil
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
