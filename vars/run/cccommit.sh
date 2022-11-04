#!/bin/bash
# Script to instantiate chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer2.distillery.supply.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/distillery.supply.com/peers/peer2.distillery.supply.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=distillery-supply-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/distillery.supply.com/users/Admin@distillery.supply.com/msp
export ORDERER_ADDRESS=orderer3.supply.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/supply.com/orderers/orderer3.supply.com/tls/ca.crt
SID=$(peer lifecycle chaincode querycommitted -C supplychain -O json \
  | jq -r '.chaincode_definitions|.[]|select(.name=="supplyflow")|.sequence' || true)

if [[ -z $SID ]]; then
  SEQUENCE=1
else
  SEQUENCE=$((1+$SID))
fi

peer lifecycle chaincode commit -o $ORDERER_ADDRESS --channelID supplychain \
  --name supplyflow --version 2.6 --sequence $SEQUENCE \
  --peerAddresses peer1.bottling.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/bottling.supply.com/peers/peer1.bottling.supply.com/tls/ca.crt \
  --peerAddresses peer2.distillery.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/distillery.supply.com/peers/peer2.distillery.supply.com/tls/ca.crt \
  --peerAddresses peer1.hmrc.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/hmrc.supply.com/peers/peer1.hmrc.supply.com/tls/ca.crt \
  --peerAddresses peer2.malting.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/malting.supply.com/peers/peer2.malting.supply.com/tls/ca.crt \
  --peerAddresses peer1.maturation.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/maturation.supply.com/peers/peer1.maturation.supply.com/tls/ca.crt \
  --peerAddresses peer1.producer1.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/producer1.supply.com/peers/peer1.producer1.supply.com/tls/ca.crt \
  --peerAddresses peer1.producer2.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/producer2.supply.com/peers/peer1.producer2.supply.com/tls/ca.crt \
  --peerAddresses peer1.retailer1.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/retailer1.supply.com/peers/peer1.retailer1.supply.com/tls/ca.crt \
  --peerAddresses peer1.retailer2.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/retailer2.supply.com/peers/peer1.retailer2.supply.com/tls/ca.crt \
  --init-required \
  --collections-config /vars/supplyflow_collection_config.json \
  --cafile $ORDERER_TLS_CA --tls
