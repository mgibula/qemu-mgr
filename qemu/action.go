package qemu

type Action interface {
	Execute(args []string)
}
