package main

import (
	"demonos/drivers/vga"
	"demonos/kernel/elf"
	"demonos/kernel/mm"
	"demonos/kernel/proc"
	"demonos/kernel/syscall"
	"demonos/kernel/trap"
	"unsafe"
)

var handlers [2]unsafe.Pointer

func main() {
	th := trapHandler
	sh := syscallHandler
	handlers[0] = *(*unsafe.Pointer)(unsafe.Pointer(&th))
	handlers[1] = *(*unsafe.Pointer)(unsafe.Pointer(&sh))
	kmain(0)
}

var initBin []byte

//go:nosplit
//go:noinline
func kmain(mbInfo uint32) {
	con := vga.DefaultConsole
	con.Clear()

	con.FG = vga.LightCyan
	con.WriteString("DemonOS booting\n")
	con.FG = vga.White

	con.WriteString("  mb2 @ ")
	con.WriteHex(uint64(mbInfo))
	con.WriteString("\n")

	mm.Init(mbInfo)
	con.WriteString("  mm OK  free=")
	con.WriteHex(mm.TotalFreeBytes())
	con.WriteString("\n")

	trap.Init()
	con.WriteString("  trap OK\n")

	if len(initBin) > 0 {
		result, err := elf.Load(initBin)
		if err != nil {
			con.FG = vga.LightRed
			con.WriteString("  elf load failed\n")
			con.FG = vga.White
			halt()
		}
		con.WriteString("  elf loaded entry=")
		con.WriteHex(uint64(result.Entry))
		con.WriteString("\n")

		t := proc.NewTask(result.Entry, result.StackTop)
		trap.SetKernelStack(t.KernelStack)

		con.FG = vga.LightGreen
		con.WriteString("  launching init\n")
		con.FG = vga.White

		t.Run()
	}

	con.FG = vga.LightGreen
	con.WriteString("DemonOS ready\n")
	con.FG = vga.White

	halt()
}

//go:nosplit
//go:noinline
func trapHandler(f *trap.Frame) {
	trap.Handle(f)
}

//go:nosplit
//go:noinline
func syscallHandler(nr, a1, a2, a3 uint64) uint64 {
	return syscall.Handle(nr, a1, a2, a3)
}
