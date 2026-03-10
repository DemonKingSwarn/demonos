package proc

import "demonos/kernel/mm"

const kernelStackSize = 16 * mm.PageSize

type State int

const (
	StateReady State = iota
	StateRunning
	StateZombie
)

type Task struct {
	PID         uint64
	State       State
	KernelStack uintptr
	UserEntry   uintptr
	UserStack   uintptr
}

var nextPID uint64 = 1

func NewTask(entry, userStack uintptr) *Task {
	kstack := allocKernelStack()
	t := &Task{
		PID:         nextPID,
		State:       StateReady,
		KernelStack: kstack,
		UserEntry:   entry,
		UserStack:   userStack,
	}
	nextPID++
	return t
}

func allocKernelStack() uintptr {
	base := uintptr(0)
	pages := kernelStackSize / mm.PageSize
	for i := 0; i < pages; i++ {
		p := mm.AllocPage()
		if base == 0 {
			base = p
		}
	}
	return base + uintptr(kernelStackSize)
}
