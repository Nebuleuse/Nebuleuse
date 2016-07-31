package main

import (
	"github.com/Nebuleuse/Nebuleuse/core"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "help" {
			log.Println("Usage: Nebuleuse [action]\nActions: help: Prints this\n\t install: execute installation script")
			return
		} else if os.Args[1] == "install" {
			core.Install()
			createInstallServer()
		}
	}

	core.Init()

	defer core.Die()

	createServer()
}
