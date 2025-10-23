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

  # Stop and disable old service if exists
  sudo systemctl stop container-${CONTAINER_NAME}.service 2>/dev/null || true
  sudo systemctl disable container-${CONTAINER_NAME}.service 2>/dev/null || true
  sudo rm -f /etc/systemd/system/container-${CONTAINER_NAME}.service

  # Remove old container if exists
  sudo podman rm -f ${CONTAINER_NAME} 2>/dev/null || true

  # Generate systemd service (--new flag means systemd will create/remove container)
  sudo podman create \
    --name ${CONTAINER_NAME} \
    --device=/dev/ttyUSB0:/dev/ttyUSB0 \
    -p 80:80 \
    ${IMAGE_NAME}:latest

  # Generate and install systemd service
  sudo podman generate systemd --name --new ${CONTAINER_NAME} > /tmp/${CONTAINER_NAME}.service
  sudo mv /tmp/${CONTAINER_NAME}.service /etc/systemd/system/container-${CONTAINER_NAME}.service

  # Enable and start service
  sudo systemctl daemon-reload
  sudo systemctl enable container-${CONTAINER_NAME}.service
  sudo systemctl start container-${CONTAINER_NAME}.service

  # Cleanup
  rm /tmp/${IMAGE_NAME}.tar

  echo "==> Service started. Status:"
  sudo systemctl status container-${CONTAINER_NAME}.service --no-pager -l
  echo ""
  echo "==> Recent logs:"
  sudo journalctl -u container-${CONTAINER_NAME}.service -n 20 --no-pager
EOF

echo "==> Deploy complete!"
