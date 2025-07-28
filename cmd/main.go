package main

import (
	"log"

	"github.com/konstantin-suspitsyn/datacomrade/cmd/api"
)

const version = "1.0.0"

func main() {
	err := api.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
