package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/context"
	"net/http"
	"strconv"
)

//Populates the context with the user struct using the request sessionId
func userBySession(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("sessionid") == "" {
			EasyResponse(w, core.NebError, "Missing sessionid")
			return
		}

		user, err := core.GetUserBySession(r.FormValue("sessionid"), core.UserMaskOnlyId)

		if err != nil {
			EasyErrorResponse(w, core.NebErrorDisconnected, err)
			return
		}

		context.Set(r, "user", user)
		next(w, r)
	}
}

// Verifies there is data being sent
func verifyFormDataExist(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("data") == "" {
			EasyResponse(w, core.NebError, "Missing data")
			return
		}
		context.Set(r, "data", r.FormValue("data"))
		next(w, r)
	}
}

// Verifies context's user rank for auth level
func mustBeAdmin(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		iusr, ok := context.GetOk(r, "user")
		if !ok {
			EasyResponse(w, core.NebError, "No User to verify admin rights on")
			return
		}
		usr := iusr.(*core.User)
		usr.FetchUserInfos(core.UserMaskBase)
		if usr.Rank < 2 {
			EasyResponse(w, core.NebError, "Unauthorized")
			return
		}
		next(w, r)
	}
}

type connectResponse struct {
	SessionId string
}

func connectUser(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("username") == "" || r.FormValue("password") == "" {
		EasyResponse(w, core.NebError, "Missing username and/or password")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	id, err := core.CreateSession(username, password)

	if err != nil {
		EasyErrorResponse(w, core.NebErrorLogin, err)
		return
	}

	response := connectResponse{id}

	res, err := json.Marshal(response)
	if err != nil {
		core.Warning.Println("Could not encode status response")
		EasyErrorResponse(w, core.NebError, err)
	} else {
		fmt.Fprint(w, string(res))
	}
}

func disconnectUser(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*core.User)
	user.Disconnect()
	EasyResponse(w, core.NebErrorNone, "User disconnected")
}

func getUserInfos(w http.ResponseWriter, r *http.Request) {
	var user *core.User
	user = new(core.User)
	var err error
	if r.FormValue("infomask") != "" {
		mask, err := strconv.ParseInt(r.FormValue("infomask"), 10, 0)

		if err != nil {
			EasyErrorResponse(w, core.NebError, err)
			return
		} else if r.FormValue("sessionid") != "" && r.FormValue("infomask") != "" {
			user, err = core.GetUserBySession(r.FormValue("sessionid"), int(mask))
			if err != nil {
				EasyResponse(w, core.NebError, "Invalid sessionid")
				return
			}
		} else if r.FormValue("userid") != "" && r.FormValue("infomask") != "" {
			var id int64
			id, err = strconv.ParseInt(r.FormValue("userid"), 10, 0)
			if err != nil {
				EasyResponse(w, core.NebError, "Invalid userid")
				return
			}

			user.Id = int(id)

			err = user.FetchUserInfos(int(mask))
			if err != nil {
				EasyErrorResponse(w, core.NebError, err)
				return
			}
		} else {
			EasyResponse(w, core.NebError, "Missing sessionid or userid and infomask")
			return
		}
	}
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
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
	user := context.Get(r, "user").(*core.User)

	data := context.Get(r, "data").([]byte)
	var req updateAchievementsRequest
	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
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
		EasyResponse(w, core.NebErrorPartialFail, "Updated "+string(count)+" Achievements")
		return
	}

	EasyResponse(w, core.NebErrorNone, "Updated Achievements")
}

type updateStatsRequest struct {
	Stats []core.UserStat
}

func updateStats(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*core.User)

	data := context.Get(r, "data").([]byte)
	var req updateStatsRequest
	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	user.UpdateStats(req.Stats)
	EasyResponse(w, core.NebErrorNone, "Updated Stats")
}

type updateComplexStatsRequest struct {
	Stats []core.ComplexStat
}

func addComplexStats(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*core.User)

	data := context.Get(r, "data").([]byte)
	var req updateComplexStatsRequest
	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	err = user.UpdateComplexStats(req.Stats)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	EasyResponse(w, core.NebErrorNone, "Inserted complex stat")
}
