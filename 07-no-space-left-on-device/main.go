// go run main.go
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	var inputFilePath = flag.String("inputFilePath", "./inputf.txt", "Input File")
	flag.Parse()

	log.Printf("inputFilePath %v\n", *inputFilePath)

	fs, err := parseFileToFileSystem(*inputFilePath)

	if err != nil {
		log.Fatalf("Error Parsing file %v: %v", *inputFilePath, err)
		return
	}

	log.Printf("> (1st Puzzle) Find all of the directories with a total size of at most 100000. What is the sum of the total sizes of those directories?")

	directoriesWhichTotalSizeIsLowerThanSizeLimit := fs.FindDirectoriesWithTotalSizeOfAtMost(PuzzleDirectorySizeLimit)

	sum := uint(0)

	for _, v := range directoriesWhichTotalSizeIsLowerThanSizeLimit {
		sum += v.TotalSize
	}

	log.Printf("Sum of diretories which total size is lower than %d is: %d", PuzzleDirectorySizeLimit, sum)
}

const (
	CommandExecutionIndicator  = "$"
	DirectoryIndicator         = "$"
	ProgramArgumentsStartIndex = 2
	PuzzleDirectorySizeLimit   = 100000
)

type DataBuffer struct {
	ByteStream []byte
}

type Directory struct {
	ParentDirectory *Directory
	Name            string
	Files           []File
	Directories     []*Directory
	TotalSize       uint
}

type File struct {
	Name string
	Size int
}

type Program struct {
	Command   string
	Arguments []string
}

func (d Directory) FindDirectoriesWithTotalSizeOfAtMost(sizeLimit uint) []Directory {
	ds := []Directory{}

	if d.TotalSize <= sizeLimit {
		ds = append(ds, d)
	}

	for _, v := range d.Directories {
		ds = append(ds, v.FindDirectoriesWithTotalSizeOfAtMost(sizeLimit)...)
	}

	return ds
}

func (d Directory) WalkToRoot() Directory {
	if d.Name == "/" {
		return d
	} else {
		return d.ParentDirectory.WalkToRoot()
	}
}

func (d *Directory) ComputeTotalSize() uint {

	ts := d.TotalSize

	if ts == 0 {
		for _, f := range d.Files {
			ts += uint(f.Size)
		}

		for _, d := range d.Directories {
			ts += d.ComputeTotalSize()
		}

		d.TotalSize = ts
	}

	return ts
}

func (p Program) IsListFilesProgram() bool {
	return p.Command == "ls" || p.Command == "dir"
}

func (p Program) IsChangeDirectoryProgram() bool {
	return p.Command == "cd"
}

func (p Program) IsChangeDirectoryToUpperLevelProgram() bool {
	return !p.IsChangeDirectoryToLowerLevelProgram()
}

func (p Program) IsChangeDirectoryToLowerLevelProgram() bool {
	return p.Arguments[len(p.Arguments)-1] == ".."
}

func isProgramExecutionLine(s string) bool {
	return strings.HasPrefix(s, CommandExecutionIndicator)
}

func isDirectoryLine(s string) bool {
	return strings.HasPrefix(s, DirectoryIndicator)
}

func lineToProgram(l string) Program {
	commandExecutionSplit := strings.Split(l, " ")
	command := commandExecutionSplit[1]
	arguments := []string{}

	if len(commandExecutionSplit) > ProgramArgumentsStartIndex {
		arguments = commandExecutionSplit[ProgramArgumentsStartIndex:]
	}

	return Program{
		Command:   command,
		Arguments: arguments,
	}
}

func lineToDirectory(l string) Directory {
	lineSplit := strings.Split(l, " ")

	return Directory{
		Name: lineSplit[1],
	}
}

func lineToFile(l string) File {
	lineSplit := strings.Split(l, " ")
	fileSizeParse, _ := strconv.Atoi(lineSplit[0])

	return File{
		Size: fileSizeParse,
		Name: lineSplit[1],
	}
}

func parseFileToFileSystem(filePath string) (*Directory, error) {
	lines, err := getFileLines(filePath)

	if err != nil {
		return nil, err
	}

	lc := len(lines)

	var currentDirectory *Directory
	collectDirectoryFiles := false

	for i := 0; i < lc; i++ {
		line := lines[i]

		if isProgramExecutionLine(line) {
			p := lineToProgram(line)

			if p.IsChangeDirectoryProgram() {
				if p.IsChangeDirectoryToUpperLevelProgram() {
					dirName := p.Arguments[len(p.Arguments)-1]

					d := Directory{
						Name: dirName,
					}

					d.ParentDirectory = currentDirectory

					currentDirectory = &d
				} else {
					d := currentDirectory
					currentDirectory = currentDirectory.ParentDirectory

					currentDirectory.Directories = append(currentDirectory.Directories, d)
				}
			} else if p.IsListFilesProgram() {
				collectDirectoryFiles = true
			}
		} else if collectDirectoryFiles {
			if isDirectoryLine(line) {
				d := lineToDirectory(line)

				currentDirectory.Directories = append(currentDirectory.Directories, &d)
			} else {
				f := lineToFile(line)

				currentDirectory.Files = append(currentDirectory.Files, f)
			}
		}
	}

	ruteDirectory := currentDirectory.WalkToRoot()
	ruteDirectory.ComputeTotalSize()

	return &ruteDirectory, nil
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

	lineBreak := "\n"

	if runtime.GOOS == "windows" {
		lineBreak = "\r"
	}

	return strings.Split(string(rawBytes), lineBreak), nil
}
