package main

import (
	"log"

	"github.com/gvx3/sportuni-book/pkg/app"
)

func main() {
	if err := app.RunApp(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
