package pci

import "demonos/arch/x86_64/pio"

const (
	pciAddrPort = 0xCF8
	pciDataPort = 0xCFC

	maxBuses     = 256
	maxDevices   = 32
	maxFunctions = 8

	regVendor     = 0x00
	regClass      = 0x08
	regHeaderType = 0x0E
)

type Device struct {
	Bus      uint8
	Dev      uint8
	Func     uint8
	VendorID uint16
	DeviceID uint16
	Class    uint8
	Subclass uint8
	ProgIF   uint8
}

const maxDevList = 64

var (
	devices [maxDevList]Device
	Count   int
)

//go:nosplit
func configAddr(bus, dev, fn, offset uint8) uint32 {
	return 0x80000000 |
		uint32(bus)<<16 |
		uint32(dev)<<11 |
		uint32(fn)<<8 |
		uint32(offset&0xFC)
}

//go:nosplit
func ReadL(bus, dev, fn, offset uint8) uint32 {
	pio.Outl(pciAddrPort, configAddr(bus, dev, fn, offset))
	return pio.Inl(pciDataPort)
}

//go:nosplit
func ReadW(bus, dev, fn, offset uint8) uint16 {
	pio.Outl(pciAddrPort, configAddr(bus, dev, fn, offset))
	v := pio.Inl(pciDataPort)
	return uint16(v >> ((offset & 2) * 8))
}

//go:nosplit
func ReadB(bus, dev, fn, offset uint8) uint8 {
	pio.Outl(pciAddrPort, configAddr(bus, dev, fn, offset))
	v := pio.Inl(pciDataPort)
	return uint8(v >> ((offset & 3) * 8))
}

//go:nosplit
func Scan() {
	Count = 0
	for bus := 0; bus < maxBuses; bus++ {
		for dev := 0; dev < maxDevices; dev++ {
			scanDevice(uint8(bus), uint8(dev))
		}
	}
}

//go:nosplit
func scanDevice(bus, dev uint8) {
	v := ReadL(bus, dev, 0, regVendor)
	if v == 0xFFFFFFFF {
		return
	}
	ht := ReadB(bus, dev, 0, regHeaderType)
	funcs := uint8(1)
	if ht&0x80 != 0 {
		funcs = maxFunctions
	}
	for fn := uint8(0); fn < funcs; fn++ {
		v = ReadL(bus, dev, fn, regVendor)
		if v == 0xFFFFFFFF {
			continue
		}
		if Count >= maxDevList {
			return
		}
		cls := ReadL(bus, dev, fn, regClass)
		devices[Count] = Device{
			Bus:      bus,
			Dev:      dev,
			Func:     fn,
			VendorID: uint16(v),
			DeviceID: uint16(v >> 16),
			Class:    uint8(cls >> 24),
			Subclass: uint8(cls >> 16),
			ProgIF:   uint8(cls >> 8),
		}
		Count++
	}
}

//go:nosplit
func Get(i int) Device {
	return devices[i]
}
