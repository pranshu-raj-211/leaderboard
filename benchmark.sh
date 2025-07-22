#!/bin/bash

IMAGE_NAME=${1:-"leaderboard"}
DOCKERFILE=${2:-"Dockerfile"}
TRIALS=${3:-5}

echo "Starting Docker build benchmarking for $IMAGE_NAME (${TRIALS} trials)..."
echo "$(date): Starting docker build benchmarking for $IMAGE_NAME (${TRIALS} trials)" >> logs/docker_benchmark.log

total_time_seconds=0
total_size_mb=0
successful_trials=0

for ((i=1; i<=TRIALS; i++)); do
    echo "Running trial $i/$TRIALS..."
    docker builder prune -f > /dev/null 2>&1
    
    start_time=$(date +%s)
    
    docker build -t $IMAGE_NAME -f $DOCKERFILE . > /dev/null 2>&1
    exit_code=$?
    
    end_time=$(date +%s)
    build_time=$((end_time - start_time))
    
    if [[ $exit_code -eq 0 ]]; then
        image_size=$(docker images $IMAGE_NAME --format "{{.Size}}")
        
        echo "Trial $i: Success - ${build_time}s, $image_size"
        echo "Trial $i: build time ${build_time}s image size $image_size time start $start_time time end $end_time" >> logs/docker_benchmark.log
        
        total_time=$((total_time + build_time))
        size_mb_total=$(awk "BEGIN {print $size_mb_total + $image_size}")
        successful_trials=$((successful_trials + 1))
    else
        echo "Trial $i: Failed (exit code: $exit_code)"
        echo "Trial $i: BUILD FAILED (exit code: $exit_code) time start $start_time time end $end_time" >> logs/docker_benchmark.log
    fi
done

echo "================================="
if [[ $successful_trials -gt 0 ]]; then
    avg_time=$(awk "BEGIN {print $total_time / $successful_trials}")
    avg_size_mb=$(awk "BEGIN {print $size_mb_total / $successful_trials}")
    
    echo "Successful trials: $successful_trials/$TRIALS"
    echo "Average time: ${avg_time}s"
    echo "Average size: $avg_size_mb"
    
    echo "$(date): Average results - time ${avg_time}s size - $avg_size_mb mb (${successful_trials}/${TRIALS} successful)" >> logs/docker_benchmark.log
else
    echo "All trials failed!"
    echo "$(date): All trials failed!" >> logs/docker_benchmark.log
fi
echo "================================="
