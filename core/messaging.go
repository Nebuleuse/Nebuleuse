package core

import (
	"encoding/json"
)

var messagePipeline map[string][]int

type messageData struct {
	Channel string
	Message interface{}
}

func initMessaging() {
	messagePipeline = make(map[string][]int)
}
func Listen(name string, userid int) {
	for _, id := range messagePipeline[name] { // prevent duplicates
		if id == userid {
			return
		}
	}
	messagePipeline[name] = append(messagePipeline[name], userid)
}
func StopListen(name string, userid int) {
	messagePipeline[name] = append(messagePipeline[name][:userid], messagePipeline[name][userid+1:]...)
	if len(messagePipeline[name]) == 0 {
		delete(messagePipeline, name)
	}
}
func Dispatch(name string, message interface{}) {
	var msg messageData
	msg.Channel = name
	msg.Message = message
	res, err := json.Marshal(msg)
	if err != nil {
		Error.Println("Could not marshal message dispatch : ", err)
		return
	}

	for _, id := range messagePipeline[name] {
		ok := SendMessageToUserId(id, string(res))
		if !ok {
			Warning.Println("Inexistant userid in connectedUser still present in Messaging pipepline: " + name + " : " + string(id))
			StopListen(name, id)
		}
	}
}
func GetMessages(userid int) chan string {
	return connectedUsers[userid].Messages
}

func SendMessageToUserId(userid int, message string) bool {
	session := GetSessionByUserId(userid)
	if session != nil {
		session.Messages <- message
		return true
	}
	return false
}
