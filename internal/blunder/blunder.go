package blunder

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/sashabaranov/go-openai"
)

const (
	Themes StoryDataType = iota
	Styles
	Modifiers
)

type StoryOptions struct {
	Theme    string
	Style    string
	Modifier string
}

type StoryResult struct {
	Prompt string
	Story  string
}

type StoryDataType int

type StoryDataProvider interface {
	GetRandomString(dataType StoryDataType) (string, error)
}

func generatePrompt(dataProvider StoryDataProvider, options StoryOptions) (string, error) {
	const promptFormatString = "Describe to me a highly comical situation stemming from a misunderstanding. " +
		"The theme should be '%v'%v. Write the description in the style of %v and limit the length to 500 characters."

	var err error = nil

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

func GenerateStory(openAIToken string, dataProvider StoryDataProvider, options StoryOptions) (StoryResult, error) {
	// Generate query
	prompt, err := generatePrompt(dataProvider, options)
	if err != nil {
		return StoryResult{}, err
	}

	// Generate story
	story, err := queryLLM(openAIToken, prompt)
	if err != nil {
		return StoryResult{}, err
	}

	return StoryResult{
		Prompt: prompt,
		Story:  story,
	}, nil
}
