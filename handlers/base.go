package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/mux"
)

func RegisterHandlers() {
	r := mux.NewRouter()

	r.HandleFunc("/", status)
	r.HandleFunc("/status", status).Methods("GET")
	r.HandleFunc("/connect", verifyFormValuesExist([]string{"username", "password"}, connectUser)).Methods("POST")
	r.HandleFunc("/disconnect", userBySession(false, disconnectUser)).Methods("POST")

	//Messaging
	r.HandleFunc("/getMessages", userBySession(false, getMessages)).Methods("POST")
	r.HandleFunc("/sendMessage", userBySession(false, verifyFormValuesExist([]string{"pipe", "channel", "message"}, sendMessage))).Methods("POST")
	r.HandleFunc("/subscribeTo", userBySession(false, verifyFormValuesExist([]string{"pipe", "channel"}, subscribeTo))).Methods("POST")
	r.HandleFunc("/unSubscribeTo", userBySession(false, verifyFormValuesExist([]string{"pipe", "channel"}, unSubscribeTo))).Methods("POST")

	//Administration
	r.PathPrefix("/admin/").Handler((http.StripPrefix("/admin/", http.FileServer(http.Dir(core.Cfg.GetSysConfig("DashboardLocation"))))))
	r.HandleFunc("/getDashboardInfos", userBySession(false, mustBeAdmin(getDashboardInfos))).Methods("POST")
	r.HandleFunc("/getLogs", userBySession(false, mustBeAdmin(getLogs))).Methods("POST")

	//User
	r.HandleFunc("/getUserInfos", userBySession(true, verifyFormValuesExist([]string{"infomask"}, getUserInfos))).Methods("POST")
	r.HandleFunc("/setUserAchievements", userBySession(false, verifyFormValuesExist([]string{"data"}, setUserAchievements))).Methods("POST")
	r.HandleFunc("/setUserStats", userBySession(false, verifyFormValuesExist([]string{"data"}, setUserStats))).Methods("POST")
	r.HandleFunc("/addComplexStats", userBySession(false, verifyFormValuesExist([]string{"data"}, addComplexStats))).Methods("POST")

	//Users
	r.HandleFunc("/getUsersInfos", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"infomask", "page"}, getUsersInfos)))).Methods("POST")
	r.HandleFunc("/getOnlineUsersList", userBySession(false, mustBeAdmin(getOnlineUsersList))).Methods("POST")

	//Achievements
	r.HandleFunc("/getAchievements", userBySession(false, mustBeAdmin(getAchievements))).Methods("POST")
	r.HandleFunc("/getAchievement", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"achievementid"}, getAchievement)))).Methods("POST")
	r.HandleFunc("/setAchievement", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"achievementid", "data"}, setAchievement)))).Methods("POST")
	r.HandleFunc("/addAchievement", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"data"}, addAchievement)))).Methods("POST")
	r.HandleFunc("/deleteAchievement", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"achievementid"}, deleteAchievement)))).Methods("POST")

	//Stats
	r.HandleFunc("/getStatTables", userBySession(false, getStatTables)).Methods("POST")
	r.HandleFunc("/getStatTable", userBySession(false, verifyFormValuesExist([]string{"name"}, getStatTable))).Methods("POST")
	r.HandleFunc("/setStatTable", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"data"}, setStatTable)))).Methods("POST")
	r.HandleFunc("/addStatTable", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"data"}, addStatTable)))).Methods("POST")
	r.HandleFunc("/deleteStatTable", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"name"}, deleteStatTable)))).Methods("POST")
	r.HandleFunc("/setUsersStatFields", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"fields"}, setUsersStatFields)))).Methods("POST")

	//Updates
	r.PathPrefix("/updates/").Handler((http.StripPrefix("/updates/", http.FileServer(http.Dir(core.Cfg.GetSysConfig("UpdatesLocation"))))))
	r.HandleFunc("/getBranchList", userBySession(false, getBranchList)).Methods("POST")
	r.HandleFunc("/getBranchUpdates", userBySession(false, verifyFormValuesExist([]string{"branch"}, getBranchUpdates))).Methods("POST")
	r.HandleFunc("/getCompleteBranchUpdates", userBySession(false, mustBeAdmin(getCompleteBranchUpdates))).Methods("POST")
	r.HandleFunc("/addBranchFromBuild", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"build", "name", "log", "semver", "accessrank"}, addBranchFromBuild)))).Methods("POST")
	r.HandleFunc("/addEmptyBranch", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"name", "accessrank"}, addEmptyBranch)))).Methods("POST")

	r.HandleFunc("/updateGitCommitCacheList", userBySession(false, mustBeAdmin(updateGitCommitCacheList))).Methods("POST")
	r.HandleFunc("/prepareGitBuild", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"commit"}, prepareGitBuild))))
	r.HandleFunc("/addGitBuild", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"commit", "log"}, createGitBuild))))
	r.HandleFunc("/addUpdate", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"build", "branch", "semver", "log"}, createUpdate))))
	r.HandleFunc("/setActiveUpdate", userBySession(false, mustBeAdmin(verifyFormValuesExist([]string{"build", "branch"}, setActiveUpdate))))
	http.Handle("/", r)
}
func RegisterInstallHandlers() {
	r := mux.NewRouter()
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
	case core.NebErrorLogin, core.NebErrorAuthFail, core.NebErrorDisconnected:
		w.WriteHeader(http.StatusUnauthorized)
	case core.NebErrorPartialFail, core.NebError:
		w.WriteHeader(http.StatusBadRequest)
	case core.NebErrorNone:
		w.WriteHeader(http.StatusOK)
	}

	fmt.Fprint(w, string(res))
}
func EasyDataResponse(w http.ResponseWriter, data interface{}) {
	res, err := json.Marshal(data)
	if err != nil {
		core.Warning.Println("Could not encode data response", data)
		EasyErrorResponse(w, core.NebError, err)
	} else {
		fmt.Fprint(w, string(res))
	}
}
func EasyErrorResponse(w http.ResponseWriter, code int, err error) {
	v, ok := err.(core.NebuleuseError)

	if ok {
		EasyResponse(w, v.Code, v.Msg)
		return
	}

	e := easyResponse{code, err.Error()}

	res, err := json.Marshal(e)
	if err != nil {
		core.Warning.Println("Could not encode easy response")
	} else {
		switch code {
		case core.NebErrorLogin, core.NebErrorAuthFail, core.NebErrorDisconnected:
			w.WriteHeader(http.StatusUnauthorized)
		case core.NebErrorPartialFail:
			w.WriteHeader(http.StatusBadRequest)
		case core.NebError:
			w.WriteHeader(http.StatusInternalServerError)
		case core.NebErrorNone:
			w.WriteHeader(http.StatusOK)
		}

		fmt.Fprint(w, string(res))
	}
}

type statusResponse struct {
	NebuleuseVersion  int
	UpdaterVersion    int
	UpdateSystem      string
	UpdatesLocation   string
	ComplexStatsTable []core.ComplexStatTableInfo
}

func status(w http.ResponseWriter, r *http.Request) {
	CStatsInfos, err := core.GetComplexStatsTablesInfos()
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	response := statusResponse{core.NebuleuseVersion, core.GetUpdaterVersion(), core.GetUpdateSystem(), core.GetUpdatesLocation(), CStatsInfos}

	EasyDataResponse(w, response)
}
