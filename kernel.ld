OUTPUT_FORMAT(elf64-x86-64)
OUTPUT_ARCH(i386:x86-64)
ENTRY(_start)

KERNEL_BASE = 0x0000000000100000;

SECTIONS {
    . = KERNEL_BASE;

    .multiboot : {
        KEEP(*(.multiboot))
    }

    .note.Xen : {
        KEEP(*(.note.Xen))
    }

    .text ALIGN(4096) : {
        *(.text .text.*)
    }

    .rodata ALIGN(4096) : {
        *(.rodata .rodata.*)
    }

    .typelinks ALIGN(8) : {
        *(.typelinks)
    }

    .itablinks ALIGN(8) : {
        *(.itablinks)
    }

    .gosymtab ALIGN(8) : {
        *(.gosymtab)
    }

    .gopclntab ALIGN(8) : {
        *(.gopclntab)
    }

    .go.buildinfo ALIGN(16) : {
        *(.go.buildinfo)
    }

    .noptrdata ALIGN(4096) : {
        *(.noptrdata .noptrdata.*)
    }

    .data ALIGN(4096) : {
        *(.data .data.*)
    }

    .bss ALIGN(4096) : {
        *(COMMON)
        *(.bss .bss.*)
    }

    .noptrbss ALIGN(4096) : {
        *(.noptrbss .noptrbss.*)
    }

    .tbss ALIGN(8) (NOLOAD) : {
        *(.tbss .tbss.*)
    }

    . = ALIGN(16);
    _kernel_end = .;

    /DISCARD/ : {
        *(.comment)
        *(.note*)
        *(.eh_frame*)
        *(.tdata)
        *(.gnu_debuglink)
    }
}
