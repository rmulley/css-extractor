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
	// Retrieve flags and set globals
	flag.StringVar(&I_FILENAME, "i", "", "'i' parameter must be specificed.")
	flag.StringVar(&O_FILENAME, "o", "", "'o' parameter must be specificed.")
	flag.Parse()

	if strings.TrimSpace(I_FILENAME) == "" {
		log.Fatalln("'i' parameter must be specificed.")
	} //if

	if strings.TrimSpace(O_FILENAME) == "" {
		log.Fatalln("'o' parameter must be specificed.")
	} //if
} //init

func main() {
	var (
		err              error
		cssToAdd         string
		lines            []string
		allCssRules      = make(map[string][]CssRule)
		findInlineCSS    = regexp.MustCompile("<(\\w+)\\s.*?style=\"(.*?)\".*?>")
		findClasses      = regexp.MustCompile("class=\"(.*?)\"")
		findId           = regexp.MustCompile("id=\"([\\w\\-]+)\"")
		replaceStyleTags = regexp.MustCompile("(.*?)style=\".*?\"(.*)")
	) //var

	// Extract each line of the file
	if lines, err = readFile(I_FILENAME); err != nil {
		log.Fatalln(err)
	} //if

	for _, line := range lines {
		temp := extractInlineCSS(line, findInlineCSS, findId, findClasses)

		for k, v := range temp {
			if _, isSet := allCssRules[k]; isSet {
				allCssRules[k] = append(allCssRules[k], v...)
			} else {
				allCssRules[k] = v
			} //else
		} //for
	} //for

	for element, rules := range allCssRules {
		cssToAdd += element + " {"

		for _, val := range rules {
			cssToAdd += "\n\t" + val.Rule + ":\t" + val.Value
		} //for

		cssToAdd += "\n}\n"
	} //for

	if err = createCssFile(O_FILENAME, cssToAdd); err != nil {
		log.Fatalln(err)
	} //if

	if err = removeStyleTags(I_FILENAME, lines, replaceStyleTags); err != nil {
		log.Fatalln(err)
	} //if
} //main

func extractInlineCSS(line string, findInlineCSS, findId, findClasses *regexp.Regexp) (cssRules map[string][]CssRule) {
	cssRules = make(map[string][]CssRule)
	styles := findInlineCSS.FindAllStringSubmatch(line, -1)

	for _, v := range styles {
		var (
			id, tag   string
			eachClass []string
		) //var

		for pos, _ := range v {
			// First position is entire capture
			if pos == 0 {
				ids := findId.FindAllStringSubmatch(v[pos], -1)

				if len(ids) > 0 && len(ids[0]) > 0 {
					id = ids[0][1]
				} else {
					classes := findClasses.FindAllStringSubmatch(v[pos], -1)

					for _, allClasses := range classes {
						for pos2, classes2 := range allClasses {
							if pos2 != 0 {
								eachClass = strings.Split(classes2, " ")
							} //for
						} //if
					} //for
				} //else
			} else if pos == 1 { // Second position is tag name
				element := strings.TrimSpace(v[pos])

				if len(id) > 0 {
					tag = element + "#" + id
				} else {
					for _, class := range eachClass {
						tag += " " + element + "." + class
					} //for
				} //else
			} else if len(tag) > 0 { // pos = 3 is the style tag.  We are ignoring tags w/out a class or ID right now
				var (
					items []string
				) //var

				items = strings.Split(v[pos], ";")

				for _, item := range items {
					if len(item) > 0 {
						rule := strings.Split(item, ":")
						rule[0] = strings.TrimSpace(rule[0])
						rule[1] = strings.TrimSpace(rule[1])

						if rule[1][len(rule[1])-1:] != ";" {
							rule[1] += ";"
						} //if

						cssRules[strings.TrimSpace(tag)] = append(cssRules[tag],
							CssRule{
								Rule:  rule[0],
								Value: rule[1],
							}) //append
					} //if
				} //for
			} //else
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

	// Write to new CSS file
	if _, err = oFile.Write([]byte(strings.TrimSpace(css))); err != nil {
		return err
	} //if

	return nil
} //writeCssFile

func removeStyleTags(filename string, lines []string, replaceStyleTags *regexp.Regexp) (err error) {
	var (
		oFile *os.File
	) //var

	// Attempt to create the file
	if oFile, err = os.OpenFile(filename+"_backup", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660); err != nil {
		return err
	} //if
	defer oFile.Close()

	for _, line := range lines {
		var (
			newLine string
			styles  [][]string
		) //var

		styles = replaceStyleTags.FindAllStringSubmatch(line, -1)

		for _, v := range styles {
			for pos, w := range v {
				// Ignore first element which is entire line
				if pos != 0 {
					newLine += w
				} //if
			} //for
		} //for

		// If no style tags were found, use original line
		if newLine == "" {
			newLine = line
		} //if

		// Append line to file
		if _, err = oFile.Write([]byte(strings.TrimSpace(newLine) + "\n")); err != nil {
			return err
		} //if
	} //for

	return err
} //removeStyleTags
