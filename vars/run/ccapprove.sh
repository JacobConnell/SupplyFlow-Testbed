#!/bin/bash
# Script to approve chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer2.distillery.supply.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/distillery.supply.com/peers/peer2.distillery.supply.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=distillery-supply-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/distillery.supply.com/users/Admin@distillery.supply.com/msp
export ORDERER_ADDRESS=orderer1.supply.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/supply.com/orderers/orderer1.supply.com/tls/ca.crt

peer lifecycle chaincode queryinstalled -O json | jq -r '.installed_chaincodes | .[] | select(.package_id|startswith("supplyflow_2.6:"))' > ccstatus.json

PKID=$(jq '.package_id' ccstatus.json | xargs)
REF=$(jq '.references.supplychain' ccstatus.json)

SID=$(peer lifecycle chaincode querycommitted -C supplychain -O json \
  | jq -r '.chaincode_definitions|.[]|select(.name=="supplyflow")|.sequence' || true)
if [[ -z $SID ]]; then
  SEQUENCE=1
elif [[ -z $REF ]]; then
  SEQUENCE=$SID
else
  SEQUENCE=$((1+$SID))
fi


export CORE_PEER_LOCALMSPID=bottling-supply-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/bottling.supply.com/peers/peer1.bottling.supply.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/bottling.supply.com/users/Admin@bottling.supply.com/msp
export CORE_PEER_ADDRESS=peer1.bottling.supply.com:7051

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID supplychain \
#   --name supplyflow --version 2.6 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.bottling-supply-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID supplychain --name supplyflow \
    --version 2.6 --package-id $PKID \
  --init-required \
    --collections-config /vars/supplyflow_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi

export CORE_PEER_LOCALMSPID=distillery-supply-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/distillery.supply.com/peers/peer1.distillery.supply.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/distillery.supply.com/users/Admin@distillery.supply.com/msp
export CORE_PEER_ADDRESS=peer1.distillery.supply.com:7051

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID supplychain \
#   --name supplyflow --version 2.6 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.distillery-supply-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID supplychain --name supplyflow \
    --version 2.6 --package-id $PKID \
  --init-required \
    --collections-config /vars/supplyflow_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi

export CORE_PEER_LOCALMSPID=hmrc-supply-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/hmrc.supply.com/peers/peer1.hmrc.supply.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/hmrc.supply.com/users/Admin@hmrc.supply.com/msp
export CORE_PEER_ADDRESS=peer1.hmrc.supply.com:7051

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID supplychain \
#   --name supplyflow --version 2.6 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.hmrc-supply-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID supplychain --name supplyflow \
    --version 2.6 --package-id $PKID \
  --init-required \
    --collections-config /vars/supplyflow_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi

export CORE_PEER_LOCALMSPID=malting-supply-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/malting.supply.com/peers/peer2.malting.supply.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/malting.supply.com/users/Admin@malting.supply.com/msp
export CORE_PEER_ADDRESS=peer2.malting.supply.com:7051

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID supplychain \
#   --name supplyflow --version 2.6 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.malting-supply-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID supplychain --name supplyflow \
    --version 2.6 --package-id $PKID \
  --init-required \
    --collections-config /vars/supplyflow_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi

export CORE_PEER_LOCALMSPID=maturation-supply-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/maturation.supply.com/peers/peer1.maturation.supply.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/maturation.supply.com/users/Admin@maturation.supply.com/msp
export CORE_PEER_ADDRESS=peer1.maturation.supply.com:7051

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID supplychain \
#   --name supplyflow --version 2.6 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.maturation-supply-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID supplychain --name supplyflow \
    --version 2.6 --package-id $PKID \
  --init-required \
    --collections-config /vars/supplyflow_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi

export CORE_PEER_LOCALMSPID=producer1-supply-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/producer1.supply.com/peers/peer1.producer1.supply.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/producer1.supply.com/users/Admin@producer1.supply.com/msp
export CORE_PEER_ADDRESS=peer1.producer1.supply.com:7051

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID supplychain \
#   --name supplyflow --version 2.6 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.producer1-supply-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID supplychain --name supplyflow \
    --version 2.6 --package-id $PKID \
  --init-required \
    --collections-config /vars/supplyflow_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi

export CORE_PEER_LOCALMSPID=producer2-supply-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/producer2.supply.com/peers/peer1.producer2.supply.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/producer2.supply.com/users/Admin@producer2.supply.com/msp
export CORE_PEER_ADDRESS=peer1.producer2.supply.com:7051

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID supplychain \
#   --name supplyflow --version 2.6 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.producer2-supply-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID supplychain --name supplyflow \
    --version 2.6 --package-id $PKID \
  --init-required \
    --collections-config /vars/supplyflow_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi

export CORE_PEER_LOCALMSPID=retailer1-supply-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/retailer1.supply.com/peers/peer1.retailer1.supply.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/retailer1.supply.com/users/Admin@retailer1.supply.com/msp
export CORE_PEER_ADDRESS=peer1.retailer1.supply.com:7051

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID supplychain \
#   --name supplyflow --version 2.6 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.retailer1-supply-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID supplychain --name supplyflow \
    --version 2.6 --package-id $PKID \
  --init-required \
    --collections-config /vars/supplyflow_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi

export CORE_PEER_LOCALMSPID=retailer2-supply-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/retailer2.supply.com/peers/peer1.retailer2.supply.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/retailer2.supply.com/users/Admin@retailer2.supply.com/msp
export CORE_PEER_ADDRESS=peer1.retailer2.supply.com:7051

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID supplychain \
#   --name supplyflow --version 2.6 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.retailer2-supply-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID supplychain --name supplyflow \
    --version 2.6 --package-id $PKID \
  --init-required \
    --collections-config /vars/supplyflow_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi
