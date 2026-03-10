#!/usr/bin/env bash
args=()
go_obj=""
skip_next=0
for arg in "$@"; do
    if [ "$skip_next" = "1" ]; then
        skip_next=0
        continue
    fi
    case "$arg" in
        -m64|-lpthread|-lm|-lc|-ldl|-lrt|-O2|-g|-rdynamic) continue ;;
        -m) skip_next=1; continue ;;
        */go.o)
            go_obj="$arg"
            args+=("$arg")
            ;;
        *) args+=("$arg") ;;
    esac
done

if [ -n "$go_obj" ]; then
    objcopy \
        --globalize-symbol="main.kmain" \
        --globalize-symbol="main.trapHandler" \
        --globalize-symbol="main.syscallHandler" \
        "$go_obj" "$go_obj.globalized" && mv "$go_obj.globalized" "$go_obj"
fi

exec /usr/bin/ld -m elf_x86_64 "${args[@]}"
