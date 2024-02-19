package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/chippolot/blunders/internal/blunder"
)

const (
	NounsFilePath       = "res/nouns.txt"
	StylesFilePath      = "res/styles.txt"
	ModifiersFilePath   = "res/modifiers.txt"
	CachedStoryFilePath = "recent_story.json"
)

type FileDataProvider struct {
}

func (f *FileDataProvider) AddStory(story string, prompt string) error {
	result := &blunder.StoryResult{
		Story:     story,
		Prompt:    prompt,
		Timestamp: time.Now().UTC(),
	}

	file, err := os.Create(CachedStoryFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(result)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileDataProvider) GetMostRecentStory() (blunder.StoryResult, error) {
	file, err := os.Open(CachedStoryFilePath)
	if err != nil {
		return blunder.StoryResult{}, err
	}
	defer file.Close()

	var result blunder.StoryResult
	if err := json.NewDecoder(file).Decode(&result); err != nil {
		return blunder.StoryResult{}, err
	}

	return result, nil
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

func (f *FileDataProvider) Close() error {
	return nil
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
