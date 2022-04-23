#!/bin/bash
cd $FABRIC_CFG_PATH
# cryptogen generate --config crypto-config.yaml --output keyfiles
configtxgen -profile OrdererGenesis -outputBlock genesis.block -channelID systemchannel

configtxgen -printOrg bottling-supply-com > JoinRequest_bottling-supply-com.json
configtxgen -printOrg distillery-supply-com > JoinRequest_distillery-supply-com.json
configtxgen -printOrg hmrc-supply-com > JoinRequest_hmrc-supply-com.json
configtxgen -printOrg malting-supply-com > JoinRequest_malting-supply-com.json
configtxgen -printOrg maturation-supply-com > JoinRequest_maturation-supply-com.json
configtxgen -printOrg producer1-supply-com > JoinRequest_producer1-supply-com.json
configtxgen -printOrg producer2-supply-com > JoinRequest_producer2-supply-com.json
configtxgen -printOrg retailer1-supply-com > JoinRequest_retailer1-supply-com.json
configtxgen -printOrg retailer2-supply-com > JoinRequest_retailer2-supply-com.json
