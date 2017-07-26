package qemu

import (
	"encoding/json"
	"log"
	"net"
	"sync"
)

const (
	MONITOR_CLOSE = iota
)

type MonitorState struct {
	sync.Mutex
	Path       string
	Instance   Instance
	control    chan int
	input      chan string
	started    bool
	ended      bool
	connection net.Conn
	command    MonitorCommand
	queue      []MonitorCommand
}

func (monitor *MonitorState) SendJson(msg interface{}) {
	b, _ := json.Marshal(msg)
	monitor.send(string(b))
}

func (monitor *MonitorState) connect() {
	monitor.Lock()
	defer monitor.Unlock()

	log.Println("Started monitoring socket " + monitor.Path)
	defer func() {
		monitor.ended = true
	}()

	c, err := net.Dial("unix", monitor.Path)
	monitor.connection = c
	monitor.Instance.MonitorPath = monitor.Path
	if err != nil {
		log.Println("Can't connect to socket " + monitor.Path)
		return
	}

	defer monitor.connection.Close()
	log.Println("Connected to " + monitor.Path)

	go monitor.reader()

	monitor.queueCommand(&ReadQmpBanner{})
	monitor.queueCommand(&ReadQmpCapabilities{})
	monitor.queueCommand(&ReadName{})
	monitor.queueCommand(&ReadUuid{})
	monitor.queueCommand(&ReadKvm{})
	monitor.queueCommand(&ReadStatus{})
	monitor.queueCommand(&ReadCpus{ParseCmdLine: true})
	monitor.queueCommand(&ReadPci{})

	monitor.started = true
	monitor.Unlock()

loop:
	for {
		if len(monitor.queue) != 0 {
			monitor.command = monitor.queue[0]
			monitor.queue = monitor.queue[1:]

			monitor.Lock()
			monitor.command.Execute(monitor)
			monitor.Unlock()
		}

		select {
		case msg := <-monitor.input:
			monitor.Lock()
			monitor.processMessage(msg)
			monitor.Unlock()
		case command := <-monitor.control:
			if command == MONITOR_CLOSE {
				break loop
			}
		}
	}

	monitor.Lock()
	log.Println("Ended monitoring socket " + monitor.Path)
}

func (monitor *MonitorState) processMessage(msg string) {
	var response MonitorResponse
	err := json.Unmarshal([]byte(msg), &response)
	if err != nil {
		monitor.control <- MONITOR_CLOSE
		return
	}

	if len(response.Error.Class) != 0 {
		// Command error
		log.Printf("Received error %v\n", response.Error)
	} else if len(response.Event) != 0 {
		// Asynchronouse notification
		log.Printf("Received notification %v\n", response.Event)
	} else if monitor.command != nil {
		// We're waiting for reply
		err := json.Unmarshal([]byte(msg), &monitor.command)
		if err != nil {
			monitor.control <- MONITOR_CLOSE
			return
		}
		monitor.command.ProcessResponse(monitor)
		monitor.command = nil
	} else {
		// No idea
		log.Printf("Received something %v", msg)
	}
}

func (monitor *MonitorState) queueCommand(command MonitorCommand) {
	monitor.queue = append(monitor.queue, command)
}

func (monitor *MonitorState) reader() {
	buf := make([]byte, 1024*64)
	for {
		n, err := monitor.connection.Read(buf[:])
		if err != nil {
			log.Println("Connection error for " + monitor.Path)
			log.Println(err)
			monitor.control <- MONITOR_CLOSE
			return
		}

		monitor.input <- string(buf[0:n])
	}
}

func (monitor *MonitorState) send(msg string) {
	_, err := monitor.connection.Write([]byte(msg))
	if err != nil {
		log.Println("Write error for " + monitor.Path)
		log.Println(err)
		return
	}
}
