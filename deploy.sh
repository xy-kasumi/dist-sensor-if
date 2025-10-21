#!/bin/bash
set -e

# Configuration
REMOTE="xyx@dist-sensor.local"
IMAGE_NAME="dist-sensor-if"
CONTAINER_NAME="dist-sensor-if"
BUILD_DIR="build"
LOCAL_TAR="${BUILD_DIR}/${IMAGE_NAME}.tar"

echo "==> Copying image to Pi..."
scp $LOCAL_TAR $REMOTE:/tmp/$IMAGE_NAME.tar

echo "==> Loading and running on Pi..."
ssh $REMOTE << EOF
  # Load the image
  sudo podman load -i /tmp/${IMAGE_NAME}.tar

  # Stop and remove old container if exists
  sudo podman stop ${CONTAINER_NAME} 2>/dev/null || true
  sudo podman rm ${CONTAINER_NAME} 2>/dev/null || true

  # Run new container
  sudo podman run -d \
    --name ${CONTAINER_NAME} \
    --restart unless-stopped \
    --device=/dev/ttyUSB0:/dev/ttyUSB0 \
    -p 80:80 \
    ${IMAGE_NAME}:latest

  # Cleanup
  rm /tmp/${IMAGE_NAME}.tar

  echo "==> Container started. Logs:"
  sudo podman logs --tail 20 ${CONTAINER_NAME}
EOF

echo "==> Deploy complete!"
