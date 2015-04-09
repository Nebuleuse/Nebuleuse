package core

import (
	"strconv"
)

var messagePipeline map[string][]int

func initMessaging() {
	messagePipeline = make(map[string][]int)
}
func Listen(name string, userid int) {

	messagePipeline[name] = append(messagePipeline[name], userid)
}
func StopListen(name string, userid int) {
	messagePipeline[name] = append(messagePipeline[name][:userid], messagePipeline[name][userid+1:]...)
	if len(messagePipeline[name]) == 0 {
		delete(messagePipeline, name)
	}
}
func Dispatch(name, message string) {
	for _, id := range messagePipeline[name] {
		Info.Println("Sending message to " + strconv.Itoa(id) + " : " + message)
		ok := SendMessageToUserId(id, message)
		if !ok {
			Warning.Println("Inexistant userid in connectedUser still present in Messaging pipepline: " + name + " : " + string(id))
			StopListen(name, id)
		}
	}
}
func GetMessages(userid int) chan string {
	return connectedUsers[userid].Messages
}
