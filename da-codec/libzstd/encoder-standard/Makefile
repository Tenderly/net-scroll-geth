.PHONY: all clean build install
.PHONY: all clean build install build-linux build-all

# Detect platform
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# Allow override of target platform for cross-compilation
TARGET_PLATFORM ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')_$(shell uname -m | sed 's/x86_64/amd64/')

ifeq ($(UNAME_S),Linux)
    ifeq ($(UNAME_M),x86_64)
        PLATFORM := linux_amd64
        RUST_TARGET := x86_64-unknown-linux-gnu
    else ifeq ($(UNAME_M),aarch64)
        PLATFORM := linux_arm64
        RUST_TARGET := aarch64-unknown-linux-gnu
    endif
else ifeq ($(UNAME_S),Darwin)
    ifeq ($(UNAME_M),arm64)
        PLATFORM := darwin_arm64
        RUST_TARGET := aarch64-apple-darwin
        # Set macOS deployment target to avoid version warnings
        export MACOSX_DEPLOYMENT_TARGET := 14.0
    else ifeq ($(UNAME_M),x86_64)
        PLATFORM := darwin_amd64
        RUST_TARGET := x86_64-apple-darwin
        # Set macOS deployment target to avoid version warnings
        export MACOSX_DEPLOYMENT_TARGET := 14.0
    endif
endif

# Cross-compilation targets
LINUX_AMD64_TARGET := x86_64-unknown-linux-gnu
LINUX_ARM64_TARGET := aarch64-unknown-linux-gnu

TARGET_DIR := ../../encoding/zstd

all: build

build:
	@echo "Building standard encoder for $(PLATFORM)..."
ifeq ($(UNAME_S),Darwin)
	@echo "Setting macOS deployment target to $(MACOSX_DEPLOYMENT_TARGET)"
	MACOSX_DEPLOYMENT_TARGET=$(MACOSX_DEPLOYMENT_TARGET) cargo build --release
	MACOSX_DEPLOYMENT_TARGET=$(MACOSX_DEPLOYMENT_TARGET) cargo build --release --target $(RUST_TARGET)
else
	cargo build --release
	cargo build --release --target $(RUST_TARGET)
endif

build-linux: build-linux-amd64 build-linux-arm64

build-linux-amd64:
	@echo "Cross-compiling standard encoder for Linux AMD64..."
	rustup target add $(LINUX_AMD64_TARGET)
	cargo build --release --target $(LINUX_AMD64_TARGET)
	cp target/$(LINUX_AMD64_TARGET)/release/libencoder_standard.a $(TARGET_DIR)/libencoder_standard_linux_amd64.a
	@echo "Linux AMD64 build complete!"

build-linux-arm64:
	@echo "Cross-compiling standard encoder for Linux ARM64..."
	rustup target add $(LINUX_ARM64_TARGET)
	cargo build --release --target $(LINUX_ARM64_TARGET)
	cp target/$(LINUX_ARM64_TARGET)/release/libencoder_standard.a $(TARGET_DIR)/libencoder_standard_linux_arm64.a
	@echo "Linux ARM64 build complete!"

build-all: build build-linux

install: build
	@echo "Installing standard library to $(TARGET_DIR)..."
	cp target/release/libencoder_standard.a $(TARGET_DIR)/libencoder_standard_$(PLATFORM).a
	cp target/$(RUST_TARGET)/release/libencoder_standard.a $(TARGET_DIR)/libencoder_standard_$(PLATFORM).a
ifeq ($(UNAME_S),Darwin)
	@echo "Prefixing symbols with scroll_standard_..."
	# Create symbol redefine file for Rust runtime symbols
@@ -54,13 +83,19 @@
endif
	@echo "Standard installation complete!"

install-linux: build-linux
	@echo "Linux libraries installed!"
	@echo "Next: Run symbol prefixing script to process Linux libraries"

install-all: install install-linux

clean:
	cargo clean
	rm -f $(TARGET_DIR)/libencoder_standard_$(PLATFORM).a

info:
	@echo "Standard encoder - Platform: $(PLATFORM)"
	@echo "Target directory: $(TARGET_DIR)"
ifeq ($(UNAME_S),Darwin)
	@echo "macOS deployment target: $(MACOSX_DEPLOYMENT_TARGET)"
endif
