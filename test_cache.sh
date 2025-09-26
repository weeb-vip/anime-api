#!/bin/bash

echo "Cache Testing Checklist"
echo "======================="
echo ""

echo "1. Starting Docker services..."
docker-compose up -d

echo ""
echo "2. Waiting for services to be ready..."
sleep 5

echo ""
echo "3. Testing Redis connection..."
if redis-cli ping > /dev/null 2>&1; then
    echo "✅ Redis is running"
else
    echo "❌ Redis is not responding. Make sure Docker services are running:"
    echo "   docker-compose up -d"
    exit 1
fi

echo ""
echo "4. Testing MySQL connection..."
if docker-compose exec mysql mysqladmin ping -h"localhost" -u"user" -p"mysecretpassword" > /dev/null 2>&1; then
    echo "✅ MySQL is running"
else
    echo "❌ MySQL is not responding. Check docker-compose services."
    exit 1
fi

echo ""
echo "5. Testing API server startup with cache logging..."
echo "Starting server for 10 seconds to check cache initialization..."
timeout 10s go run cmd/main.go serve 2>&1 | grep -E "(cache|Cache|redis|Redis)" || echo "No cache logs found"

echo ""
echo "6. Testing a GraphQL query to see cache behavior..."
echo "Make a GraphQL request and watch the logs for cache activity:"
echo ""
echo "curl -X POST http://localhost:3000/graphql \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -H 'Accept-Encoding: gzip' \\"
echo "  -d '{\"query\":\"{ currentlyAiring { id title { english } } }\"}'"
echo ""
echo "Look for these log messages:"
echo "- 'Cache MISS for AiringAnimeWithEpisodes' (first request)"
echo "- 'Successfully SET cache for AiringAnimeWithEpisodes' (cache storage)"
echo "- 'Cache HIT for AiringAnimeWithEpisodes' (subsequent requests)"