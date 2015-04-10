package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Nebuleuse/Nebuleuse/core"
	"net/http"
	"strconv"
)

type connectResponse struct {
	SessionId string
}

func connectUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["username"] == nil || r.PostForm["password"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing username and/or password"))
		return
	}

	username := r.PostForm["username"][0] // For some reason r.PostForm[i] is String[]
	password := r.PostForm["password"][0]
	id, err := core.CreateSession(username, password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid"))
		return
	}

	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskOnlyId)

	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
		return
	}

	user.Disconnect()
	fmt.Fprint(w, EasyResponse(core.NebErrorNone, "User disconnected"))
}

func getUserInfos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var user *core.User
	user = new(core.User)
	var err error
	if r.PostForm["infomask"] != nil {
		mask, err := strconv.ParseInt(r.PostForm["infomask"][0], 10, 0)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, EasyErrorResponse(core.NebError, err))
			return
		} else if r.PostForm["sessionid"] != nil && r.PostForm["infomask"] != nil {
			user, err = core.GetUserBySession(r.FormValue("sessionid"), int(mask))

		} else if r.PostForm["userid"] != nil && r.PostForm["infomask"] != nil {
			var id int64
			id, err = strconv.ParseInt(r.FormValue("userid"), 10, 0)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, EasyResponse(core.NebError, "Invalid userid"))
				return
			}

			user.Id = int(id)

			err = user.FetchUserInfos(int(mask))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, EasyErrorResponse(core.NebError, err))
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid or userid and infomask"))
			return
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, EasyErrorResponse(core.NebError, err))
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid or data"))
		return
	}

	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskOnlyId)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
		return
	}

	data := r.PostForm["data"][0]
	var req updateAchievementsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid or data"))
		return
	}

	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskOnlyId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
		return
	}

	data := r.PostForm["data"][0]
	var req updateStatsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid or data"))
		return
	}

	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskOnlyId)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
		return
	}

	data := r.PostForm["data"][0]
	var req updateComplexStatsRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, EasyErrorResponse(core.NebError, err))
		return
	}

	err = user.UpdateComplexStats(req.Stats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, EasyErrorResponse(core.NebError, err))
		return
	}

	fmt.Fprint(w, EasyResponse(core.NebErrorNone, "Inserted complex stat"))
}
