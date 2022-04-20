# Install the chaincode
minifab install -n privatemarbles -r true

# Modify the vars/privatemarbles_collection_config.json with the following content

Order=$( echo '{"BarleyOrderID":"222", "Status":"Will Do", "Price":2, "InvoiceID":1234}' | base64 | tr -d \\n )

./minifab query -p '"ReadBarleyOrder","222"' -t ''

Order=$( echo '{"BarleyOrderID":"222","Size":"5OKilo", "Producer":"NA", "QCPass":"NA", "Status":"NA", "Accepted":"NA"}' | base64 | tr -d \\n )

producer1-supply-com
producer2-supply-com
malting-supply-com
distillery-supply-com
maturation-supply-com
bottling-supply-com
retailer1-supply-com
retailer2-supply-com

Order

Order=$( echo '{"BarleyOrderID”:”222”,”Size":"5OKilo", "Producer":"NA", "QCPass":"NA", "Status":"NA", "Accepted":"NA"}' | base64 | tr -d \\n )


Order=$( echo '{"BarleyOrderID”:”33445”,”Size":"51Kilo", "Producer":"NA”, "Status":"NA”}’ | base64 | tr -d \\n )

Order=$( echo '{"BarleyOrderID”:”12345”,”Size":"51Kilo", "Producer”:”Producer1”}' | base64 | tr -d \\n )
    

./minifab invoke -p '"InitBarleyOrder"' -t '{"InputJSON":"'$Order'"}'

./minifab query -p '”ReadBarleyOrder”,”222”’ -t ''



Order=$( echo '{"BarleyOrderID”:”444553343”, "Status":"NA", “Price”:2, “InvoiceID”:”1234"}’ | base64 | tr -d \\n )
./minifab invoke -p '"ConfirmBarleyOrder"' -t '{"InputJSON":"'$Order'"}'

[
 {
    "name": "collectionMarbles",
    "policy": "OR( 'org0examplecom.member', 'org1examplecom.member' )",
    "requiredPeerCount": 0,
    "maxPeerCount": 3,
    "blockToLive":1000000,
    "memberOnlyRead": true
 },
 {
    "name": "collectionMarblePrivateDetails",
    "policy": "OR( 'org0examplecom.member' )",
    "requiredPeerCount": 0,
    "maxPeerCount": 3,
    "blockToLive":3,
    "memberOnlyRead": true
 }
]
```
# Approve,commit,initialize the chaincode
    ./minifab approve,commit,initialize -p ''


	ObjectType string `json:"docType"` 
	BarleyOrderID       string `json:"BarleyOrderID"`  
	Producer      string `json:"Producer"`
	QCPass      string `json:"QCPass"`
	Status       string    `json:"Status"`
	Size      string `json:"Size"`
	Accepted      string `json:"Accepted"`

# To init marble
    Order=$( echo '{"BarleyOrderID":"444553343","Size":"5OKilo", "Producer":"NA", "QCPass":"NA", "Status":"NA", "Accepted":"NA"}' | base64 | tr -d \\n )
    ./minifab invoke -p '"InitBarleyOrder"' -t '{"InputJSON":"'$Order'"}'

    MARBLE=$( echo '{"name":"marble2","color":"red","size":50,"owner":"tom","price":102}' | base64 | tr -d \\n )
    minifab invoke -p '"initMarble"' -t '{"marble":"'$MARBLE'"}'

    MARBLE=$( echo '{"name":"marble3","color":"blue","size":70,"owner":"tom","price":103}' | base64 | tr -d \\n )
    minifab invoke -p '"initMarble"' -t '{"marble":"'$MARBLE'"}'

# To transfer marble
    MARBLE_OWNER=$( echo '{"name":"marble2","owner":"jerry"}' | base64 | tr -d \\n )
    minifab invoke -p '"transferMarble"' -t '{"marble_owner":"'$MARBLE_OWNER'"}'

# To query marble
    minifab query -p '"readMarble","marble1"' -t ''
    minifab query -p '"readMarblePrivateDetails","marble1"' -t ''
    minifab query -p '"getMarblesByRange","marble1","marble4"' -t ''

# To delete marble
    MARBLE_ID=$( echo '{"name":"marble1"}' | base64 | tr -d \\n )
    minifab invoke -p '"delete"' -t '{"marble_delete":"'$MARBLE_ID'"}'
