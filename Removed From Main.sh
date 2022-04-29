Removed From Main
#if [ ! -d "$(pwd)/vars/chaincode" ]; then
#  cp -r $(pwd)/chaincode $(pwd)/vars/
#fi
#if [ ! -d "$(pwd)/vars/app" ]; then
#  cp -r $(pwd)/app $(pwd)/vars/
#fi

#if [ ! -f "$(pwd)/vars/collection_config.json" ]; then
#  cp $(pwd)/collection_config.json $(pwd)/vars/collection_config.json
#fi