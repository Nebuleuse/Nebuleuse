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
func getAchievement(w http.ResponseWriter, r *http.Request) {
	id := context.Get(r, "achievementid").(string)
	/*id, err := strconv.ParseInt(sid, 10, 0)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}*/

	ach, err := core.GetAchievementData(id)
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

	var ach core.AchievementData
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

	var ach core.AchievementData
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

func getLogs(w http.ResponseWriter, r *http.Request) {
	res := core.GetPastLogs(5000)
	fmt.Fprint(w, string(res))
}

func getUserStatsList(w http.ResponseWriter, r *http.Request) {
	fields, err := core.GetUserStatsFields()
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	res, err := json.Marshal(fields)
	if err != nil {
		core.Warning.Println("Could not encode status response")
	} else {
		fmt.Fprint(w, string(res))
	}
}

func getStatTables(w http.ResponseWriter, r *http.Request) {
	fields, err := core.GetComplexStatsTablesInfos()
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	res, err := json.Marshal(fields)
	if err != nil {
		core.Warning.Println("Could not encode status response")
	} else {
		fmt.Fprint(w, string(res))
	}
}

func getStatTable(w http.ResponseWriter, r *http.Request) {
	name := context.Get(r, "name").(string)

	table, err := core.GetComplexStatsTableInfos(name)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}

	res, err := json.Marshal(table)
	if err != nil {
		core.Warning.Println("Could not encode status response")
	} else {
		fmt.Fprint(w, string(res))
	}
}

func addStatTable(w http.ResponseWriter, r *http.Request) {
	data := []byte(context.Get(r, "data").(string))
	var table core.ComplexStatTableInfo

	err := json.Unmarshal(data, &table)

	if err != nil {
		core.Warning.Println("Could not encode status response")
		return
	}

	err = core.AddStatTable(table)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyResponse(w, core.NebErrorNone, "added table")
	}
}

func setStatTable(w http.ResponseWriter, r *http.Request) {
	data := []byte(context.Get(r, "data").(string))
	var table core.ComplexStatTableInfo

	err := json.Unmarshal(data, &table)

	if err != nil {
		core.Warning.Println("Could not encode status response")
		return
	}

	err = core.SetStatTable(table)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyResponse(w, core.NebErrorNone, "updated table")
	}
}

func deleteStatTable(w http.ResponseWriter, r *http.Request) {
	name := context.Get(r, "name").(string)

	err := core.DeleteStatTable(name)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyResponse(w, core.NebErrorNone, "updated table")
	}
}
