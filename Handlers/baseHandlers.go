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
	r.HandleFunc("/disconnect", disconnectUser).Methods("POST")
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

type ComplexStatTableInfo struct {
	Name string
	Fields []string
	AutoCount bool
}
type statusResponse struct {
	Maintenance      bool
	NebuleuseVersion int
	GameVersion      int
	UpdaterVersion   int
	ComplexStatTable []ComplexStatTableInfo
}

func status(w http.ResponseWriter, r *http.Request) {
	var CStatsInfos, err := getComplexStatsTablesInfos()
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebError, err))
	}
	response := statusResponse{false, nebuleuseVersion, getGameVersion(), getUpdaterVersion(), CStatsInfos}

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