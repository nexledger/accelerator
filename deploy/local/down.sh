#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

DOCKER_COMPOSE_FILE=$1
if [ -z ${DOCKER_COMPOSE_FILE} ]; then
    DOCKER_COMPOSE_FILE="docker-compose.yml"
fi

cd ${SCRIPT_DIR}
echo docker-compose -f ${DOCKER_COMPOSE_FILE} down
docker-compose -f ${DOCKER_COMPOSE_FILE} down --remove-orphans
docker rmi -f $(docker images dev* -q)
cd -

exit 0