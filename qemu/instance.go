package qemu

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Instance struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	MonitorPath string `json:"monitor-path"`
	QemuVersion string `json:"qemu-version"`
	Cpus        uint   `json:"cpus"`
	CpuTopology struct {
		Threads uint `json:"threads"`
		Cores   uint `json:"cores"`
		Sockets uint `json:"sockets"`
	} `json:"cpu-topology"`
	Memory struct {
		Current uint `json:"current"`
		Max     uint `json:"max"`
	}
	MemorySlots struct {
		Available uint `json:"slots"`
	}
	Paused     bool   `json:"paused"`
	KvmEnabled bool   `json:"kvm_enabled"`
	Machine    string `json:"machine-type"`
	Pidfile    string `json:"-"`
}

func (instance *Instance) ConstructCommandLineArgs() []string {
	var cmd []string

	cmd = append(cmd, "-daemonize")
	cmd = append(cmd, "-no-user-config")
	cmd = append(cmd, "-nodefaults")
	cmd = append(cmd, "-no-shutdown")
	cmd = append(cmd, "-display")
	cmd = append(cmd, "none")
	cmd = append(cmd, "-msg")
	cmd = append(cmd, "timestamp=on")
	cmd = append(cmd, "-name")
	cmd = append(cmd, instance.Name)
	cmd = append(cmd, "-pidfile")
	cmd = append(cmd, instance.Pidfile)
	cmd = append(cmd, "-machine")
	cmd = append(cmd, instance.constructMachineConfig())
	cmd = append(cmd, "-smp")
	cmd = append(cmd, instance.constructSmpConfig())
	cmd = append(cmd, "-m")
	cmd = append(cmd, instance.constructMemoryConfig())

	if len(instance.MonitorPath) != 0 {
		cmd = append(cmd, "-chardev")
		cmd = append(cmd, fmt.Sprintf("socket,id=charmonitor,path=%s,server,nowait", instance.MonitorPath))
		cmd = append(cmd, "-mon")
		cmd = append(cmd, "chardev=charmonitor,id=monitor,mode=control")
	}

	return cmd
}

func (instance *Instance) constructMachineConfig() string {
	var cmd string

	cmd = fmt.Sprintf("type=%s", instance.Machine)
	return cmd
}

func (instance *Instance) constructSmpConfig() string {
	var cmd string

	cmd = fmt.Sprintf("cpus=%d", instance.Cpus)
	if instance.CpuTopology.Threads > 0 || instance.CpuTopology.Cores > 0 || instance.CpuTopology.Sockets > 0 {
		if instance.CpuTopology.Threads == 0 || instance.CpuTopology.Cores == 0 || instance.CpuTopology.Sockets == 0 {
			log.Fatalln("Params '-threads', '-cores' and '-sockets' must be specified if one of it is")
		}

		if instance.CpuTopology.Sockets*instance.CpuTopology.Cores*instance.CpuTopology.Threads < instance.Cpus {
			log.Fatalln("'-threads', '-cores' and '-sockets' multiplied must by greater than or equal '-cpus'")
		}

		cmd += fmt.Sprintf(",maxcpus=%d", instance.CpuTopology.Sockets*instance.CpuTopology.Cores*instance.CpuTopology.Threads)
		cmd += fmt.Sprintf(",cores=%d", instance.CpuTopology.Cores)
		cmd += fmt.Sprintf(",threads=%d", instance.CpuTopology.Threads)
		cmd += fmt.Sprintf(",sockets=%d", instance.CpuTopology.Sockets)
	}

	return cmd
}

func (instance *Instance) constructMemoryConfig() string {
	var cmd string

	cmd = fmt.Sprintf("size=%dM", instance.Memory.Current)
	if instance.Memory.Max > 0 || instance.MemorySlots.Available > 0 {
		if instance.Memory.Max == 0 || instance.MemorySlots.Available == 0 {
			log.Fatalln("Both '-max-memory' and '-memory-slots' must be defined")
		}

		cmd += fmt.Sprintf(",maxmem=%dM", instance.Memory.Max)
		cmd += fmt.Sprintf(",slots=%d", instance.MemorySlots.Available)
	}

	return cmd
}

// Since not everything can be read from QMP, we need to parse
// cmdline to read what we can
func (instance *Instance) ParseCommandLine(cmdline string) {
	parts := strings.Split(cmdline, "\x00")

	// Try to determine CPU topology
	smp_pos := findOption(parts, "-smp")
	if smp_pos != 0 {
		instance.readCpuTopology(parts[smp_pos+1])
	}
}

func (instance *Instance) readCpuTopology(options string) {
	_, cores := findInteger("cores=([0-9]+)", options)
	_, threads := findInteger("threads=([0-9]+)", options)
	_, sockets := findInteger("sockets=([0-9]+)", options)

	if cores == 0 && threads == 0 && sockets == 0 {
		return
	}

	log.Printf("[%v] Found CPU topology %v sockets, %v cores, %v threads\n", instance.MonitorPath, sockets, cores, threads)
	instance.CpuTopology.Threads = uint(threads)
	instance.CpuTopology.Cores = uint(cores)
	instance.CpuTopology.Sockets = uint(sockets)
}

func findInteger(expression string, line string) (bool, int) {
	re := regexp.MustCompile(expression)
	found := re.FindStringSubmatch(line)

	if len(found) < 2 {
		return false, 0
	}

	retval, err := strconv.Atoi(found[1])
	if err != nil {
		return false, 0
	}

	return true, retval
}

func findOption(parts []string, part string) int {
	for i, v := range parts {
		if v == part {
			return i
		}
	}

	// 0 is safe, because first element of cmdline is binary itself
	return 0
}
