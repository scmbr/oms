package main

import (
	"log"

	"github.com/scmbr/oms/user-service/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
