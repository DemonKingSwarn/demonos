package pio

//go:nosplit
func Outb(port uint16, val uint8)

//go:nosplit
func Inb(port uint16) uint8

//go:nosplit
func Outw(port uint16, val uint16)

//go:nosplit
func Inw(port uint16) uint16

//go:nosplit
func Outl(port uint16, val uint32)

//go:nosplit
func Inl(port uint16) uint32
