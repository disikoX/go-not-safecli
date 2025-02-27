package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "go-safecli",
		Usage: " simple CLI tool for securely storing and managing passwords and email credentials ",
		Action: func(*cli.Context) error {
			fmt.Println("Hello, this is my first cli-app")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal((err))
	}
}
