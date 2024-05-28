#!/bin/bash

# Build the Go server image
docker build -t my-go-server ./goserver

# Start all services using Docker Compose
docker-compose up -d
