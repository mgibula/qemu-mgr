package main

import "github.com/mgibula/qemu-mgr/qemu"
import "log"
import "os"

func main() {
	var action qemu.Action

	if len(os.Args) == 1 {
		log.Fatalln("Parameter 'mode' is required")
	}

	switch os.Args[1] {
	case "create":
		action = new(qemu.StartAction)
	case "monitor":
		action = new(qemu.MonitorAction)
	default:
		log.Fatalf("Mode '%v' is invalid\n", os.Args[1])
	}

	action.Execute(os.Args[2:])
}
