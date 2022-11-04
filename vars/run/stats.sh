#!/bin/bash
# Script to check network status

let oked=0
let total=0
declare -a allpeernodes=(peer1.distillery.supply.com peer2.distillery.supply.com peer1.producer1.supply.com peer1.producer2.supply.com peer1.malting.supply.com peer2.malting.supply.com peer1.bottling.supply.com peer1.maturation.supply.com peer1.retailer1.supply.com peer1.retailer2.supply.com peer1.hmrc.supply.com)
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