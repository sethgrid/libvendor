package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var SRCDIR string
var CWD string

func init() {
	var err error
	CWD, err = os.Getwd() // does not work with `go run <script>`; only works on a build
	if err != nil {
		fmt.Println("Error getting current working directory!")
		os.Exit(1)
	}
	flag.StringVar(&SRCDIR, "src", CWD+"/src", "Usually `./src`. The directory that we are going to vendorize.")
}

func main() {
	flag.Parse()

	fmt.Printf("Starting in %s\n\nYou need to: mv src _vendor, and make the changes below.\n\n", CWD)
	filepath.Walk(SRCDIR, ParseFile)
}

func IsGoExt(filename string) bool {
	return filename[len(filename)-3:] == ".go"
}

type CandidateLine struct {
	LineNum  int
	Text     string
	Filepath string
	Prepend  string
}

func ParseFile(path string, info os.FileInfo, err error) error {

	if !info.IsDir() && IsGoExt(info.Name()) {
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("error opening file: %v\n", err)
			os.Exit(1)
		}

		s := bufio.NewScanner(f)

		candidates := make([]CandidateLine, 0)
		lineNum := 0
		// a candidate line is eligible for rewriting
		candidateLine := false
		for s.Scan() {
			line := s.Text()
			lineNum++

			if strings.Contains(line, "import \"") {
				// single line import statement
				candidateLine = true

			} else if strings.Contains(line, "import (") {
				candidateLine = true

			} else if strings.Contains(line, ")") {
				// done importing
				candidateLine = false
			}

			// if the candidate line is importing with .com (google.com, github.com,...)
			if candidateLine && strings.Contains(line, ".com") {
				candidates = append(candidates, CandidateLine{LineNum: lineNum, Text: s.Text(), Filepath: path})
			}
		}

		// if this file has lines to be updated, report them
		if len(candidates) > 0 {
			fmt.Println(path, ":")
			index := strings.Index(path, "/src/")
			for _, c := range candidates {
				nestedDirs := c.Filepath[index+4:]
				c.Prepend = relativeDirs(strings.Count(nestedDirs, "/")) + "_vendor/"
				fmt.Printf("\tLine %d: %s; prepend with %s\n", c.LineNum, c.Text, c.Prepend)
			}
		}
	}
	return nil
}

func relativeDirs(count int) string {
	str := ""
	for i := 0; i < count; i++ {
		str += "../"
	}
	return str
}
