#!/bin/bash
set -e

# Configuration
IMAGE_NAME="dist-sensor-if"
BUILD_DIR="build"
OUTPUT_BINARY="${BUILD_DIR}/${IMAGE_NAME}"
OUTPUT_TAR="${BUILD_DIR}/${IMAGE_NAME}.tar"

# Create build directory
mkdir -p $BUILD_DIR

echo "==> Cross-compiling for arm64..."
GOOS=linux GOARCH=arm64 go build -o $OUTPUT_BINARY

echo "==> Building container image..."
podman build --platform linux/arm64 -t $IMAGE_NAME:latest .

echo "==> Saving image to ${OUTPUT_TAR}..."
rm -f $OUTPUT_TAR
podman save $IMAGE_NAME:latest -o $OUTPUT_TAR

echo "==> Build complete!"
echo "    Image (arm64):  ${OUTPUT_TAR}"
