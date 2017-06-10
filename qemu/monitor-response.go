package qemu

type MonitorRequest struct {
	Execute string `json:"execute"`
}

type MonitorCommand interface {
	Execute(monitor *MonitorState)
	ProcessResponse(monitor *MonitorState)
}

type MonitorResponse struct {
	Id string `json:"id"`

	// Errors
	Error struct {
		Class string `json:"class"`
		Desc  string `json:"desc"`
	} `json:"error"`

	// Asynchronous events
	Event     string `json:"event"`
	Timestamp struct {
		Seconds      int `json:"seconds"`
		Microseconds int `json:"microseconds"`
	} `json:"timestamp"`
}
