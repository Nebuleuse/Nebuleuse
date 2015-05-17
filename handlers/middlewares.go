package handlers

import (
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/context"
	"net/http"
	"strconv"
)

type middleWare func(http.ResponseWriter, *http.Request)

//Populates the context with the user struct using the request sessionId
func userBySession(allowTarget bool, next middleWare) middleWare {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("sessionid") == "" {
			EasyResponse(w, core.NebError, "Missing sessionid")
			return
		}

		requester, err := core.GetUserBySession(r.FormValue("sessionid"), core.UserMaskOnlyId)

		if err != nil {
			EasyErrorResponse(w, core.NebErrorDisconnected, err)
			return
		}

		context.Set(r, "sessionid", r.FormValue("sessionid"))
		context.Set(r, "requester", requester)
		context.Set(r, "session", core.GetSessionBySessionId(r.FormValue("sessionid")))

		if r.FormValue("userid") == "" {
			context.Set(r, "user", requester)
			next(w, r)
		} else {
			user := new(core.User)
			user.Id, err = strconv.Atoi(r.FormValue("userid"))
			context.Set(r, "user", user)
			if !allowTarget {
				// If we don't allow targeting other users, only admins can target
				n := mustBeAdmin(next)
				n(w, r)
			} else {
				next(w, r)
			}
		}
	}
}

// Verifies the existence of specified form fields
func verifyFormValuesExist(vals []string, next middleWare) middleWare {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, val := range vals {
			if r.FormValue(val) == "" {
				EasyResponse(w, core.NebError, "Missing input : "+val)
				return
			}
			context.Set(r, val, r.FormValue(val))
		}
		next(w, r)
	}
}

// Verifies context's requester rank for auth level
func mustBeAdmin(next middleWare) middleWare {
	return authRank(2, next)
}
func authRank(rank int, next middleWare) middleWare {
	return func(w http.ResponseWriter, r *http.Request) {
		irqst, ok := context.GetOk(r, "requester")
		if !ok {
			EasyResponse(w, core.NebErrorAuthFail, "No User to verify admin rights on")
			return
		}
		rqst := irqst.(*core.User)
		rqst.FetchUserInfos(core.UserMaskBase)
		if rqst.Rank < 2 {
			EasyResponse(w, core.NebErrorAuthFail, "Unauthorized")
			return
		}
		next(w, r)
	}
}
