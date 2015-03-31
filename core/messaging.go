package core

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
	for val := range messagePipeline[name] {
		ok := SendMessageToUserId(val, message)
		if !ok {
			Warning.Println("Inexistant userid in connectedUser still present in Messaging pipepline: " + name + " : " + string(val))
			StopListen(name, val)
		}
	}
}
func GetMessages(userid int) chan string {
	return connectedUsers[userid].Messages
}
