package qemu

import "log"

// { "execute": "qmp_capabilities" }
// {"return": {}}

type QmpCapabilitiesMessage struct {
	MonitorRequest
}

type ReadQmpCapabilities struct {
	MonitorResponse

	Return struct {
	} `json:"return"`
}

func (command *ReadQmpCapabilities) Execute(monitor *MonitorState) {
	request := QmpCapabilitiesMessage{}
	request.Execute = "qmp_capabilities"

	monitor.SendJson(request)
}

func (command *ReadQmpCapabilities) ProcessResponse(monitor *MonitorState) {
	log.Printf("[%v] Initial communication completed with no error", monitor.Path)
}
