set shell := ["bash", "-c"]

build_dir  := "build"
kernel_elf := build_dir / "demonos.elf"
iso_dir    := build_dir / "iso"
iso        := build_dir / "demonos.iso"
disk_img   := build_dir / "disk.img"
disk_size  := "256M"

nasm_flags := "-f elf64"
ld_script  := "kernel.ld"
ld_wrapper := "ld_wrapper.sh"

go_gcflags := "-e -B -wb=false"

asm_srcs   := "boot/boot.s arch/x86_64/gdt.s arch/x86_64/idt.s arch/x86_64/cpu.s arch/x86_64/syscall.s arch/x86_64/usermode.s"

default: build

build: _dirs _asm _link
    @echo "kernel: {{kernel_elf}}"

_dirs:
    mkdir -p {{build_dir}}

_asm:
    #!/usr/bin/env bash
    set -e
    for src in {{asm_srcs}}; do
        obj={{build_dir}}/$(basename ${src%.s}).o
        nasm {{nasm_flags}} -o "$obj" "$src"
    done

_link: _asm
    #!/usr/bin/env bash
    set -e
    asm_objs=$(ls {{build_dir}}/*.o | tr '\n' ' ')
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build \
        -gcflags="{{go_gcflags}}" \
        -ldflags="-linkmode external -extld $(pwd)/{{ld_wrapper}} -extldflags '-T {{ld_script}} -nostdlib -static $asm_objs'" \
        -o {{kernel_elf}} \
        ./kernel

iso: build
    #!/usr/bin/env bash
    set -e
    if ! command -v limine &>/dev/null; then
        echo "limine not found; install with: sudo pacman -S limine" >&2
        exit 1
    fi

    LIMINE_DATA=$(limine --print-datadir 2>/dev/null || echo /usr/share/limine)
    ISO_DIR={{iso_dir}}
    ISO={{iso}}

    mkdir -p "$ISO_DIR/boot/limine" "$ISO_DIR/EFI/BOOT"

    cp {{kernel_elf}} "$ISO_DIR/boot/demonos.elf"

    printf 'timeout: 0\n\n/DemonOS\n    protocol: multiboot2\n    kernel_path: boot():/boot/demonos.elf\n' \
        > "$ISO_DIR/boot/limine/limine.conf"

    cp "$LIMINE_DATA/limine-bios.sys"       "$ISO_DIR/boot/limine/"
    cp "$LIMINE_DATA/limine-bios-cd.bin"    "$ISO_DIR/boot/limine/"
    cp "$LIMINE_DATA/limine-uefi-cd.bin"    "$ISO_DIR/boot/limine/"
    cp "$LIMINE_DATA/BOOTX64.EFI"           "$ISO_DIR/EFI/BOOT/"

    xorriso -as mkisofs \
        -b boot/limine/limine-bios-cd.bin \
        -no-emul-boot -boot-load-size 4 -boot-info-table \
        --efi-boot boot/limine/limine-uefi-cd.bin \
        -efi-boot-part --efi-boot-image \
        --protective-msdos-label \
        -o "$ISO" "$ISO_DIR"

    limine bios-install "$ISO"

    echo "iso: $ISO"

disk:
    #!/usr/bin/env bash
    set -e
    IMG={{disk_img}}
    MNT=$(mktemp -d)
    trap "sudo umount '$MNT' 2>/dev/null; rmdir '$MNT'" EXIT

    dd if=/dev/zero of="$IMG" bs=1M count=0 seek=$(echo "{{disk_size}}" | sed 's/M//') 2>/dev/null
    mkfs.ext4 -q -L "DemonOS" "$IMG"
    sudo mount -o loop "$IMG" "$MNT"

    sudo mkdir -p \
        "$MNT/Applications" \
        "$MNT/Library/Frameworks" \
        "$MNT/Library/Extensions" \
        "$MNT/Library/Preferences" \
        "$MNT/System/Library/CoreServices" \
        "$MNT/Users/Shared" \
        "$MNT/bin" \
        "$MNT/sbin" \
        "$MNT/dev" \
        "$MNT/cores" \
        "$MNT/Volumes" \
        "$MNT/opt/local/bin" \
        "$MNT/opt/local/lib" \
        "$MNT/private/etc" \
        "$MNT/private/var/log" \
        "$MNT/private/var/run" \
        "$MNT/private/tmp" \
        "$MNT/usr/bin" \
        "$MNT/usr/lib" \
        "$MNT/usr/libexec" \
        "$MNT/usr/local/bin" \
        "$MNT/usr/local/lib" \
        "$MNT/usr/share"

    sudo ln -s private/etc  "$MNT/etc"
    sudo ln -s private/var  "$MNT/var"
    sudo ln -s private/tmp  "$MNT/tmp"

    sudo chmod 1777 "$MNT/private/tmp"
    sudo chmod 755  "$MNT/private/var/log"

    echo "disk: $IMG"

run: build
    qemu-system-x86_64 \
        -kernel {{kernel_elf}} \
        -m 256M \
        -serial stdio \
        -no-reboot \
        -no-shutdown \
        -display none

run-disk: build disk
    qemu-system-x86_64 \
        -kernel {{kernel_elf}} \
        -m 256M \
        -serial stdio \
        -no-reboot \
        -no-shutdown \
        -display none \
        -drive file={{disk_img}},format=raw,if=virtio

run-iso: iso
    qemu-system-x86_64 \
        -cdrom {{iso}} \
        -m 256M \
        -serial stdio \
        -no-reboot \
        -no-shutdown

run-kvm: build
    qemu-system-x86_64 \
        -kernel {{kernel_elf}} \
        -m 256M \
        -enable-kvm \
        -serial stdio \
        -no-reboot \
        -no-shutdown \
        -display none

debug: build
    #!/usr/bin/env bash
    set -e
    qemu-system-x86_64 \
        -kernel {{kernel_elf}} \
        -m 256M \
        -serial stdio \
        -no-reboot \
        -no-shutdown \
        -display none \
        -s -S &
    gdb -ex "target remote :1234" \
        -ex "symbol-file {{kernel_elf}}"

clean:
    rm -rf {{build_dir}}
