# Symbol Prefix Guide for da-codec Libraries

This document describes how to manage symbol conflicts between multiple Rust static libraries in the da-codec package.

## Problem

The da-codec package uses multiple Rust static libraries that contain identical Rust runtime symbols, causing duplicate symbol errors during linking on macOS:

- `libstylus_darwin_arm64.a` (from nitro package)
- `libencoder_legacy_darwin_arm64.a` (custom zstd-rs fork)  
- `libencoder_standard_darwin_arm64.a` (standard zstd 0.13.3)

Common conflicting symbols include:
- `___rg_oom`
- `___rust_foreign_exception`
- `___rust_drop_panic`
- `___rdl_alloc`
- `_rust_begin_unwind`
- `___rdl_dealloc`
- `___rdl_oom`
- `___rdl_realloc`
- `_rust_panic`
- `_rust_eh_personality`
- `___rdl_alloc_zeroed`
- `___rust_start_panic`
- `___rust_panic_cleanup`

## Solution Architecture

The solution uses symbol prefixing to make each library's Rust runtime symbols unique:

1. **libstylus**: Uses `nitro` prefix
2. **libencoder_legacy**: Uses `scroll_legacy_` prefix  
3. **libencoder_standard**: Uses `scroll_standard_` prefix

## Files and Structure

```
da-codec/encoding/
├── zstd/
│   ├── add_symbol_prefix.sh           # Main symbol prefixing script
│   ├── symbol_wrapper.c               # Common symbol implementations
│   ├── libstylus_darwin_arm64.a       # Nitro library (prefixed)
│   ├── libencoder_legacy_darwin_arm64.a    # Legacy encoder (prefixed)
│   ├── libencoder_standard_darwin_arm64.a  # Standard encoder (prefixed)
│   └── zstd.go                        # Go bindings
└── libzstd/
    ├── encoder-legacy/                # Legacy encoder source
    │   ├── Cargo.toml
    │   ├── Makefile
    │   └── src/lib.rs
    └── encoder-standard/              # Standard encoder source  
        ├── Cargo.toml
        ├── Makefile
        └── src/lib.rs
```

## Rebuilding Libraries

### Prerequisites

- Homebrew LLVM tools: `brew install llvm`
- Rust toolchain
- Access to scroll-tech/zstd-rs repository (for legacy encoder)

### Step 1: Rebuild Standard Encoder

```bash
cd da-codec/encoding/libzstd/encoder-standard
make clean
make build
make install
```

### Step 2: Rebuild Legacy Encoder (if source is available)

```bash
cd da-codec/encoding/libzstd/encoder-legacy
make clean
make build  
make install
```

**Note**: Legacy encoder uses a private repository that may require authentication.

### Step 3: Apply Symbol Prefixing

```bash
cd da-codec/encoding/zstd
./add_symbol_prefix.sh
```

If the script detects already-processed libraries, force processing:

```bash
FORCE_PROCESS_STANDARD=1 ./add_symbol_prefix.sh
```

### Step 4: Verify Symbol Prefixing

Check that symbols have been prefixed:

```bash
# Check legacy encoder
/opt/homebrew/opt/llvm/bin/llvm-nm libencoder_legacy_darwin_arm64.a | grep "scroll_legacy_"

# Check standard encoder  
/opt/homebrew/opt/llvm/bin/llvm-nm libencoder_standard_darwin_arm64.a | grep "scroll_standard_"

# Check stylus library
/opt/homebrew/opt/llvm/bin/llvm-nm libstylus_darwin_arm64.a | grep "nitro"
```

### Step 5: Test Build

```bash
cd da-codec/encoding
go build ./zstd
```

Should complete without duplicate symbol errors.

## Manual Symbol Prefixing (if script fails)

If the automatic script fails, use manual symbol prefixing:

### For Standard Encoder:

```bash
cd da-codec/encoding/zstd

# Create symbol redefine file
echo '___rg_oom scroll_standard____rg_oom
___rust_foreign_exception scroll_standard____rust_foreign_exception
___rust_drop_panic scroll_standard____rust_drop_panic
___rdl_alloc scroll_standard____rdl_alloc
_rust_begin_unwind scroll_standard__rust_begin_unwind
___rdl_dealloc scroll_standard____rdl_dealloc
___rdl_oom scroll_standard____rdl_oom
___rdl_realloc scroll_standard____rdl_realloc
_rust_panic scroll_standard__rust_panic
_rust_eh_personality scroll_standard__rust_eh_personality
___rdl_alloc_zeroed scroll_standard____rdl_alloc_zeroed
___rust_start_panic scroll_standard____rust_start_panic
___rust_panic_cleanup scroll_standard____rust_panic_cleanup' > redefine_standard.syms

# Apply prefixing
/opt/homebrew/opt/llvm/bin/llvm-objcopy --redefine-syms=redefine_standard.syms \
    libencoder_standard_darwin_arm64.a libencoder_standard_darwin_arm64_new.a
mv libencoder_standard_darwin_arm64_new.a libencoder_standard_darwin_arm64.a
ranlib libencoder_standard_darwin_arm64.a
```

### For Legacy Encoder:

```bash
# Create symbol redefine file  
echo '___rg_oom scroll_legacy____rg_oom
___rust_foreign_exception scroll_legacy____rust_foreign_exception
___rust_drop_panic scroll_legacy____rust_drop_panic
___rdl_alloc scroll_legacy____rdl_alloc
_rust_begin_unwind scroll_legacy__rust_begin_unwind
___rdl_dealloc scroll_legacy____rdl_dealloc
___rdl_oom scroll_legacy____rdl_oom
___rdl_realloc scroll_legacy____rdl_realloc
_rust_panic scroll_legacy__rust_panic
_rust_eh_personality scroll_legacy__rust_eh_personality
___rdl_alloc_zeroed scroll_legacy____rdl_alloc_zeroed
___rust_start_panic scroll_legacy____rust_start_panic
___rust_panic_cleanup scroll_legacy____rust_panic_cleanup' > redefine_legacy.syms

# Apply prefixing
/opt/homebrew/opt/llvm/bin/llvm-objcopy --redefine-syms=redefine_legacy.syms \
    libencoder_legacy_darwin_arm64.a libencoder_legacy_darwin_arm64_new.a  
mv libencoder_legacy_darwin_arm64_new.a libencoder_legacy_darwin_arm64.a
ranlib libencoder_legacy_darwin_arm64.a
```

### For Stylus Library:

```bash
# Create symbol redefine file
echo '___rg_oom nitro___rg_oom
___rust_foreign_exception nitro___rust_foreign_exception
___rust_drop_panic nitro___rust_drop_panic
___rdl_alloc nitro___rdl_alloc
_rust_begin_unwind nitro_rust_begin_unwind
___rdl_dealloc nitro___rdl_dealloc
___rdl_oom nitro___rdl_oom
___rdl_realloc nitro___rdl_realloc
_rust_panic nitro_rust_panic
_rust_eh_personality nitro_rust_eh_personality
___rdl_alloc_zeroed nitro___rdl_alloc_zeroed
___rust_start_panic nitro___rust_start_panic
___rust_panic_cleanup nitro___rust_panic_cleanup' > redefine_stylus.syms

# Apply prefixing
/opt/homebrew/opt/llvm/bin/llvm-objcopy --redefine-syms=redefine_stylus.syms \
    libstylus_darwin_arm64.a libstylus_darwin_arm64_new.a
mv libstylus_darwin_arm64_new.a libstylus_darwin_arm64.a
ranlib libstylus_darwin_arm64.a
```

## Understanding the add_symbol_prefix.sh Script

The script automatically:

1. **Detects platform**: darwin_arm64, linux_amd64, linux_arm64
2. **Processes each library** with its designated prefix
3. **Identifies conflicting symbols** using pattern matching
4. **Applies symbol prefixing** using llvm-objcopy
5. **Verifies no conflicts** remain between libraries
6. **Preserves original functions** like `compress_scroll_batch_bytes_*`

Key configurations in the script:

```bash
LIBRARIES=(
    "libstylus_darwin_arm64.a:nitro"
    "libencoder_legacy_darwin_arm64.a:scroll_legacy_"
    "libencoder_legacy_linux_amd64.a:scroll_legacy_"
    "libencoder_legacy_linux_arm64.a:scroll_legacy_"
    "libencoder_standard_darwin_arm64.a:scroll_standard_"
    "libencoder_standard_linux_amd64.a:scroll_standard_"
    "libencoder_standard_linux_arm64.a:scroll_standard_"
)
```

## Troubleshooting

### Build fails with duplicate symbols

1. **Check if libraries have been processed**:
   ```bash
   ./add_symbol_prefix.sh
   ```

2. **Force reprocessing if needed**:
   ```bash
   FORCE_PROCESS_STANDARD=1 ./add_symbol_prefix.sh
   ```

3. **Manually verify symbol prefixing** using the commands in Step 4

4. **Use manual prefixing** if automatic script fails

### llvm-objcopy not found

Install LLVM tools:
```bash
brew install llvm
```

The script uses `/opt/homebrew/opt/llvm/bin/llvm-objcopy`.

### Legacy encoder build fails

The legacy encoder uses a private repository. If you can't access it:

1. **Skip legacy encoder rebuild** - use existing library
2. **Apply symbol prefixing to existing library**:
   ```bash
   FORCE_PROCESS_LEGACY=1 ./add_symbol_prefix.sh
   ```

### Symbol prefixing appears to fail

On macOS, `llvm-nm` might still show original symbol names even after successful prefixing due to the `symbol_wrapper.c` providing common implementations. If the build succeeds without duplicate symbol errors, the prefixing is working correctly.

## Maintenance Notes

- **Always backup libraries** before applying symbol prefixing
- **Test builds after any changes** to ensure no regressions
- **Update this guide** when adding new libraries or changing prefixes
- **Document any script modifications** for future reference

## Platform Support

Currently tested and supported on:
- **macOS darwin_arm64** (primary development platform)
- **Linux amd64** (build scripts support, limited testing)
- **Linux arm64** (build scripts support, limited testing)

Windows support would require adapting the build system and tools.