package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
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
	flag.StringVar(&O_FILENAME, "ofile", "", "Unable to set O_FILENAME flag variable.")
	flag.Parse()

	if I_FILENAME == "" {
		log.Fatalln("ifile parameter must be specificed.")
	} //if
} //init

func main() {
	var (
		err           error
		cssToAdd      string
		lines         []string
		allCssRules   = make(map[string][]CssRule)
		findInlineCSS = regexp.MustCompile("<(\\w+)\\s(.*?)style=\"([\\w\\-]+):([\\w]+);\"")
		findClasses   = regexp.MustCompile("class=\"(.*?)\"")
	) //var

	// Extract each line of the file
	if lines, err = readFile(I_FILENAME); err != nil {
		log.Fatalln(err)
	} //if
	ctr := 0
	for _, line := range lines {
		if ctr == 28 {
			// Determine if
			if allCssRules = extractInlineCSS(line, findInlineCSS, findClasses); err != nil {
				log.Println(err)
			} //if
		} //if

		ctr++
	} //for

	for element, rules := range allCssRules {
		log.Println(element)
		log.Println(rules)
		log.Println("")

		cssToAdd += element + "{"

		for _, val := range rules {
			cssToAdd += "\n\t" + val.Rule + ":\t" + val.Value
		} //for

		cssToAdd += "\n}\n"
	} //for

	log.Println("CSS to add:")
	log.Println(cssToAdd)

	if err = createCssFile(O_FILENAME, cssToAdd); err != nil {

	} //if
} //main

func extractInlineCSS(line string, findInlineCSS, findClasses *regexp.Regexp) (cssRules map[string][]CssRule) {
	//log.Println(line)
	cssRules = make(map[string][]CssRule)
	styles := findInlineCSS.FindAllStringSubmatch(line, -1)

	for _, v := range styles {
		var tag string

		for pos, _ := range v {
			// First position is entire capture
			if pos != 0 {
				// Second position is tag name
				if pos == 1 {
					tag = strings.TrimSpace(v[pos])
				} else if pos == 2 {
					log.Println(v[pos])
					classes := findClasses.FindAllStringSubmatch(v[pos], -1)

					for _, allClasses := range classes {
						for pos2, classes2 := range allClasses {
							if pos2 != 0 {
								var origTag string

								origTag = tag
								tag = ""

								eachClass := strings.Split(classes2, " ")

								for _, class := range eachClass {
									tag += " " + origTag + "." + class
								} //for
							} //for
						} //if
					} //for
				} else if pos%2 != 0 {
					// Odd positions are rules
					// Even positions are values
					cssRules[tag] = append(cssRules[tag],
						CssRule{
							Rule:  strings.TrimSpace(v[pos]),
							Value: strings.TrimSpace(v[pos+1]),
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

func createCssFile(filename, css string) (err error) {
	var (
		oFile *os.File
	) //var

	// Attempt to create the file
	if oFile, err = os.Create(filename); err != nil {
		return err
	} //if
	defer oFile.Close()

	if _, err = oFile.Write([]byte(strings.TrimSpace(css))); err != nil {
		return err
	} //if

	return nil
} //writeCssFile

func addCssToHead(filename string) (err error) {

	return err
} //addCssToHead

func removeStyleTags(filename string) (err error) {

	return err
} //removeStyleTags
