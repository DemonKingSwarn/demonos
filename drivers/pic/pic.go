package pic

import "demonos/arch/x86_64/pio"

const (
	pic1Cmd  = 0x20
	pic1Data = 0x21
	pic2Cmd  = 0xA0
	pic2Data = 0xA1

	icw1Init = 0x10
	icw1ICW4 = 0x01
	icw4x86  = 0x01
	picEOI   = 0x20

	irqBasemaster = 0x20
	irqBaseSlave  = 0x28
)

//go:nosplit
func Init() {
	pio.Outb(pic1Cmd, icw1Init|icw1ICW4)
	pio.Outb(pic2Cmd, icw1Init|icw1ICW4)

	pio.Outb(pic1Data, irqBasemaster)
	pio.Outb(pic2Data, irqBaseSlave)

	pio.Outb(pic1Data, 0x04)
	pio.Outb(pic2Data, 0x02)

	pio.Outb(pic1Data, icw4x86)
	pio.Outb(pic2Data, icw4x86)

	pio.Outb(pic1Data, 0xFF)
	pio.Outb(pic2Data, 0xFF)
}

//go:nosplit
func EnableIRQ(irq uint8) {
	if irq < 8 {
		mask := pio.Inb(pic1Data) &^ (1 << irq)
		pio.Outb(pic1Data, mask)
	} else {
		mask := pio.Inb(pic2Data) &^ (1 << (irq - 8))
		pio.Outb(pic2Data, mask)
		mask = pio.Inb(pic1Data) &^ (1 << 2)
		pio.Outb(pic1Data, mask)
	}
}

//go:nosplit
func DisableIRQ(irq uint8) {
	if irq < 8 {
		mask := pio.Inb(pic1Data) | (1 << irq)
		pio.Outb(pic1Data, mask)
	} else {
		mask := pio.Inb(pic2Data) | (1 << (irq - 8))
		pio.Outb(pic2Data, mask)
	}
}

//go:nosplit
func EOI(irq uint8) {
	if irq >= 8 {
		pio.Outb(pic2Cmd, picEOI)
	}
	pio.Outb(pic1Cmd, picEOI)
}

//go:nosplit
func IRQVector(irq uint8) uint64 {
	if irq < 8 {
		return uint64(irqBasemaster) + uint64(irq)
	}
	return uint64(irqBaseSlave) + uint64(irq-8)
}
