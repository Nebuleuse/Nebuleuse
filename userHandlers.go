package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type connectResponse struct {
	SessionId string
}

func connectUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["username"] == nil || r.PostForm["password"] == nil {
		fmt.Fprint(w, EasyResponse(NebError, "Missing username and/or password"))
		return
	}

	username := r.PostForm["username"][0] // For some reason r.PostForm[i] is String[]
	password := r.PostForm["password"][0]
	id, err := CreateSession(username, password)

	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebErrorLogin, err))
		return
	}

	response := connectResponse{id}

	res, err := json.Marshal(response)
	if err != nil {
		log.Println("Could not encode status response")
	} else {
		fmt.Fprint(w, string(res))
	}
}

func disconnectUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil {
		fmt.Fprint(w, EasyResponse(NebError, "Missing sessionid"))
		return
	}

	user, err := GetUserBySession(r.PostForm["sessionid"][0], UserMaskOnlyId)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebErrorDisconnected, err))
		return
	}

	user.Disconnect()
	fmt.Fprint(w, EasyResponse(NebErrorNone, "User disconnected"))
}

func getUserInfos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var user *User
	var err error
	if r.PostForm["sessionid"] != nil {
		user, err = GetUserBySession(r.PostForm["sessionid"][0], UserMaskAll)

	} else if r.PostForm["userid"] != nil && r.PostForm["infomask"] != nil {
		id, err := strconv.ParseInt(r.PostForm["userid"][0], 10, 8)
		user.Id = int(id)
		mask, err := strconv.ParseInt(r.PostForm["infomask"][0], 10, 8)

		if err != nil {
			fmt.Fprint(w, EasyErrorResponse(NebError, err))
			return
		}

		err = user.FetchUserInfos(int(mask))
	} else {
		fmt.Fprint(w, EasyResponse(NebError, "Missing sessionid or userid and infomask"))
		return
	}
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebError, err))
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
		fmt.Fprint(w, EasyResponse(NebError, "Missing sessionid or data"))
		return
	}

	user, err := GetUserBySession(r.PostForm["sessionid"][0], UserMaskOnlyId)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebErrorDisconnected, err))
		return
	}

	data := r.PostForm["data"][0]
	var req updateAchievementsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebError, err))
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
		fmt.Fprint(w, EasyResponse(NebErrorPartialFail, "Updated "+string(count)+" Achievements"))
		return
	}

	fmt.Fprint(w, EasyResponse(NebErrorNone, "Updated Achievements"))
}

type updateStatsRequest struct {
	Stats []UserStat
}

func updateStats(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil || r.PostForm["data"] == nil {
		fmt.Fprint(w, EasyResponse(NebError, "Missing sessionid or data"))
		return
	}

	user, err := GetUserBySession(r.PostForm["sessionid"][0], UserMaskOnlyId)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebErrorDisconnected, err))
		return
	}

	data := r.PostForm["data"][0]
	var req updateStatsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebError, err))
		return
	}

	user.UpdateStats(req.Stats)

	fmt.Fprint(w, EasyResponse(NebErrorNone, "Updated Stats"))
}

type updateComplexStatsRequest struct {
	Stats []ComplexStat
}

func addComplexStats(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil || r.PostForm["data"] == nil {
		fmt.Fprint(w, EasyResponse(NebError, "Missing sessionid or data"))
		return
	}

	user, err := GetUserBySession(r.PostForm["sessionid"][0], UserMaskOnlyId)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebErrorDisconnected, err))
		return
	}

	data := r.PostForm["data"][0]
	var req updateComplexStatsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebError, err))
		return
	}

	err = user.updateComplexStats(req.Stats)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(NebError, err))
		return
	}

	fmt.Fprint(w, EasyResponse(NebErrorNone, "Inserted complex stat"))
}
