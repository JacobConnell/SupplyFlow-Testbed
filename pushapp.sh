#!/bin/bash

if [ -d "$(pwd)/vars/app" ]; then
    rm -r $(pwd)/vars/app
fi
  cp -r $(pwd)/app $(pwd)/vars/app

