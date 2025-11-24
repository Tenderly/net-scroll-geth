.PHONY: fix-symbols fix-symbols-darwin-arm64 backup restore clean help

LEGACY_LIB = libencoder_legacy_darwin_arm64.a
STANDARD_LIB = libencoder_standard_darwin_arm64.a
OBJCONV_CONF = objconv.conf

help:
	@echo "Scroll-geth Symbol Fix Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  fix-symbols           - Fix duplicate Rust symbols in Darwin ARM64 libraries"
	@echo "  fix-symbols-darwin-arm64 - Same as fix-symbols"
	@echo "  backup                - Backup original libraries"
	@echo "  restore               - Restore libraries from backup"
	@echo "  clean                 - Remove temporary files"
	@echo "  help                  - Show this help message"
	@echo ""
	@echo "Usage:"
	@echo "  make fix-symbols      # After updating or linking scroll-geth"

fix-symbols: fix-symbols-darwin-arm64

fix-symbols-darwin-arm64: check-objconv backup
	@echo "==> Fixing symbols in $(LEGACY_LIB)..."
	@objconv @$(OBJCONV_CONF) $(LEGACY_LIB) $(LEGACY_LIB).fixed
	@mv $(LEGACY_LIB).fixed $(LEGACY_LIB)
	@echo "==> Fixing symbols in $(STANDARD_LIB)..."
	@objconv @$(OBJCONV_CONF) $(STANDARD_LIB) $(STANDARD_LIB).fixed
	@mv $(STANDARD_LIB).fixed $(STANDARD_LIB)
	@echo "==> Symbol renaming complete!"
	@echo ""
	@echo "Verifying symbols:"
	@nm -g $(LEGACY_LIB) | grep -c scroll_ffi || true
	@nm -g $(STANDARD_LIB) | grep -c scroll_ffi || true
	@echo "✓ Symbols successfully renamed with scroll_ffi_ prefix"

check-objconv:
	@which objconv > /dev/null || (echo "Error: objconv not found. Install with: brew install objconv" && exit 1)

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
