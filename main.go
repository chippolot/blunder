package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
)

const (
	NounsFilePath            = "res/nouns.txt"
	StyleModifiersFilePath   = "res/style_modifiers.txt"
	ContentModifiersFilePath = "res/content_modifiers.txt"
)

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

func randomLinesFromFile(filePath string, num int) []string {
	lines, err := readLines(filePath)
	if err != nil {
		panic(err)
	}

	rand.Shuffle(len(lines), func(i, j int) {
		lines[i], lines[j] = lines[j], lines[i]
	})

	return lines[:num]
}

func generatePrompt(theme string, style string, modifier string) string {
	const promptFormatString = "Describe to me a highly comical situation stemming from a misunderstanding. " +
		"The theme should be '%v'%v. Write the description in the style of %v and limit the length to 500 characters."

	// Get a random theme
	if theme == "" {
		theme = randomLinesFromFile(NounsFilePath, 1)[0]
	}

	// Get a random style
	if style == "" {
		style = randomLinesFromFile(StyleModifiersFilePath, 1)[0]
	}

	// Get a random content modifier
	if modifier == "" && rand.Float32() > 0.5 {
		modifier = randomLinesFromFile(ContentModifiersFilePath, 1)[0]
	}

	// Build and output query
	return fmt.Sprintf(promptFormatString, theme, modifier, style)
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

func generateStory(openAIToken string, theme string, style string, modifier string, showPrompt bool) (string, error) {

	// Generate query
	prompt := generatePrompt(theme, style, modifier)
	if showPrompt {
		fmt.Println("Prompt:")
		fmt.Println(prompt)
		fmt.Println()
	}

	// Generate result
	result, err := queryLLM(openAIToken, prompt)
	if err != nil {
		return "", err
	}

	return result, nil
}

func main() {

	app := &cli.App{
		Name:  "blunderbuddy",
		Usage: "generate a comical story of a misundertanding!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Value:    "",
				Usage:    "OpenAI token",
				Required: false,
				EnvVars:  []string{"OPEN_AI_API_KEY"},
			},
			&cli.StringFlag{
				Name:     "theme",
				Value:    "",
				Usage:    "Theme for story",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "style",
				Value:    "",
				Usage:    "Style for story",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "modifier",
				Value:    "",
				Usage:    "Modifier for story",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "showPrompt",
				Value:    false,
				Usage:    "If true, shows the generated prompt alongside the story",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			token := ctx.String("token")
			theme := ctx.String("theme")
			style := ctx.String("style")
			modifier := ctx.String("modifier")
			showPrompt := ctx.Bool("showPrompt")
			story, err := generateStory(token, theme, style, modifier, showPrompt)
			if err != nil {
				return err
			}
			fmt.Println(story)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
