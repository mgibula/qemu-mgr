package qemu

import "log"

// -> { "execute": "query-pci" }
// <- { "return": [
//          {
//             "bus": 0,
//             "devices": [
//                {
//                   "bus": 0,
//                   "qdev_id": "",
//                   "slot": 0,
//                   "class_info": {
//                      "class": 1536,
//                      "desc": "Host bridge"
//                   },
//                   "id": {
//                      "device": 32902,
//                      "vendor": 4663
//                   },
//                   "function": 0,
//                   "regions": [
//                   ]
//                },

type QueryPciMessage struct {
	MonitorRequest
}

type ReadPci struct {
	MonitorResponse

	Return []struct {
		bus     uint `json:"bus"`
		Devices []struct {
			Bus       uint `json:"bus"`
			Slot      uint `json:"slot"`
			Function  uint `json:"function"`
			Irq       uint `json:"irq"`
			ClassInfo struct {
				Class uint   `json:"class"`
				Desc  string `json:"desc"`
			} `json:"class_info"`
			Id struct {
				Device     uint   `json:"device"`
				Vendor     uint   `json:"vendor"`
				DeviceDesc string `json:"device-desc"`
				VendorDesc string `json:"vendor-desc"`
			} `json:"id"`
		} `json:"devices"`
	} `json:"return"`
}

func (command *ReadPci) Execute(monitor *MonitorState) {
	request := QueryPciMessage{}
	request.Execute = "query-pci"
	monitor.SendJson(request)
}

func (command *ReadPci) ProcessResponse(monitor *MonitorState) {

	for i, bus := range command.Return {
		for j, device := range bus.Devices {
			command.Return[i].Devices[j].Id.DeviceDesc = GetDeviceName(device.Id.Device)
			command.Return[i].Devices[j].Id.VendorDesc = GetVendorName(device.Id.Vendor)
		}
	}

	log.Printf("[%v] PCI info %v)\n", monitor.Path, command.Return)
}
