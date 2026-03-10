package pit

import "demonos/arch/x86_64/pio"

const (
	pitChan0 = 0x40
	pitCmd   = 0x43

	pitFreq    = 1193182
	targetHz   = 100
	pitDivisor = pitFreq / targetHz

	cmdChan0     = 0x00
	cmdAccessLow = 0x30
	cmdModeRate  = 0x04
	cmdBinary    = 0x00
)

var Ticks uint64

//go:nosplit
func Init() {
	cmd := uint8(cmdChan0 | cmdAccessLow | cmdModeRate | cmdBinary)
	pio.Outb(pitCmd, cmd)
	pio.Outb(pitChan0, uint8(pitDivisor&0xFF))
	pio.Outb(pitChan0, uint8(pitDivisor>>8))
}

//go:nosplit
func Tick() {
	Ticks++
}
