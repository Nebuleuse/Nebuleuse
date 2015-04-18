package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/context"
	"net/http"
	"strconv"
)

type dashboardInfosResponse struct {
	UserCount   int
	OnlineUsers int
	UpdateCount int
}

func getDashboardInfos(w http.ResponseWriter, r *http.Request) {
	var dashRes dashboardInfosResponse
	dashRes.OnlineUsers = core.CountOnlineUsers()
	dashRes.UserCount = core.GetUserCount()
	dashRes.UpdateCount = core.GetUpdateCount()
	res, err := json.Marshal(dashRes)
	if err != nil {
		core.Warning.Println("Could not encode status response")
	} else {
		fmt.Fprint(w, string(res))
	}
}

func getAchievements(w http.ResponseWriter, r *http.Request) {
	ach, err := core.GetAchievementsData()
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	res, err := json.Marshal(ach)
	if err != nil {
		core.Warning.Println("Could not encode status response")
	} else {
		fmt.Fprint(w, string(res))
	}
}

func setAchievement(w http.ResponseWriter, r *http.Request) {
	data := context.Get(r, "data").(string)
	sid := context.Get(r, "achievementid").(string)

	id, err := strconv.ParseInt(sid, 10, 0)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	var ach core.AchievementsTable
	json.Unmarshal([]byte(data), &ach)
	err = core.SetAchievementData(int(id), ach)

	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	EasyResponse(w, core.NebErrorNone, "Updated achievement table")
}

func addAchievement(w http.ResponseWriter, r *http.Request) {
	data := context.Get(r, "data").(string)

	var ach core.AchievementsTable
	json.Unmarshal([]byte(data), &ach)
	value, err := core.AddAchievementData(ach)

	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	res, err := json.Marshal(value)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	fmt.Fprint(w, string(res))
}

func deleteAchievement(w http.ResponseWriter, r *http.Request) {
	sid := context.Get(r, "achievementid").(string)

	id, err := strconv.ParseInt(sid, 10, 0)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	err = core.DeleteAchievementData(int(id))

	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	EasyResponse(w, core.NebErrorNone, "deleted achievement table")
}
