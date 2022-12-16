package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"sync"
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
	listCannotContain := make([][]int, max+1)
	p := Position{}

	var wg sync.WaitGroup
	wg.Add(max)
	for i := min; i < max; i++ {
		go func(i int) {
			println(i)
			listCannotContain[i] = z.positionsCannotContainBeacon(i)

			for _, beacon := range z.Beacons {
				if beacon.Y == i {
					// beacon is in the same row
					listCannotContain[i] = append(listCannotContain[i], beacon.X)
				}
			}
			// listCannotContain[i] = removeOutOfBounds(listCannotContain[i], min, max)

			// if len(listCannotContain[i]) <= max-min {

			x := findValuesNotPresentBetween(listCannotContain[i], min, max)

			if len(x) == 1 {
				println("here")
				log.Printf("%v %v", x, i)
				p.X = x[0]
				p.Y = i
			}
			wg.Done()
			// }
		}(i)
	}
	wg.Wait()

	return p
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
