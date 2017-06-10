package qemu

import "log"
import "fmt"

// {"QMP": {"version": {"qemu": {"micro": 0, "minor": 9, "major": 2}, "package": ""}, "capabilities": []}}

type ReadQmpBanner struct {
	QMP struct {
		Version struct {
			Qemu struct {
				Micro int `json:"micro"`
				Minor int `json:"minor"`
				Major int `json:"major"`
			} `json:"qemu"`
		} `json:"version"`
	} `json:"QMP"`
}

func (response *ReadQmpBanner) Execute(monitor *MonitorState) {

}

func (response *ReadQmpBanner) ProcessResponse(monitor *MonitorState) {
	log.Printf("[%v] Connected to QEMU %v.%v.%v\n", monitor.Path, response.QMP.Version.Qemu.Major, response.QMP.Version.Qemu.Minor, response.QMP.Version.Qemu.Micro)
	monitor.Instance.QemuVersion = fmt.Sprintf("%v.%v.%v", response.QMP.Version.Qemu.Major, response.QMP.Version.Qemu.Minor, response.QMP.Version.Qemu.Micro)
}
