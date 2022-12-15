package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
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

	log.Printf("> (1st Puzzle) Using your scan, simulate the falling sand. How many units of sand come to rest before sand starts flowing into the abyss below?")

	log.Printf("%v\n", zoneMap)
}

const (
	SensorLabel = 'S'
	BeaconLabel = 'B'
)

const (
	positionsParseRegex = `(x=re(.+), y=re(.+).*)`
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

func lineToSensorAndBeacon(line string) (Sensor, Beacon) {
	positionsRegex := regexp.MustCompile(positionsParseRegex)

	xPositionString := positionsRegex.FindAllString(line, 1)
	yPositionString := positionsRegex.FindAllString(line, 2)

	log.Printf("%v\n", xPositionString)
	log.Printf("%v\n", yPositionString)

	sensorPositionXParse, _ := strconv.Atoi(xPositionString[0])
	sensorPositionYParse, _ := strconv.Atoi(yPositionString[0])

	beaconPositionXParse, _ := strconv.Atoi(xPositionString[1])
	beaconPositionYParse, _ := strconv.Atoi(yPositionString[1])

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

		sensors = append(sensors, sensor)
		beacons = append(beacons, beacon)
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
