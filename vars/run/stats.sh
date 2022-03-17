#!/bin/bash
# Script to check network status

let oked=0
let total=0
declare -a allpeernodes=(peer1.distillery.supply.com peer2.distillery.supply.com peer1.shop.supply.com peer2.shop.supply.com peer1.supplier.supply.com peer2.supplier.supply.com)
for anode in ${allpeernodes[@]}; do
  let total=1+$total
  ss=$(wget -O- -S ${anode}:7061/healthz | jq '.status')
  printf "%20s %s\n" $anode $ss
  if [ $ss == '"OK"' ]; then
    let oked=1+$oked
  fi
done

declare -a allorderernodes=(orderer1.supply.com orderer2.supply.com orderer3.supply.com)
for anode in ${allorderernodes[@]}; do
  let total=1+$total
  ss=$(wget -O- -S ${anode}:7060/healthz | jq '.status')
  printf "%20s %s\n" $anode $ss
  if [ $ss == '"OK"' ]; then
    let oked=1+$oked
  fi
done

let percent=$oked*100/$total
echo "Network Status: $percent%"