package internal

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Grep struct {
	CaseInSensitive bool
	Recursive       bool
	ReadFromStdIn   bool
	Path            string
	Pattern         string
	OutputFile      string
}

func (g *Grep) MatchPattern() {
	if g.Recursive {
		g.matchPatternInDir(g.Path, g.Pattern)
	} else if g.ReadFromStdIn {
		g.matchPatternInStdIn(g.Pattern)
	} else {
		g.matchPatternInFile(g.Path, g.Pattern)
	}
}

func (g *Grep) CheckMatch(lines []string, pattern string) []string {
	var matchedLines []string
	for _, line := range lines {
		if g.CaseInSensitive {
			lineLowercase := strings.ToLower(line)
			patternLowerCase := strings.ToLower(pattern)
			if strings.Contains(lineLowercase, patternLowerCase) {
				matchedLines = append(matchedLines, line)
			}
		} else {
			if strings.Contains(line, pattern) {
				matchedLines = append(matchedLines, line)
			}
		}
	}
	return matchedLines
}

func (g *Grep) matchPatternInFile(path string, pattern string) {
	content, err := os.ReadFile(path)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	con := string(content)
	lines := strings.Split(con, "\n")

	matchedLines := g.CheckMatch(lines, pattern)
	g.WriteOutput(path, lines, matchedLines)
}

func (g *Grep) matchPatternInDir(dirName string, pattern string) {

	if dirName == "" || dirName == "." {
		dirName, _ = os.Getwd()
	}
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if info != nil {
			if !info.IsDir() {
				g.matchPatternInFile(path, pattern)
			}
		}
		return errors.New("check the directory name")
	})
	fmt.Println(err)
}

func (g *Grep) matchPatternInStdIn(pattern string) {

	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "^D" {
			break
		}
		lines = append(lines, line)
	}
	matchedLines := g.CheckMatch(lines, pattern)
	g.WriteOutput("", lines, matchedLines)

}

func (g *Grep) WriteOutput(path string, lines []string, matchedLines []string) {

	if g.OutputFile != "" {
		g.writeOutputToFile(path, lines, matchedLines)
		return
	}

	for _, line := range matchedLines {
		if g.Recursive {
			fmt.Print(path + "  ")
		}
		fmt.Println(line)
	}
}

func (g *Grep) writeOutputToFile(path string, lines []string, matchedLines []string) {
	outputFileName := g.OutputFile
	f, err := os.Create(outputFileName)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer f.Close()
	for _, line := range matchedLines {
		stringToWrite := ""
		if g.Recursive {
			stringToWrite = path + "   " + line + "\n"
		} else {
			stringToWrite = line + "\n"
		}
		_, err2 := f.WriteString(stringToWrite)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
}
