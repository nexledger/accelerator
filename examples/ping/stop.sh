#!/bin/bash

cd ../network

#Shut down the containers
docker-compose -f docker-compose.yml down

#Clean up chaincode images
docker rm $(docker ps -aq)
docker rmi $(docker images dev-* -q)

cd -