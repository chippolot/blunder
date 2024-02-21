package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/chippolot/jokegen"
)

const (
	ThemesFilePath      = "res/nouns.txt"
	StylesFilePath      = "res/%s/styles.txt"
	ModifiersFilePath   = "res/%s/modifiers.txt"
	CachedStoryFilePath = "recent_story.json"
)

type FileDataProvider struct {
}

func (f *FileDataProvider) AddStory(story, prompt string, storyType jokegen.StoryType) error {
	result := &jokegen.StoryResult{
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

func (f *FileDataProvider) GetMostRecentStory(storyType jokegen.StoryType) (jokegen.StoryResult, error) {
	file, err := os.Open(CachedStoryFilePath)
	if err != nil {
		return jokegen.StoryResult{}, err
	}
	defer file.Close()

	var result jokegen.StoryResult
	if err := json.NewDecoder(file).Decode(&result); err != nil {
		return jokegen.StoryResult{}, err
	}

	return result, nil
}

func (f *FileDataProvider) GetRandomString(dataType jokegen.StoryDataType, storyType jokegen.StoryType) (string, error) {
	filePath, err := getFilePath(dataType, storyType)
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

func getFilePath(dataType jokegen.StoryDataType, storyType jokegen.StoryType) (string, error) {
	storyTypeString, err := storyType.ToString()
	if err != nil {
		return "", err
	}

	switch dataType {
	case jokegen.Themes:
		return ThemesFilePath, nil
	case jokegen.Styles:
		return fmt.Sprintf(StylesFilePath, storyTypeString), nil
	case jokegen.Modifiers:
		return fmt.Sprintf(ModifiersFilePath, storyTypeString), nil
	}
	return "", fmt.Errorf("unknown data type: %v", dataType)
}

func readLines(path string) ([]string, error) {
	file, err := jokegen.ResourcesFS.Open(path)
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
