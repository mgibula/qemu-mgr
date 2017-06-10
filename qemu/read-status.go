package qemu

import "log"

// -> { "execute": "query-status" }
// <- { "return": { "running": true,
//                  "singlestep": false,
//                  "status": "running" } }

type QueryStatusMessage struct {
	MonitorRequest
}

type ReadStatus struct {
	MonitorResponse

	Return struct {
		Running    bool   `json:"running"`
		Singlestep bool   `json:"singlestep"`
		Status     string `json:"status"`
	} `json:"return"`
}

func (command *ReadStatus) Execute(monitor *MonitorState) {
	request := QueryStatusMessage{}
	request.Execute = "query-status"
	monitor.SendJson(request)
}

func (command *ReadStatus) ProcessResponse(monitor *MonitorState) {
	log.Printf("[%v] Instance status is: %v. (singlestep: %v, running: %v)\n", monitor.Path, command.Return.Status, command.Return.Singlestep, command.Return.Running)
}
