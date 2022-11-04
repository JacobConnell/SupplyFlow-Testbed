#!/bin/bash
# Script to join a peer to a channel
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer1.retailer1.supply.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/retailer1.supply.com/peers/peer1.retailer1.supply.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=retailer1-supply-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/retailer1.supply.com/users/Admin@retailer1.supply.com/msp
export ORDERER_ADDRESS=orderer2.supply.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/supply.com/orderers/orderer2.supply.com/tls/ca.crt
if [ ! -f "supplychain.genesis.block" ]; then
  peer channel fetch oldest -o $ORDERER_ADDRESS --cafile $ORDERER_TLS_CA \
  --tls -c supplychain /vars/supplychain.genesis.block
fi

peer channel join -b /vars/supplychain.genesis.block \
  -o $ORDERER_ADDRESS --cafile $ORDERER_TLS_CA --tls
