package qemu

import (
	"fmt"
	"io/ioutil"
	"log"
)

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
	ParseCmdLine bool

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
	request := QueryCpusMessage{}
	request.Execute = "query-cpus"
	monitor.SendJson(request)
}

func (command *ReadCpus) ProcessResponse(monitor *MonitorState) {
	log.Printf("[%v] Found %v CPUs\n", monitor.Path, len(command.Return))
	monitor.Instance.Cpus = uint(len(command.Return))

	if command.ParseCmdLine {
		log.Printf("[%v] Parsing command line from /proc/%d/cmdline", monitor.Path, command.Return[0].ThreadId)
		data, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", command.Return[0].ThreadId))
		if err != nil {
			log.Printf("Error %v while reading /proc/%d/cmdline", err, command.Return[0].ThreadId)
		}

		monitor.Instance.ParseCommandLine(string(data))
	}
}
