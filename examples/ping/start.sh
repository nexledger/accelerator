#!/bin/bash

cd ../network
docker-compose -f docker-compose.yml up -d

# Create the channel
docker exec cli peer channel create -c accelerator -f /etc/hyperledger/channel-artifacts/accelerator.tx --outputBlock /etc/hyperledger/channel-artifacts/accelerator.block -o orderer1.ordererorg1:7050

# Join the channel
docker exec -e "CORE_PEER_ADDRESS=peer1.peerorg1:7051" cli peer channel join -b /etc/hyperledger/channel-artifacts/accelerator.block -o orderer1.ordererorg1:7050
docker exec -e "CORE_PEER_ADDRESS=peer2.peerorg1:7051" cli peer channel join -b /etc/hyperledger/channel-artifacts/accelerator.block -o orderer1.ordererorg1:7050

# Install the ping chaincode
docker exec -e "CORE_PEER_ADDRESS=peer1.peerorg1:7051" cli peer chaincode install -n ping -p ping -v 1.0
docker exec -e "CORE_PEER_ADDRESS=peer2.peerorg1:7051" cli peer chaincode install -n ping -p ping -v 1.0

# Instantiate the ping chaincode
docker exec -e "CORE_PEER_ADDRESS=peer1.peerorg1:7051" cli peer chaincode instantiate -C accelerator -n ping -v 1.0 -c '{"Args":[]}' -o orderer1.ordererorg1:7050
docker exec -e "CORE_PEER_ADDRESS=peer2.peerorg1:7051" cli peer chaincode instantiate -C accelerator -n ping -v 1.0 -c '{"Args":[]}' -o orderer1.ordererorg1:7050

cd -