package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chippolot/blunders/internal/blunder"
)

var dataProvider *SQLiteDataProvider

func storyHandler(w http.ResponseWriter, r *http.Request) {
	// Resolve API key
	token := os.Getenv("OPEN_AI_API_KEY")
	if token == "" {
		log.Fatal("OpenAI API key not found in environment variables")
	}

	// Generate Story
	options := blunder.StoryOptions{}
	result, err := blunder.GenerateStory(token, dataProvider, options)
	if err != nil {
		panic(err)
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		panic(err)
	}
}

func main() {
	dataProvider = MakeSQLiteDataProvider()
	defer dataProvider.Close()

	http.HandleFunc("/story", storyHandler)

	port := 8080
	fmt.Printf("Server is running on http://localhost:%v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
