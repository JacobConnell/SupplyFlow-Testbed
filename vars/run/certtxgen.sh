#!/bin/bash
cd $FABRIC_CFG_PATH
# cryptogen generate --config crypto-config.yaml --output keyfiles
configtxgen -profile OrdererGenesis -outputBlock genesis.block -channelID systemchannel

configtxgen -printOrg distillery-supply-com > JoinRequest_distillery-supply-com.json
configtxgen -printOrg producer1-supply-com > JoinRequest_producer1-supply-com.json
