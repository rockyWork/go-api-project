#!/usr/bin/env bash
# wait-for-it.sh - 等待服务就绪
# 来源: https://github.com/vishnubob/wait-for-it

set -e

host="$1"
port="$2"
shift 2
cmd="$@"

timeout=30
echo "Waiting for $host:$port..."

for i in $(seq $timeout); do
    if nc -z "$host" "$port" 2>/dev/null; then
        echo "$host:$port is available!"
        exec $cmd
    fi
    sleep 1
done

echo "Timeout waiting for $host:$port"
exit 1
