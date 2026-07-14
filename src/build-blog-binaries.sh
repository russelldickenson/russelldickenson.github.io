#!/bin/bash

# Define the application name based on the current directory or a fixed string
APP_NAME="blog"
OUTPUT_DIR=".."

# Define target platforms: "OS/ARCH"
PLATFORMS=("linux/amd64" "darwin/arm64")

echo "Starting build process..."

for PLATFORM in "${PLATFORMS[@]}"; do
    # Split the platform string into OS and ARCH
    IFS="/" read -r OS ARCH <<< "$PLATFORM"

    # Set the output filename
    OUTPUT_FILE="${OUTPUT_DIR}/${APP_NAME}-${OS}-${ARCH}"

    # Add .exe extension if building for Windows (optional)
    if [ "$OS" == "windows" ]; then
        OUTPUT_FILE="${OUTPUT_FILE}.exe"
    fi

    echo "Building for ${OS}/${ARCH}..."

    # GOOS and GOARCH tell the Go compiler what to target
    # CGO_ENABLED=0 ensures a static binary for better portability
    # -ldflags="-s -w" strips debug information and symbols to reduce binary size
    env CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -ldflags="-s -w" -o "$OUTPUT_FILE" main.go

    if [ $? -eq 0 ]; then
        echo "Successfully built: $OUTPUT_FILE"
    else
        echo "Error: Build failed for ${OS}/${ARCH}"
        exit 1
    fi
done

echo "All builds completed."
ls -lh ${OUTPUT_DIR}/blog-*