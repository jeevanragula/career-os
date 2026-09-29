#!/usr/bin/env bash
set -euo pipefail

command -v docker >/dev/null || { echo "Docker is required"; exit 1; }

if [ ! -f .env ]; then cp .env.example .env; fi

echo "Starting CareerOS..."
docker compose up --build -d

echo
echo "CareerOS: http://localhost:8080"
echo "Health:   http://localhost:8080/healthz"
echo
echo "Edit .env to configure an AI provider, then:"
echo "  docker compose restart careeros"
