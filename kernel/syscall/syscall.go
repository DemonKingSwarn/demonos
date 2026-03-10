package syscall

import (
	"demonos/drivers/kbd"
	"demonos/drivers/vga"
	"unsafe"
)

const (
	SysRead      = 0
	SysWrite     = 1
	SysOpen      = 2
	SysClose     = 3
	SysExit      = 60
	SysExitGroup = 231
)

const (
	ENOSYS = 38
	EBADF  = 9
	EFAULT = 14
	EAGAIN = 11
)

const (
	FdStdin  = 0
	FdStdout = 1
	FdStderr = 2
)

//go:nosplit
func syscallHandler(nr, a1, a2, a3 uint64) uint64 {
	switch nr {
	case SysRead:
		return sysRead(a1, uintptr(a2), a3)
	case SysWrite:
		return sysWrite(a1, uintptr(a2), a3)
	case SysExit, SysExitGroup:
		sysExit(int(a1))
	}
	return ^uint64(ENOSYS - 1)
}

//go:nosplit
func Handle(nr, a1, a2, a3 uint64) uint64 {
	return syscallHandler(nr, a1, a2, a3)
}

//go:nosplit
func sysRead(fd uint64, buf uintptr, count uint64) uint64 {
	if fd != FdStdin {
		return ^uint64(EBADF - 1)
	}
	if count == 0 {
		return 0
	}
	b := unsafe.Slice((*byte)(unsafe.Pointer(buf)), count)
	n := uint64(0)
	for n < count {
		ch, ok := kbd.Read()
		if !ok {
			if n == 0 {
				return ^uint64(EAGAIN - 1)
			}
			break
		}
		b[n] = ch
		n++
		if ch == '\n' {
			break
		}
	}
	return n
}

//go:nosplit
func sysWrite(fd uint64, buf uintptr, count uint64) uint64 {
	if fd != FdStdout && fd != FdStderr {
		return ^uint64(EBADF - 1)
	}
	if count == 0 {
		return 0
	}
	con := vga.DefaultConsole
	b := unsafe.Slice((*byte)(unsafe.Pointer(buf)), count)
	for _, c := range b {
		con.PutChar(c)
	}
	return count
}

//go:nosplit
func sysExit(code int) {
	_ = code
	for {
	}
}
