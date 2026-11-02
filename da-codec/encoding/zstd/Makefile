.PHONY: fix-symbols fix-symbols-darwin-arm64 check-llvm-objcopy backup restore clean help

LEGACY_LIB = libencoder_legacy_darwin_arm64.a
STANDARD_LIB = libencoder_standard_darwin_arm64.a
# Same default path as add_symbol_prefix.sh. Override with: make fix-symbols LLVM_OBJCOPY=<path>
LLVM_OBJCOPY ?= /opt/homebrew/opt/llvm/bin/llvm-objcopy
RENAME_SYMBOLS = \
	--redefine-sym _rust_eh_personality=scroll_ffi_rust_eh_personality \
	--redefine-sym ___rust_no_alloc_shim_is_unstable=scroll_ffi___rust_no_alloc_shim_is_unstable

help:
	@echo "Scroll-geth Symbol Fix Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  fix-symbols           - Make the Darwin ARM64 libraries from their .backup files, with the duplicate Rust symbols renamed"
	@echo "  fix-symbols-darwin-arm64 - Same as fix-symbols"
	@echo "  backup                - Backup original libraries"
	@echo "  restore               - Restore libraries from backup"
	@echo "  clean                 - Remove temporary files"
	@echo "  help                  - Show this help message"
	@echo ""
	@echo "Usage:"
	@echo "  make fix-symbols      # After rebuilding the encoder libraries"
	@echo "  make fix-symbols LLVM_OBJCOPY=<path>  # Use another llvm-objcopy"

fix-symbols: fix-symbols-darwin-arm64

fix-symbols-darwin-arm64: check-llvm-objcopy backup
	@echo "==> Making $(LEGACY_LIB) from $(LEGACY_LIB).backup..."
	@"$(LLVM_OBJCOPY)" $(RENAME_SYMBOLS) $(LEGACY_LIB).backup $(LEGACY_LIB)
	@echo "==> Making $(STANDARD_LIB) from $(STANDARD_LIB).backup..."
	@"$(LLVM_OBJCOPY)" $(RENAME_SYMBOLS) $(STANDARD_LIB).backup $(STANDARD_LIB)
	@echo "==> Symbol renaming complete!"
	@echo ""
	@echo "Verifying symbols:"
	@nm -g $(LEGACY_LIB) | grep -c scroll_ffi || true
	@nm -g $(STANDARD_LIB) | grep -c scroll_ffi || true
	@echo "✓ Symbols successfully renamed with scroll_ffi_ prefix"

check-llvm-objcopy:
	@test -x "$(LLVM_OBJCOPY)" || (echo "Error: llvm-objcopy not found at $(LLVM_OBJCOPY). Install with: brew install llvm, or run: make fix-symbols LLVM_OBJCOPY=<path>" && exit 1)

backup:
	@echo "==> Creating backups..."
	@cp -n $(LEGACY_LIB) $(LEGACY_LIB).backup 2>/dev/null || true
	@cp -n $(STANDARD_LIB) $(STANDARD_LIB).backup 2>/dev/null || true
	@echo "✓ Backups created (if they didn't exist)"

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

clean:
	@echo "==> Cleaning temporary files..."
	@rm -f $(LEGACY_LIB).fixed $(STANDARD_LIB).fixed
	@echo "✓ Clean complete"
