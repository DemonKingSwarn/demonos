package elf

import (
	"demonos/kernel/mm"
	"unsafe"
)

const (
	elfMagic = 0x464C457F

	etExec = 2

	ptLoad = 1

	pfX = 0x1
	pfW = 0x2
	pfR = 0x4

	pageSize = mm.PageSize
)

type Header64 struct {
	Ident     [16]byte
	Type      uint16
	Machine   uint16
	Version   uint32
	Entry     uint64
	PhOff     uint64
	ShOff     uint64
	Flags     uint32
	EhSize    uint16
	PhEntSize uint16
	PhNum     uint16
	ShEntSize uint16
	ShNum     uint16
	ShStrNdx  uint16
}

type Phdr64 struct {
	Type   uint32
	Flags  uint32
	Offset uint64
	VAddr  uint64
	PAddr  uint64
	FileSz uint64
	MemSz  uint64
	Align  uint64
}

type LoadResult struct {
	Entry    uintptr
	StackTop uintptr
}

var ErrBadMagic = loadErr("bad ELF magic")
var ErrNotExec = loadErr("not an executable ELF")
var ErrBadArch = loadErr("not an x86-64 ELF")

type loadErr string

func (e loadErr) Error() string { return string(e) }

func Load(image []byte) (LoadResult, error) {
	if len(image) < int(unsafe.Sizeof(Header64{})) {
		return LoadResult{}, ErrBadMagic
	}

	hdr := (*Header64)(unsafe.Pointer(&image[0]))

	magic := *(*uint32)(unsafe.Pointer(&hdr.Ident[0]))
	if magic != elfMagic {
		return LoadResult{}, ErrBadMagic
	}
	if hdr.Type != etExec {
		return LoadResult{}, ErrNotExec
	}
	if hdr.Machine != 0x3E {
		return LoadResult{}, ErrBadArch
	}

	phdrBase := uintptr(unsafe.Pointer(&image[0])) + uintptr(hdr.PhOff)

	for i := uint16(0); i < hdr.PhNum; i++ {
		ph := (*Phdr64)(unsafe.Pointer(phdrBase + uintptr(i)*uintptr(hdr.PhEntSize)))
		if ph.Type != ptLoad {
			continue
		}
		if err := loadSegment(image, ph); err != nil {
			return LoadResult{}, err
		}
	}

	stackPages := uintptr(8)
	stackBase := uintptr(0)
	for i := uintptr(0); i < stackPages; i++ {
		p := mm.AllocPage()
		if p == 0 {
			return LoadResult{}, loadErr("out of memory for stack")
		}
		if stackBase == 0 {
			stackBase = p
		}
	}

	result := LoadResult{
		Entry:    uintptr(hdr.Entry),
		StackTop: stackBase + stackPages*uintptr(pageSize),
	}
	return result, nil
}

func loadSegment(image []byte, ph *Phdr64) error {
	if ph.MemSz == 0 {
		return nil
	}

	numPages := (ph.MemSz + uint64(pageSize) - 1) / uint64(pageSize)
	destBase := uintptr(ph.VAddr) &^ uintptr(pageSize-1)

	for i := uint64(0); i < numPages; i++ {
		p := mm.AllocPage()
		if p == 0 {
			return loadErr("out of memory")
		}
		dest := destBase + uintptr(i)*uintptr(pageSize)
		zeroPage(dest)
		_ = dest
	}

	if ph.FileSz > 0 {
		src := image[ph.Offset : ph.Offset+ph.FileSz]
		dst := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ph.VAddr))), ph.FileSz)
		copy(dst, src)
	}

	return nil
}

func zeroPage(addr uintptr) {
	p := unsafe.Slice((*byte)(unsafe.Pointer(addr)), pageSize)
	for i := range p {
		p[i] = 0
	}
}
