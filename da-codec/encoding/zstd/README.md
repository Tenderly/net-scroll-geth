# Scroll-geth Rust Symbol Conflict Fix

## Problem

When linking both `scroll-geth` and `nitro` in the same Go binary, there are duplicate Rust runtime symbols that cause linker errors:
- `_rust_eh_personality`
- `___rust_no_alloc_shim_is_unstable`

Both libraries export these standard Rust symbols, causing a conflict during final linking.

## Solution

We use `objconv` (a tool by Agner Fog) to rename the conflicting symbols in scroll-geth's libraries with a `scroll_ffi_` prefix. This is the same approach used successfully in the [zksync-go](https://github.com/Tenderly/zksync-go) project.

### Prerequisites

Install objconv:
```bash
brew install objconv
```

### Manual Process

1. Backup the original libraries:
```bash
cp libencoder_legacy_darwin_arm64.a libencoder_legacy_darwin_arm64.a.backup
cp libencoder_standard_darwin_arm64.a libencoder_standard_darwin_arm64.a.backup
```

2. Run objconv with the configuration file:
```bash
objconv @objconv.conf libencoder_legacy_darwin_arm64.a libencoder_legacy_darwin_arm64_fixed.a
objconv @objconv.conf libencoder_standard_darwin_arm64.a libencoder_standard_darwin_arm64_fixed.a
```

3. Replace the original files:
```bash
mv libencoder_legacy_darwin_arm64_fixed.a libencoder_legacy_darwin_arm64.a
mv libencoder_standard_darwin_arm64_fixed.a libencoder_standard_darwin_arm64.a
```

### Automated Process

Run the Makefile target:
```bash
make fix-symbols
```

## Configuration

The `objconv.conf` file contains the symbol renaming rules:
```
-nr:_rust_eh_personality:scroll_ffi_rust_eh_personality
-nr:___rust_no_alloc_shim_is_unstable:scroll_ffi___rust_no_alloc_shim_is_unstable
```

The `-nr` flag tells objconv to rename symbols: `-nr:old_name:new_name`

## Verification

You can verify the symbols were renamed correctly:
```bash
nm -g libencoder_legacy_darwin_arm64.a | grep scroll_ffi
nm -g libencoder_standard_darwin_arm64.a | grep scroll_ffi
```

You should see the renamed symbols with the `scroll_ffi_` prefix.

## When to Run

This needs to be run:
- After updating scroll-geth from upstream
- After rebuilding the Rust encoder libraries
- On fresh checkouts/links of scroll-geth

## Platform Support

Currently, this fix is only applied to Darwin ARM64 libraries (`darwin_arm64`). If building on other platforms (Linux AMD64/ARM64), similar fixes would be needed for those library files.

## References

- Original issue: Duplicate Rust symbols between scroll-geth and nitro
- Similar fix in zksync-go: https://github.com/Tenderly/zksync-go/blob/master/Makefile#L15
- objconv documentation: https://www.agner.org/optimize/objconv-instructions.pdf
