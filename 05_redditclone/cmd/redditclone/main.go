package main

import (
	"log"
	"redditclone/internal/app"
)

const configPath = "./configs/config.yml"

func main() {
	if err := app.Run(configPath); err != nil {
		log.Fatalln(err)
	}
}
