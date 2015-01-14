package main

import (
	"log"
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

func registerHandlers(){
	r := mux.NewRouter()
    r.HandleFunc("/", homeHandler)
    r.HandleFunc("/status", status)
    r.HandleFunc("/connect", connectUser)
    http.Handle("/", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request){
	user, err := GetUser("64dcb8cf-0820-4714-4bc2-e885566d54f9")
	if(err != nil){
		log.Println(err)
		return
	}
	res, err := json.Marshal(user)
	if err != nil {
		log.Println("Could not encode error response")
	} 
	fmt.Fprint(w, string(res))
}

type errorResponse struct{
	ErrorCode int
	ErrorMessage string
}
func prepareErrorResponse(code int, message string) string{
	e := errorResponse{code, message}
	res, err := json.Marshal(e)
	if err != nil {
		log.Println("Could not encode error response")
	} 

	return string(res)
}

type statusResponse struct{
	Maintenance bool
	NebuleuseVersion int
	GameVersion int
	UpdaterVersion int
	Motd string
}
func status(w http.ResponseWriter, r *http.Request){
	response := statusResponse{false, nebuleuseVersion, getGameVersion(), getUpdaterVersion(), "abc"}

	if(_cfg["maintenance"] != "0"){
		response.Maintenance = true
	}

	res, err := json.Marshal(response)
	if err != nil {
		log.Println("Could not encode status response")
	} else {
		fmt.Fprint(w, string(res))
	}
}

type connectRequest struct{
	Username string
	Password string
}
type connectResponse struct{
	SessionId string
}
func connectUser(w http.ResponseWriter, r *http.Request){
	id, err := CreateSession("test", "test")
	if err != nil && err.Error() == "" {
		fmt.Fprint(w, prepareErrorResponse(NebErrorLogin, "Wrong login information"))
		return
	} else if err != nil {
		fmt.Fprint(w, prepareErrorResponse(NebError, err.Error()))
		return
	}

	response := connectResponse{id}

	res, err := json.Marshal(response)
	if err != nil {
		log.Println("Could not encode status response")
	} else {
		fmt.Fprint(w, string(res))
	}
}