#!/bin/bash

# Define the number of concurrent requests
parallel_max=10

# Define total requests
total_requests=10000

# Function to make requests
make_requests() {
    local url_file=$1
    local -i total=$2
    local -i max_concurrent=$3

    # Calculate how many times to call curl based on total requests and parallelism
    local -i num_calls=$((total / max_concurrent))

    for ((i=1; i<=num_calls; i++)); do
        curl --silent --parallel --parallel-immediate --parallel-max "$max_concurrent" --config "$url_file" > /dev/null &
    done

    # Wait for all background jobs to finish
    wait
}

# Hit the load balancer URL 100 times and measure the time
echo "Hitting the load balancer URLs"
time make_requests urls.txt $total_requests $parallel_max

# Hit the single URL 100 times and measure the time
echo "Hitting the single server URL"
time make_requests url.txt $total_requests $parallel_max

