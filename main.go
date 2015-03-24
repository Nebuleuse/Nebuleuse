package main

import (
	"github.com/Nebuleuse/Nebuleuse/core"
)

func main() {
	core.Init()

	defer core.Die()

	CreateServer()
}
