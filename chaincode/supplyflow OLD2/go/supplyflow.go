package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//Barley Orders

type BarleyOrder struct {
	ObjectType string `json:"docType"` 
	BarleyOrderID       string `json:"BarleyOrderID"`  
	Producer      string `json:"Producer"`
	QCPass      string `json:"QCPass"`
	Status       string    `json:"Status"`
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
	Status       string    `json:"Status"`
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
	Status       string    `json:"Status"`
	Size      string `json:"Size"`
	InitialProof      string `json:"InitalProof"`
}

type Maturation struct {
	ObjectType string `json:"docType"` 
	CaskID       string `json:"CaskID"`  
	BatchID       string `json:"BatchID"`  
	QCPass      string `json:"QCPass"`
	Status       string    `json:"Status"`
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
	Status       string    `json:"Status"`
	Size      string `json:"Size"`
	PalletID      string `json:"PalletID"`
}

type RetailerOrder struct {
	ObjectType string `json:"docType"` 
	RetailerOrderID       string `json:"RetailerOrderID"`
	Shop       string `json:"Shop"`  
	PalletID       string `json:"PalletID"` 
	Status       string    `json:"Status"`
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

//Distrillery Starts New Barley Order

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
	} else if orderAsBytes != nil {
		fmt.Println("This order already exists: " + orderInput.BarleyOrderID)
		return fmt.Errorf("This order already exists: " + orderInput.BarleyOrderID)
	}

	// ==== Create order object, marshal to JSON, and save to state ====
	order := &BarleyOrder{
		ObjectType: "BarleyOrder",
		BarleyOrderID:       orderInput.BarleyOrderID,
		Size:       orderInput.Size,
		Producer:	orderInput.Producer,
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

	// ==== BarleyOrder saved and indexed. Return success ====

	return nil
}

Producer Confirms and Agrees Price


// ===========================================================
// transfer a marble by setting a new owner name on the marble
// ===========================================================
func (s *SmartContract) ConfirmBarleyOrder(ctx contractapi.TransactionContextInterface) error {

	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}

	// BarleyOrder properties are private, therefore they get passed in transient field
	transientJSON, ok := transMap["InputJSON"]
	if !ok {
		return fmt.Errorf("Order not found in the transient map")
	}

	type OrderTransientInput struct {
		BarleyOrderID  string `json:"BarleyOrderID"`
		Status	string	`json:"Status"`
		Price int `json:"Price"`
		InvoiceID int `json:"InvoiceID"`
	}

	var OrderInput OrderTransientInput
	err = json.Unmarshal(transientJSON, &OrderInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	if len(OrderInput.BarleyOrderID) == 0 {
		return fmt.Errorf("ID field must be a non-empty string")
	}
	if len(OrderInput.Price) == 0 {
		return fmt.Errorf("price field must be a non-empty string")
	}
	if len(OrderInput.InvoiceID) == 0 {
		return fmt.Errorf("InvoiceID field must be a non-empty string")
	}

	orderAsBytes, err := ctx.GetStub().GetPrivateData("collectionBarleyOrders", OrderInput.BarleyOrderID)
	if err != nil {
		return fmt.Errorf("Failed to get order:" + err.Error())
	} else if orderAsBytes == nil {
		return fmt.Errorf("BarleyOrder does not exist: " + OrderInput.Name)
	}

	orderToUpdate := BarleyOrder{}
	err = json.Unmarshal(orderAsBytes, &orderToUpdate) 
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	orderToUpdate.Status = "OrderInput.Status"

	orderJSONasBytes, _ := json.Marshal(orderToUpdate)
	err = ctx.GetStub().PutPrivateData("collectionBarleyOrders", orderToUpdate.BarleyOrderID, orderJSONasBytes)
	if err != nil {
		return err
	}

	privOrder := &BarleyPrivateOrder{
		ObjectType: "BarleyPrivateOrder",
		BarleyOrderID:       orderInput.BarleyOrderID,
		Price:       orderInput.Price,
		InvocieID:	orderInput.InvoiceID,
	}

	orderJSONasBytes, err := json.Marshal(privOrder)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	//TODO SETUP IF STATEMENT FOR OTHER PRODUCER
	
	err = ctx.GetStub().PutPrivateData("collectionPrivateProducer1-Orders", orderInput.BarleyOrderID, orderJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put Order: %s", err.Error())
	}

	return nil


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

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating supplychain chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting supplychain chaincode: %s", err.Error())
	}
}