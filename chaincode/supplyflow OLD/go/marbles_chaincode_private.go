package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

//Barley Orders

type BarleyOrder struct {
	ObjectType string `json:"docType"` 
	BarleyOrderID       string `json:"BarleyOrderID"`  
	Producer      string `json:"Producer"`
	QCPass      string `json:"QCPass"`
	Status       int    `json:"Status"`
	Size      string `json:"Size"`
	Accepted      string `json:"Accepted"`
}

type BarleyPrivateOrder struct {
	ObjectType string `json:"docType"` 
	BarleyOrderID       string `json:"BarleyOrderID"` 
	Price      int    `json:"price"`
	InvoiceID	int		`json:"InvoiceID"`
}

//Malting Orders

type MaltOrder struct {
	ObjectType string `json:"docType"` 
	MaltOrderID       string `json:"MaltOrderID"`  
	BarleyOrderID       string `json:"BarleyOrderID"`  
	QCPass      string `json:"QCPass"`
	Status       int    `json:"Status"`
	Size      string `json:"Size"`
	Accepted      string `json:"Accepted"`
}

type MaltPrivateOrder struct {
	ObjectType string `json:"docType"` 
	MaltOrderID       string `json:"BarleyOrderID"` 
	Price      int    `json:"price"`
	InvoiceID	int		`json:"InvoiceID"`
}

type Distillation struct {
	ObjectType string `json:"docType"` 
	BatchID       string `json:"BatchID"`  
	MaltOrderID       string `json:"MaltOrderID"`  
	QCPass      string `json:"QCPass"`
	Status       int    `json:"Status"`
	Size      string `json:"Size"`
	InitialProof      string `json:"InitalProof"`
}

type Maturation struct {
	ObjectType string `json:"docType"` 
	CaskID       string `json:"CaskID"`  
	BatchID       string `json:"BatchID"`  
	QCPass      string `json:"QCPass"`
	Status       int    `json:"Status"`
	Size      string `json:"Size"`
	StartDate	string `json:"StartDate"`
	EndDate	string `json:"EndDate"`
	FinalProof      string `json:"InitalProof"`
	Notes	string `json:"Notes"`
	Taste	string `json:"Taste"`
}

type Bottling struct {
	ObjectType string `json:"docType"` 
	BottleID       string `json:"BottleID"`  
	CaskID       string `json:"CaskID"`  
	Status       int    `json:"Status"`
	Size      string `json:"Size"`
	PalletID      string `json:"PalletID"`
}

type RetailerOrder struct {
	ObjectType string `json:"docType"` 
	RetailerOrderID       string `json:"RetailerOrderID"`
	Shop       string `json:"Shop"`  
	PalletID       string `json:"PalletID"` 
	Status       int    `json:"Status"`
	Size      string `json:"Size"`
}

type RetailerPrivateOrder struct {
	ObjectType string `json:"docType"` 
	RetailerOrderID       string `json:"RetailerOrderID"` 
	Price      int    `json:"price"`
	InvoiceID	int		`json:"InvoiceID"`
}



type SmartContract struct {
	contractapi.Contract
}


func (s *SmartContract) InitBarleyOrder(ctx contractapi.TransactionContextInterface) error {

	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}

	transientBarleyOrderJSON, ok := transMap["InputJSON"]
	if !ok {
		return fmt.Errorf("Barley Order not found in the transient map")
	}

	var orderInput BarleyOrder

	err = json.Unmarshal(transientBarleyOrderJSON, &orderInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	if len(orderInput.BarleyOrderID) == 0 {
		return fmt.Errorf("Order ID field must be a non-empty string")
	}
	if len(orderInput.Producer) == 0 {
		return fmt.Errorf("Producer field must be a non-empty string")
	}
	if len(orderInput.Size) == 0 {
		return fmt.Errorf("Producer field must be a non-empty string")
	}

	// ==== Check if order already exists ====
	orderAsBytes, err := ctx.GetStub().GetPrivateData("collectionBarleyOrders", orderInput.BarleyOrderID)
	if err != nil {
		return fmt.Errorf("Failed to get order: " + err.Error())
	} else if marbleAsBytes != nil {
		fmt.Println("This order already exists: " + orderInput.BarleyOrderID)
		return fmt.Errorf("This order already exists: " + orderInput.BarleyOrderID)
	}

	// ==== Create marble object, marshal to JSON, and save to state ====
	order := &BarleyOrder{
		ObjectType: "BarleyOrder",
		BarleyOrderID:       orderInput.BarleyOrderID,
		Size:       orderInput.Size,
		Status:      "Ordered",
	}

	orderJSONasBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// === Save order to state ===
	err = ctx.GetStub().PutPrivateData("collectionBarleyOrders", orderInput.BarleyOrderID, orderJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put Order: %s", err.Error())
	}

	// ==== Marble saved and indexed. Return success ====

	return nil

}

// ===============================================
// readMarble - read a marble from chaincode state
// ===============================================

func (s *SmartContract) ReadBarleyOrder(ctx contractapi.TransactionContextInterface, BarleyOrderID string) (*BarleyOrder, error) {

	orderJSON, err := ctx.GetStub().GetPrivateData("collectionBarleyOrders", BarleyOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from order %s", err.Error())
	}
	if orderJSON == nil {
		return nil, fmt.Errorf("%s does not exist", BarleyOrderID)
	}

	order := new(BarleyOrder)
	_ = json.Unmarshal(orderJSON, order)

	return order, nil

}

// ===============================================
// ReadMarblePrivateDetails - read a marble private details from chaincode state
// ===============================================
func (s *SmartContract) ReadMarblePrivateDetails(ctx contractapi.TransactionContextInterface, marbleID string) (*MarblePrivateDetails, error) {

	marbleDetailsJSON, err := ctx.GetStub().GetPrivateData("collectionMarblePrivateDetails", marbleID) //get the marble from chaincode state
	if err != nil {
		return nil, fmt.Errorf("failed to read from marble details %s", err.Error())
	}
	if marbleDetailsJSON == nil {
		return nil, fmt.Errorf("%s does not exist", marbleID)
	}

	marbleDetails := new(MarblePrivateDetails)
	_ = json.Unmarshal(marbleDetailsJSON, marbleDetails)

	return marbleDetails, nil
}

// ==================================================
// delete - remove a marble key/value pair from state
// ==================================================
func (s *SmartContract) Delete(ctx contractapi.TransactionContextInterface) error {

	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}

	// Marble properties are private, therefore they get passed in transient field
	transientDeleteMarbleJSON, ok := transMap["marble_delete"]
	if !ok {
		return fmt.Errorf("marble to delete not found in the transient map")
	}

	type marbleDelete struct {
		Name string `json:"name"`
	}

	var marbleDeleteInput marbleDelete
	err = json.Unmarshal(transientDeleteMarbleJSON, &marbleDeleteInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	if len(marbleDeleteInput.Name) == 0 {
		return fmt.Errorf("name field must be a non-empty string")
	}

	// to maintain the color~name index, we need to read the marble first and get its color
	valAsbytes, err := ctx.GetStub().GetPrivateData("collectionMarbles", marbleDeleteInput.Name) //get the marble from chaincode state
	if err != nil {
		return fmt.Errorf("failed to read marble: %s", err.Error())
	}
	if valAsbytes == nil {
		return fmt.Errorf("marble private details does not exist: %s", marbleDeleteInput.Name)
	}

	var marbleToDelete Marble
	err = json.Unmarshal([]byte(valAsbytes), &marbleToDelete)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	// delete the marble from state
	err = ctx.GetStub().DelPrivateData("collectionMarbles", marbleDeleteInput.Name)
	if err != nil {
		return fmt.Errorf("Failed to delete state:" + err.Error())
	}

	// Also delete the marble from the color~name index
	indexName := "color~name"
	colorNameIndexKey, err := ctx.GetStub().CreateCompositeKey(indexName, []string{marbleToDelete.Color, marbleToDelete.Name})
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelPrivateData("collectionMarbles", colorNameIndexKey)
	if err != nil {
		return fmt.Errorf("Failed to delete marble:" + err.Error())
	}

	// Finally, delete private details of marble
	err = ctx.GetStub().DelPrivateData("collectionMarblePrivateDetails", marbleDeleteInput.Name)
	if err != nil {
		return err
	}

	return nil

}

// ===========================================================
// transfer a marble by setting a new owner name on the marble
// ===========================================================
func (s *SmartContract) TransferMarble(ctx contractapi.TransactionContextInterface) error {

	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}

	// Marble properties are private, therefore they get passed in transient field
	transientTransferMarbleJSON, ok := transMap["marble_owner"]
	if !ok {
		return fmt.Errorf("marble owner not found in the transient map")
	}

	type marbleTransferTransientInput struct {
		Name  string `json:"name"`
		Owner string `json:"owner"`
	}

	var marbleTransferInput marbleTransferTransientInput
	err = json.Unmarshal(transientTransferMarbleJSON, &marbleTransferInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	if len(marbleTransferInput.Name) == 0 {
		return fmt.Errorf("name field must be a non-empty string")
	}
	if len(marbleTransferInput.Owner) == 0 {
		return fmt.Errorf("owner field must be a non-empty string")
	}

	marbleAsBytes, err := ctx.GetStub().GetPrivateData("collectionMarbles", marbleTransferInput.Name)
	if err != nil {
		return fmt.Errorf("Failed to get marble:" + err.Error())
	} else if marbleAsBytes == nil {
		return fmt.Errorf("Marble does not exist: " + marbleTransferInput.Name)
	}

	marbleToTransfer := Marble{}
	err = json.Unmarshal(marbleAsBytes, &marbleToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	marbleToTransfer.Owner = marbleTransferInput.Owner //change the owner

	marbleJSONasBytes, _ := json.Marshal(marbleToTransfer)
	err = ctx.GetStub().PutPrivateData("collectionMarbles", marbleToTransfer.Name, marbleJSONasBytes) //rewrite the marble
	if err != nil {
		return err
	}

	return nil

}

// ===========================================================================================
// getMarblesByRange performs a range query based on the start and end keys provided.

// Read-only function results are not typically submitted to ordering. If the read-only
// results are submitted to ordering, or if the query is used in an update transaction
// and submitted to ordering, then the committing peers will re-execute to guarantee that
// result sets are stable between endorsement time and commit time. The transaction is
// invalidated by the committing peers if the result set has changed between endorsement
// time and commit time.
// Therefore, range queries are a safe option for performing update transactions based on query results.
// ===========================================================================================
func (s *SmartContract) GetMarblesByRange(ctx contractapi.TransactionContextInterface, startKey string, endKey string) ([]Marble, error) {

	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange("collectionMarbles", startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []Marble{}

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newMarble := new(Marble)

		err = json.Unmarshal(response.Value, newMarble)
		if err != nil {
			return nil, err
		}

		results = append(results, *newMarble)
	}

	return results, nil

}

// =======Rich queries =========================================================================
// Two examples of rich queries are provided below (parameterized query and ad hoc query).
// Rich queries pass a query string to the state database.
// Rich queries are only supported by state database implementations
//  that support rich query (e.g. CouchDB).
// The query string is in the syntax of the underlying state database.
// With rich queries there is no guarantee that the result set hasn't changed between
//  endorsement time and commit time, aka 'phantom reads'.
// Therefore, rich queries should not be used in update transactions, unless the
// application handles the possibility of result set changes between endorsement and commit time.
// Rich queries can be used for point-in-time queries against a peer.
// ============================================================================================

// ===== Example: Parameterized rich query =================================================
// queryMarblesByOwner queries for marbles based on a passed in owner.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (owner).
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (s *SmartContract) QueryMarblesByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]Marble, error) {

	ownerString := strings.ToLower(owner)

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"marble\",\"owner\":\"%s\"}}", ownerString)

	queryResults, err := s.getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// ===== Example: Ad hoc rich query ========================================================
// queryMarbles uses a query string to perform a query for marbles.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the queryMarblesForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (s *SmartContract) QueryMarbles(ctx contractapi.TransactionContextInterface, queryString string) ([]Marble, error) {

	queryResults, err := s.getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func (s *SmartContract) getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]Marble, error) {

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult("collectionMarbles", queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []Marble{}

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newMarble := new(Marble)

		err = json.Unmarshal(response.Value, newMarble)
		if err != nil {
			return nil, err
		}

		results = append(results, *newMarble)
	}
	return results, nil
}

// ===============================================
// getMarbleHash - use the public data hash to verify a private marble
// Result is the hash on the public ledger of a marble stored a private data collection
// ===============================================
func (s *SmartContract) GetMarbleHash(ctx contractapi.TransactionContextInterface, collection string, marbleID string) (string, error) {

	// GetPrivateDataHash can use any collection deployed with the chaincode as input
	hashAsbytes, err := ctx.GetStub().GetPrivateDataHash(collection, marbleID)
	if err != nil {
		return "", fmt.Errorf("Failed to get public data hash for marble:" + err.Error())
	} else if hashAsbytes == nil {
		return "", fmt.Errorf("Marble does not exist: " + marbleID)
	}

	return string(hashAsbytes), nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating private mables chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting private mables chaincode: %s", err.Error())
	}
}
