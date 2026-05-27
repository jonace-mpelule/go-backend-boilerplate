package main

import (
	"log"

	"github.com/username/project-name/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal("error starting app:", err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
