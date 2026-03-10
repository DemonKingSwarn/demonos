package main

import (
	"demonos/drivers/kbd"
	"demonos/drivers/pci"
	"demonos/drivers/pic"
	"demonos/drivers/pit"
	"demonos/drivers/serial"
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
	serial.Init()
	serial.WriteString("DemonOS kmain\n")

	con := vga.DefaultConsole
	con.Clear()

	con.FG = vga.LightCyan
	con.WriteString("DemonOS booting\n")
	con.FG = vga.White

	con.WriteString("  mb2 @ ")
	con.WriteHex(uint64(mbInfo))
	con.WriteString("\n")
	serial.WriteString("  mb2 @ ")
	serial.WriteHex(uint64(mbInfo))
	serial.WriteString("\n")

	mm.Init(mbInfo)
	con.WriteString("  mm OK  free=")
	con.WriteHex(mm.TotalFreeBytes())
	con.WriteString("\n")
	serial.WriteString("  mm OK\n")

	trap.Init()
	con.WriteString("  trap OK\n")
	serial.WriteString("  trap OK\n")

	pic.Init()
	pic.EnableIRQ(0)
	pic.EnableIRQ(1)
	pit.Init()
	con.WriteString("  pic/pit OK\n")
	serial.WriteString("  pic/pit OK\n")

	pci.Scan()
	con.WriteString("  pci OK  devs=")
	con.WriteHex(uint64(pci.Count))
	con.WriteString("\n")
	serial.WriteString("  pci OK  devs=")
	serial.WriteHex(uint64(pci.Count))
	serial.WriteString("\n")

	_ = kbd.Pending

	if len(initBin) > 0 {
		result, err := elf.Load(initBin)
		if err != nil {
			con.FG = vga.LightRed
			con.WriteString("  elf load failed\n")
			con.FG = vga.White
			serial.WriteString("  elf load failed\n")
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
		serial.WriteString("  launching init\n")

		t.Run()
	}

	con.FG = vga.LightGreen
	con.WriteString("DemonOS ready\n")
	con.FG = vga.White
	serial.WriteString("DemonOS ready\n")

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
