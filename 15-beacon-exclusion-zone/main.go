package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
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

type Range struct {
	min int
	max int
}

func (r *Range) Len() int {
	return r.max - r.min
}

type Ranges []Range

// Implementing these because of sort function
func (r Ranges) Len() int           { return len(r) }
func (r Ranges) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Ranges) Less(i, j int) bool { return r[i].min < r[j].min }

func (z ZoneMap) ComputeTuningFrequency() int {
	p := z.findPossiblePositionForDistressBeacon()
	return p.X*4000000 + p.Y
}

func (z ZoneMap) findPossiblePositionForDistressBeacon() Position {
	max := 4000000 + 1
	min := 0

	p := Position{}
	for row := min; row <= max; row++ {
		listCannotContain := z.positionsCannotContainBeacon(row)

		if len(listCannotContain) > 1 {
			interval := listCannotContain[1]
			p.X = interval.min - 1
			p.Y = row
			break
		}
	}

	return p
}

func (z ZoneMap) ComputeNumberPositionsCannotContainBeacon(row int) int {
	ranges := z.positionsCannotContainBeacon(row)
	sum := 0

	for _, interval := range ranges {
		sum += interval.Len()
	}

	return sum
}

func (z ZoneMap) positionsCannotContainBeacon(row int) Ranges {
	list := Ranges{}

	for s, sensor := range z.Sensors {
		beacon := z.Beacons[s]

		diffX := math.Abs(float64(sensor.X) - float64(beacon.X))
		diffY := math.Abs(float64(sensor.Y) - float64(beacon.Y))
		dist := int(diffX + diffY)

		// distance between the row and the sensor.Y
		distY := int(math.Abs(float64(row - sensor.Y)))

		distX := dist - distY
		if distX < 0 {
			// We are not able to conclude that the row cannot contain beacon with this sensor/beacon
			continue
		}

		interval := Range{
			min: sensor.X - distX,
			max: sensor.X + distX,
		}

		list = append(list, interval)

	}

	sort.Sort(list)

	count := len(list)
	cleanList := []Range{}
	interval := list[0]

	for i := 1; i < count; i++ {
		if list[i].min <= interval.max {
			interval.max = getBigger(interval.max, list[i].max)
			interval.min = getSmaller(interval.min, list[i].min)
			continue
		}
		cleanList = append(cleanList, interval)
		interval = list[i]
	}
	cleanList = append(cleanList, interval)

	return cleanList
}

func getBigger(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func getSmaller(a int, b int) int {
	if a < b {
		return a
	}
	return b
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
