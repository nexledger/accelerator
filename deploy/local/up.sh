#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

DOCKER_COMPOSE_FILE=$1
if [ -z ${DOCKER_COMPOSE_FILE} ]; then
    DOCKER_COMPOSE_FILE="docker-compose.yml"
fi

cd ${SCRIPT_DIR}

echo docker-compose -f ${DOCKER_COMPOSE_FILE} up -d
docker-compose -f ${DOCKER_COMPOSE_FILE} up -d

cd -

exit 0