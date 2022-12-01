// go run main.go
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	log.Printf("> (1st Puzzle) Find the Elf carrying the most Calories. How many total Calories is that Elf carrying?")

	start := time.Now()

	log.Println("Start First Solution")
	elfsMap, err := parseFile(*inputFilePath)
	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	elfs := Elfs{
		List: elfsMap,
	}

	elfMostCalories := elfs.FindElfMostCalories()

	log.Printf("The elf with most calories is %v, with %v calories", elfMostCalories.Id, elfMostCalories.GetTotalCalories())

	elapsed := time.Since(start)
	log.Printf("End First Solution - %s", elapsed)

	log.Println("Start Second Solution")
	start = time.Now()
	part1, err := getMostCalories(*inputFilePath)
	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}
	log.Printf("The most calories is %v", part1)
	elapsed = time.Since(start)
	log.Printf("End Second Solution - %s", elapsed)

	log.Printf("> (2nd Puzzle) Find the top three Elves carrying the most Calories. How many Calories are those Elves carrying in total?")

	elvesThatCarryMostCaloriesCount := 3

	elves := elfs.FindElvesThatCarryMostCalories(elvesThatCarryMostCaloriesCount)

	elvesCarriedCaloriesTotal := 0

	for _, v := range elves {
		totalCalores := v.GetTotalCalories()
		elvesCarriedCaloriesTotal += totalCalores

		log.Printf("Elf (%d): %d", v.Id, totalCalores)
	}

	log.Printf("Total calories carried by elves: %d", elvesCarriedCaloriesTotal)
}

func parseFile(filePath string) (map[int]Elf, error) {
	lines, err := getFileLines(filePath)
	if err != nil {
		return nil, err
	}

	var result = make(map[int]Elf)
	elfNumber := 0
	result[elfNumber] = Elf{Id: elfNumber}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			elfNumber = elfNumber + 1
			result[elfNumber] = Elf{Id: elfNumber}
			continue
		}
		calorie, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalf("Error while parsing line to int: %v", err)
			return nil, err
		}

		result[elfNumber] = result[elfNumber].AddCalorie(calorie)
	}

	return result, nil
}

func getMostCalories(filePath string) (int, error) {
	lines, err := getFileLines(filePath)
	if err != nil {
		return -1, err
	}

	var result = 0
	sum := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			sum = 0
			continue
		}
		calorie, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalf("Error while parsing line to int: %v", err)
			return -1, err
		}
		sum = sum + calorie
		if sum > result {
			result = sum
		}
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

type Elfs struct {
	List map[int]Elf
}

func (this Elfs) FindElfMostCalories() Elf {
	most := this.List[0]
	for i := 1; i < len(this.List); i++ {
		if most.GetTotalCalories() < this.List[i].GetTotalCalories() {
			most = this.List[i]
		}

	}
	return most
}

func (this Elfs) FindElvesThatCarryMostCalories(elvesCount int) []Elf {
	elves := make([]Elf, len(this.List))

	for i, elf := range this.List {
		elves[i] = elf
	}

	sort.Slice(elves, func(i, j int) bool {
		return elves[i].GetTotalCalories() >= elves[j].GetTotalCalories()
	})

	topElves := elves[0:elvesCount]

	return topElves
}

type Elf struct {
	Id       int
	Calories []int
}

func (this Elf) GetTotalCalories() int {
	sum := 0
	for _, v := range this.Calories {
		sum += v
	}
	return sum
}

func (this Elf) AddCalorie(calorie int) Elf {
	this.Calories = append(this.Calories, calorie)
	return this
}
