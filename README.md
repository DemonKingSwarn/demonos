# DemonOS

A kernel written in Go and x86-64 assembly, with Linux syscall ABI compatibility.

## Requirements

- `go` 1.21+
- `nasm` 2.15+
- `ld` (binutils)
- `limine` + `xorriso` (for ISO)
- `qemu-system-x86_64` (for running)
- `just`

## Build

```sh
just build
```

## Run in QEMU

```sh
just run
```

## Debug with GDB

```sh
just debug
```

This starts QEMU paused on port 1234 and opens GDB connected to it.

## How it works

1. GRUB loads the kernel ELF via Multiboot2 and jumps to `_start` in 32-bit protected mode.
2. `boot/boot.s` sets up identity-mapped page tables (first 1 GiB as 2 MiB pages), enables PAE, sets the LME bit in EFER, enables paging, then far-jumps into 64-bit code.
3. `kmain` (Go) initialises the physical memory allocator from the Multiboot2 memory map, loads the full GDT (including TSS), loads the IDT, and enables `SYSCALL`/`SYSRET` via MSRs.
4. If an init ELF binary is linked in (`kernel/kernel.go:initBin`), it is parsed, loaded into physical memory, and executed in ring 3.
5. Syscalls from user space land in `arch/x86_64/syscall.s`, which saves caller-saved registers and calls `syscallHandler` in Go. Currently implemented: `write` (fd 1/2 go to VGA), `exit`, `exit_group`.
