package main

import (
	"log"
	"moe.two.bgmi-gls/cmd/app"
)

func main() {
	if err := app.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}