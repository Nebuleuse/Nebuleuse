package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Nebuleuse/Nebuleuse/core"
	"net/http"
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

type getAchievementsResponse struct {
	Achievements []core.AchievementsTable
}

func getAchievements(w http.ResponseWriter, r *http.Request) {
	var achievements getAchievementsResponse
	ach, err := core.GetAchievements()
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	achievements.Achievements = ach
	res, err := json.Marshal(achievements)
	if err != nil {
		core.Warning.Println("Could not encode status response")
	} else {
		core.Info.Println(string(res))
		fmt.Fprint(w, string(res))
	}
}
