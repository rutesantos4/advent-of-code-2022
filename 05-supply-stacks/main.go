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

	log.Printf("> (1st Puzzle) After the rearrangement procedure completes, what crate ends up on top of each stack (Mover 9000)?")

	rearrangement, err := parseFile(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	rearrangementWithMover9000 := rearrangement.Copy()

	rearrangementWithMover9000.ProcessRearrangementWithCrateMover9000()
	topCratesOfStacksWithMover9000 := rearrangementWithMover9000.GetTopCratesStacks()

	log.Printf("The crates that end up on top of each stack are: %v", topCratesOfStacksWithMover9000)

	log.Printf("> (2nd Puzzle) After the rearrangement procedure completes, what crate ends up on top of each stack (Mover 9001)?")

	rearrangementWithMover9001 := rearrangement.Copy()

	rearrangementWithMover9001.ProcessRearrangementWithCrateMover9001()
	topCratesOfStacksWithMover9001 := rearrangementWithMover9001.GetTopCratesStacks()

	log.Printf("The crates that end up on top of each stack are: %v", topCratesOfStacksWithMover9001)
}

const (
	MovesRegularExpression string = `move (.*?) from (.*?) to (.?)`
)

type Rearrangement struct {
	CratesStack  map[int]Crates //int - stack number | Crates - list if crates
	Moves        []Move
	IsRearranged bool
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
	cratesRange := len(r.CratesStack) + 1

	for i := 1; i < cratesRange; i++ {
		crate := r.CratesStack[i][0]
		result += fmt.Sprintf("%c", crate)
	}

	return result
}

func (r Rearrangement) ProcessRearrangementWithCrateMover9000() {
	for _, move := range r.Moves {
		for i := 0; i < move.Count; i++ {
			crates := r.CratesStack[move.From]
			value := crates[0]
			r.CratesStack[move.From] = crates[1:]
			r.CratesStack[move.To] = append([]byte{value}, r.CratesStack[move.To]...)
		}
	}

	r.IsRearranged = true
}

func (r Rearrangement) ProcessRearrangementWithCrateMover9001() {
	for _, move := range r.Moves {
		orgCrates := r.CratesStack[move.From]
		dstCrates := r.CratesStack[move.To]

		// In Go, append can either me mutable or immutable
		// A slice of a slice ([:...]) is a view of the sliced slice
		// Hence we need to create a copy of the moving crates, else we will
		// end up swapping both crates stacks
		// Example:
		// a := []int {1, 2, 3}
		// b := a[2:]
		// b[0] = 5
		// Slice a is now = [1, 2 ,5]

		cratesToMove := orgCrates[:move.Count]
		cratesToMoveCopy := make(Crates, len(cratesToMove))
		copy(cratesToMoveCopy, cratesToMove)

		dstCrates = append(cratesToMoveCopy, dstCrates...)
		orgCrates = orgCrates[move.Count:]

		r.CratesStack[move.From] = orgCrates
		r.CratesStack[move.To] = dstCrates

	}

	r.IsRearranged = true
}

func (r Rearrangement) Copy() Rearrangement {
	mc := make([]Move, len(r.Moves))
	copy(mc, r.Moves)

	csc := make(map[int]Crates, len(r.CratesStack))

	for k, v := range r.CratesStack {
		vc := make(Crates, len(v))
		copy(vc, v)

		csc[k] = vc
	}

	return Rearrangement{
		Moves:        mc,
		CratesStack:  csc,
		IsRearranged: r.IsRearranged,
	}
}

func linesToCratesStack(lines []string) map[int]Crates {
	indexOfLineOfStacksIds := len(lines) - 1
	result := make(map[int]Crates)
	lineStacksNumbers := lines[indexOfLineOfStacksIds]
	stacksNumbers := strings.Split(lineStacksNumbers, " ")

	for _, stackNumberString := range stacksNumbers {

		if strings.TrimSpace(stackNumberString) == "" {
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
	re := regexp.MustCompile(MovesRegularExpression)
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
	lineNumber := 0
	lineCount := len(lines)

	for lineNumber = 0; lineNumber < lineCount; lineNumber++ {
		line := lines[lineNumber]

		if strings.TrimSpace(line) == "" {
			lineNumber++
			break
		}

		cratesStackLines = append(cratesStackLines, line)
	}

	cratesStack := linesToCratesStack(cratesStackLines)

	numberMoves := lineCount - lineNumber
	moves := make([]Move, lineCount-lineNumber)

	for i := 0; i < numberMoves; i++ {
		moves[i] = lineToMove(lines[i+lineNumber])
	}

	return Rearrangement{
		CratesStack:  cratesStack,
		Moves:        moves,
		IsRearranged: false,
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
