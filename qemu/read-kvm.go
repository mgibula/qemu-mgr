package qemu

import "log"

// -> { "execute": "query-kvm" }
// <- { "return": { "enabled": true, "present": true } }

type QueryKvmMessage struct {
	MonitorRequest
}

type ReadKvm struct {
	MonitorResponse

	Return struct {
		Enabled bool `json:"enabled"`
		Present bool `json:"present"`
	} `json:"return"`
}

func (command *ReadKvm) Execute(monitor *MonitorState) {
	request := QueryKvmMessage{}
	request.Execute = "query-kvm"
	monitor.SendJson(request)
}

func (command *ReadKvm) ProcessResponse(monitor *MonitorState) {
	var present, enabled string
	if command.Return.Present {
		present = "present"
	} else {
		present = "not present"
	}

	if command.Return.Enabled {
		enabled = "enabled"
	} else {
		enabled = "not enabled"
	}

	log.Printf("[%v] KVM is %v and %v", monitor.Path, present, enabled)
}
