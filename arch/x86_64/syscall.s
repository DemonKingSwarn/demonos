bits 64
default rel

section .text

global syscallEntryAsm
syscallEntryAsm:
    swapgs

    mov   [gs:16], rsp
    mov   rsp, [gs:8]

    push  r11
    push  rcx
    push  rbp
    push  rbx
    push  r12
    push  r13
    push  r14
    push  r15

    mov   rdi, rax
    mov   rsi, rbx
    mov   rdx, rcx
    mov   rcx, r10

    extern main.syscallHandler
    call  main.syscallHandler

    pop   r15
    pop   r14
    pop   r13
    pop   r12
    pop   rbx
    pop   rbp
    pop   rcx
    pop   r11

    mov   rsp, [gs:16]
    swapgs
    o64 sysret
