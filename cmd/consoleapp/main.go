package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chippolot/jokegen"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "jokegen",
		Usage: "generates a comical story",
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
				Name:     "storyType",
				Aliases:  []string{"st"},
				Value:    "",
				Usage:    "Story type",
				Required: true,
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
				Aliases:  []string{"p"},
				Value:    false,
				Usage:    "If true, shows the generated prompt alongside the story",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "forceRegenerate",
				Aliases:  []string{"f"},
				Value:    false,
				Usage:    "If true, the story will always be regerated, even if a valid cached story exists",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			token := ctx.String("token")
			showPrompt := ctx.Bool("showPrompt")
			options := jokegen.StoryOptions{
				Theme:           ctx.String("theme"),
				Style:           ctx.String("style"),
				Modifier:        ctx.String("modifier"),
				ForceRegenerate: ctx.Bool("forceRegenerate"),
			}

			storyTypeString := ctx.String("storyType")
			storyType, err := jokegen.ParseStoryType(storyTypeString)
			if err != nil {
				return err
			}

			dataProvider := &FileDataProvider{}
			defer dataProvider.Close()

			result, err := jokegen.GenerateStory(token, storyType, dataProvider, options)
			if err != nil {
				return err
			}
			if showPrompt {
				fmt.Println("Prompt:")
				fmt.Println(result.Prompt)
				fmt.Println()
			}
			fmt.Println(result.Story)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
