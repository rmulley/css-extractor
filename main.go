package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
) //import

type CssRule struct {
	Rule  string
	Value string
} //struct

// Global variables
var (
	I_FILENAME string
	O_FILENAME string
) //var

func init() {
	// Retrieve flags & set globals
	flag.StringVar(&I_FILENAME, "ifile", "", "Unable to set I_FILENAME flag variable.")
	flag.Parse()

	if I_FILENAME == "" {
		log.Fatalln("ifile parameter must be specificed.")
	} //if
} //init

func main() {
	var (
		err           error
		lines         []string
		allCssRules   = make(map[string][]CssRule)
		findInlineCSS = regexp.MustCompile("<(\\w+)\\s.*?style=\"([\\w\\-]+):([\\w]+);\"")
	) //var

	// Extract each line of the file
	if lines, err = readFile(I_FILENAME); err != nil {
		log.Fatalln(err)
	} //if
	ctr := 0
	for _, line := range lines {
		if ctr == 28 {
			// Determine if
			if allCssRules = extractInlineCSS(line, findInlineCSS); err != nil {
				log.Println(err)
			} //if
		} //if

		ctr++
	} //for

	for k, v := range allCssRules {
		log.Println(k)
		log.Println(v)
		log.Println("")
	} //for
} //main

func extractInlineCSS(line string, findInlineCSS *regexp.Regexp) (cssRules map[string][]CssRule) {
	log.Println(line)
	cssRules = make(map[string][]CssRule)
	styles := findInlineCSS.FindAllStringSubmatch(line, -1)

	for _, v := range styles {
		var tag string

		for pos, _ := range v {
			// First position is entire capture
			if pos != 0 {
				// Second position is tag name
				if pos == 1 {
					tag = v[pos]
				} else if pos%2 == 0 {
					// Even positions are rules
					// Odd positions are values
					cssRules[tag] = append(cssRules[tag],
						CssRule{
							Rule:  v[pos],
							Value: v[pos+1],
						}) //append
				} //elseif
			} //if
		} //for
	} //for

	return cssRules
} //extractInlineCSS

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
