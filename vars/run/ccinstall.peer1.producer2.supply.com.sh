#!/bin/bash
# Script to install chaincode onto a peer node
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer1.producer2.supply.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/producer2.supply.com/peers/peer1.producer2.supply.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=producer2-supply-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/producer2.supply.com/users/Admin@producer2.supply.com/msp
cd /go/src/github.com/chaincode/supplyflow


if [ ! -f "supplyflow_go_2.6.tar.gz" ]; then
  cd go
  GO111MODULE=on
  go mod vendor
  cd -
  peer lifecycle chaincode package supplyflow_go_2.6.tar.gz \
    -p /go/src/github.com/chaincode/supplyflow/go/ \
    --lang golang --label supplyflow_2.6
fi

peer lifecycle chaincode install supplyflow_go_2.6.tar.gz
