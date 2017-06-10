package qemu

import "log"

// -> { "execute": "query-cpus" }
// <- { "return": [
//          {
//             "CPU":0,
//             "current":true,
//             "halted":false,
//             "qom_path":"/machine/unattached/device[0]",
//             "arch":"x86",
//             "pc":3227107138,
//             "thread_id":3134
//          },
//          {
//             "CPU":1,
//             "current":false,
//             "halted":true,
//             "qom_path":"/machine/unattached/device[2]",
//             "arch":"x86",
//             "pc":7108165,
//             "thread_id":3135
//          }
//       ]
//    }

type QueryCpusMessage struct {
	MonitorRequest
}

type ReadCpus struct {
	MonitorResponse

	Return []struct {
		Cpu      int    `json:"CPU"`
		Current  bool   `json:"current"`
		Halted   bool   `json:"halted"`
		QomPath  string `json:"qom_path"`
		Arch     string `json:"arch"`
		Pc       uint   `json:"pc"`
		ThreadId uint   `json:"thread_id"`
	} `json:"return"`
}

func (command *ReadCpus) Execute(monitor *MonitorState) {
	request := QueryStatusMessage{}
	request.Execute = "query-cpus"
	monitor.SendJson(request)
}

func (command *ReadCpus) ProcessResponse(monitor *MonitorState) {
	log.Printf("[%v] Found %v CPUs\n", monitor.Path, len(command.Return))
}
