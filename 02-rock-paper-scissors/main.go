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

	log.Printf("> (1st Puzzle) What would your total score be if everything goes exactly according to your strategy guide?")

	games, err := parseFile(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	finalPlayerScore := games.ComputePlayerScore()

	log.Printf("Final Player Score is: %d", finalPlayerScore)
}

const (
	OpponentPlayerMoveDistance int8 = 23
	Rock                            = 65
	WinPoints                       = 6
	DrawPoints                      = 3
	OpponentMoveLineIndex           = 0
	PlayerMoveLineIndex             = 2
)

type Game struct {
	PlayerMove   int8
	OpponentMove int8
	Score        []int
}

type Games []Game

func lineToGame(l string) Game {
	opponentMove := int8(l[OpponentMoveLineIndex])
	playerMove := int8(l[PlayerMoveLineIndex]) - OpponentPlayerMoveDistance
	gameScore := computeScore(playerMove, opponentMove)

	return Game{
		PlayerMove:   playerMove,
		OpponentMove: opponentMove,
		Score:        gameScore,
	}
}

func computeScore(pm int8, om int8) []int {
	score := []int{int(pm - (Rock - 1)), int(om - (Rock - 1))}

	movesDiff := pm - om

	// scissors - rock = 2
	// paper/scissors - rock/paper = 1
	// rock/paper/scissors - rock/paper/scissors = 0
	// rock/paper - paper/scissors = -1
	// rock - scissors = -2

	switch movesDiff {
	case 2:
		score[1] += WinPoints
		break
	case 1:
		score[0] += WinPoints
		break
	case 0:
		score[0] += DrawPoints
		score[1] += DrawPoints
		break
	case -1:
		score[1] += WinPoints
		break
	case -2:
		score[0] += WinPoints
		break
	}

	return score
}

func (gs Games) ComputePlayerScore() uint {
	sum := uint(0)

	for _, g := range gs {
		sum += uint(g.Score[0])
	}

	return sum
}

func parseFile(filePath string) (Games, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	result := make([]Game, len(lines))

	for i, line := range lines {
		result[i] = lineToGame(line)
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
