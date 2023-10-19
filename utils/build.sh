#!/bin/bash
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 linux/arm64 -t louvandtech/quote-bot:1.01 --push .