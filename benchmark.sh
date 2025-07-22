#!/bin/bash

IMAGE_NAME=${1:-"leaderboard"}
DOCKERFILE=${2:-"Dockerfile"}

# TODO: multiple runs - take average
# TODO: in case of build fail, log that instead of image size and build time

echo "Benchmarking build for $IMAGE_NAME"

docker builder prune -f

start_time=$(date +%s)
docker buildx build -t $IMAGE_NAME -f $DOCKERFILE .
end_time=$(date +%s)

build_time=$((end_time - start_time))

image_size=$(docker images $IMAGE_NAME --format "{{.Size}}")

echo "================================="
echo "Build Time: ${build_time}s"
echo "Image Size: $image_size"
echo "================================="

echo "$(date): $IMAGE_NAME - Build: ${build_time}s, Size: $image_size" >> logs/docker_benchmark.log