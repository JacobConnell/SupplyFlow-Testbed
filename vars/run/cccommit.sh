#!/bin/bash
# Script to instantiate chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer1.supplier.supply.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/supplier.supply.com/peers/peer1.supplier.supply.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=supplier-supply-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/supplier.supply.com/users/Admin@supplier.supply.com/msp
export ORDERER_ADDRESS=orderer1.supply.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/supply.com/orderers/orderer1.supply.com/tls/ca.crt
SID=$(peer lifecycle chaincode querycommitted -C supplychain -O json \
  | jq -r '.chaincode_definitions|.[]|select(.name=="supplyflow")|.sequence' || true)

if [[ -z $SID ]]; then
  SEQUENCE=1
else
  SEQUENCE=$((1+$SID))
fi

peer lifecycle chaincode commit -o $ORDERER_ADDRESS --channelID supplychain \
  --name supplyflow --version 4.8 --sequence $SEQUENCE \
  --peerAddresses peer1.distillery.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/distillery.supply.com/peers/peer1.distillery.supply.com/tls/ca.crt \
  --peerAddresses peer2.shop.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/shop.supply.com/peers/peer2.shop.supply.com/tls/ca.crt \
  --peerAddresses peer1.supplier.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/supplier.supply.com/peers/peer1.supplier.supply.com/tls/ca.crt \
  --init-required \
  --collections-config /vars/supplyflow_collection_config.json \
  --cafile $ORDERER_TLS_CA --tls
