package qemu

import "log"

// -> { "execute": "query-uuid" }
// <- { "return": { "UUID": "550e8400-e29b-41d4-a716-446655440000" } }

type QueryUuidMessage struct {
	MonitorRequest
}

type ReadUuid struct {
	MonitorResponse

	Return struct {
		Uuid string `json:"UUID"`
	} `json:"return"`
}

func (command *ReadUuid) Execute(monitor *MonitorState) {
	request := QueryUuidMessage{}
	request.Execute = "query-uuid"
	monitor.SendJson(request)
}

func (command *ReadUuid) ProcessResponse(monitor *MonitorState) {
	log.Printf("[%v] Instance UUID is %v", monitor.Path, command.Return.Uuid)
	monitor.Instance.Uuid = command.Return.Uuid
}
