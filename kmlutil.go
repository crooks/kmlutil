package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// Make a <tag> regexp global so as not to recompile within a loop
var tagRE = regexp.MustCompile("<([^>]+)>")

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Too many arguments.")
		os.Exit(0)
	} else if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "No filename specified.")
		os.Exit(0)
	}
	f, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
	defer f.Close()
	readKML(f, os.Stdout)
}

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func reverseCoords(s string) (revStr string) {
	coords := strings.Split(s, " ")
	reverse(coords)
	return strings.Join(coords, " ")
}

func readKML(r io.Reader, w io.Writer) {
	var flagCoords bool
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if getTag(line) == "/coordinates" {
			flagCoords = false
		}
		if flagCoords {
			line = reverseCoords(line)
		}
		io.WriteString(w, line+"\n")
		if getTag(line) == "coordinates" {
			flagCoords = true
		}
	}
}

func getTag(s string) string {
	matchTags := tagRE.FindStringSubmatch(s)
	if len(matchTags) == 0 {
		return ""
	}
	return matchTags[1]
}
