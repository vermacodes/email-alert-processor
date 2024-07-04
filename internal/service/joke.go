package service

import (
	"log/slog"

	"github.com/icelain/jokeapi"
)

func GetJoke() string {
	jokeType := "single"
	blacklist := []string{"nsfw", "religious", "political", "racist", "sexist", "explicit"}
	categories := []string{"Programming"}

	api := jokeapi.New()

	api.Set(jokeapi.Params{Blacklist: blacklist, JokeType: jokeType, Categories: categories})

	response, err := api.Fetch()
	if err != nil {
		slog.Error("Error fetching joke: ", err)
		return "Oops! No jokes available right now."
	}

	// Return the first (only one) in the list.
	return response.Joke[0]
}
