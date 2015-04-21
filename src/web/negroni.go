package main

import (
	"fmt"
	"os"
	"web/controllers/frontend"
	"web/helpers/mynegroni"
)

func init() {
	env := os.Getenv("ENV")

	if env == "" {
		panic("ENV not set")
	}

	fmt.Printf("[negroni] %s\n", env)
}

func main() {

	// Negroni
	n := mynegroni.New()

	// Frontend
	n.Use(frontend.Restricted)
	n.Use(frontend.Controller)
	n.Use(frontend.NotFound)

	// Run negroni run!
	n.Run(":3000")
}
