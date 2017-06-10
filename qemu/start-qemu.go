package qemu

import "log"
import "os"
import "os/exec"
import "flag"

var qemuPath string

type StartAction struct {
}

func (action *StartAction) readParameters(args []string) Instance {
	var instance Instance

	startCommand := flag.NewFlagSet("create", flag.ExitOnError)

	startCommand.StringVar(&qemuPath, "qemu-path", "", "path to QEMU binary")
	startCommand.StringVar(&instance.Name, "name", "", "instance name")
	startCommand.BoolVar(&instance.Paused, "paused", false, "instance paused at start")
	startCommand.UintVar(&instance.Cpus, "cpus", 1, "CPU count")
	startCommand.UintVar(&instance.CpuTopology.Threads, "threads", 0, "CPU topology - threads per core")
	startCommand.UintVar(&instance.CpuTopology.Cores, "cores", 0, "CPU topology - cores per socket")
	startCommand.UintVar(&instance.CpuTopology.Sockets, "sockets", 0, "CPU topology - sockets number")
	startCommand.UintVar(&instance.Memory.Current, "memory", 512, "instance memory (in MB)")
	startCommand.UintVar(&instance.Memory.Max, "max-memory", 0, "maximum instance memory (in MB)")
	startCommand.UintVar(&instance.MemorySlots.Available, "memory-slots", 0, "additional memory slots")
	startCommand.BoolVar(&instance.KvmEnabled, "kvm", false, "use KVM acceleration")
	startCommand.StringVar(&instance.Machine, "machine-type", "pc", "machine type")
	startCommand.StringVar(&instance.MonitorPath, "monitor", "", "monitor socket path")
	startCommand.StringVar(&instance.Pidfile, "pidfile", "", "PID file")

	startCommand.Parse(args)

	if len(instance.Name) == 0 {
		log.Fatalln("Error: nonempty -name parameter is required")
	}

	return instance
}

func (action *StartAction) Execute(args []string) {
	instance := action.readParameters(args)

	if len(qemuPath) == 0 {
		searchForQemu()
	}

	var cmd exec.Cmd
	cmd.Path = qemuPath
	cmd.Args = append(cmd.Args, qemuPath)
	cmd.Args = append(cmd.Args, instance.ConstructCommandLineArgs()...)

	log.Printf("Starting QEMU as %v\n", cmd.Args)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error while starting QEMU: %v\n", err)
	}
}

func searchForQemu() {
	guesses := []string{
		"/usr/bin/qemu-system-x86_64",
		"/usr/bin/qemu-kvm",
	}

	for _, path := range guesses {
		_, err := os.Stat(path)
		if err == nil {
			qemuPath = path
			log.Printf("Found QEMU in %v\n", path)
			return
		}
	}

	log.Fatalln("Can't find qemu binary (try using -qemu-path)")
}
