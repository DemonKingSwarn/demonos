package vga

import "unsafe"

const (
	Width  = 80
	Height = 25
	Base   = uintptr(0xB8000)
)

type Colour uint8

const (
	Black        Colour = 0
	Blue         Colour = 1
	Green        Colour = 2
	Cyan         Colour = 3
	Red          Colour = 4
	Magenta      Colour = 5
	Brown        Colour = 6
	LightGrey    Colour = 7
	DarkGrey     Colour = 8
	LightBlue    Colour = 9
	LightGreen   Colour = 10
	LightCyan    Colour = 11
	LightRed     Colour = 12
	LightMagenta Colour = 13
	Yellow       Colour = 14
	White        Colour = 15
)

//go:nosplit
func entry(ch byte, fg, bg Colour) uint16 {
	attr := uint16(bg)<<12 | uint16(fg)<<8
	return attr | uint16(ch)
}

//go:nosplit
func cell(col, row int) *uint16 {
	offset := uintptr(row*Width+col) * 2
	return (*uint16)(unsafe.Pointer(Base + offset))
}

type Console struct {
	Col, Row int
	FG, BG   Colour
}

var DefaultConsole = &Console{FG: White, BG: Black}

//go:nosplit
func (c *Console) Clear() {
	blank := entry(' ', c.FG, c.BG)
	for row := 0; row < Height; row++ {
		for col := 0; col < Width; col++ {
			*cell(col, row) = blank
		}
	}
	c.Col, c.Row = 0, 0
}

//go:nosplit
func (c *Console) scroll() {
	blank := entry(' ', c.FG, c.BG)
	for row := 1; row < Height; row++ {
		for col := 0; col < Width; col++ {
			*cell(col, row-1) = *cell(col, row)
		}
	}
	for col := 0; col < Width; col++ {
		*cell(col, Height-1) = blank
	}
	c.Row = Height - 1
}

//go:nosplit
func (c *Console) PutChar(ch byte) {
	switch ch {
	case '\n':
		c.Col = 0
		c.Row++
	case '\r':
		c.Col = 0
	case '\t':
		c.Col = (c.Col + 8) &^ 7
		if c.Col >= Width {
			c.Col = 0
			c.Row++
		}
	default:
		*cell(c.Col, c.Row) = entry(ch, c.FG, c.BG)
		c.Col++
		if c.Col >= Width {
			c.Col = 0
			c.Row++
		}
	}
	if c.Row >= Height {
		c.scroll()
	}
}

//go:nosplit
func (c *Console) WriteString(s string) {
	for i := 0; i < len(s); i++ {
		c.PutChar(s[i])
	}
}

//go:nosplit
func (c *Console) WriteHex(v uint64) {
	const digits = "0123456789abcdef"
	c.WriteString("0x")
	for i := 60; i >= 0; i -= 4 {
		c.PutChar(digits[(v>>uint(i))&0xf])
	}
}
