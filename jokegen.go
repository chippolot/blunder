package jokegen

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/sashabaranov/go-openai"
)

const (
	// Story data types
	Themes StoryDataType = iota
	Styles
	Modifiers

	// Story types
	Misunderstanding StoryType = iota
	Slapstick
)

type StoryOptions struct {
	Theme           string
	Style           string
	Modifier        string
	ForceRegenerate bool
}

type StoryResult struct {
	Timestamp time.Time
	Prompt    string
	Story     string
}

type StoryDataType int
type StoryType int

func ParseStoryType(str string) (StoryType, error) {
	switch str {
	case "misunderstanding":
		return Misunderstanding, nil
	case "slapstick":
		return Slapstick, nil
	}
	return -1, fmt.Errorf("unknown story type: %v", str)
}

type StoryDataProvider interface {
	AddStory(story string, prompt string, storyType StoryType) error
	GetMostRecentStory(storyType StoryType) (StoryResult, error)
	GetRandomString(dataType StoryDataType) (string, error)
	Close() error
}

func getPrompt(storyType StoryType) (string, error) {
	const postfix string = "The theme should be '%v'%v. Write the description in the style of %v and limit the length to 500 characters."

	switch storyType {
	case Misunderstanding:
		return "Describe to me a highly comical situation stemming from a misunderstanding. " + postfix, nil
	case Slapstick:
		return "Describe to me a highly comical situation revolving around slapstick humor, using florid language to describe the action. " + postfix, nil
	}
	return "", fmt.Errorf("unknown story type %v", storyType)
}

func generatePrompt(storyType StoryType, dataProvider StoryDataProvider, options StoryOptions) (string, error) {
	promptFormatString, err := getPrompt(storyType)
	if err != nil {
		return "", err
	}

	// Get a random theme
	theme := options.Theme
	if theme == "" {
		theme, err = dataProvider.GetRandomString(Themes)
		if err != nil {
			return "", err
		}
	}

	// Get a random style
	style := options.Style
	if style == "" {
		style, err = dataProvider.GetRandomString(Styles)
		if err != nil {
			return "", err
		}
	}

	// Get a random content modifier
	modifier := options.Modifier
	if modifier == "" && rand.Float32() > 0.5 {
		modifier, err = dataProvider.GetRandomString(Modifiers)
		if err != nil {
			return "", err
		}
	}
	if modifier != "" {
		modifier = " " + modifier
	}

	// Build and output query
	return fmt.Sprintf(promptFormatString, theme, modifier, style), nil
}

func queryLLM(token string, prompt string) (string, error) {
	client := openai.NewClient(token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func GenerateStory(openAIToken string, storyType StoryType, dataProvider StoryDataProvider, options StoryOptions) (StoryResult, error) {
	// Check for cached story
	if !options.ForceRegenerate {
		now := time.Now().UTC()
		cached, err := dataProvider.GetMostRecentStory(storyType)
		if err == nil {
			cacheDuration := now.Sub(cached.Timestamp)
			if cacheDuration < time.Hour*24 {
				return cached, nil
			}
		}
	}

	// Generate query
	prompt, err := generatePrompt(storyType, dataProvider, options)
	if err != nil {
		return StoryResult{}, err
	}

	// Generate story
	story, err := queryLLM(openAIToken, prompt)
	if err != nil {
		return StoryResult{}, err
	}

	// Cache story
	err = dataProvider.AddStory(story, prompt, storyType)
	if err != nil {
		return StoryResult{}, err
	}

	return StoryResult{
		Prompt:    prompt,
		Story:     story,
		Timestamp: time.Now().UTC(),
	}, nil
}
