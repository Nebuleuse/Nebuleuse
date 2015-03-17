package main

import (
	"Nebuleuse/core"
)

func main() {
	core.Init()

	defer core.Die()

	CreateServer()
}
