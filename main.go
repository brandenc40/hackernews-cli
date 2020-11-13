package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brandenc40/hackernews-cli/cli"
)

func main() {
	app := cli.HNCliApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
}
