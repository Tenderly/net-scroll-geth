# Scroll-geth Rust Symbol Conflict Fix

## Problem

When linking both `scroll-geth` and `nitro` in the same Go binary, there are duplicate Rust runtime symbols that cause linker errors:
- `_rust_eh_personality`
- `___rust_no_alloc_shim_is_unstable`

Both libraries export these standard Rust symbols, causing a conflict during final linking.

## Solution

We use `llvm-objcopy --redefine-sym` to rename the conflicting symbols in scroll-geth's Darwin ARM64 libraries with a `scroll_ffi_` prefix:

| Old name | New name |
|---|---|
| `_rust_eh_personality` | `scroll_ffi_rust_eh_personality` |
| `___rust_no_alloc_shim_is_unstable` | `scroll_ffi___rust_no_alloc_shim_is_unstable` |

The names are raw Mach-O names: a C name gets a leading `_` in Mach-O, the new names do not.

The inputs are the `.backup` files. They hold the libraries after `add_symbol_prefix.sh` and before the rename. Keep them in the repo. `make fix-symbols` always makes the committed libraries from them, so the output is the same each time.

We used `objconv` for this step before. It wrote broken archives: junk bytes in many member names, and member data that did not start on an 8-byte boundary. Apple's linker skips those members ("not 8-byte aligned"), and rustc refuses the names ("invalid utf-8"). `llvm-objcopy` writes clean archives.

### Prerequisites

Install LLVM:
```bash
brew install llvm
```

The Makefile uses `/opt/homebrew/opt/llvm/bin/llvm-objcopy`, the same path as `add_symbol_prefix.sh`. To use another copy, give its path. For example, the copy in a rustup toolchain that has the `llvm-tools` component:
```bash
make fix-symbols LLVM_OBJCOPY=$HOME/.rustup/toolchains/<toolchain>/lib/rustlib/aarch64-apple-darwin/bin/llvm-objcopy
```

### Manual Process

```bash
/opt/homebrew/opt/llvm/bin/llvm-objcopy \
  --redefine-sym _rust_eh_personality=scroll_ffi_rust_eh_personality \
  --redefine-sym ___rust_no_alloc_shim_is_unstable=scroll_ffi___rust_no_alloc_shim_is_unstable \
  libencoder_legacy_darwin_arm64.a.backup libencoder_legacy_darwin_arm64.a
/opt/homebrew/opt/llvm/bin/llvm-objcopy \
  --redefine-sym _rust_eh_personality=scroll_ffi_rust_eh_personality \
  --redefine-sym ___rust_no_alloc_shim_is_unstable=scroll_ffi___rust_no_alloc_shim_is_unstable \
  libencoder_standard_darwin_arm64.a.backup libencoder_standard_darwin_arm64.a
```

### Automated Process

Run the Makefile target:
```bash
make fix-symbols
```

## Verification

You can verify the symbols were renamed correctly:
```bash
llvm-nm -g --no-llvm-bc libencoder_legacy_darwin_arm64.a | grep scroll_ffi
llvm-nm -g --no-llvm-bc libencoder_standard_darwin_arm64.a | grep scroll_ffi
```

You should see the renamed symbols with the `scroll_ffi_` prefix.

Some members (`std`, `core`, `alloc` and other Rust standard library members) also carry LLVM bitcode. Plain `nm` shows the bitcode symbols for those members, and the rename does not change the bitcode. `--no-llvm-bc` shows the Mach-O symbols, which are the ones the linker uses.

## When to Run

Run this after you rebuild the Rust encoder libraries:
1. Run `add_symbol_prefix.sh`.
2. Copy the two new Darwin ARM64 libraries over the `.backup` files.
3. Run `make fix-symbols`.

## Platform Support

Currently, this fix is only applied to Darwin ARM64 libraries (`darwin_arm64`). If building on other platforms (Linux AMD64/ARM64), similar fixes would be needed for those library files.

## References

- Original issue: Duplicate Rust symbols between scroll-geth and nitro
- Similar fix in zksync-go: https://github.com/Tenderly/zksync-go/blob/master/Makefile#L15
- llvm-objcopy documentation: https://llvm.org/docs/CommandGuide/llvm-objcopy.html
