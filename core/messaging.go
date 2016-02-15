package core

import (
	"encoding/json"
)

type MessagePipeline struct {
	name                         string
	rank                         int
	canUserJoin, canUserDispatch bool
	pipes                        map[string][]*UserSession
}

func createPipeline(name string, rank int, canUserJoin, canDispatch bool) MessagePipeline {
	var p MessagePipeline
	p.pipes = make(map[string][]*UserSession)
	p.name = name
	p.canUserJoin = canUserJoin
	p.canUserDispatch = canDispatch
	p.rank = rank
	return p
}

var messagePipelines map[string]MessagePipeline

type messageData struct {
	Channel string
	Message interface{}
}

func initMessaging() {
	messagePipelines = make(map[string]MessagePipeline)
	messagePipelines["system"] = createPipeline("system", 1, false, false)
	messagePipelines["admin"] = createPipeline("admin", 2, true, false)
}

func Listen(pipe, name string, user *UserSession) {
	pipeline := messagePipelines[pipe]
	if !pipeline.canUserJoin {
		return
	} else if pipeline.rank > user.UserRank {
		return
	}
	for _, list := range pipeline.pipes[name] { // prevent duplicates
		if user.UserId == list.UserId {
			return
		}
	}
	pipeline.pipes[name] = append(pipeline.pipes[name], user)
}
func StopListen(pipe, name string, user *UserSession) {
	pipeline := messagePipelines[pipe]

	list := pipeline.pipes[name]
	for i, usr := range list {
		if usr.UserId == user.UserId { // remove the userSession from the slice
			list, list[len(list)-1] = append(list[:i], list[i+1:]...), nil
		}
	}

	if len(pipeline.pipes[name]) == 0 {
		delete(pipeline.pipes, name)
	}
}
func UserStopListen(user *UserSession) {
	for _, pipeline := range messagePipelines {
		for _, pipe := range pipeline.pipes {
			for i, usr := range pipe {
				if usr.UserId == user.UserId {
					pipe, pipe[len(pipe)-1] = append(pipe[:i], pipe[i+1:]...), nil
				}
			}
		}
	}
}
func CanUserDispatch(pipe string) bool {
	return messagePipelines[pipe].canUserDispatch
}
func CanUserJoin(pipe string, user *UserSession) bool {
	return messagePipelines[pipe].canUserJoin && messagePipelines[pipe].rank <= user.UserId
}
func Dispatch(pipe, name string, message interface{}) {
	var msg messageData
	msg.Channel = name
	msg.Message = message
	res, err := json.Marshal(msg)
	if err != nil {
		Error.Println("Could not marshal message dispatch : ", err)
		return
	}

	for _, id := range messagePipelines[pipe].pipes[name] {
		SendMessageToUserId(id, string(res))
	}
}
func GetMessages(userid int) chan string {
	return connectedUsers[userid].Messages
}

func SendMessageToUserId(user *UserSession, message string) {
	user.Messages <- message
}
