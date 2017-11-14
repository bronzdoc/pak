package main

import (
	"github.com/bronzdoc/pak/api"
)

func main() {
	pakfile := api.Parse("pakfile.json")
	api.Build(&pakfile)
}
