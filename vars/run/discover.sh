#!/bin/bash
# Script to discover endorsers and channel config
cd /vars

export PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/distillery.supply.com/users/Admin@distillery.supply.com/tls/ca.crt
export ADMINPRIVATEKEY=/vars/keyfiles/peerOrganizations/distillery.supply.com/users/Admin@distillery.supply.com/msp/keystore/priv_sk
export ADMINCERT=/vars/keyfiles/peerOrganizations/distillery.supply.com/users/Admin@distillery.supply.com/msp/signcerts/Admin@distillery.supply.com-cert.pem

discover endorsers --peerTLSCA $PEER_TLS_ROOTCERT_FILE \
  --userKey $ADMINPRIVATEKEY \
  --userCert $ADMINCERT \
  --MSP distillery-supply-com --channel supplychain \
  --server peer2.distillery.supply.com:7051 \
  --chaincode simple | jq '.[0]' | \
  jq 'del(.. | .Identity?)' | jq 'del(.. | .LedgerHeight?)' \
  > /vars/discover/supplychain_simple_endorsers.json

discover config --peerTLSCA $PEER_TLS_ROOTCERT_FILE \
  --userKey $ADMINPRIVATEKEY \
  --userCert $ADMINCERT \
  --MSP distillery-supply-com --channel supplychain \
  --server peer2.distillery.supply.com:7051 > /vars/discover/supplychain_config.json
