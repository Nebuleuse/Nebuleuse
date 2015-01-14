package main

import (
	"log"
	"net/http"
)

const (  
	NebErrorNone = iota
	NebError = iota
	NebErrorDisconnected = iota
	NebErrorLogin = iota
)

func createServer(){
	registerHandlers()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
