#!/bin/bash
cd $FABRIC_CFG_PATH
# cryptogen generate --config crypto-config.yaml --output keyfiles
configtxgen -profile OrdererGenesis -outputBlock genesis.block -channelID systemchannel

configtxgen -printOrg distillery-supply-com > JoinRequest_distillery-supply-com.json
configtxgen -printOrg shop-supply-com > JoinRequest_shop-supply-com.json
configtxgen -printOrg supplier-supply-com > JoinRequest_supplier-supply-com.json
