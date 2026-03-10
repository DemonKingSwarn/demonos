package trap

import (
	_ "unsafe"
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
	}
}

//go:nosplit
func Handle(f *Frame) {
	trapHandler(f)
}

//go:nosplit
func panicTrap(f *Frame) {
	name := "unknown"
	if f.Vector < 32 {
		name = exceptionNames[f.Vector]
	}
	_ = name
	for {
	}
}
