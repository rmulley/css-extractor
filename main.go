package main

import (
	"bufio"
	"flag"
	"log"
	"os"
) //import

// Global variables
var (
	I_FILENAME string
	O_FILENAME string
) //var

func init() {
	// Retrieve flags & set globals
	flag.StringVar(&I_FILENAME, "ifile", "", "Unable to set I_FILENAME flag variable.")
	flag.Parse()

	//HOST = strings.TrimSpace(strings.ToLower(HOST))
} //init

func main() {
	var (
		err   error
		lines []string
	) //var

	// Extract each line of the file
	if lines, err = readFile(I_FILENAME); err != nil {
		log.Fatalln(err)
	} //if

	for _, line := range lines {

	} //for
} //main

func readFile(filename string) (lines []string, err error) {
	var (
		iFile   *os.File
		scanner *bufio.Scanner
	) //var

	// Attempt to open the file
	if iFile, err = os.Open(I_FILENAME); err != nil {
		return lines, err
	} //if
	defer iFile.Close()

	scanner = bufio.NewScanner(iFile)

	// Read in each line of file
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	} //for

	return lines, scanner.Err()
} //readFile
