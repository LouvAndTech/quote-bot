#!/bin/bash
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t louvandtech/quote-bot:2.1 -t louvandtech/quote-bot:latest --push .