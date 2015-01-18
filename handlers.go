package main

import (
	"strconv"
	"log"
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

func registerHandlers(){
	r := mux.NewRouter()
    r.HandleFunc("/", homeHandler)
    r.HandleFunc("/status", status).Methods("GET")
    r.HandleFunc("/connect", connectUser).Methods("POST")
    r.HandleFunc("/getUserInfos", updateStats).Methods("POST")
    r.HandleFunc("/updateAchievement", updateAchievement).Methods("POST")
    r.HandleFunc("/updateStats", updateStats).Methods("POST")
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

type easyResponse struct{
	Code int
	Message string
}
func EasyResponse(code int, message string) string{
	e := easyResponse{code, message}
	res, err := json.Marshal(e)
	if err != nil {
		log.Println("Could not encode easy response")
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
	r.ParseForm()
	if(r.PostForm["username"] == nil || r.PostForm["password"] == nil){
		fmt.Fprint(w, EasyResponse(NebError, "Missing username and/or password"))
		return
	}

	username := r.PostForm["username"][0] // For some reason r.PostForm[i] is String[]
	password := r.PostForm["password"][0]
	id, err := CreateSession(username, password)
	
	if err != nil {
		fmt.Fprint(w, EasyResponse(NebErrorLogin, err.Error()))
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
func updateAchievement(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if(r.PostForm["sessionid"] == nil || r.PostForm["achievementid"] == nil || r.PostForm["value"] == nil){
		fmt.Fprint(w, EasyResponse(NebError, "Missing sessionid, achievementid or value"))
		return
	}

	user,err := GetUser(r.PostForm["sessionid"][0])
	if(err != nil){
		fmt.Fprint(w, EasyResponse(NebErrorDisconnected, err.Error()))
		return
	}

	aid, err := strconv.Atoi(r.PostForm["achievementid"][0])
	if(err != nil){
		fmt.Fprint(w, EasyResponse(NebError, err.Error()))
		return
	}

	value, err := strconv.Atoi(r.PostForm["value"][0])
	if(err != nil){
		fmt.Fprint(w, EasyResponse(NebError, err.Error()))
		return
	}

	err = user.UpdateAchievementProgress(aid, value)
	if err != nil{
		fmt.Fprint(w, EasyResponse(NebError, err.Error()))
		return
	}
	
	fmt.Fprint(w, EasyResponse(NebErrorNone, "Updated Achievement"))

	go user.Heartbeat()
}
type updateStatsRequest struct {
	Map string
	Stats []Stat
	Kills []Kill
}
func updateStats(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if(r.PostForm["sessionid"] == nil || r.PostForm["data"] == nil ){
		fmt.Fprint(w, EasyResponse(NebError, "Missing sessionid or data"))
		return
	}

	user,err := GetUser(r.PostForm["sessionid"][0])
	if(err != nil){
		fmt.Fprint(w, EasyResponse(NebErrorDisconnected, err.Error()))
		return
	}

	data := r.PostForm["data"][0]
	var req updateStatsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil{
		fmt.Fprint(w, EasyResponse(NebError, err.Error()))
		return
	}

	user.UpdateStats(req.Stats)
	user.InsertKills(req.Kills, req.Map)
	
	fmt.Fprint(w, EasyResponse(NebErrorNone, "Updated Stats"))
	
	go user.Heartbeat()
}

func getUserInfos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if(r.PostForm["sessionid"] == nil || r.PostForm["data"] == nil ){
		fmt.Fprint(w, EasyResponse(NebError, "Missing sessionid or data"))
		return
	}

	user,err := GetUser(r.PostForm["sessionid"][0])
	if(err != nil){
		fmt.Fprint(w, EasyResponse(NebErrorDisconnected, err.Error()))
		return
	}

	res, err := json.Marshal(user)
	fmt.Fprint(w, string(res))

	go user.Heartbeat()
}