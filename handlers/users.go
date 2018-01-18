package handlers

import (
	"encoding/json"
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/context"
	"net/http"
	"strconv"
)

type connectResponse struct {
	SessionId string
}

func connectUser(w http.ResponseWriter, r *http.Request) {
	username := context.Get(r, "username").(string)
	password := context.Get(r, "password").(string)
	id, err := core.CreateSession(username, password)

	if err != nil {
		EasyErrorResponse(w, core.NebErrorLogin, err)
		return
	}

	response := connectResponse{id}

	EasyDataResponse(w, response)
}

func disconnectUser(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*core.User)
	user.Disconnect()
	EasyResponse(w, core.NebErrorNone, "User disconnected")
}

//User connected, must be admin
func getOnlineUsersList(w http.ResponseWriter, r *http.Request) {
	list := core.GetOnlineUsersList()

	EasyDataResponse(w, list)
}

//User connected, can target someone else, form value: infomask
func getUserInfos(w http.ResponseWriter, r *http.Request) {
	smask := context.Get(r, "infomask").(string)
	mask, err := strconv.ParseInt(smask, 10, 0)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	user := context.Get(r, "user").(*core.User)

	err = user.FetchUserInfos(int(mask))
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	EasyDataResponse(w, user)
}

//User connected, must be admin, form value: infomask, page
func getUsersInfos(w http.ResponseWriter, r *http.Request) {
	smask := context.Get(r, "infomask").(string)
	mask, err := strconv.ParseInt(smask, 10, 0)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	spage := context.Get(r, "page").(string)
	page, err := strconv.ParseInt(spage, 10, 0)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	// Pages start at 1
	page -= 1
	users, err := core.GetUsersInfos(int(page)*30, 30, int(mask))
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	EasyDataResponse(w, users)
}

type setAchievementsRequest struct {
	Achievements []struct {
		Id    int
		Value int
	}
}

//User connected, form value: data
func setUserAchievements(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*core.User)

	data := context.Get(r, "data").(string)
	var req setAchievementsRequest
	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	count := len(req.Achievements)
	for _, ach := range req.Achievements {
		err = user.SetAchievementProgress(ach.Id, ach.Value)
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

type setStatsRequest struct {
	Stats []core.UserStat
}

//User connected, form value: data
func setUserStats(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*core.User)

	data := context.Get(r, "data").(string)
	var req setStatsRequest
	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	user.SetStats(req.Stats)
	EasyResponse(w, core.NebErrorNone, "Updated Stats")
}

type setComplexStatsRequest struct {
	Stats []core.ComplexStat
}

//User connected, form value: data
func addComplexStats(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*core.User)

	data := context.Get(r, "data").(string)
	var req setComplexStatsRequest
	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	err = user.SetComplexStats(req.Stats)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	EasyResponse(w, core.NebErrorNone, "Inserted complex stat")
}
