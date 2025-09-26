#!/bin/bash

echo "Redis Connection Pool Performance Test"
echo "======================================"
echo ""

echo "1. Starting Redis (if not running)..."
if ! redis-cli ping > /dev/null 2>&1; then
    echo "Starting Redis with Docker..."
    docker run -d --name redis-test -p 6379:6379 redis:alpine
    sleep 3
fi

echo "✅ Redis is running"

echo ""
echo "2. Testing Redis latency directly..."
echo "Direct Redis latency (should be sub-millisecond):"
redis-cli --latency -i 1 -c 5

echo ""
echo "3. Connection Pool Configuration:"
echo "- Pool Size: 10 connections"
echo "- Min Idle: 3 connections"
echo "- Max Idle: 6 connections"
echo "- Read Timeout: 1000ms"
echo "- Write Timeout: 1000ms"
echo "- Dial Timeout: 2000ms"

echo ""
echo "4. Expected improvements:"
echo "- Redis.Get operations should be < 1ms (was 4.23ms)"
echo "- Connection reuse should show in pool stats"
echo "- Pool hits should increase with repeated requests"

echo ""
echo "5. Monitor your traces for these new attributes:"
echo "- redis.pool.hits"
echo "- redis.pool.misses"
echo "- redis.pool.total_conns"
echo "- redis.pool.idle_conns"
echo "- redis.pool.timeouts"

echo ""
echo "6. To test with your app:"
echo "   go run cmd/main.go serve"
echo ""
echo "7. Make multiple GraphQL requests and watch for:"
echo "   - Decreasing Redis.Get times"
echo "   - Increasing pool.hits in traces"
echo "   - Stable pool.total_conns"

echo ""
echo "Pool performance should show:"
echo "- First request: pool.misses > 0 (creating connections)"
echo "- Subsequent requests: pool.hits > 0 (reusing connections)"
echo "- Redis operations: < 1ms (down from 4.23ms)"