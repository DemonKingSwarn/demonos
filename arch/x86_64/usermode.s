bits 64
default rel

section .text

global jumpUserMode
jumpUserMode:
    mov   rcx, rdi
    mov   rsp, rsi
    mov   r11, 0x202
    xor   rax, rax
    xor   rbx, rbx
    xor   rdx, rdx
    xor   rbp, rbp
    xor   rdi, rdi
    xor   rsi, rsi
    xor   r8,  r8
    xor   r9,  r9
    xor   r10, r10
    xor   r12, r12
    xor   r13, r13
    xor   r14, r14
    xor   r15, r15
    swapgs
    o64 sysret
