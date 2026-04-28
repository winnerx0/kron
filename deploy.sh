#!/usr/bin/env bash
set -euo pipefail

REPO="git@github.com:winnerx0/kron.git"
BRANCH="${1:-main}"
DEPLOY_DIR="/opt/kron"

log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"; }

log "Deploying kron (branch: $BRANCH)"

if [ -d "$DEPLOY_DIR/.git" ]; then
    log "Pulling latest changes..."
    git -C "$DEPLOY_DIR" fetch origin
    git -C "$DEPLOY_DIR" checkout "$BRANCH"
    git -C "$DEPLOY_DIR" reset --hard "origin/$BRANCH"
else
    log "Cloning repository..."
    git clone --branch "$BRANCH" "$REPO" "$DEPLOY_DIR"
fi

cd "$DEPLOY_DIR"

log "Stopping existing containers..."
docker compose down --remove-orphans

log "Building images..."
docker compose build --no-cache

log "Starting services..."
docker compose up -d

log "Waiting for services to be healthy..."
sleep 5
docker compose ps

log "Deploy complete. Logs:"
docker compose logs --tail=20
