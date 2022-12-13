package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	pairs, err := parseFileToDistressSignal(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Printf("> (1st Puzzle) Determine which pairs of packets are already in the right order. What is the sum of the indices of those pairs?")

	indicesSumOfOrdered := pairs.IndicesSumOfOrdered()

	log.Printf("The sum of the indices of the pairs of packets are already in the right order is: %d\n", indicesSumOfOrdered)

	log.Printf("> (2nd Puzzle) Organize all of the packets into the correct order. What is the decoder key for the distress signal?")

	decoderKey := pairs.FindDecoderKey()

	log.Printf("The decoder key for the distress signal is: %d\n", decoderKey)
}

const (
	StartList        = '['
	EndList          = ']'
	SeparateElements = ','
)

const (
	List   = 1
	Number = 2
)

const (
	Equals  = 0
	Bigger  = 1
	Smaller = -1
)

const (
	Divider2 = "[[2]]"
	Divider6 = "[[6]]"
)

type Packet struct {
	Type      int
	Value     int
	Children  []*Packet
	Parent    *Packet
	isDivider bool
}

type Pair struct {
	Left  Packet
	Right Packet
}

type DistressSignal []Pair

func (d DistressSignal) IndicesSumOfOrdered() int {
	count := len(d)
	sum := 0

	for i := 1; i < count+1; i++ {
		if d[i-1].isRightOrder() {
			sum += i
		}
	}

	return sum
}

func (d DistressSignal) FindDecoderKey() int {

	divider2 := lineToPacket(Divider2)
	divider6 := lineToPacket(Divider6)
	divider2.isDivider = true
	divider6.isDivider = true
	pair := Pair{
		Left:  divider2,
		Right: divider6,
	}
	d = append(d, pair)
	packetsOrdered := d.orderPackets()

	count := len(packetsOrdered)
	mul := 1
	for i := 1; i < count+1; i++ {
		if packetsOrdered[i-1].isDivider {
			mul *= i
		}
	}

	return mul
}

func (d DistressSignal) orderPackets() []Packet {
	var packets []Packet = []Packet{}
	for _, v := range d {
		packets = append(packets, v.Left)
		packets = append(packets, v.Right)
	}

	sort.SliceStable(packets, func(i, j int) bool {
		return comparePairs(packets[i], packets[j]) == Smaller
	})
	return packets
}

func (p *Packet) changeToList() {
	p.Type = List
	packet := Packet{
		Type:   Number,
		Value:  p.Value,
		Parent: p,
	}
	p.Children = []*Packet{&packet}
}

func (p Packet) compareNumberTo(right Packet) int {

	if p.Value == right.Value {
		return Equals
	}

	if p.Value < right.Value {
		return Smaller
	}

	return Bigger
}

func (p Pair) isRightOrder() bool {
	return comparePairs(p.Left, p.Right) == Smaller
}

func comparePairs(left Packet, right Packet) int {

	if left.Type == Number && right.Type == Number {
		return left.compareNumberTo(right)
	}

	// if left/rigth is number then change to list
	if left.Type == Number {
		left.changeToList()
	}

	if right.Type == Number {
		right.changeToList()
	}

	smallerLen := len(left.Children)
	if len(right.Children) < smallerLen {
		smallerLen = len(right.Children)
	}

	// compare children (for i -> comparePairs)
	for i := 0; i < smallerLen; i++ {
		if check := comparePairs(*left.Children[i], *right.Children[i]); check != Equals {
			return check
		}
	}

	// if len(left) is bigger then len(rigth) then inputs are not in the right order
	if len(left.Children) > len(right.Children) {
		return Bigger
	}

	if len(left.Children) < len(right.Children) {
		return Smaller
	}

	return Equals
}

func lineToPacket(line string) Packet {
	count := len(line)
	var packet *Packet = &Packet{
		Type:     List,
		Children: []*Packet{},
	}

	parent := packet
	open := 0

	for i := 1; i < count; i++ {
		char := line[i]

		switch char {
		case StartList:
			newPacket := &Packet{
				Type:     List,
				Children: []*Packet{},
				Parent:   packet,
			}
			packet.Children = append(packet.Children, newPacket)
			packet = newPacket
			open++

		case EndList:
			if open > 0 {
				packet = packet.Parent
			}
			open--

		case SeparateElements:
			continue

		default:
			//number can be more than one digit
			number := string(char)
			digits := numberOfDigits(line[i:])
			if digits > 1 {
				number = line[i : i+digits]
				i = i + digits - 1
			}
			packet.Children = append(packet.Children, parseToNumber(number))
		}
	}

	return *parent
}

func numberOfDigits(chars string) int {
	count := len(chars)

	for i := 0; i < count; i++ {
		char := chars[i]
		if char == SeparateElements || char == EndList || char == StartList {
			return i
		}
	}

	return 0
}

func parseToNumber(s string) *Packet {
	value, _ := strconv.Atoi(s)

	return &Packet{
		Type:  Number,
		Value: value,
	}
}

func linesToPair(lines []string) Pair {
	left := lineToPacket(lines[0])
	right := lineToPacket(lines[1])

	return Pair{
		Left:  left,
		Right: right,
	}
}

func parseFileToDistressSignal(filePath string) (DistressSignal, error) {
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

	distressSignalLinesGroupLen := 2
	pairsLen := len(cleanLines) / distressSignalLinesGroupLen

	distressSignal := make(DistressSignal, pairsLen)

	for i := 0; i < pairsLen; i++ {
		si := i * distressSignalLinesGroupLen
		ei := si + distressSignalLinesGroupLen

		distressSignal[i] = linesToPair(cleanLines[si:ei])
	}

	return distressSignal, nil
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
