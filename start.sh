#!/bin/bash

# Stop and remove any existing container with the same name
docker rm -f spy-bot 2>/dev/null || true

# Build the Docker image
echo "Building Docker image..."
docker build -t spy-bot .

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Error: .env file not found!"
    exit 1
fi

# Run the container with the .env file
echo "Starting container..."
docker run -d \
    --env-file .env \
    --name spy-bot \
    --restart unless-stopped \
    spy-bot

echo "Container started successfully!"
