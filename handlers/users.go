package handlers

import (
	"Nebuleuse/core"
	"encoding/json"
	"fmt"
	"net/http"
)

type connectResponse struct {
	SessionId string
}

func connectUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["username"] == nil || r.PostForm["password"] == nil {
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing username and/or password"))
		return
	}

	username := r.PostForm["username"][0] // For some reason r.PostForm[i] is String[]
	password := r.PostForm["password"][0]
	id, err := core.CreateSession(username, password)

	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorLogin, err))
		return
	}

	response := connectResponse{id}

	res, err := json.Marshal(response)
	if err != nil {
		core.Warning.Println("Could not encode status response")
	} else {
		fmt.Fprint(w, string(res))
	}
}

func disconnectUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil {
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid"))
		return
	}

	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskNone)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
		return
	}

	user.Disconnect()
	fmt.Fprint(w, EasyResponse(core.NebErrorNone, "User disconnected"))
}

func getUserInfos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil {
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid"))
		return
	}

	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskStats|core.UserMaskAchievements)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
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
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid or data"))
		return
	}

	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskNone)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
		return
	}

	data := r.PostForm["data"][0]
	var req updateAchievementsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebError, err))
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
		fmt.Fprint(w, EasyResponse(core.NebErrorPartialFail, "Updated "+string(count)+" Achievements"))
		return
	}

	fmt.Fprint(w, EasyResponse(core.NebErrorNone, "Updated Achievements"))
}

type updateStatsRequest struct {
	Stats []core.UserStat
}

func updateStats(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil || r.PostForm["data"] == nil {
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid or data"))
		return
	}

	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskNone)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
		return
	}

	data := r.PostForm["data"][0]
	var req updateStatsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebError, err))
		return
	}

	user.UpdateStats(req.Stats)

	fmt.Fprint(w, EasyResponse(core.NebErrorNone, "Updated Stats"))
}

type updateComplexStatsRequest struct {
	Stats []core.ComplexStat
}

func addComplexStats(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil || r.PostForm["data"] == nil {
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid or data"))
		return
	}

	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskNone)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
		return
	}

	data := r.PostForm["data"][0]
	var req updateComplexStatsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebError, err))
		return
	}

	err = user.UpdateComplexStats(req.Stats)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebError, err))
		return
	}

	fmt.Fprint(w, EasyResponse(core.NebErrorNone, "Inserted complex stat"))
}
