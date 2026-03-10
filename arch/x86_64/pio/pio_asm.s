#include "textflag.h"

TEXT ·Outb(SB),NOSPLIT,$0-3
    MOVWLZX port+0(FP), DX
    MOVBLZX val+2(FP), AX
    BYTE $0xEE
    RET

TEXT ·Inb(SB),NOSPLIT,$0-9
    MOVWLZX port+0(FP), DX
    BYTE $0xEC
    MOVB AX, ret+8(FP)
    RET

TEXT ·Outw(SB),NOSPLIT,$0-4
    MOVWLZX port+0(FP), DX
    MOVWLZX val+2(FP), AX
    BYTE $0x66; BYTE $0xEF
    RET

TEXT ·Inw(SB),NOSPLIT,$0-10
    MOVWLZX port+0(FP), DX
    BYTE $0x66; BYTE $0xED
    MOVW AX, ret+8(FP)
    RET

TEXT ·Outl(SB),NOSPLIT,$0-8
    MOVWLZX port+0(FP), DX
    MOVLQZX val+4(FP), AX
    BYTE $0xEF
    RET

TEXT ·Inl(SB),NOSPLIT,$0-12
    MOVWLZX port+0(FP), DX
    BYTE $0xED
    MOVL AX, ret+8(FP)
    RET
