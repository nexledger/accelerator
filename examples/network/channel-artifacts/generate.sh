#!/bin/bash

rm genesis.block
rm accelerator.tx

FABRIC_CFG_PATH=. ./configtxgen -profile genesis -outputBlock genesis.block -channelID default
FABRIC_CFG_PATH=. ./configtxgen -profile accelerator -outputCreateChannelTx accelerator.tx -outputAnchorPeersUpdate accelerator-anchor.tx -asOrg peerorg1msp -channelID accelerator