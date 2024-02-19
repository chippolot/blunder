package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chippolot/blunders/internal/blunder"
	"github.com/urfave/cli/v2"
)

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
			options := blunder.StoryOptions{
				Theme:           ctx.String("theme"),
				Style:           ctx.String("style"),
				Modifier:        ctx.String("modifier"),
				ForceRegenerate: ctx.Bool("forceRegenerate"),
			}
			dataProvider := &FileDataProvider{}
			result, err := blunder.GenerateStory(token, dataProvider, options)
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
