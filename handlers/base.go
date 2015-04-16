package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHandlers() {
	r := mux.NewRouter()

	r.HandleFunc("/", status)
	r.HandleFunc("/status", status).Methods("GET")
	r.HandleFunc("/connect", connectUser).Methods("POST")
	r.HandleFunc("/disconnect", userBySession(false, disconnectUser)).Methods("POST")
	r.HandleFunc("/getUserInfos", userBySession(true, verifyFormValuesExist([]string{"infomask"}, getUserInfos))).Methods("POST")

	r.HandleFunc("/seAchievements", userBySession(false, verifyFormValuesExist([]string{"data"}, setAchievements))).Methods("POST")
	r.HandleFunc("/seStats", userBySession(false, verifyFormValuesExist([]string{"data"}, setStats))).Methods("POST")
	r.HandleFunc("/addComplexStats", userBySession(false, verifyFormValuesExist([]string{"data"}, addComplexStats))).Methods("POST")

	r.HandleFunc("/fetchMessage", userBySession(false, fetchMessage)).Methods("POST")
	r.HandleFunc("/sendMessage", userBySession(false, verifyFormValuesExist([]string{"channel", "message"}, sendMessage))).Methods("POST")
	r.HandleFunc("/subscribeTo", userBySession(false, verifyFormValuesExist([]string{"channel"}, subscribeTo))).Methods("POST")
	r.HandleFunc("/unSubscribeTo", userBySession(false, verifyFormValuesExist([]string{"channel"}, unSubscribeTo))).Methods("POST")

	r.PathPrefix("/admin/").Handler((http.StripPrefix("/admin/", http.FileServer(http.Dir("./admin/dist/")))))
	r.HandleFunc("/getDashboardInfos", userBySession(false, mustBeAdmin(getDashboardInfos))).Methods("POST")
	http.Handle("/", r)
}

type easyResponse struct {
	Code    int
	Message string
}

func EasyResponse(w http.ResponseWriter, code int, message string) {
	e := easyResponse{code, message}
	res, err := json.Marshal(e)
	if err != nil {
		core.Warning.Println("Could not encode easy response")
	}

	switch code {
	case core.NebErrorLogin:
		w.WriteHeader(http.StatusUnauthorized)
	case core.NebErrorAuthFail:
		w.WriteHeader(http.StatusUnauthorized)
	case core.NebErrorPartialFail, core.NebError:
		w.WriteHeader(http.StatusBadRequest)
	case core.NebErrorNone:
		w.WriteHeader(http.StatusOK)
	}

	fmt.Fprint(w, res)
}
func EasyErrorResponse(w http.ResponseWriter, code int, err error) {
	v, ok := err.(core.NebuleuseError)
	var e easyResponse

	if ok {
		e = easyResponse{v.Code, v.Msg}
	} else {
		e = easyResponse{code, err.Error()}
	}

	res, err := json.Marshal(e)
	if err != nil {
		core.Warning.Println("Could not encode easy response")
	}

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, res)
}

type statusResponse struct {
	Maintenance       bool
	NebuleuseVersion  int
	GameVersion       int
	UpdaterVersion    int
	ComplexStatsTable []core.ComplexStatTableInfo
}

func status(w http.ResponseWriter, r *http.Request) {
	CStatsInfos, err := core.GetComplexStatsTablesInfos()
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	response := statusResponse{false, core.NebuleuseVersion, core.GetGameVersion(), core.GetUpdaterVersion(), CStatsInfos}

	if core.Cfg["maintenance"] != "0" {
		response.Maintenance = true
	}

	res, err := json.Marshal(response)
	if err != nil {
		core.Warning.Println("Could not encode status response")
		EasyErrorResponse(w, core.NebError, err)
	} else {
		fmt.Fprint(w, string(res))
	}
}
