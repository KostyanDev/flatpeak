package main

import (
	"log"
	"os"

	"app/internal/app"

	_ "app/docs"
	_ "app/internal/domain"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
