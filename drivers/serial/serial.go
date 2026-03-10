package serial

import "demonos/arch/x86_64/pio"

const (
	com1 = 0x3F8

	regData  = com1 + 0
	regIER   = com1 + 1
	regFCR   = com1 + 2
	regLCR   = com1 + 3
	regMCR   = com1 + 4
	regLSR   = com1 + 5
	regDivLo = com1 + 0
	regDivHi = com1 + 1

	lsrTHRE = 1 << 5
	lsrDR   = 1 << 0

	lcrDLAB       = 1 << 7
	lcr8N1  uint8 = 0x03
	divisor       = 1
)

//go:nosplit
func Init() {
	pio.Outb(regIER, 0x00)
	pio.Outb(regLCR, lcrDLAB)
	pio.Outb(regDivLo, divisor)
	pio.Outb(regDivHi, 0x00)
	pio.Outb(regLCR, lcr8N1)
	pio.Outb(regFCR, 0xC7)
	pio.Outb(regMCR, 0x0B)
	pio.Outb(regIER, 0x00)
}

//go:nosplit
func WriteChar(c byte) {
	for pio.Inb(regLSR)&lsrTHRE == 0 {
	}
	pio.Outb(regData, c)
}

//go:nosplit
func WriteString(s string) {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			WriteChar('\r')
		}
		WriteChar(s[i])
	}
}

//go:nosplit
func WriteHex(v uint64) {
	const digits = "0123456789abcdef"
	WriteChar('0')
	WriteChar('x')
	for i := 60; i >= 0; i -= 4 {
		WriteChar(digits[(v>>uint(i))&0xf])
	}
}

//go:nosplit
func ReadChar() (byte, bool) {
	if pio.Inb(regLSR)&lsrDR == 0 {
		return 0, false
	}
	return pio.Inb(regData), true
}
