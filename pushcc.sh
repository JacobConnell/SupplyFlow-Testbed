#!/bin/bash

if [ -d "$(pwd)/vars/chaincode/supplyflow" ]; then
    rm -r $(pwd)/vars/chaincode/supplyflow
fi
  cp -r $(pwd)/chaincode/supplyflow $(pwd)/vars/chaincode/supplyflow

