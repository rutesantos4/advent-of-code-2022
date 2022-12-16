package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	zoneMap, err := parseFileToZoneMap(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Printf("> (1st Puzzle) Consult the report from the sensors you just deployed. In the row where y=2000000, how many positions cannot contain a beacon?")

	numberPositionsCannotContainBeacon := zoneMap.ComputeNumberPositionsCannotContainBeacon(RowPuzzle1)

	log.Printf("The number of positions that cannot contain a beacon is: %v", numberPositionsCannotContainBeacon)

	log.Printf("> (2nd Puzzle) Find the only possible position for the distress beacon. What is its tuning frequency?")

	tuningFrequency := zoneMap.ComputeTuningFrequency()

	log.Printf("The tuning frequency is: %v", tuningFrequency)
}

const (
	// RowPuzzle1 = 10
	RowPuzzle1 = 2000000
)

const (
	SensorLabel = 'S'
	BeaconLabel = 'B'
)

const (
	positionsParseRegex = `(x=(-?\d+), y=(-?\d+))`
)

type Position struct {
	X int
	Y int
}

type Sensor = Position
type Beacon = Position

type Sensors = []Sensor
type Beacons = []Beacon

type ZoneMap struct {
	Sensors Sensors
	Beacons Beacons
}

func (z ZoneMap) ComputeTuningFrequency() int {
	p := z.findPossiblePositionForDistressBeacon()
	println(p.X)
	println(p.Y)
	return p.X*4000000 + p.Y
}

func (z ZoneMap) findPossiblePositionForDistressBeacon() Position {
	max := 4000000 + 1
	min := 0
	pempty := Position{}

	start := time.Now()
	defer func() {
		fmt.Println("Part2:", time.Since(start))
	}()

	// c1 := make(chan Position)
	// c2 := make(chan Position)
	// c3 := make(chan Position)
	// c4 := make(chan Position)
	// c5 := make(chan Position)
	// c6 := make(chan Position)
	// c7 := make(chan Position)
	// c8 := make(chan Position)
	// c9 := make(chan Position)
	// c10 := make(chan Position)
	// for i := min; i < max/10; i++ {
	// 	go findPositionPerLine(z, i, min, max, c1)

	// 	a := i + (max / 10)
	// 	go findPositionPerLine(z, a, min, max, c2)

	// 	b := a + (max / 10)
	// 	go findPositionPerLine(z, b, min, max, c3)

	// 	c := b + (max / 10)
	// 	go findPositionPerLine(z, c, min, max, c4)

	// 	d := c + (max / 10)
	// 	go findPositionPerLine(z, d, min, max, c5)

	// 	e := d + (max / 10)
	// 	go findPositionPerLine(z, e, min, max, c6)

	// 	f := e + (max / 10)
	// 	go findPositionPerLine(z, f, min, max, c7)

	// 	g := f + (max / 10)
	// 	go findPositionPerLine(z, g, min, max, c8)

	// 	h := g + (max / 10)
	// 	go findPositionPerLine(z, h, min, max, c9)

	// 	j := h + (max / 10)
	// 	go findPositionPerLine(z, j, min, max, c10)

	// }

	// for p := range c1 {
	// 	if p != pempty {
	// 		return p
	// 	}
	// }

	// for p := range c2 {
	// 	if p != pempty {
	// 		return p
	// 	}
	// }

	// for p := range c3 {
	// 	if p != pempty {
	// 		return p
	// 	}
	// }

	// for p := range c4 {
	// 	if p != pempty {
	// 		return p
	// 	}
	// }

	// for p := range c5 {
	// 	if p != pempty {
	// 		return p
	// 	}
	// }

	ch := make(chan Position)
	for row := min; row <= max; row++ {
		go findPositionPerLine(z, row, min, max, ch)
	}

	pempty = <-ch

	return pempty
}

func findPositionPerLine(z ZoneMap, i int, min int, max int, c chan Position) {
	println(i)

	p := Position{}
	listCannotContain := z.positionsCannotContainBeacon(i)

	for _, beacon := range z.Beacons {
		if beacon.Y == i {
			// beacon is in the same row
			listCannotContain = append(listCannotContain, beacon.X)
		}
	}

	listCannotContain = removeOutOfBounds(listCannotContain, min, max)

	if len(listCannotContain) <= max-min {

		x := findValuesNotPresentBetween(listCannotContain, min, max)

		if len(x) == 1 {
			println("here")
			log.Printf("%v %v", x, i)
			p.X = x[0]
			p.Y = i
			c <- p
		}
	}
}

func (z ZoneMap) ComputeNumberPositionsCannotContainBeacon(row int) int {
	return len(z.positionsCannotContainBeacon(row))
}

func (z ZoneMap) positionsCannotContainBeacon(row int) []int {
	p := []int{}

	for s, sensor := range z.Sensors {
		beacon := z.Beacons[s]

		diffX := math.Abs(float64(sensor.X) - float64(beacon.X))
		diffY := math.Abs(float64(sensor.Y) - float64(beacon.Y))
		dist := int(diffX + diffY)

		for distY := 0; distY <= dist; distY++ {

			y := sensor.Y + distY
			p = sensor.addPositionsCannotContainBeacon(y, p, (dist - distY), row)

			y = sensor.Y - distY
			p = sensor.addPositionsCannotContainBeacon(y, p, (dist - distY), row)
		}

	}

	//remove duplicates from p
	p = removeDuplicateInt(p)

	for _, beacon := range z.Beacons {
		if beacon.Y == row {
			// beacon is in the same row
			p = removeElementInt(p, beacon.X)
		}
	}

	return p
}

func (s Sensor) addPositionsCannotContainBeacon(y int, p []int, dist int, row int) []int {
	if y == row {
		p = append(p, s.X)

		xRight := s.X + dist

		for i := s.X; i <= xRight; i++ {
			p = append(p, i)
		}

		xLeft := s.X - dist

		for i := xLeft; i <= s.X; i++ {
			p = append(p, i)
		}
	}
	return p
}

func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func removeElementInt(intSlice []int, el int) []int {

	for i, item := range intSlice {

		if item == el {
			return append(intSlice[:i], intSlice[i+1:]...)
		}
	}

	return intSlice
}

func findValuesNotPresentBetween(intSlice []int, min int, max int) []int {
	result := []int{}

	for i := min; i < max; i++ {

		if !contains(intSlice, i) {
			result = append(result, i)
		}
	}

	return result
}

func contains(intSlice []int, el int) bool {
	for _, v := range intSlice {
		if v == el {
			return true
		}
	}

	return false
}

func removeOutOfBounds(intSlice []int, min int, max int) []int {
	for i := 0; i < len(intSlice); i++ {

		if intSlice[i] < min {
			intSlice = append(intSlice[:i], intSlice[i+1:]...)
			i--
		} else if intSlice[i] > max {
			intSlice = append(intSlice[:i], intSlice[i+1:]...)
			i--
		}
	}

	return intSlice
}

func lineToSensorAndBeacon(line string) (Sensor, Beacon) {
	positionsRegex := regexp.MustCompile(positionsParseRegex)

	positionString := positionsRegex.FindAllStringSubmatch(line, 2)
	sPositionString := positionString[0][2:]
	bPositionString := positionString[1][2:]

	sensorPositionXParse, _ := strconv.Atoi(sPositionString[0])
	sensorPositionYParse, _ := strconv.Atoi(sPositionString[1])

	beaconPositionXParse, _ := strconv.Atoi(bPositionString[0])
	beaconPositionYParse, _ := strconv.Atoi(bPositionString[1])

	return Sensor{
			X: sensorPositionXParse,
			Y: sensorPositionYParse,
		}, Beacon{
			X: beaconPositionXParse,
			Y: beaconPositionYParse,
		}
}

func parseFileToZoneMap(filePath string) (*ZoneMap, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	ll := len(lines)

	sensors := make(Sensors, ll)
	beacons := make(Beacons, ll)

	for i := 0; i < ll; i++ {
		l := lines[i]

		sensor, beacon := lineToSensorAndBeacon(l)

		sensors[i] = sensor
		beacons[i] = beacon
	}

	return &ZoneMap{
		Sensors: sensors,
		Beacons: beacons,
	}, nil
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
