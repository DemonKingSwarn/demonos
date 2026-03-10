#include "textflag.h"

TEXT ·gdtLoad(SB),NOSPLIT,$0
    RET

TEXT ·idtLoad(SB),NOSPLIT,$0
    RET

TEXT ·setKernelStack(SB),NOSPLIT,$0-8
    RET

TEXT ·syscallEntryAsm(SB),NOSPLIT,$0
    RET

TEXT ·writeMSR(SB),NOSPLIT,$0-12
    RET

TEXT ·readMSR(SB),NOSPLIT,$0-12
    RET

TEXT ·cpuSetGSBase(SB),NOSPLIT,$0-8
    RET

TEXT ·cpuSetKernelGSBase(SB),NOSPLIT,$0-8
    RET
