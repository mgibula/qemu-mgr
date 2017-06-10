package qemu

import "log"

// -> { "execute": "query-name" }
// <- { "return": { "name": "qemu-name" } }

type QueryNameMessage struct {
	MonitorRequest
}

type ReadName struct {
	MonitorResponse

	Return struct {
		Name string `json:"name"`
	} `json:"return"`
}

func (command *ReadName) Execute(monitor *MonitorState) {
	request := QueryNameMessage{}
	request.Execute = "query-name"
	monitor.SendJson(request)
}

func (command *ReadName) ProcessResponse(monitor *MonitorState) {
	log.Printf("[%v] Instance name is %v", monitor.Path, command.Return.Name)
	monitor.Instance.Name = command.Return.Name
}
