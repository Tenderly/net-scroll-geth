.PHONY: all clean build install

# Detect platform
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Linux)
    ifeq ($(UNAME_M),x86_64)
        PLATFORM := linux_amd64
    else ifeq ($(UNAME_M),aarch64)
        PLATFORM := linux_arm64
    endif
else ifeq ($(UNAME_S),Darwin)
    ifeq ($(UNAME_M),arm64)
        PLATFORM := darwin_arm64
        # Set macOS deployment target to avoid version warnings
        export MACOSX_DEPLOYMENT_TARGET := 14.0
    else ifeq ($(UNAME_M),x86_64)
        PLATFORM := darwin_amd64
        # Set macOS deployment target to avoid version warnings
        export MACOSX_DEPLOYMENT_TARGET := 14.0
    endif
endif

TARGET_DIR := ../../encoding/zstd

all: build

build:
	@echo "Building standard encoder for $(PLATFORM)..."
ifeq ($(UNAME_S),Darwin)
	@echo "Setting macOS deployment target to $(MACOSX_DEPLOYMENT_TARGET)"
	MACOSX_DEPLOYMENT_TARGET=$(MACOSX_DEPLOYMENT_TARGET) cargo build --release
else
	cargo build --release
endif

install: build
	@echo "Installing standard library to $(TARGET_DIR)..."
	cp target/release/libencoder_standard.a $(TARGET_DIR)/libencoder_standard_$(PLATFORM).a
ifeq ($(UNAME_S),Darwin)
	@echo "Prefixing symbols with scroll_standard_..."
	# Create symbol redefine file for Rust runtime symbols
	/opt/homebrew/opt/llvm/bin/llvm-nm $(TARGET_DIR)/libencoder_standard_$(PLATFORM).a | \
	awk '/^[0-9a-fA-F]+ [TDBS] / { if ($$3 ~ /^(_rust_|___r[dg]l_)/ || $$3 == "___rust_foreign_exception" || $$3 == "___rust_drop_panic" || $$3 == "___rust_panic_cleanup") print $$3 " scroll_standard_" $$3 }' > redefine_standard.syms
	# Apply symbol prefixing
	if [ -s redefine_standard.syms ]; then \
		/opt/homebrew/opt/llvm/bin/llvm-objcopy --redefine-syms=redefine_standard.syms $(TARGET_DIR)/libencoder_standard_$(PLATFORM).a $(TARGET_DIR)/libencoder_standard_$(PLATFORM)_new.a && \
		mv $(TARGET_DIR)/libencoder_standard_$(PLATFORM)_new.a $(TARGET_DIR)/libencoder_standard_$(PLATFORM).a; \
	fi
	rm -f redefine_standard.syms
	@echo "Fixing library symbol table..."
	ranlib $(TARGET_DIR)/libencoder_standard_$(PLATFORM).a
endif
	@echo "Standard installation complete!"

clean:
	cargo clean
	rm -f $(TARGET_DIR)/libencoder_standard_$(PLATFORM).a

info:
	@echo "Standard encoder - Platform: $(PLATFORM)"
	@echo "Target directory: $(TARGET_DIR)"
ifeq ($(UNAME_S),Darwin)
	@echo "macOS deployment target: $(MACOSX_DEPLOYMENT_TARGET)"
endif
