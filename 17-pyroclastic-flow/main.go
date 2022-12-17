package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	moves, err := parseFileToMoves(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Println(moves)

	log.Printf("> (1st Puzzle) How many units tall will the tower of rocks be after 2022 rocks have stopped falling?")
	println(Tetris(moves, 20).String())
	// numberSandBeforeAbyss, _ := moves.DrawSand()

	// log.Printf("The units tall tower of rocks will be after 2022 rocks is: %d\n", numberSandBeforeAbyss)
}

const (
	Left  = "<"
	Right = ">"
	Down  = "v"
)

const (
	ChamberWidth         = 7
	RockEdgeAwayLeftWall = 2
	RockEdgeAwayBottom   = 3
)

type Element string
type Chamber [][]Element
type Position struct {
	X int
	Y int
}
type Rock [][]Element

var AllRocks []Rock = []Rock{Rock1(), Rock2(), Rock3(), Rock4(), Rock5()}

type Move string
type Moves []Move

func Tetris(m Moves, plays int) Chamber {
	chamber := BuildChamper()
	numberMoves := len(m)
	numberRocks := len(AllRocks)
	currentMove := 0
	currentRock := 0
	height := 0

	for plays > 0 {

		if currentRock == numberRocks {
			currentRock = 0
		}

		rock := AllRocks[currentRock]
		currentPosition := rock.getFistPosition(height + RockEdgeAwayBottom)
		goNextRock := false

		for !goNextRock {

			if currentMove == numberMoves {
				currentMove = 0
			}

			move := m[currentMove]

			currentPosition = currentPosition.MovePosition(move, rock, height)

			// Not working properly, can still go down, but it stops
			if chamber.isStopFalling(rock, currentPosition, height) {
				// if rock stopped then substract play
				plays--
				currentMove++
				move = m[currentMove]
				currentPosition = currentPosition.MovePosition(move, rock, height)
				chamber.drawRock(rock, currentPosition)
				currentRock++
				height += len(rock)
				goNextRock = true
			}
			currentMove++

		}

	}

	return chamber
}

func (p Position) MovePosition(move Move, rock Rock, height int) Position {

	width := rock.getWidth()
	switch move {
	case Left:
		if x := p.X - 1; x >= height {
			p.X = x
		}
	case Right:
		if x := p.X + 1; x+width < ChamberWidth {
			p.X = x
		}
	case Down:
		p.Y = p.Y - 1
	}

	return p
}

func (r Rock) getWidth() int {
	l := len(r)
	max := 0

	for i := 0; i < l; i++ {
		col := len(r[i])

		for j := 0; j < col; j++ {

			if r[i][j] == "#" && j > max {
				max = j
			}
		}
	}

	return max
}

func (r Rock) getFistPosition(height int) Position {

	l := len(r)

	p := Position{
		X: RockEdgeAwayLeftWall,
		Y: height,
	}

	for i := 0; i < l; i++ {
		col := len(r[i])

		for j := 0; j < col; j++ {

			if r[i][j] == "#" {
				p.X += j
				return p
			}
		}
	}

	return p
}

func (c Chamber) isStopFalling(rock Rock, position Position, height int) bool {
	l := len(rock)
	col := len(rock[0])

	if position.Y <= height {
		return true
	}

	for i := 0; i < l; i++ {
		for j := 0; j < col; j++ {

			r := rock[i][j]
			m := c[i+position.Y][j+position.X]

			if m == Element("#") && r == Element("#") {
				return true
			}
		}
	}

	return false
}

func (c Chamber) drawRock(rock Rock, position Position) {
	l := len(rock)
	col := len(rock[0])

	for i := 0; i < l; i++ {
		for j := 0; j < col; j++ {

			r := rock[i][j]
			c[i+position.Y][j+position.X] = r
		}
	}
}

func (c Chamber) String() string {
	lines := len(c)
	columns := len(c[0])
	result := "\n"
	for l := 0; l < lines; l++ {
		for col := 0; col < columns; col++ {
			result += string(c[l][col])
		}
		result += "\n"
	}

	return result
}

func BuildChamper() Chamber {
	size := 60
	chamber := make(Chamber, size)

	for i := 0; i < size; i++ {
		chamber[i] = make([]Element, ChamberWidth)
		for j := 0; j < ChamberWidth; j++ {
			chamber[i][j] = "."
		}
	}

	return chamber
}

func Rock1() Rock {
	rock := make(Rock, 1)
	rock[0] = make([]Element, 4)

	for i := 0; i < 4; i++ {
		rock[0][i] = "#"
	}

	return rock
}

func Rock2() Rock {
	rock := make(Rock, 3)

	for i := 0; i < 3; i++ {
		rock[i] = make([]Element, 3)
	}

	rock[0][0] = "."
	rock[0][1] = "#"
	rock[0][2] = "."

	rock[1][0] = "#"
	rock[1][1] = "#"
	rock[1][2] = "#"

	rock[2][0] = "."
	rock[2][1] = "#"
	rock[2][2] = "."

	return rock
}

func Rock3() Rock {
	rock := make(Rock, 3)

	for i := 0; i < 3; i++ {
		rock[i] = make([]Element, 3)
	}

	rock[0][0] = "#"
	rock[0][1] = "#"
	rock[0][2] = "#"

	rock[1][0] = "."
	rock[1][1] = "."
	rock[1][2] = "#"

	rock[2][0] = "."
	rock[2][1] = "."
	rock[2][2] = "#"

	return rock
}

func Rock4() Rock {
	rock := make(Rock, 4)

	for i := 0; i < 4; i++ {
		rock[i] = make([]Element, 1)
		rock[i][0] = "#"
	}

	return rock
}

func Rock5() Rock {
	rock := make(Rock, 2)

	for i := 0; i < 2; i++ {
		rock[i] = make([]Element, 2)
		rock[i][0] = "#"
		rock[i][1] = "#"
	}

	return rock
}

func parseFileToMoves(filePath string) (Moves, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	// It should have only one line
	line := lines[0]
	nMoves := len(line) * 2

	var m = make(Moves, nMoves)
	i := 0

	for _, v := range line {
		m[i] = Move(string(v))
		i++
		m[i] = Move(Down)
		i++
	}

	return m, nil
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
