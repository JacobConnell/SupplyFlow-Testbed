#!/bin/bash
# Script to approve chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer2.distillery.supply.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/distillery.supply.com/peers/peer2.distillery.supply.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=distillery-supply-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/distillery.supply.com/users/Admin@distillery.supply.com/msp
export ORDERER_ADDRESS=orderer3.supply.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/supply.com/orderers/orderer3.supply.com/tls/ca.crt

peer lifecycle chaincode queryinstalled -O json | jq -r '.installed_chaincodes | .[] | select(.package_id|startswith("supplyflow_4.3:"))' > ccstatus.json

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


export CORE_PEER_LOCALMSPID=distillery-supply-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/distillery.supply.com/peers/peer1.distillery.supply.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/distillery.supply.com/users/Admin@distillery.supply.com/msp
export CORE_PEER_ADDRESS=peer1.distillery.supply.com:7051

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID supplychain \
#   --name supplyflow --version 4.3 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.distillery-supply-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID supplychain --name supplyflow \
    --version 4.3 --package-id $PKID \
  --init-required \
    --collections-config /vars/supplyflow_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi

export CORE_PEER_LOCALMSPID=shop-supply-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/shop.supply.com/peers/peer2.shop.supply.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/shop.supply.com/users/Admin@shop.supply.com/msp
export CORE_PEER_ADDRESS=peer2.shop.supply.com:7051

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID supplychain \
#   --name supplyflow --version 4.3 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.shop-supply-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID supplychain --name supplyflow \
    --version 4.3 --package-id $PKID \
  --init-required \
    --collections-config /vars/supplyflow_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi

export CORE_PEER_LOCALMSPID=supplier-supply-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/supplier.supply.com/peers/peer1.supplier.supply.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/supplier.supply.com/users/Admin@supplier.supply.com/msp
export CORE_PEER_ADDRESS=peer1.supplier.supply.com:7051

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID supplychain \
#   --name supplyflow --version 4.3 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.supplier-supply-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID supplychain --name supplyflow \
    --version 4.3 --package-id $PKID \
  --init-required \
    --collections-config /vars/supplyflow_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi
