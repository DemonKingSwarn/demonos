bits 32
default rel

MB2_MAGIC       equ 0xE85250D6
MB2_ARCH_I386   equ 0
MB2_HDR_LEN     equ (mb2_end - mb2_start)
MB2_CHECKSUM    equ -(MB2_MAGIC + MB2_ARCH_I386 + MB2_HDR_LEN)

XEN_ELFNOTE_PHYS32_ENTRY equ 18

section .multiboot
align 8
mb2_start:
    dd MB2_MAGIC
    dd MB2_ARCH_I386
    dd MB2_HDR_LEN
    dd MB2_CHECKSUM
    dw 0
    dw 0
    dd 8
mb2_end:

section .note.Xen note
align 4
    dd pvh_name_end - pvh_name_start
    dd 4
    dd XEN_ELFNOTE_PHYS32_ENTRY
pvh_name_start:
    db "Xen", 0
pvh_name_end:
align 4
    dd pvh_start

section .bss
align 4096
pml4_table:   resb 4096
pdpt_table:   resb 4096
pd_table:     resb 4096

align 16
stack_bottom:
    resb 65536
stack_top:

section .rodata
align 8
gdt64:
    dq 0
.code: equ $ - gdt64
    dq (1<<43)|(1<<44)|(1<<47)|(1<<53)
.data: equ $ - gdt64
    dq (1<<41)|(1<<44)|(1<<47)
gdt64_ptr:
    dw $ - gdt64 - 1
    dq gdt64

section .text
global pvh_start
pvh_start:
    mov   edi, 0
    jmp   setup_paging

global _start
_start:
    mov   edi, ebx

    mov   eax, 0x80000000
    cpuid
    cmp   eax, 0x80000001
    jb    .no_longmode

    mov   eax, 0x80000001
    cpuid
    test  edx, (1 << 29)
    jz    .no_longmode
    jmp   setup_paging

.no_longmode:
    mov   word [0xb8000], 0x4F4E
    mov   word [0xb8002], 0x4F4C
    mov   word [0xb8004], 0x4F4D
.hang: hlt
    jmp   .hang

setup_paging:
    mov   eax, pdpt_table
    or    eax, 0x3
    mov   [pml4_table], eax

    mov   eax, pd_table
    or    eax, 0x3
    mov   [pdpt_table], eax

    mov   ecx, 0
.pd_loop:
    mov   eax, 0x200000
    mul   ecx
    or    eax, 0x83
    mov   [pd_table + ecx * 8], eax
    inc   ecx
    cmp   ecx, 512
    jl    .pd_loop

    mov   eax, pml4_table
    mov   cr3, eax

    mov   eax, cr4
    or    eax, (1 << 5)
    mov   cr4, eax

    mov   ecx, 0xC0000080
    rdmsr
    or    eax, (1 << 8)
    wrmsr

    mov   eax, cr0
    or    eax, (1 << 31) | 1
    mov   cr0, eax

    lgdt  [gdt64_ptr]
    jmp   gdt64.code:long_mode_entry

bits 64
long_mode_entry:
    mov   ax, gdt64.data
    mov   ds, ax
    mov   es, ax
    mov   fs, ax
    mov   gs, ax
    mov   ss, ax

    mov   rsp, stack_top

    mov   word [0xB8000], 0x0A4B

    call  serial_init

    lea   rsi, [rel boot_msg]
    call  serial_puts

    mov   edi, edi

    extern main.kmain
    call  main.kmain

    lea   rsi, [rel halt_msg]
    call  serial_puts

.halt:
    cli
    hlt
    jmp   .halt

serial_init:
    mov   dx, 0x3F8 + 1
    mov   al, 0x00
    out   dx, al
    mov   dx, 0x3F8 + 3
    mov   al, 0x80
    out   dx, al
    mov   dx, 0x3F8 + 0
    mov   al, 0x01
    out   dx, al
    mov   dx, 0x3F8 + 1
    mov   al, 0x00
    out   dx, al
    mov   dx, 0x3F8 + 3
    mov   al, 0x03
    out   dx, al
    mov   dx, 0x3F8 + 2
    mov   al, 0xC7
    out   dx, al
    mov   dx, 0x3F8 + 4
    mov   al, 0x0B
    out   dx, al
    ret

serial_putchar:
    mov   bl, al
.wait:
    mov   dx, 0x3F8 + 5
    in    al, dx
    test  al, 0x20
    jz    .wait
    mov   al, bl
    mov   dx, 0x3F8
    out   dx, al
    ret

serial_puts:
.loop:
    mov   al, [rsi]
    test  al, al
    jz    .done
    call  serial_putchar
    inc   rsi
    jmp   .loop
.done:
    ret

section .rodata
boot_msg: db "DemonOS boot64", 13, 10, 0
halt_msg: db "DemonOS halt", 13, 10, 0
