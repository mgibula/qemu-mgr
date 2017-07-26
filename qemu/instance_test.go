package qemu_test

import (
	"fmt"
	"github.com/mgibula/qemu-mgr/qemu"
)

func ExampleInstance1() {
	var i qemu.Instance
	i.Name = "test1"
	i.Cpus = 2
	i.Memory.Current = 128

	fmt.Println(i.ConstructCommandLineArgs())
	// Output: [-daemonize -no-user-config -nodefaults -no-shutdown -display none -msg timestamp=on -name test1 -smp cpus=2 -m size=128M]
}

func ExampleInstance2() {
	var i qemu.Instance
	i.Name = "test2"
	i.Cpus = 2
	i.Memory.Current = 128
	i.Pidfile = "/tmp/pid.file"

	fmt.Println(i.ConstructCommandLineArgs())
	// Output: [-daemonize -no-user-config -nodefaults -no-shutdown -display none -msg timestamp=on -name test2 -smp cpus=2 -m size=128M -pidfile /tmp/pid.file]
}

func ExampleInstance3() {
	var i qemu.Instance
	i.Name = "test2"
	i.Cpus = 2
	i.Memory.Current = 128
	i.MonitorPath = "/tmp/vm.sock"

	fmt.Println(i.ConstructCommandLineArgs())
	// Output: [-daemonize -no-user-config -nodefaults -no-shutdown -display none -msg timestamp=on -name test2 -smp cpus=2 -m size=128M -chardev socket,id=charmonitor,path=/tmp/vm.sock,server,nowait -mon chardev=charmonitor,id=monitor,mode=control]
}

func ExampleInstance4() {
	var i qemu.Instance
	i.Name = "test1"
	i.Cpus = 2
	i.Memory.Current = 128
	i.Machine = "pc-i440fx-2.3"

	fmt.Println(i.ConstructCommandLineArgs())
	// Output: [-daemonize -no-user-config -nodefaults -no-shutdown -display none -msg timestamp=on -name test1 -machine type=pc-i440fx-2.3 -smp cpus=2 -m size=128M]
}

func ExampleInstance5() {
	var i qemu.Instance
	i.Name = "test1"
	i.Cpus = 2
	i.CpuTopology.Threads = 2
	i.CpuTopology.Cores = 1
	i.CpuTopology.Sockets = 1
	i.Memory.Current = 128

	fmt.Println(i.ConstructCommandLineArgs())
	// Output: [-daemonize -no-user-config -nodefaults -no-shutdown -display none -msg timestamp=on -name test1 -smp cpus=2,maxcpus=2,cores=1,threads=2,sockets=1 -m size=128M]
}

func ExampleInstance6() {
	var i qemu.Instance
	i.Name = "test1"
	i.Cpus = 4
	i.CpuTopology.Threads = 2
	i.CpuTopology.Cores = 2
	i.CpuTopology.Sockets = 1
	i.Memory.Current = 128

	fmt.Println(i.ConstructCommandLineArgs())
	// Output: [-daemonize -no-user-config -nodefaults -no-shutdown -display none -msg timestamp=on -name test1 -smp cpus=4,maxcpus=4,cores=2,threads=2,sockets=1 -m size=128M]
}

func ExampleInstance7() {
	var i qemu.Instance
	i.Name = "test1"
	i.Cpus = 8
	i.CpuTopology.Threads = 2
	i.CpuTopology.Cores = 2
	i.CpuTopology.Sockets = 2
	i.Memory.Current = 128

	fmt.Println(i.ConstructCommandLineArgs())
	// Output: [-daemonize -no-user-config -nodefaults -no-shutdown -display none -msg timestamp=on -name test1 -smp cpus=8,maxcpus=8,cores=2,threads=2,sockets=2 -m size=128M]
}

func ExampleInstance8() {
	var i qemu.Instance
	i.Name = "test1"
	i.Cpus = 2
	i.CpuTopology.Threads = 2
	i.CpuTopology.Cores = 2
	i.CpuTopology.Sockets = 2
	i.Memory.Current = 128

	fmt.Println(i.ConstructCommandLineArgs())
	// Output: [-daemonize -no-user-config -nodefaults -no-shutdown -display none -msg timestamp=on -name test1 -smp cpus=2,maxcpus=8,cores=2,threads=2,sockets=2 -m size=128M]
}

func ExampleInstance9() {
	var i qemu.Instance
	i.Name = "test1"
	i.Cpus = 2
	i.Memory.Current = 128
	i.Memory.Max = 1024
	i.MemorySlots.Available = 10

	fmt.Println(i.ConstructCommandLineArgs())
	// Output: [-daemonize -no-user-config -nodefaults -no-shutdown -display none -msg timestamp=on -name test1 -smp cpus=2 -m size=128M,maxmem=1024M,slots=10]
}
