#!/bin/bash
# Script to create channel block 0 and then create channel
cp $FABRIC_CFG_PATH/core.yaml /vars/core.yaml
cd /vars
export FABRIC_CFG_PATH=/vars
configtxgen -profile OrgChannel \
  -outputCreateChannelTx supplychain.tx -channelID supplychain

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer1.supplier.supply.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/supplier.supply.com/peers/peer1.supplier.supply.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=supplier-supply-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/supplier.supply.com/users/Admin@supplier.supply.com/msp
export ORDERER_ADDRESS=orderer3.supply.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/supply.com/orderers/orderer3.supply.com/tls/ca.crt
peer channel create -c supplychain -f supplychain.tx -o $ORDERER_ADDRESS \
  --cafile $ORDERER_TLS_CA --tls
