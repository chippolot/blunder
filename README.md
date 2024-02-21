# JokeGen
Simple Go library and console application which use the OpenAI API to randomly generate different types of comical short stories.

## Overview

JokeGen can be used in a few different ways:
* The `github.com/chippolot/jokegen` package can be used in other Go projects for story generation.
* The `jokegen` console application can be used to test story generation directly.

JokeGen required an active OpenAI API key in order to generate new stories.

## Story Types

JokeGen supports a few different story types:
* Misunderstandings (`storyType: misunderstanding`)
  * Comical stories stemming from misunderstandings.
* Slapsticks (`storyType: slapstick`)
  * Stories involving slapstick humor.
* Curses (`storyType: curse`)
  * Comical stories revolving around curses.
* Curses (`storyType: creature`)
  * Comical stories revolving around mythical creatures.

## Installing Package

To include JokeGen in another Go application just use `go get`:
```
go get -u github.com/chippolot/jokegen@latest
```

Next, import JokeGen in your application:
```go
import "github.com/chippolot/jokegen"
```

You will need to create a new data provider implementing the following interface:
```go
type StoryDataProvider interface {
	AddStory(story string, prompt string, storyType StoryType) error
	GetMostRecentStory(storyType StoryType) (StoryResult, error)
	GetRandomString(dataType StoryDataType, storyType StoryType) (string, error)
	Close() error
}
```

Finally, generate new stories:
```go
openAIToken := "YOUR_OPEN_AI_TOKEN"
storyType := "misunderstanding"
dataProvider := YOUR_DATA_PROVIDER_INSTANCE
options := jokegen.StoryOptions{}
jokegen.GenerateStory(openAIToken, storyType, dataProvider, options)
```

## Executing CLI program

Build the CLI program with:
```
go build -o jokegen ./cmd/consoleapp
```

And run with:
```
./consoleapp -st STORY_TYPE
```

## License

This project is licensed under the MIT License - see the LICENSE.md file for details
