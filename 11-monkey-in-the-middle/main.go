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

	monkeys, err := parseFileToMonkeys(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Printf("> (1st Puzzle) What is the level of monkey business after 20 rounds of stuff-slinging simian shenanigans?")

	monkeysAfterPuzzleRounds := monkeys.PlayMonkeyInTheMiddleFor(PuzzleGameRounds, reduceWorryLevelDivisionBy3())
	monkeyBusinessLevel := monkeysAfterPuzzleRounds.ComputeMonkeyBusinessLevel()

	log.Printf("The monkey level business after %d rounds is: %d\n", PuzzleGameRounds, monkeyBusinessLevel)

	log.Printf("> (2nd Puzzle) Starting again from the initial state in your puzzle input, what is the level of monkey business after 10000 rounds?")

	monkeysAfterPuzzle2Rounds := monkeys.PlayMonkeyInTheMiddleFor(Puzzle2GameRounds, reduceWorryLevelPuzzleModularArithmetic(monkeys))
	monkeyBusinessLevelPuzzle2 := monkeysAfterPuzzle2Rounds.ComputeMonkeyBusinessLevel()

	log.Printf("The monkey level business after %d rounds is: %d\n", Puzzle2GameRounds, monkeyBusinessLevelPuzzle2)

}

const (
	PuzzleGameRounds  = 20
	Puzzle2GameRounds = 10000
)

const (
	AddOp      = "+"
	SubtractOp = "-"
	MultiplyOp = "*"
	DivideOp   = "/"
)

type Monkey struct {
	Id                        MonkeyId
	ItemsWorryLevel           Items
	WorryLevelUpdateOperation WorryLevelUpdateOperationCallback
	MonkeyPassTestOperation   MonkeyPassTestOperationCallback
	DivisionNumber            int
	PassedItemsCount          int
}

type MonkeyId byte
type ItemWorryLevel uint
type WorryLevelUpdateOperationCallback func(owl ItemWorryLevel) ItemWorryLevel
type MonkeyPassTestOperationCallback func(wl ItemWorryLevel, rwl WorryLevelReductionCallback) (MonkeyId, ItemWorryLevel)

type Items []ItemWorryLevel
type Monkeys []Monkey

type WorryLevelReductionCallback func(owl ItemWorryLevel) ItemWorryLevel

func reduceWorryLevelDivisionBy3() WorryLevelReductionCallback {
	return func(owl ItemWorryLevel) ItemWorryLevel {
		return owl / 3
	}
}

// Reduce the worry by applying Modular arithmetic
func reduceWorryLevelPuzzleModularArithmetic(m Monkeys) WorryLevelReductionCallback {
	prod := 1
	for _, v := range m {
		prod *= v.DivisionNumber
	}

	return func(owl ItemWorryLevel) ItemWorryLevel {
		return owl % ItemWorryLevel(prod)
	}
}

func (ms Monkeys) ComputeMonkeyBusinessLevel() uint {
	var firstMaxPassedItemsCount, secondMaxPassedItemsCount int
	firstMaxPassedItemsCount = 0
	secondMaxPassedItemsCount = 0

	for _, m := range ms {
		if m.PassedItemsCount > firstMaxPassedItemsCount {
			secondMaxPassedItemsCount = firstMaxPassedItemsCount
			firstMaxPassedItemsCount = m.PassedItemsCount
		} else if m.PassedItemsCount > secondMaxPassedItemsCount {
			secondMaxPassedItemsCount = m.PassedItemsCount
		}
	}

	return uint(firstMaxPassedItemsCount) * uint(secondMaxPassedItemsCount)
}

func (ms Monkeys) PlayMonkeyInTheMiddleFor(rounds int, rwl WorryLevelReductionCallback) *Monkeys {
	nms := make(Monkeys, len(ms))

	nmsMap := map[MonkeyId]*Monkey{}

	for i, m := range ms {
		isc := make(Items, len(m.ItemsWorryLevel))
		copy(isc, m.ItemsWorryLevel)

		mc := &Monkey{
			Id:                        m.Id,
			WorryLevelUpdateOperation: m.WorryLevelUpdateOperation,
			ItemsWorryLevel:           isc,
			MonkeyPassTestOperation:   m.MonkeyPassTestOperation,
			PassedItemsCount:          0,
		}

		nmsMap[m.Id] = mc
		nms[i] = *mc
	}

	for i := 0; i < rounds; i++ {
		for _, mh := range nms {
			m := nmsMap[mh.Id]
			for _, iwl := range m.ItemsWorryLevel {
				niwl := iwl
				niwlu := m.WorryLevelUpdateOperation(niwl)
				pmi, pmiwl := m.MonkeyPassTestOperation(niwlu, rwl)
				pm := nmsMap[pmi]

				pm.ItemsWorryLevel = append(pm.ItemsWorryLevel, pmiwl)
				m.PassedItemsCount++
			}

			m.ItemsWorryLevel = Items{}
		}

	}

	for i, m := range nms {
		nms[i] = *nmsMap[m.Id]
	}

	return &nms
}

func parseWorryLevelUpdateOperation(s string) WorryLevelUpdateOperationCallback {
	sSplit := strings.Split(s, " ")
	sSplitLen := len(sSplit)

	rightOperandParse, err := strconv.Atoi(sSplit[sSplitLen-1])
	rightOperand := ItemWorryLevel(rightOperandParse)
	operation := sSplit[sSplitLen-2]

	return func(owl ItemWorryLevel) ItemWorryLevel {
		var nwl ItemWorryLevel

		if err != nil {
			rightOperand = owl
		}

		switch operation {
		case AddOp:
			nwl = owl + rightOperand
		case SubtractOp:
			nwl = owl - rightOperand
		case MultiplyOp:
			nwl = owl * rightOperand
		case DivideOp:
			nwl = owl / rightOperand
		}

		return nwl
	}
}

func parseMonkeyPassTestOperation(ts, itt, ift string) (MonkeyPassTestOperationCallback, int) {
	tsSplit := strings.Split(ts, " ")
	tSplitLen := len(tsSplit)

	rightOperandParse, _ := strconv.Atoi(tsSplit[tSplitLen-1])
	rightOperand := ItemWorryLevel(rightOperandParse)

	return func(wl ItemWorryLevel, rwl WorryLevelReductionCallback) (MonkeyId, ItemWorryLevel) {
		div := rwl(wl)
		res := div % rightOperand

		if res == 0 {
			return MonkeyId(itt[len(itt)-1]), div
		} else {
			return MonkeyId(ift[len(ift)-1]), div
		}
	}, rightOperandParse
}

func linesToMonkey(ls []string) *Monkey {
	monkeyInfoMap := map[string]string{}

	for _, l := range ls {
		lSplit := strings.Split(l, ":")
		monkeyInfoMap[strings.TrimSpace(lSplit[0])] = strings.TrimSpace(lSplit[1])
	}

	monkeyId := MonkeyId(ls[0][len(ls[0])-2])

	itemsSplit := strings.Split(monkeyInfoMap["Starting items"], ", ")
	items := make(Items, len(itemsSplit))

	for i, is := range itemsSplit {
		isParse, _ := strconv.Atoi(is)

		items[i] = ItemWorryLevel(isParse)
	}

	worryLevelUpdateOperation := parseWorryLevelUpdateOperation(monkeyInfoMap["Operation"])

	monkeyPassTestOperation, divisionNumber := parseMonkeyPassTestOperation(monkeyInfoMap["Test"], monkeyInfoMap["If true"], monkeyInfoMap["If false"])

	return &Monkey{
		Id:                        monkeyId,
		ItemsWorryLevel:           items,
		WorryLevelUpdateOperation: worryLevelUpdateOperation,
		MonkeyPassTestOperation:   monkeyPassTestOperation,
		DivisionNumber:            divisionNumber,
		PassedItemsCount:          0,
	}
}

func parseFileToMonkeys(filePath string) (Monkeys, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	var cleanLines []string

	for _, l := range lines {
		cl := strings.TrimSpace(l)
		if cl != "" {
			cleanLines = append(cleanLines, cl)
		}
	}

	monkeyLinesGroupLen := 6
	monkeysLen := len(cleanLines) / monkeyLinesGroupLen

	monkeys := make(Monkeys, monkeysLen)

	for i := 0; i < monkeysLen; i++ {
		si := i * monkeyLinesGroupLen
		ei := si + monkeyLinesGroupLen

		monkeys[i] = *linesToMonkey(cleanLines[si:ei])
	}

	return monkeys, nil
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
