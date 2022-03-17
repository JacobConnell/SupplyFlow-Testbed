#!/bin/bash
# Script to instantiate chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer1.supplier.supply.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/supplier.supply.com/peers/peer1.supplier.supply.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=supplier-supply-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/supplier.supply.com/users/Admin@supplier.supply.com/msp
export ORDERER_ADDRESS=orderer2.supply.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/supply.com/orderers/orderer2.supply.com/tls/ca.crt
peer chaincode query -o $ORDERER_ADDRESS --cafile $ORDERER_TLS_CA \
  --tls -C supplychain -n supplyflow  \
  --peerAddresses peer1.supplier.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/supplier.supply.com/peers/peer1.supplier.supply.com/tls/ca.crt \
  -c '{"Args":["ReadBarleyOrder","222"]}'
