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

	if O_FILENAME == "" {
		log.Fatalln("ofile parameter must be specificed.")
	} //if
} //init

func main() {
	var (
		err         error
		cssToAdd    string
		lines       []string
		allCssRules = make(map[string][]CssRule)
		//findInlineCSS    = regexp.MustCompile("<(\\w+)\\s(.*?)style=\"([\\w\\-]+):\\s*([\\w]+);\"")
		findInlineCSS    = regexp.MustCompile("<(\\w+)\\s(.*?)style=\"(.*?)\"")
		findClasses      = regexp.MustCompile("class=\"(.*?)\"")
		replaceStyleTags = regexp.MustCompile("(.*?)style=\".*?\"(.*)")
	) //var

	// Extract each line of the file
	if lines, err = readFile(I_FILENAME); err != nil {
		log.Fatalln(err)
	} //if

	for _, line := range lines {
		temp := extractInlineCSS(line, findInlineCSS, findClasses)

		for k, v := range temp {
			if _, isSet := allCssRules[k]; isSet {
				allCssRules[k] = append(allCssRules[k], v...)
			} else {
				allCssRules[k] = v
			} //else
		} //for
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
		log.Fatalln(err)
	} //if

	if err = removeStyleTags(I_FILENAME, lines, replaceStyleTags); err != nil {
		log.Fatalln(err)
	} //if
} //main

func extractInlineCSS(line string, findInlineCSS, findClasses *regexp.Regexp) (cssRules map[string][]CssRule) {
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
				} else { // pos = 4 is the style tag
					var (
						items []string
					) //var

					items = strings.Split(v[pos], ";")

					for _, item := range items {
						if len(item) > 0 {
							log.Println("item:", item)

							rule := strings.Split(item, ":")
							log.Println("rule: ", rule)

							rule[0] = strings.TrimSpace(rule[0])
							rule[1] = strings.TrimSpace(rule[1])

							if rule[1][len(rule[1])-1:] != ";" {
								rule[1] += ";"
							} //if

							cssRules[tag] = append(cssRules[tag],
								CssRule{
									Rule:  rule[0],
									Value: rule[1],
								}) //append
						} //if
					} //for
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
	if oFile, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0660); err != nil {
		return err
	} //if
	defer oFile.Close()

	if _, err = oFile.Write([]byte(strings.TrimSpace(css))); err != nil {
		return err
	} //if

	return nil
} //writeCssFile

func removeStyleTags(filename string, lines []string, replaceStyleTags *regexp.Regexp) (err error) {
	var (
		oFile *os.File
	) //var

	log.Println("Output:", filename+"_backup")

	// Attempt to create the file
	if oFile, err = os.OpenFile(filename+"_backup", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660); err != nil {
		return err
	} //if
	defer oFile.Close()

	for _, line := range lines {
		styles := replaceStyleTags.FindAllStringSubmatch(line, -1)

		for _, v := range styles {
			for _, w := range v {
				line = w
			} //for
		} //for

		log.Println("LINE:", line)

		if len(line) > 1 {
			if line[len(line)-2:len(line)-1] != ">" {
				line += ">"
			} //if
		} //if

		if _, err = oFile.Write([]byte(strings.TrimSpace(line) + "\n")); err != nil {
			return err
		} //if
	} //for

	return err
} //removeStyleTags
