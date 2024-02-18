package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/sashabaranov/go-openai"
)

const (
	NounsFilePath            = "res/nouns.txt"
	StyleModifiersFilePath   = "res/style_modifiers.txt"
	ContentModifiersFilePath = "res/content_modifiers.txt"
)

func getEnvVariable(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", errors.New("environment variable not found: " + key)
	}
	return value, nil
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

func generateQuery() string {
	const promptFormatString = "Describe to me a highly comical situation stemming from a misunderstanding. " +
		"The theme should be '%v'%v. Write the description in the style of %v and limit the length to 500 characters."

	// Get a random word
	word := randomLinesFromFile(NounsFilePath, 1)[0]

	// Get some random modifiers
	styleModifier := randomLinesFromFile(StyleModifiersFilePath, 1)[0]
	contentModifier := ""
	if rand.Float32() > 0.5 {
		contentModifier = randomLinesFromFile(ContentModifiersFilePath, 1)[0]
	}

	// Build and output query
	return fmt.Sprintf(promptFormatString, word, contentModifier, styleModifier)
}

func queryLLM(token string, query string) (string, error) {
	client := openai.NewClient(token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func main() {

	// Generate query
	query := generateQuery()

	// Print the query
	// fmt.Println("Query:")
	// fmt.Println(query)

	// Get secret key
	openAIToken, err := getEnvVariable("OPEN_AI_API_KEY")
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// Generate result
	result, err := queryLLM(openAIToken, query)
	if err != nil {
		fmt.Printf("LLM generation error: %v\n", err)
		return
	}

	// Print the result
	// fmt.Println()
	// fmt.Println("Story:")
	fmt.Println(result)
}
