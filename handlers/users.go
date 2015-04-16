package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/context"
	"net/http"
	"strconv"
)

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
	smask := context.Get(r, "infomask").(string)
	mask, err := strconv.Atoi(smask)
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

	res, err := json.Marshal(user)
	fmt.Fprint(w, string(res))
}

type setAchievementsRequest struct {
	Achievements []struct {
		Id    int
		Value int
	}
}

func setAchievements(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*core.User)

	data := context.Get(r, "data").([]byte)
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

func setStats(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*core.User)

	data := context.Get(r, "data").([]byte)
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

func addComplexStats(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*core.User)

	data := context.Get(r, "data").([]byte)
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
