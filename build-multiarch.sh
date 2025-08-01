#!/bin/bash

set -e

echo "Building multi-architecture node_exporter with corrected approach..."

# Create build directory
mkdir -p .build

# Build for Linux AMD64 using local cross-compilation
echo "Building for linux/amd64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o .build/linux-amd64/node_exporter .

# Build for Linux ARM64 using local cross-compilation
echo "Building for linux/arm64..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o .build/linux-arm64/node_exporter .

# Create and use a new builder instance for multi-architecture builds
echo "Setting up Docker Buildx..."
docker buildx create --name multiarch-builder --use || true

# Build multi-architecture image using buildx with the corrected Dockerfile
echo "Building multi-architecture image with buildx..."
docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --tag ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${AWS_ECR_NAMESPACE}/node-exporter:multiarch \
    --file Dockerfile.multiarch \
    --push \
    .

echo ""
echo "Multi-architecture build completed successfully!"
echo ""
echo "Generated binaries:"
echo "  - .build/linux-amd64/node_exporter (Linux AMD64)"
echo "  - .build/linux-arm64/node_exporter (Linux ARM64)"
echo ""
echo "Generated multi-architecture Docker image:"
echo "  - ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${AWS_ECR_NAMESPACE}/node-exporter:multiarch (supports both AMD64 and ARM64)"
echo ""
echo "Architecture verification:"
file .build/linux-amd64/node_exporter
file .build/linux-arm64/node_exporter
echo ""
echo "To test the multi-architecture image:"
echo "  docker pull ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${AWS_ECR_NAMESPACE}/node-exporter:multiarch"
echo "  docker run --rm ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${AWS_ECR_NAMESPACE}/node-exporter:multiarch --version" 