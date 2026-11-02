.PHONY: fix-symbols fix-symbols-darwin-arm64 check-llvm-tools check-backup-inputs restore help

LEGACY_LIB = libencoder_legacy_darwin_arm64.a
STANDARD_LIB = libencoder_standard_darwin_arm64.a
# Same default paths as add_symbol_prefix.sh. Override with: make fix-symbols LLVM_OBJCOPY=<path> LLVM_NM=<path>
LLVM_OBJCOPY ?= /opt/homebrew/opt/llvm/bin/llvm-objcopy
LLVM_NM ?= /opt/homebrew/opt/llvm/bin/llvm-nm
RENAME_SYMBOLS = \
	--redefine-sym _rust_eh_personality=scroll_ffi_rust_eh_personality \
	--redefine-sym ___rust_no_alloc_shim_is_unstable=scroll_ffi___rust_no_alloc_shim_is_unstable

help:
	@echo "Scroll-geth Symbol Fix Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  fix-symbols           - Make the Darwin ARM64 libraries from their .backup files, with the duplicate Rust symbols renamed"
	@echo "  fix-symbols-darwin-arm64 - Same as fix-symbols"
	@echo "  restore               - Copy the .backup files over the libraries (undo the rename)"
	@echo "  help                  - Show this help message"
	@echo ""
	@echo "Usage:"
	@echo "  make fix-symbols      # After rebuilding the encoder libraries"
	@echo "  make fix-symbols LLVM_OBJCOPY=<path> LLVM_NM=<path>  # Use other LLVM tools"

fix-symbols: fix-symbols-darwin-arm64

fix-symbols-darwin-arm64: check-llvm-tools check-backup-inputs
	@echo "==> Making $(LEGACY_LIB) from $(LEGACY_LIB).backup..."
	@"$(LLVM_OBJCOPY)" $(RENAME_SYMBOLS) $(LEGACY_LIB).backup $(LEGACY_LIB)
	@echo "==> Making $(STANDARD_LIB) from $(STANDARD_LIB).backup..."
	@"$(LLVM_OBJCOPY)" $(RENAME_SYMBOLS) $(STANDARD_LIB).backup $(STANDARD_LIB)
	@echo "==> Symbol renaming complete!"
	@echo ""
	@echo "Verifying symbols (Mach-O symbol table, not the LLVM bitcode):"
	@"$(LLVM_NM)" -g --no-llvm-bc $(LEGACY_LIB) | grep scroll_ffi | sort | uniq -c || true
	@"$(LLVM_NM)" -g --no-llvm-bc $(STANDARD_LIB) | grep scroll_ffi | sort | uniq -c || true
	@echo "✓ Symbols successfully renamed with scroll_ffi_ prefix"

check-llvm-tools:
	@test -x "$(LLVM_OBJCOPY)" || (echo "Error: llvm-objcopy not found at $(LLVM_OBJCOPY). Install with: brew install llvm, or run: make fix-symbols LLVM_OBJCOPY=<path>" && exit 1)
	@test -x "$(LLVM_NM)" || (echo "Error: llvm-nm not found at $(LLVM_NM). Install with: brew install llvm, or run: make fix-symbols LLVM_NM=<path>" && exit 1)

# The .backup files are the inputs: the libraries after add_symbol_prefix.sh and before the rename.
# make never writes them.
check-backup-inputs:
	@for f in $(LEGACY_LIB).backup $(STANDARD_LIB).backup; do \
		test -f "$$f" || { echo "Error: $$f not found. It is the input before the rename: run add_symbol_prefix.sh on the new libraries, then copy each Darwin ARM64 library to its .backup file."; exit 1; }; \
	done

restore:
	@echo "==> Restoring from backups..."
	@if [ -f $(LEGACY_LIB).backup ]; then \
		cp $(LEGACY_LIB).backup $(LEGACY_LIB); \
		echo "✓ Restored $(LEGACY_LIB)"; \
	else \
		echo "✗ No backup found for $(LEGACY_LIB)"; \
	fi
	@if [ -f $(STANDARD_LIB).backup ]; then \
		cp $(STANDARD_LIB).backup $(STANDARD_LIB); \
		echo "✓ Restored $(STANDARD_LIB)"; \
	else \
		echo "✗ No backup found for $(STANDARD_LIB)"; \
	fi
