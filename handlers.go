package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func registerHandlers() {
	r := mux.NewRouter()
	r.HandleFunc("/", status)
	r.HandleFunc("/status", status).Methods("GET")
	r.HandleFunc("/connect", connectUser).Methods("POST")
	r.HandleFunc("/getUserInfos", getUserInfos).Methods("POST")
	r.HandleFunc("/updateAchievements", updateAchievements).Methods("POST")
	r.HandleFunc("/updateStats", updateStats).Methods("POST")
	r.HandleFunc("/addComplexStats", addComplexStats).Methods("POST")
	http.Handle("/", r)
}

type easyResponse struct {
	Code    int
	Message string
}

func EasyResponse(code int, message string) string {
	e := easyResponse{code, message}
	res, err := json.Marshal(e)
	if err != nil {
		log.Println("Could not encode easy response")
	}

	return string(res)
}
func EasyErrorResponse(code int, err error) string {
	v, ok := err.(NebuleuseError)
	var e easyResponse
	if ok {
		e = easyResponse{v.code, v.msg}
	} else {
		e = easyResponse{code, err.Error()}
	}
	res, err := json.Marshal(e)
	if err != nil {
		log.Println("Could not encode easy response")
	}

	return string(res)
}

type statusResponse struct {
	Maintenance      bool
	NebuleuseVersion int
	GameVersion      int
	UpdaterVersion   int
	Motd             string
}

func status(w http.ResponseWriter, r *http.Request) {
	response := statusResponse{false, nebuleuseVersion, getGameVersion(), getUpdaterVersion(), "abc"}

	if _cfg["maintenance"] != "0" {
		response.Maintenance = true
	}

	res, err := json.Marshal(response)
	if err != nil {
		log.Println("Could not encode status response")
		fmt.Fprint(w, EasyErrorResponse(NebError, err))
	} else {
		fmt.Fprint(w, string(res))
	}
}

type connectResponse struct {
	SessionId string
}

func connectUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["username"] == nil || r.PostForm["password"] == nil {
		fmt.Fprint(w, EasyResponse(NebError, "Missing username and/or password"))
		return
	}

	username := r.PostForm["username"][0] // For some reason r.PostForm[i] is String[]
	password := r.PostForm["password"][0]
	id, err := CreateSession(username, password)

	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebErrorLogin, err))
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

func getUserInfos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil {
		fmt.Fprint(w, EasyResponse(NebError, "Missing sessionid"))
		return
	}

	user, err := GetUser(r.PostForm["sessionid"][0])
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebErrorDisconnected, err))
		return
	}

	res, err := json.Marshal(user)
	fmt.Fprint(w, string(res))
}

type achievementRequest struct {
	Id    int
	Value int
}
type updateAchievementsRequest struct {
	Achievements []achievementRequest
}

func updateAchievements(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil || r.PostForm["data"] == nil {
		fmt.Fprint(w, EasyResponse(NebError, "Missing sessionid or data"))
		return
	}

	user, err := GetUser(r.PostForm["sessionid"][0])
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebErrorDisconnected, err))
		return
	}

	data := r.PostForm["data"][0]
	var req updateAchievementsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebError, err))
		return
	}

	count := len(req.Achievements)
	for _, ach := range req.Achievements {
		err = user.UpdateAchievementProgress(ach.Id, ach.Value)
		if err != nil {
			count--
			continue
		}
	}
	if count != len(req.Achievements) {
		fmt.Fprint(w, EasyResponse(NebErrorPartialFail, "Updated "+string(count)+" Achievements"))
		return
	}

	fmt.Fprint(w, EasyResponse(NebErrorNone, "Updated Achievements"))
}

type updateStatsRequest struct {
	Stats []UserStat
}

func updateStats(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil || r.PostForm["data"] == nil {
		fmt.Fprint(w, EasyResponse(NebError, "Missing sessionid or data"))
		return
	}

	user, err := GetUser(r.PostForm["sessionid"][0])
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebErrorDisconnected, err))
		return
	}

	data := r.PostForm["data"][0]
	var req updateStatsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebError, err))
		return
	}

	user.UpdateStats(req.Stats)

	fmt.Fprint(w, EasyResponse(NebErrorNone, "Updated Stats"))
}

type updateComplexStatsRequest struct {
	Stats []ComplexStat
}

func addComplexStats(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil || r.PostForm["data"] == nil {
		fmt.Fprint(w, EasyResponse(NebError, "Missing sessionid or data"))
		return
	}

	user, err := GetUser(r.PostForm["sessionid"][0])
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebErrorDisconnected, err))
		return
	}

	data := r.PostForm["data"][0]
	var req updateComplexStatsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebError, err))
		return
	}

	err = user.updateComplexStats(req.Stats)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebError, err))
		return
	}

	fmt.Fprint(w, EasyResponse(NebErrorNone, "Inserted complex stat"))
}
