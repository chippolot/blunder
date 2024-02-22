package jokegen

import (
	"context"
	"embed"
	"fmt"
	"math/rand"
	"time"

	utils "github.com/chippolot/jokegen/internal"
	"github.com/sashabaranov/go-openai"
)

//go:embed res/nouns.txt
//go:embed res/misunderstanding/styles.txt
//go:embed res/misunderstanding/modifiers.txt
//go:embed res/slapstick/styles.txt
//go:embed res/slapstick/modifiers.txt
var ResourcesFS embed.FS

type StoryDataType int

const (
	// Story data types
	Themes StoryDataType = iota
	Styles
	Modifiers
)

type StoryType int

const (
	// Story types
	Misunderstanding StoryType = 1 << iota
	Slapstick
	Curse
	Creature
	AntiHumor
)

var storyTypeStringMapping = utils.NewStringMapping[StoryType](map[StoryType]string{
	Misunderstanding: "misunderstanding",
	Slapstick:        "slapstick",
	Curse:            "curse",
	Creature:         "creature",
	AntiHumor:        "antihumor",
})

func ParseStoryType(str string) (StoryType, error) {
	if val, ok := storyTypeStringMapping.ToValue[str]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("unknown story type: %s", str)
}

func (s StoryType) ToString() (string, error) {
	if str, ok := storyTypeStringMapping.ToString[s]; ok {
		return str, nil
	}
	return "", fmt.Errorf("unknown story type: %v", s)
}

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

type StoryDataProvider interface {
	AddStory(story string, prompt string, storyType StoryType) error
	GetMostRecentStory(storyType StoryType) (StoryResult, error)
	GetRandomString(dataType StoryDataType, storyType StoryType) (string, error)
	Close() error
}

func getPrompt(storyType StoryType) (string, error) {
	const prefix string = "Describe to me a highly comical situation "
	const postfix string = " The theme should be '%v'%v. Write the description in the style of %v and limit the length to 500 characters."

	switch storyType {
	case Misunderstanding:
		return prefix + "stemming from a misunderstanding." + postfix, nil
	case Slapstick:
		return prefix + "revolving around slapstick humor, using florid language to describe the action." + postfix, nil
	case Curse:
		return prefix + "revolving around a curse." + postfix, nil
	case Creature:
		return prefix + "revolving around a newly created mythical creature." + postfix, nil
	case AntiHumor:
		return "Describe a story using antihumor. Nothing funny should happen and the story should neither acknowledge that it is not funny nor that there was the expectation of humor." + postfix, nil
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
		theme, err = dataProvider.GetRandomString(Themes, storyType)
		if err != nil {
			return "", err
		}
	}

	// Get a random style
	style := options.Style
	if style == "" {
		style, err = dataProvider.GetRandomString(Styles, storyType)
		if err != nil {
			return "", err
		}
	}

	// Get a random content modifier
	modifier := options.Modifier
	if modifier == "" && rand.Float32() > 0.5 {
		modifier, err = dataProvider.GetRandomString(Modifiers, storyType)
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
