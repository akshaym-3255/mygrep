package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

func (g *Grep) MatchPattern() ([]string, error) {
	var err error
	var matchedLines []string
	if g.Recursive {
		matchedLines, err = g.matchPatternInDir(g.Path, g.Pattern)
	} else if g.ReadFromStdIn {
		matchedLines = g.matchPatternInStdIn(g.Pattern, os.Stdin)
	} else {
		matchedLines, err = g.matchPatternInFile(g.Path, g.Pattern)
	}
	return matchedLines, err
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

func (g *Grep) matchPatternInFile(path string, pattern string) ([]string, error) {
	content, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}
	con := string(content)
	lines := strings.Split(con, "\n")

	matchedLines := g.CheckMatch(lines, pattern)
	return matchedLines, nil
}

func (g *Grep) matchPatternInDir(dirName string, pattern string) ([]string, error) {

	if dirName == "" || dirName == "." {
		dirName, _ = os.Getwd()
	}
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		return nil, errors.New("directory not present")
	}

	var matchedLinesInDir []string
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			matchedLines, err := g.matchPatternInFile(path, pattern)
			if err != nil {
				return err
			}
			for _, line := range matchedLines {
				matchWithPath := path + "   " + line
				matchedLinesInDir = append(matchedLinesInDir, matchWithPath)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matchedLinesInDir, nil
}

func (g *Grep) matchPatternInStdIn(pattern string, reader io.Reader) []string {

	scanner := bufio.NewScanner(reader)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "^D" {
			break
		}
		lines = append(lines, line)
	}
	matchedLines := g.CheckMatch(lines, pattern)
	return matchedLines
}

func (g *Grep) WriteOutput(matchedLines []string) {

	if g.OutputFile != "" {
		g.writeOutputToFile(matchedLines)
		return
	}

	for _, line := range matchedLines {
		fmt.Println(line)
	}
}

func (g *Grep) writeOutputToFile(matchedLines []string) {
	outputFileName := g.OutputFile
	f, err := os.Create(outputFileName)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer f.Close()
	for _, line := range matchedLines {
		stringToWrite := line + "\n"

		_, err2 := f.WriteString(stringToWrite)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
}
