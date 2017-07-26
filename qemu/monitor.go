package qemu

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// MonitorAction struct provide global structure holding
// information of all monitored qemu instances
type MonitorAction struct {
	Directory string
	terminate chan os.Signal
	sockets   map[string]*MonitorState
	sync.Mutex
}

// Execute starts periodic monitoring of provided directory
// For each socket found, a new InstanceMonitor is launch
func (action *MonitorAction) Execute(args []string) {
	action.readParameters(args)

	log.Println("Starting monitoring of directory " + action.Directory)
	log.Println("Send SIGINT (Ctrl-C) to interrupt")

	action.terminate = make(chan os.Signal)
	action.sockets = make(map[string]*MonitorState)
	signal.Notify(action.terminate, syscall.SIGINT, syscall.SIGTERM)

	action.scanDirectory()

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			action.scanDirectory()
		}
	}()

	go action.startHttpServer()

	<-action.terminate
	ticker.Stop()
}

func (action *MonitorAction) readParameters(args []string) {
	monitorCommand := flag.NewFlagSet("monitor", flag.ExitOnError)
	monitorCommand.StringVar(&action.Directory, "directory", "", "directory with monitor sockets")

	monitorCommand.Parse(args)

	if len(action.Directory) == 0 {
		log.Fatalln("Error: nonempty -directory parameter is required")
	}
}

func (action *MonitorAction) startHttpServer() {
	http.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		action.listInstances(w, r)
	})
	log.Fatalln(http.ListenAndServe("127.0.0.1:8080", nil))
}

func (action *MonitorAction) listInstances(w http.ResponseWriter, r *http.Request) {
	action.Lock()
	defer action.Unlock()

	var list ListInstances

	for _, m := range action.sockets {
		if m.started && !m.ended {
			m.Lock()
			list.Instances = append(list.Instances, m.Instance)
			m.Unlock()
		}
	}

	response, _ := json.MarshalIndent(list, "", "    ")
	fmt.Fprint(w, string(response))
}

func (action *MonitorAction) scanDirectory() {
	action.Lock()
	defer action.Unlock()

	files, _ := filepath.Glob(action.Directory)

	for _, file := range files {
		if _, exists := action.sockets[file]; !exists {
			log.Println("Adding monitor to socket " + file)
			action.sockets[file] = &MonitorState{started: false, ended: false, Path: file, control: make(chan int), input: make(chan string)}
			go action.sockets[file].connect()
		}
	}

	for file := range action.sockets {
		if action.sockets[file].ended {
			log.Println("Removing monitor from socket " + file)
			delete(action.sockets, file)
		}
	}
}
