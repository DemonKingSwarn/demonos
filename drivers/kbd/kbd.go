package kbd

import "demonos/arch/x86_64/pio"

const (
	kbdData   = 0x60
	kbdStatus = 0x64

	bufSize = 256
)

var (
	buf  [bufSize]byte
	head uint64
	tail uint64
)

var scancodeSet1 = [128]byte{
	0, 0x1B, '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-', '=', '\b',
	'\t', 'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', '[', ']', '\n',
	0, 'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l', ';', '\'', '`',
	0, '\\', 'z', 'x', 'c', 'v', 'b', 'n', 'm', ',', '.', '/', 0,
	'*', 0, ' ', 0,
}

//go:nosplit
func Handle() {
	sc := pio.Inb(kbdData)
	if sc&0x80 != 0 {
		return
	}
	if int(sc) >= len(scancodeSet1) {
		return
	}
	ch := scancodeSet1[sc]
	if ch == 0 {
		return
	}
	next := (tail + 1) % bufSize
	if next != head {
		buf[tail] = ch
		tail = next
	}
}

//go:nosplit
func Read() (byte, bool) {
	if head == tail {
		return 0, false
	}
	ch := buf[head]
	head = (head + 1) % bufSize
	return ch, true
}

//go:nosplit
func Pending() bool {
	return head != tail
}
