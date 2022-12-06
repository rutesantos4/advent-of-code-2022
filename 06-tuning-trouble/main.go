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

	pmp := dataBuffer.FindPacketMarkersPosition()

	fmt.Printf("%v\n", pmp)
}

const (
	SequenceOfDifferentBytesUntilPacketMarker = 4
)

type DataBuffer struct {
	ByteStream []byte
}

func (b DataBuffer) FindPacketMarkersPosition() map[int]byte {
	bs := b.ByteStream
	bslen := len(bs)

	pmp := map[int]byte{}

	for i := 0; i < bslen; i++ {
		ni := i + 1
		si := ni + SequenceOfDifferentBytesUntilPacketMarker + 1
		if ni < bslen && si < bslen {
			seq := bs[ni:si]
			set := map[byte]bool{}

			for _, v := range seq {
				set[v] = true
			}

			if len(set) != len(seq) {
				pmp[si] = seq[SequenceOfDifferentBytesUntilPacketMarker]

				i += SequenceOfDifferentBytesUntilPacketMarker
			}
		}
	}

	return pmp
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
