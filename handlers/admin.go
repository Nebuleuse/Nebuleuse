package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Nebuleuse/Nebuleuse/core"
	"net/http"
)

type dashboardInfosResponse struct {
	OnlineUsers int
	TotalUsers  int
}

func getDashboardInfos(w http.ResponseWriter, r *http.Request) {
	//usr := context.Get(r, "user").(*core.User)
	var dashRes dashboardInfosResponse
	dashRes.OnlineUsers = 5
	dashRes.TotalUsers = 10
	res, err := json.Marshal(dashRes)
	if err != nil {
		core.Warning.Println("Could not encode status response")
	} else {
		fmt.Fprint(w, string(res))
	}
}
