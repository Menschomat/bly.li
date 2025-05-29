#!/bin/bash
services=(shortn blowup dasher perso)
base_dir="src/services"

pids=()
for service in "${services[@]}"; do
    (
        echo "========== STARTING $service =========="
        cd "$base_dir/$service" || exit 1
        go run .
    ) &
    pids+=($!)
done

trap "echo 'Stopping all services...'; kill ${pids[*]}; exit" INT
wait