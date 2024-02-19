package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/chippolot/blunders/internal/blunder"
)

const (
	NounsFilePath     = "res/nouns.txt"
	StylesFilePath    = "res/styles.txt"
	ModifiersFilePath = "res/modifiers.txt"
)

type FileDataProvider struct {
}

func (f *FileDataProvider) GetRandomString(dataType blunder.StoryDataType) (string, error) {
	filePath, err := getFilePath(dataType)
	if err != nil {
		return "", err
	}
	lines, err := readLines(filePath)
	if err != nil {
		return "", err
	}
	randomIndex := rand.Intn(len(lines))
	return lines[randomIndex], nil
}

func getFilePath(dataType blunder.StoryDataType) (string, error) {
	switch dataType {
	case blunder.Themes:
		return NounsFilePath, nil
	case blunder.Styles:
		return StylesFilePath, nil
	case blunder.Modifiers:
		return ModifiersFilePath, nil
	}
	return "", fmt.Errorf("unknown data type: %v", dataType)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
