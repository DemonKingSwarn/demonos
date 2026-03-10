package trap

import "unsafe"

func gdtLoad()
func idtLoad()
func setKernelStack(rsp uintptr)
func syscallEntryAsm()
func writeMSR(msr uint32, val uint64)
func readMSR(msr uint32) uint64
func cpuSetGSBase(base uintptr)
func cpuSetKernelGSBase(base uintptr)

const (
	msrEFER   = uint32(0xC0000080)
	msrSTAR   = uint32(0xC0000081)
	msrLSTAR  = uint32(0xC0000082)
	msrSFMASK = uint32(0xC0000084)
	eferSCE   = uint64(1 << 0)
)

func setupSyscall() {
	writeMSR(msrEFER, readMSR(msrEFER)|eferSCE)

	star := (uint64(0x08) << 32) | (uint64(0x18) << 48)
	writeMSR(msrSTAR, star)

	ep := *(*uint64)(unsafe.Pointer(&struct{ f func() }{syscallEntryAsm}))
	writeMSR(msrLSTAR, ep)

	writeMSR(msrSFMASK, 0x200)
}

func SetKernelStack(rsp uintptr) {
	setKernelStack(rsp)
}

func Init() {
	gdtLoad()
	idtLoad()
	setupSyscall()
}
