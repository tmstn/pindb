package main

import (
	"fmt"
	"os"
)

func main() {
	app := createApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("\n%s", err.Error())
	}
}
