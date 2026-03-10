bits 64
default rel

%macro ISR_NOERR 1
isr%1:
    push  qword 0
    push  qword %1
    jmp   isr_common
%endmacro

%macro ISR_ERR 1
isr%1:
    push  qword %1
    jmp   isr_common
%endmacro

ISR_NOERR 0
ISR_NOERR 1
ISR_NOERR 2
ISR_NOERR 3
ISR_NOERR 4
ISR_NOERR 5
ISR_NOERR 6
ISR_NOERR 7
ISR_ERR   8
ISR_NOERR 9
ISR_ERR   10
ISR_ERR   11
ISR_ERR   12
ISR_ERR   13
ISR_ERR   14
ISR_NOERR 15
ISR_NOERR 16
ISR_ERR   17
ISR_NOERR 18
ISR_NOERR 19
ISR_NOERR 20
ISR_NOERR 21
ISR_NOERR 22
ISR_NOERR 23
ISR_NOERR 24
ISR_NOERR 25
ISR_NOERR 26
ISR_NOERR 27
ISR_NOERR 28
ISR_NOERR 29
ISR_ERR   30
ISR_NOERR 31

section .text

isr_stub_generic:
    push  qword 0
    push  qword 255
    jmp   isr_common

isr_common:
    push  rax
    push  rbx
    push  rcx
    push  rdx
    push  rsi
    push  rdi
    push  rbp
    push  r8
    push  r9
    push  r10
    push  r11
    push  r12
    push  r13
    push  r14
    push  r15

    mov   rdi, rsp
    extern main.trapHandler
    call  main.trapHandler

    pop   r15
    pop   r14
    pop   r13
    pop   r12
    pop   r11
    pop   r10
    pop   r9
    pop   r8
    pop   rbp
    pop   rdi
    pop   rsi
    pop   rdx
    pop   rcx
    pop   rbx
    pop   rax

    add   rsp, 16
    iretq

section .data
align 16

global idt
idt:
    times 256*2 dq 0

global idt_ptr
idt_ptr:
    dw 256 * 16 - 1
    dq idt

global isrTable
isrTable:
%assign i 0
%rep 32
    dq isr%+i
%assign i i+1
%endrep
%rep 224
    dq isr_stub_generic
%endrep

section .text

global idtLoad
idtLoad:
    xor   rcx, rcx
.fill:
    cmp   rcx, 256
    jge   .done

    lea   rbx, [isrTable]
    mov   rax, [rbx + rcx * 8]

    mov   rdx, rax
    and   rdx, 0xFFFF
    or    rdx, (0x08 << 16)
    mov   r9, (0x8E << 40)
    or    rdx, r9
    mov   rbx, rax
    mov   r10, 0xFFFF0000
    and   rbx, r10
    shl   rbx, 32
    or    rdx, rbx

    mov   rbx, rax
    shr   rbx, 32

    lea   rdi, [idt]
    mov   r8, rcx
    shl   r8, 4
    add   rdi, r8
    mov   [rdi],     rdx
    mov   [rdi + 8], rbx

    inc   rcx
    jmp   .fill
.done:
    lidt  [idt_ptr]
    ret
