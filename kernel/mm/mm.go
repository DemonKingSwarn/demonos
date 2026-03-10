package mm

import (
	"unsafe"
)

const PageSize = 4096

type mb2Header struct {
	totalSize uint32
	reserved  uint32
}

type mb2Tag struct {
	tagType uint32
	size    uint32
}

type mb2MemMapTag struct {
	tagType   uint32
	size      uint32
	entrySize uint32
	entryVer  uint32
}

type mb2MemEntry struct {
	baseAddr uint64
	length   uint64
	memType  uint32
	reserved uint32
}

const maxPages = 1024 * 1024

var (
	bitmap    [maxPages / 64]uint64
	totalFree uint64
	totalMem  uint64
)

//go:nosplit
func markUsed(page uint64) {
	bitmap[page/64] |= 1 << (page % 64)
}

//go:nosplit
func markFree(page uint64) {
	bitmap[page/64] &^= 1 << (page % 64)
}

//go:nosplit
func isFree(page uint64) bool {
	return bitmap[page/64]&(1<<(page%64)) == 0
}

//go:nosplit
func Init(mbInfo uint32) {
	for i := range bitmap {
		bitmap[i] = ^uint64(0)
	}

	if mbInfo == 0 {
		firstPage := uint64(0x100000) / PageSize
		lastPage := uint64(256*1024*1024) / PageSize
		if lastPage > maxPages {
			lastPage = maxPages
		}
		for p := firstPage; p < lastPage; p++ {
			markFree(p)
			totalFree++
			totalMem += PageSize
		}
		for p := uint64(0); p < 512; p++ {
			if isFree(p) {
				markUsed(p)
				totalFree--
			}
		}
		return
	}

	ptr := uintptr(mbInfo) + 8
	end := uintptr(mbInfo) + uintptr((*mb2Header)(unsafe.Pointer(uintptr(mbInfo))).totalSize)

	for ptr < end {
		tag := (*mb2Tag)(unsafe.Pointer(ptr))
		if tag.tagType == 0 {
			break
		}

		if tag.tagType == 6 {
			mmTag := (*mb2MemMapTag)(unsafe.Pointer(ptr))
			entryPtr := ptr + unsafe.Sizeof(*mmTag)
			numEntries := (uintptr(mmTag.size) - unsafe.Sizeof(*mmTag)) / uintptr(mmTag.entrySize)

			for i := uintptr(0); i < numEntries; i++ {
				e := (*mb2MemEntry)(unsafe.Pointer(entryPtr + i*uintptr(mmTag.entrySize)))
				totalMem += e.length
				if e.memType != 1 {
					continue
				}
				firstPage := e.baseAddr / PageSize
				lastPage := (e.baseAddr + e.length) / PageSize
				if lastPage > maxPages {
					lastPage = maxPages
				}
				for p := firstPage; p < lastPage; p++ {
					markFree(p)
					totalFree++
				}
			}
		}

		ptr += uintptr((tag.size + 7) &^ 7)
	}

	for p := uint64(0); p < 512; p++ {
		if isFree(p) {
			markUsed(p)
			totalFree--
		}
	}
}

func AllocPage() uintptr {
	for i, word := range bitmap {
		if word == ^uint64(0) {
			continue
		}
		for bit := uint64(0); bit < 64; bit++ {
			if word&(1<<bit) == 0 {
				page := uint64(i)*64 + bit
				markUsed(page)
				totalFree--
				return uintptr(page * PageSize)
			}
		}
	}
	return 0
}

func FreePage(physAddr uintptr) {
	page := uint64(physAddr) / PageSize
	if page < maxPages && !isFree(page) {
		markFree(page)
		totalFree++
	}
}

//go:nosplit
func TotalFreeBytes() uint64 { return totalFree * PageSize }

//go:nosplit
func TotalMemBytes() uint64 { return totalMem }
