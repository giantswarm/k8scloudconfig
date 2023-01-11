#!/bin/bash
env_file="/etc/apiserver-environment"

cpus="$(grep -c ^processor /proc/cpuinfo)"
memory="$(awk '/MemTotal/ { printf "%d \n", $2/1024/1024 }' /proc/meminfo)"

declare -i cpus
declare -i memory

# convert to millicores
cpus=$((cpus*1000))

# subtract dedicated millicores (see 'systemReserved' and 'kubeReserved' fields in the kubelet config file)
cpus=$((cpus-600))

# we dedicate 3/4 of all millicores to the api server
dedicated=$((cpus * 3 / 4))

# for every core we have dedicated to api server we allow 200 requests for fairness.
# 1/3 of such number is for mutating requests, 2/3 are for all requests.
echo "MAX_REQUESTS_INFLIGHT=$((dedicated / 1000 * 200 * 2 / 3))" >$env_file
echo "MAX_MUTATING_REQUESTS_INFLIGHT=$((dedicated / 1000 * 200 * 1 / 3 ))" >>$env_file
echo "CPU_LIMIT=${dedicated}m" >>$env_file
echo "MEMORY_LIMIT=$((memory * 3 / 4))Gi" >>$env_file

# We request 1/3 of the available CPUs and 1/2 of the system memory just for api server.
echo "CPU_REQUEST=$((cpus / 3))m" >>$env_file
echo "MEMORY_REQUEST=$((memory / 2))Gi" >>$env_file
