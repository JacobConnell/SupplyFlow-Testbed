#!/bin/bash
# Script to instantiate chaincode
cp $FABRIC_CFG_PATH/core.yaml /vars/core.yaml
cd /vars
export FABRIC_CFG_PATH=/vars

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer2.distillery.supply.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/distillery.supply.com/peers/peer2.distillery.supply.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=distillery-supply-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/distillery.supply.com/users/Admin@distillery.supply.com/msp
export ORDERER_ADDRESS=orderer1.supply.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/supply.com/orderers/orderer1.supply.com/tls/ca.crt

peer channel fetch config config_block.pb -o $ORDERER_ADDRESS \
  --cafile $ORDERER_TLS_CA --tls -c supplychain

configtxlator proto_decode --input config_block.pb --type common.Block \
  | jq .data.data[0].payload.data.config > supplychain_config.json