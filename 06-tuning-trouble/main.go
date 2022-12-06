// go run main.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	dataBuffer, err := parseFile(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Printf("> (1st Puzzle) How many characters need to be processed before the first start-of-packet marker is detected?")

	firstPacketMarkersPosition := dataBuffer.FindFirstPacketMarkersPosition()

	fmt.Printf("The number of characters that need to be processed before the first start-of-packet marker is detected is: %v\n", firstPacketMarkersPosition)

	log.Printf("> (2st Puzzle) How many characters need to be processed before the first start-of-message marker is detected?")

	firstMessageMarkersPosition := dataBuffer.FindFirstMessageMarkersPosition()

	fmt.Printf("The number of characters that need to be processed before the first start-of-message marker is detected is: %v\n", firstMessageMarkersPosition)

}

const (
	SequenceOfDifferentBytesUntilPacketMarker  = 4
	SequenceOfDifferentBytesUntilMessageMarker = 14
)

type DataBuffer struct {
	ByteStream []byte
}

func (b DataBuffer) FindFirstPacketMarkersPosition() int {
	packetMarkersPosition := b.FindMarkersPosition(SequenceOfDifferentBytesUntilPacketMarker)

	return getSmallestKey(packetMarkersPosition)
}

func (b DataBuffer) FindFirstMessageMarkersPosition() int {
	packetMarkersPosition := b.FindMarkersPosition(SequenceOfDifferentBytesUntilMessageMarker)

	return getSmallestKey(packetMarkersPosition)
}

func (b DataBuffer) FindMarkersPosition(incremental int) map[int]byte {
	bs := b.ByteStream
	bslen := len(bs)

	pmp := map[int]byte{}

	for i := 0; i < bslen; i++ {
		ni := i
		si := ni + incremental
		if ni < bslen && si < bslen {
			seq := bs[ni:si]
			set := map[byte]bool{}

			for _, v := range seq {
				set[v] = true
			}

			if len(set) == len(seq) {
				pmp[si] = seq[incremental-1]
			}
		}
	}

	return pmp
}

func getSmallestKey(m map[int]byte) int {

	keys := make([]int, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	minKey := keys[0]
	for _, v := range keys {
		if v < minKey {
			minKey = v
		}
	}

	return minKey
}

func parseFile(filePath string) (*DataBuffer, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	bs := []byte(lines[0])

	return &DataBuffer{ByteStream: bs}, nil
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
