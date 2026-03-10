bits 64
default rel

section .data
align 8

gdt64_full:
    dq 0
    dq (1<<43)|(1<<44)|(1<<47)|(1<<53)
    dq (1<<41)|(1<<44)|(1<<47)
    dq (1<<41)|(1<<44)|(1<<45)|(1<<46)|(1<<47)
    dq (1<<43)|(1<<44)|(1<<45)|(1<<46)|(1<<47)|(1<<53)
tss_desc_lo: dq 0
tss_desc_hi: dq 0

gdt64_full_ptr:
    dw $ - gdt64_full - 1
    dq gdt64_full

TSS_SIZE equ 108

align 16
global tss
tss:
    dd 0
    dq 0
    dq 0
    dq 0
    dq 0
    times 7 dq 0
    dq 0
    dw 0
    dw TSS_SIZE

section .text
global gdtLoad
gdtLoad:
    mov   rax, tss
    mov   rcx, rax

    mov   rdx, (TSS_SIZE - 1)
    and   rdx, 0xFFFF
    mov   rbx, rax
    shl   rbx, 16
    or    rdx, rbx
    mov   rbx, rax
    shr   rbx, 16
    and   rbx, 0xFF
    shl   rbx, 32
    or    rdx, rbx
    mov   r9, (0x89 << 40)
    or    rdx, r9
    mov   rbx, rax
    mov   r9, 0xFF000000
    and   rbx, r9
    shl   rbx, 32
    or    rdx, rbx
    mov   [tss_desc_lo], rdx

    mov   rbx, rax
    shr   rbx, 32
    mov   [tss_desc_hi], rbx

    lgdt  [gdt64_full_ptr]

    push  qword 0x08
    lea   rax, [.reload_cs]
    push  rax
    retfq

.reload_cs:
    mov   ax, 0x10
    mov   ds, ax
    mov   es, ax
    mov   fs, ax
    mov   gs, ax
    mov   ss, ax

    mov   ax, 0x28
    ltr   ax
    ret

global setKernelStack
setKernelStack:
    mov   [tss + 4], rdi
    ret
