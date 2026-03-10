bits 64

section .text

global writeMSR
writeMSR:
    mov   ecx, edi
    mov   eax, esi
    mov   rdx, rsi
    shr   rdx, 32
    wrmsr
    ret

global readMSR
readMSR:
    mov   ecx, edi
    rdmsr
    shl   rdx, 32
    or    rax, rdx
    ret

global cpuSetGSBase
cpuSetGSBase:
    mov   rcx, 0xC0000101
    mov   rax, rdi
    mov   rdx, rdi
    shr   rdx, 32
    wrmsr
    ret

global cpuSetKernelGSBase
cpuSetKernelGSBase:
    mov   rcx, 0xC0000102
    mov   rax, rdi
    mov   rdx, rdi
    shr   rdx, 32
    wrmsr
    ret

global halt
halt:
    cli
    hlt
    jmp   halt
