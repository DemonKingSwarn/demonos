package trap

import (
	"demonos/drivers/kbd"
	"demonos/drivers/pic"
	"demonos/drivers/pit"
	"demonos/drivers/serial"
)

type Frame struct {
	R15     uint64
	R14     uint64
	R13     uint64
	R12     uint64
	R11     uint64
	R10     uint64
	R9      uint64
	R8      uint64
	RBP     uint64
	RDI     uint64
	RSI     uint64
	RDX     uint64
	RCX     uint64
	RBX     uint64
	RAX     uint64
	Vector  uint64
	ErrCode uint64
	RIP     uint64
	CS      uint64
	RFLAGS  uint64
	RSP     uint64
	SS      uint64
}

var exceptionNames = [32]string{
	"divide error",
	"debug",
	"nmi",
	"breakpoint",
	"overflow",
	"bound range exceeded",
	"invalid opcode",
	"device not available",
	"double fault",
	"coprocessor segment overrun",
	"invalid tss",
	"segment not present",
	"stack fault",
	"general protection",
	"page fault",
	"reserved",
	"x87 fpu error",
	"alignment check",
	"machine check",
	"simd fpu exception",
	"virtualisation exception",
	"reserved",
	"reserved",
	"reserved",
	"reserved",
	"reserved",
	"reserved",
	"reserved",
	"reserved",
	"reserved",
	"security exception",
	"reserved",
}

//go:nosplit
func trapHandler(f *Frame) {
	if f.Vector < 32 {
		panicTrap(f)
		return
	}
	switch f.Vector {
	case pic.IRQVector(0):
		pit.Tick()
		pic.EOI(0)
	case pic.IRQVector(1):
		kbd.Handle()
		pic.EOI(1)
	default:
		if f.Vector >= 0x20 && f.Vector < 0x30 {
			pic.EOI(uint8(f.Vector - 0x20))
		}
	}
}

//go:nosplit
func Handle(f *Frame) {
	trapHandler(f)
}

//go:nosplit
func panicTrap(f *Frame) {
	serial.WriteString("\r\nKERNEL PANIC: ")
	if f.Vector < 32 {
		serial.WriteString(exceptionNames[f.Vector])
	} else {
		serial.WriteString("unknown")
	}
	serial.WriteString(" (vec=")
	serial.WriteHex(f.Vector)
	serial.WriteString(" rip=")
	serial.WriteHex(f.RIP)
	serial.WriteString(" rsp=")
	serial.WriteHex(f.RSP)
	serial.WriteString(" err=")
	serial.WriteHex(f.ErrCode)
	serial.WriteString(")\r\n")
	for {
	}
}
