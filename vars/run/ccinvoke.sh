#!/bin/bash
# Script to invoke chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer2.supplier.supply.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/supplier.supply.com/peers/peer2.supplier.supply.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=supplier-supply-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/supplier.supply.com/users/Admin@supplier.supply.com/msp
export ORDERER_ADDRESS=orderer3.supply.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/supply.com/orderers/orderer3.supply.com/tls/ca.crt
peer chaincode invoke -o $ORDERER_ADDRESS --cafile $ORDERER_TLS_CA \
  --tls -C supplychain -n supplyflow  \
  --peerAddresses peer2.distillery.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/distillery.supply.com/peers/peer2.distillery.supply.com/tls/ca.crt \
  --peerAddresses peer1.shop.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/shop.supply.com/peers/peer1.shop.supply.com/tls/ca.crt \
  --peerAddresses peer2.supplier.supply.com:7051 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/supplier.supply.com/peers/peer2.supplier.supply.com/tls/ca.crt \
  --transient '{"InputJSON":"eyJCYXJsZXlPcmRlcklEIjoiMjIyIiwgIlN0YXR1cyI6IldpbGwgRG8iLCAiUHJpY2UiOjIsICJJbnZvaWNlSUQiOjEyMzR9Cg=="}' \
  -c '{"Args":["ConfirmBarleyOrder"]}'
