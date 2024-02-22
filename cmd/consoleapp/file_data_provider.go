package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"slices"

	"github.com/chippolot/jokegen"
)

const (
	ThemesFileName    = "nouns.txt"
	StylesFileName    = "styles.txt"
	ModifiersFileName = "modifiers.txt"
)

type FileDataProvider struct {
}

func (f *FileDataProvider) AddStory(story, prompt string, storyType jokegen.StoryType) error {
	// No support for story caching
	return nil
}

func (f *FileDataProvider) GetMostRecentStory(storyType jokegen.StoryType) (jokegen.StoryResult, error) {
	// No support for story caching
	return jokegen.StoryResult{}, fmt.Errorf("no recent story available")
}

func (f *FileDataProvider) GetRandomString(dataType jokegen.StoryDataType, storyType jokegen.StoryType) (string, error) {
	filename, err := getFilename(dataType)
	if err != nil {
		return "", nil
	}
	storyTypeStr, _ := storyType.ToString()

	lines := slices.Concat(readLines(fmt.Sprintf("res/%s", filename)), readLines(fmt.Sprintf("res/%s/%s", storyTypeStr, filename)))
	if len(lines) == 0 {
		return "", nil
	}
	randomIndex := rand.Intn(len(lines))
	return lines[randomIndex], nil
}

func (f *FileDataProvider) Close() error {
	return nil
}

func getFilename(dataType jokegen.StoryDataType) (string, error) {
	switch dataType {
	case jokegen.Themes:
		return ThemesFileName, nil
	case jokegen.Styles:
		return StylesFileName, nil
	case jokegen.Modifiers:
		return ModifiersFileName, nil
	}
	return "", fmt.Errorf("unknown data type: %v", dataType)
}

func readLines(path string) []string {
	file, err := jokegen.ResourcesFS.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
