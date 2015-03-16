package main

import (
	"Nebuleuse/core"
	"Nebuleuse/gitUpdater"
)

func main() {
	core.Init()

	defer core.Die()

	gitUpdater.Init(".")

	CreateServer()
}
