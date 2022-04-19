package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"time"
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
	GeoLocation      string `json:"GeoLocation"`
	SoilPH      string `json:"SoilPH"`
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
	InitialProof      string `json:"InitialProof"`
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
	FinalProof      string `json:"FinalProof"`
	DutyPaid	string `json:"DutyPaid"`
	Notes	string `json:"Notes"`
	Taste	string `json:"Taste"`
	Age		string `json:"Age"`
}

type Bottling struct {
	ObjectType string `json:"docType"` 
	BottleID       string `json:"BottleID"`  
	CaskID       []string `json:"CaskID"`  
	Status       string    `json:"Status"`
	Duty	string    `json:"Duty"` 
	Age	string    `json:"Age"`
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

type CaskLifeModel struct {
    ObjectType string `json:"docType"` 
    CaskID string `json:”CaskID”`
    QCPass  string `json:"QCPass"`
    BatchID  string `json:”BatchID”`
	BatchQCPass  string `json:"BatchQCPass"`
    MaltID  string `json:”MaltID”`
	MaltQCPass  string `json:"MaltQCPass"`
    BarleyOrderID  string `json:”BarleyOrderID”`
	Producer    string `json:"Producer"`
    SoilPH  string `json:"SoilPH"`
    GeoLocation string `json:"GeoLocation"`
}

type BottleLifeModel struct {
    ObjectType string `json:"docType"` 
    BottleID    string `json:"BottleID"`
	Age		string `json:"Age"`
	Duty	string `json:"Duty"`
    Casks []CaskLifeModel `json:"Casks"`
}

//Malting Starts New Barley Order
func (s *SmartContract) InitBarleyOrder(ctx contractapi.TransactionContextInterface) error {
	//TODO: CHECK MSP IS OF MALTING AND NO ONE ELSE
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
	OrderID := ("BARLEY"+orderInput.BarleyOrderID) 
	orderAsBytes, err := ctx.GetStub().GetState(OrderID)
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
	err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put Order: %s", err.Error())
	}

	// ==== BarleyOrder saved and indexed. Return success ====

	return nil
}

func (s *SmartContract) ConfirmBarleyOrder(ctx contractapi.TransactionContextInterface) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}
	fmt.Println("Failed: ", msp)
	//TODO: Limit to check the producer field of the Supplier.
	if (msp == "producer1-supply-com" || msp == "producer2-supply-com") {
		transMap, err := ctx.GetStub().GetTransient()
		if err != nil {
			return fmt.Errorf("Error getting transient: " + err.Error())
		}

		transientJSON, ok := transMap["InputJSON"]
		if !ok {
			return fmt.Errorf("Order not found in the transient map")
		}

		type OrderTransientInput struct {
			BarleyOrderID  string `json:"BarleyOrderID"`
			Status	string	`json:"Status"`
			GeoLocation      string `json:"GeoLocation"`
			SoilPH      string `json:"SoilPH"`
			Price int `json:"Price"`
			InvoiceID int `json:"InvoiceID"`
		}

		var OrderInput OrderTransientInput
		err = json.Unmarshal(transientJSON, &OrderInput)
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		OrderID := ("BARLEY"+ OrderInput.BarleyOrderID) 

		if len(OrderInput.BarleyOrderID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}
		if OrderInput.Price == 0 {
			return fmt.Errorf("price field must not be nil")
		}
		if OrderInput.InvoiceID == 0 {
			return fmt.Errorf("InvoiceID field must not be nil")
		}

		orderAsBytes, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("Failed to get order:" + err.Error())
		} else if orderAsBytes == nil {
			return fmt.Errorf("BarleyOrder does not exist: " + OrderInput.BarleyOrderID)
		}

		orderToUpdate := BarleyOrder{}
		err = json.Unmarshal(orderAsBytes, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		orderToUpdate.Status = (OrderInput.Status)
		orderToUpdate.GeoLocation =	OrderInput.GeoLocation
		orderToUpdate.SoilPH	=	OrderInput.SoilPH

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

		privOrder := &BarleyPrivateOrder{
			ObjectType: "BarleyPrivateOrder",
			BarleyOrderID:       OrderInput.BarleyOrderID,
			Price:       OrderInput.Price,
			InvoiceID:	OrderInput.InvoiceID,
		}

		orderPrivJSONasBytes, err := json.Marshal(privOrder)
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		//TODO: SETUP IF STATEMENT FOR OTHER PRODUCER TO CONTROL ACCESS TO BOTH PRIVATE STATES
		//If Prod 1
		err = ctx.GetStub().PutPrivateData("collectionPrivateProducer1-Orders", OrderID, orderPrivJSONasBytes)
		//If Prod 2

		//err = ctx.GetStub().PutPrivateData("collectionPrivateProducer2-Orders", OrderID, orderPrivJSONasBytes)
		if err != nil {
			return fmt.Errorf("failed to put Order: %s", err.Error())
		}

		//ELSE ERROR

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) ShipBarleyOrder(ctx contractapi.TransactionContextInterface, BarleyOrderID string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)
	//TODO: Limit to check the producer field of the Supplier.
	if (msp == "producer1-supply-com" || msp == "producer2-supply-com") {
		
		if len(BarleyOrderID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}

		OrderID := ("BARLEY"+BarleyOrderID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", BarleyOrderID)
		}

		orderToUpdate := BarleyOrder{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		orderToUpdate.Status = "Shipped"

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) AcceptBarleyOrder(ctx contractapi.TransactionContextInterface, BarleyOrderID string, Accepted string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)
	if msp == "malting-supply-com" {
		
		if len(BarleyOrderID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}
		
		OrderID := ("BARLEY"+BarleyOrderID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", BarleyOrderID)
		}

		orderToUpdate := BarleyOrder{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		orderToUpdate.Status = "Delivered"
		orderToUpdate.Accepted = Accepted

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) ReadBarleyOrder(ctx contractapi.TransactionContextInterface, BarleyOrderID string) (*BarleyOrder, error) {

	OrderID := ("BARLEY"+BarleyOrderID) 
	orderJSON, err := ctx.GetStub().GetState(OrderID)
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

func (s *SmartContract) ReadPrivateBarleyOrder(ctx contractapi.TransactionContextInterface, BarleyOrderID string) (*BarleyPrivateOrder, error) {

	OrderID := ("BARLEY"+BarleyOrderID) 

	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)

	if (msp == "malting-supply-com") {
		orderJSON, err := ctx.GetStub().GetPrivateData("collectionPrivateProducer1-Orders", OrderID)
		if orderJSON == nil {
			orderJSON, err = ctx.GetStub().GetPrivateData("collectionPrivateProducer2-Orders", OrderID)
		}
	}

	if (msp == "producer1-supply-com") {
		orderJSON, err := ctx.GetStub().GetPrivateData("collectionPrivateProducer1-Orders", OrderID)
	}
	if (msp == "producer2-supply-com") {
		orderJSON, err := ctx.GetStub().GetPrivateData("collectionPrivateProducer2-Orders", OrderID)
	}


	if err != nil {
		return nil, fmt.Errorf("failed to read from order %s", err.Error())
	}
	if orderJSON == nil {
		return nil, fmt.Errorf("%s does not exist", BarleyOrderID)
	}

	Privorder := new(BarleyPrivateOrder)
	_ = json.Unmarshal(orderJSON, Privorder)

	return Privorder, nil
}

//Distillery Starts New Malt Order
func (s *SmartContract) InitMaltOrder(ctx contractapi.TransactionContextInterface) error {
	//TODO: MSP MUST BE DISTILLERY
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}

	transientMaltOrderJSON, ok := transMap["InputJSON"]
	if !ok {
		return fmt.Errorf("Malt Order not found in the transient map")
	}

	var orderInput MaltOrder

	err = json.Unmarshal(transientMaltOrderJSON, &orderInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	if len(orderInput.MaltOrderID) == 0 {
		return fmt.Errorf("Order ID field must be a non-empty string")
	}
	if len(orderInput.Size) == 0 {
		return fmt.Errorf("Producer field must be a non-empty string")
	}

	// ==== Check if order already exists ====
	OrderID := ("MALT"+orderInput.MaltOrderID) 
	orderAsBytes, err := ctx.GetStub().GetState(OrderID)
	if err != nil {
		return fmt.Errorf("Failed to get order: " + err.Error())
	} else if orderAsBytes != nil {
		fmt.Println("This order already exists: " + orderInput.MaltOrderID)
		return fmt.Errorf("This order already exists: " + orderInput.MaltOrderID)
	}

	// ==== Create order object, marshal to JSON, and save to state ====
	order := &MaltOrder{
		ObjectType: "MaltOrder",
		MaltOrderID:       orderInput.MaltOrderID,
		Size:       orderInput.Size,
		Status:      "Ordered",
	}

	orderJSONasBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// === Save order to state ===
	err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put Order: %s", err.Error())
	}

	return nil
}

func (s *SmartContract) ConfirmMaltOrder(ctx contractapi.TransactionContextInterface) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}
	fmt.Println("Failed: ", msp)
	if msp == "malting-supply-com" {
		transMap, err := ctx.GetStub().GetTransient()
		if err != nil {
			return fmt.Errorf("Error getting transient: " + err.Error())
		}

		transientJSON, ok := transMap["InputJSON"]
		if !ok {
			return fmt.Errorf("Order not found in the transient map")
		}

		type OrderTransientInput struct {
			MaltOrderID  string `json:"MaltOrderID"`
			BarleyOrderID	string	`json:"BarleyOrderID"`
			Price int `json:"Price"`
			InvoiceID int `json:"InvoiceID"`
		}

		var OrderInput OrderTransientInput
		err = json.Unmarshal(transientJSON, &OrderInput)
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		OrderID := ("MALT"+ OrderInput.MaltOrderID) 

		if len(OrderInput.MaltOrderID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}
		if len(OrderInput.BarleyOrderID) == 0 {
			return fmt.Errorf("Barley ID field must be a non-empty string")
		}
		if OrderInput.Price == 0 {
			return fmt.Errorf("price field must not be nil")
		}
		if OrderInput.InvoiceID == 0 {
			return fmt.Errorf("InvoiceID field must not be nil")
		}

		//Test Barley Batch Exists
		BarleyAsBytes, err := ctx.GetStub().GetState(("BARLEY"+OrderInput.BarleyOrderID))
		if err != nil {
			return fmt.Errorf("Failed to Check Barley Order:" + err.Error())
		} else if BarleyAsBytes == nil {
			return fmt.Errorf("BarleyOrder does not exist: " + OrderInput.BarleyOrderID)
		}

		orderAsBytes, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("Failed to get order:" + err.Error())
		} else if orderAsBytes == nil {
			return fmt.Errorf("MaltOrder does not exist: " + OrderInput.BarleyOrderID)
		}

		orderToUpdate := MaltOrder{}
		err = json.Unmarshal(orderAsBytes, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		orderToUpdate.Status = "Confirmed"
		orderToUpdate.BarleyOrderID = OrderInput.BarleyOrderID

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

		privOrder := &MaltPrivateOrder{
			ObjectType: "MaltPrivateOrder",
			MaltOrderID:       OrderInput.MaltOrderID,
			Price:       OrderInput.Price,
			InvoiceID:	OrderInput.InvoiceID,
		}

		orderPrivJSONasBytes, err := json.Marshal(privOrder)
		if err != nil {
			return fmt.Errorf(err.Error())
		}

	
		err = ctx.GetStub().PutPrivateData("collectionPrivateMalt-Orders", OrderID, orderPrivJSONasBytes)

		if err != nil {
			return fmt.Errorf("failed to put Order: %s", err.Error())
		}

		//ELSE ERROR

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) ShipMaltOrder(ctx contractapi.TransactionContextInterface, MaltOrderID string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)
	if msp == "malting-supply-com" {
		
		if len(MaltOrderID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}

		OrderID := ("MALT"+MaltOrderID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", MaltOrderID)
		}

		orderToUpdate := MaltOrder{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		orderToUpdate.Status = "Shipped"

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) AcceptMaltOrder(ctx contractapi.TransactionContextInterface, MaltOrderID string, Accepted string, QC string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)
	if msp == "distillery-supply-com" {
		
		if len(MaltOrderID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}
		
		OrderID := ("MALT"+MaltOrderID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", MaltOrderID)
		}

		orderToUpdate := MaltOrder{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		orderToUpdate.Status = "Delivered"
		orderToUpdate.Accepted = Accepted
		orderToUpdate.QCPass = QC

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) ReadMaltOrder(ctx contractapi.TransactionContextInterface, MaltOrderID string) (*MaltOrder, error) {

	OrderID := ("MALT"+MaltOrderID) 
	orderJSON, err := ctx.GetStub().GetState(OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from order %s", err.Error())
	}
	if orderJSON == nil {
		return nil, fmt.Errorf("%s does not exist", MaltOrderID)
	}

	order := new(MaltOrder)
	_ = json.Unmarshal(orderJSON, order)

	return order, nil
}

func (s *SmartContract) ReadPrivateMaltOrder(ctx contractapi.TransactionContextInterface, MaltOrderID string) (*MaltPrivateOrder, error) {

	OrderID := ("MALT"+MaltOrderID) 
	orderJSON, err := ctx.GetStub().GetPrivateData("collectionPrivateMalt-Orders", OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from order %s", err.Error())
	}
	if orderJSON == nil {
		return nil, fmt.Errorf("%s does not exist", MaltOrderID)
	}

	Privorder := new(MaltPrivateOrder)
	_ = json.Unmarshal(orderJSON, Privorder)

	return Privorder, nil
}

//Distillery Start New Batch
func (s *SmartContract) InitBatch(ctx contractapi.TransactionContextInterface) error {
	//TODO: MSP MUST BE DISTILLERY
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}

	transientMaltOrderJSON, ok := transMap["InputJSON"]
	if !ok {
		return fmt.Errorf("Malt Order not found in the transient map")
	}

	var orderInput Distillation

	err = json.Unmarshal(transientMaltOrderJSON, &orderInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	if len(orderInput.BatchID) == 0 {
		return fmt.Errorf("Batch ID field must be a non-empty string")
	}
	if len(orderInput.MaltOrderID) == 0 {
		return fmt.Errorf("Batch ID field must be a non-empty string")
	}
	if len(orderInput.Size) == 0 {
		return fmt.Errorf("Producer field must be a non-empty string")
	}

	//Test Malt Batch Exists
	MaltAsBytes, err := ctx.GetStub().GetState(("MALT"+orderInput.MaltOrderID))
	if err != nil {
		return fmt.Errorf("Failed to Check Malt Order:" + err.Error())
	} else if MaltAsBytes == nil {
		return fmt.Errorf("MaltOrder does not exist: " + orderInput.MaltOrderID)
	}

	// ==== Check if order already exists ====
	OrderID := ("BATCH"+orderInput.BatchID) 
	orderAsBytes, err := ctx.GetStub().GetState(OrderID)
	if err != nil {
		return fmt.Errorf("Failed to get order: " + err.Error())
	} else if orderAsBytes != nil {
		fmt.Println("This order already exists: " + orderInput.BatchID)
		return fmt.Errorf("This order already exists: " + orderInput.BatchID)
	}

	// ==== Create order object, marshal to JSON, and save to state ====
	order := &Distillation{
		ObjectType: "Distillation",
		BatchID:	orderInput.BatchID,
		MaltOrderID:       orderInput.MaltOrderID,
		Size:       orderInput.Size,
		Status:      "Started",
	}

	orderJSONasBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// === Save order to state ===
	err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put Order: %s", err.Error())
	}

	return nil
}

func (s *SmartContract) UpdateBatchStatus(ctx contractapi.TransactionContextInterface, BatchID string, Status string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)
	//TODO: SETUP TO LIMIT ONLY Distillery

	if msp == "distillery-supply-com" {
		
		if len(BatchID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}
		if len(Status) == 0 {
			return fmt.Errorf("Status field must be a non-empty string")
		}

		OrderID := ("BATCH"+BatchID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", BatchID)
		}

		orderToUpdate := Distillation{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		orderToUpdate.Status = Status

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}


func (s *SmartContract) SetInitialProof(ctx contractapi.TransactionContextInterface, BatchID string, Proof string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)

	if msp == "distillery-supply-com" {
		
		if len(BatchID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}
		if len(Proof) == 0 {
			return fmt.Errorf("Status field must be a non-empty string")
		}

		OrderID := ("BATCH"+BatchID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", BatchID)
		}

		orderToUpdate := Distillation{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		if len(orderToUpdate.InitialProof) != 0 {
			return fmt.Errorf("Initial Proof Already Set")
		}
		orderToUpdate.InitialProof = Proof
		orderToUpdate.Status = "Batch Proofed"

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) SendToWarehouse(ctx contractapi.TransactionContextInterface, BatchID string, QC string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)

	if msp == "distillery-supply-com" {
		
		if len(BatchID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}

		if len(QC) == 0 {
			return fmt.Errorf("QC field must be a non-empty string")
		}

		OrderID := ("BATCH"+BatchID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", BatchID)
		}

		orderToUpdate := Distillation{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		if len(orderToUpdate.InitialProof) == 0 {
			return fmt.Errorf("Not Proofed Yet - Disallowed")
		}
		orderToUpdate.Status = "Shipped to Warehouse"
		orderToUpdate.QCPass = QC

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) AcceptAtWarehouse(ctx contractapi.TransactionContextInterface, BatchID string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)


	if msp == "maturation-supply-com" {
		
		if len(BatchID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}

		OrderID := ("BATCH"+BatchID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", BatchID)
		}

		orderToUpdate := Distillation{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		if (orderToUpdate.Status) != "Shipped to Warehouse" {
			return fmt.Errorf("Batch Not Shipped to Warehouse!")
		}
		orderToUpdate.Status = "Accepted at Warehouse"

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) ReadBatch(ctx contractapi.TransactionContextInterface, BatchID string) (*Distillation, error) {

	OrderID := ("BATCH"+BatchID) 
	orderJSON, err := ctx.GetStub().GetState(OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from order %s", err.Error())
	}
	if orderJSON == nil {
		return nil, fmt.Errorf("%s does not exist", BatchID)
	}

	order := new(Distillation)
	_ = json.Unmarshal(orderJSON, order)

	return order, nil
}

func (s *SmartContract) InitMaturation(ctx contractapi.TransactionContextInterface) error {
	//TODO: MSP MUST BE Maturation
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}

	transientMaturationJSON, ok := transMap["InputJSON"]
	if !ok {
		return fmt.Errorf("Malt Order not found in the transient map")
	}

	var orderInput Maturation

	err = json.Unmarshal(transientMaturationJSON, &orderInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	if len(orderInput.CaskID) == 0 {
		return fmt.Errorf("Batch ID field must be a non-empty string")
	}
	if len(orderInput.BatchID) == 0 {
		return fmt.Errorf("Batch ID field must be a non-empty string")
	}
	if len(orderInput.Size) == 0 {
		return fmt.Errorf("Producer field must be a non-empty string")
	}

	//Test Batch Exists
	MaltAsBytes, err := ctx.GetStub().GetState(("BATCH"+orderInput.BatchID))
	if err != nil {
		return fmt.Errorf("Failed to Check Malt Order:" + err.Error())
	} else if MaltAsBytes == nil {
		return fmt.Errorf("Batch does not exist: " + orderInput.BatchID)
	}

	// ==== Check if order already exists ====
	OrderID := ("CASK"+orderInput.CaskID) 
	orderAsBytes, err := ctx.GetStub().GetState(OrderID)
	if err != nil {
		return fmt.Errorf("Failed to get order: " + err.Error())
	} else if orderAsBytes != nil {
		fmt.Println("This order already exists: " + orderInput.CaskID)
		return fmt.Errorf("This order already exists: " + orderInput.CaskID)
	}

	// ==== Create order object, marshal to JSON, and save to state ====
	order := &Maturation{
		ObjectType: "Maturation",
		CaskID:		orderInput.CaskID,
		BatchID:	orderInput.BatchID,
		Size:       orderInput.Size,
		Status:      "Casked",
		StartDate:		"2019-04-03", //time.Now().Format("2006-01-02"), 
	}

	//TODO SETUP TEST DATA THEN REMOVE THE FIXED DATE VALUE

	orderJSONasBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// === Save order to state ===
	err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put Order: %s", err.Error())
	}

	return nil
}

func (s *SmartContract) SetFinalProof(ctx contractapi.TransactionContextInterface, CaskID string, Proof string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)

	if msp == "maturation-supply-com" {
		
		if len(CaskID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}
		if len(Proof) == 0 {
			return fmt.Errorf("Status field must be a non-empty string")
		}

		OrderID := ("CASK"+CaskID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", CaskID)
		}

		orderToUpdate := Maturation{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		if len(orderToUpdate.FinalProof) != 0 {
			return fmt.Errorf("Final Proof Already Set")
		}
		orderToUpdate.FinalProof = Proof
		orderToUpdate.Status = "Batch Proofed"

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) QualityControl(ctx contractapi.TransactionContextInterface, CaskID string, QC string, Notes string, Taste string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)

	if msp == "maturation-supply-com" {
		
		if len(CaskID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}
		if len(QC) == 0 {
			return fmt.Errorf("QC field must be a non-empty string")
		}
		if len(Notes) == 0 {
			return fmt.Errorf("Notes field must be a non-empty string")
		}
		if len(Taste) == 0 {
			return fmt.Errorf("Taste field must be a non-empty string")
		}

		OrderID := ("CASK"+CaskID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", CaskID)
		}

		orderToUpdate := Maturation{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		if len(orderToUpdate.QCPass) != 0 {
			return fmt.Errorf("Cannot Overwrite QC Record")
		}
		orderToUpdate.Status = "Quality Controlled"
		orderToUpdate.QCPass = QC
		orderToUpdate.Notes = Notes
		orderToUpdate.Taste = Taste

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) SendToBottling(ctx contractapi.TransactionContextInterface, CaskID string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)

	if msp == "maturation-supply-com" {
		
		if len(CaskID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}

		OrderID := ("CASK"+CaskID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", CaskID)
		}

		orderToUpdate := Maturation{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		if len(orderToUpdate.QCPass) == 0 {
			return fmt.Errorf("Quality Control Not Completed")
		}
		orderToUpdate.Status = "Ready For Bottling"
		orderToUpdate.EndDate = time.Now().Format("2006-01-02")

		end, err := time.Parse("2006-01-02", orderToUpdate.EndDate)
		if err != nil {
			return err
		}
		start, err := time.Parse("2006-01-02", orderToUpdate.StartDate)
		if err != nil {
			return err
		}
		hours := end.Sub(start)
		orderToUpdate.Age = fmt.Sprint("%.2f", (hours.Hours()/24/365))+" Years"

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) AcceptAtBottling(ctx contractapi.TransactionContextInterface, CaskID string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)

	if msp == "bottling-supply-com" {
		
		if len(CaskID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}

		OrderID := ("CASK"+CaskID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", CaskID)
		}

		orderToUpdate := Maturation{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		if orderToUpdate.Status != "Ready For Bottling" {
			return fmt.Errorf("Not sent for bottling")
		}
		orderToUpdate.Status = "Accepted at Bottling"

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) ReadCask(ctx contractapi.TransactionContextInterface, CaskID string) (*Maturation, error) {

	OrderID := ("CASK"+CaskID) 
	orderJSON, err := ctx.GetStub().GetState(OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from order %s", err.Error())
	}
	if orderJSON == nil {
		return nil, fmt.Errorf("%s does not exist", CaskID)
	}

	order := new(Maturation)
	_ = json.Unmarshal(orderJSON, order)

	return order, nil
}

//BOTTLING

func (s *SmartContract) InitBottling(ctx contractapi.TransactionContextInterface) error {
	//TODO: MSP MUST BE Bottling
	//TODO: Private Cask Data
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}

	transientBottlingJSON, ok := transMap["InputJSON"]
	if !ok {
		return fmt.Errorf("Malt Order not found in the transient map")
	}

	var orderInput Bottling

	err = json.Unmarshal(transientBottlingJSON, &orderInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	if len(orderInput.BottleID) == 0 {
		return fmt.Errorf("Bottle ID field must be a non-empty string")
	}
	if len(orderInput.CaskID) == 0 {
		return fmt.Errorf("Batch ID field must be a non-empty string")
	}
	if len(orderInput.Size) == 0 {
		return fmt.Errorf("Producer field must be a non-empty string")
	}

	//Test Batch Exists
	for i, s := range orderInput.CaskID{
		CaskAsBytes, err := ctx.GetStub().GetState(("CASK"+s))
		if err != nil {
			return fmt.Errorf("Failed to Check Cask ID:" + err.Error())
		} else if CaskAsBytes == nil {
			return fmt.Errorf("Cask does not exist: " + s, i)
		}
	}

	// ==== Check if order already exists ====
	OrderID := ("BOTTLE"+orderInput.BottleID) 
	orderAsBytes, err := ctx.GetStub().GetState(OrderID)
	if err != nil {
		return fmt.Errorf("Failed to get order: " + err.Error())
	} else if orderAsBytes != nil {
		fmt.Println("This order already exists: " + orderInput.BottleID)
		return fmt.Errorf("This order already exists: " + orderInput.BottleID)
	}

	// ==== Create order object, marshal to JSON, and save to state ====
	order := &Bottling{
		ObjectType: "Bottling",
		BottleID:		orderInput.BottleID,
		CaskID:	orderInput.CaskID,
		Size:       orderInput.Size,
		Status:      "Bottled",
	}

	//TODO SETUP TEST DATA THEN REMOVE THE FIXED DATE VALUE

	orderJSONasBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// === Save order to state ===
	err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put Order: %s", err.Error())
	}

	return nil
}

func (s *SmartContract) SetPallet(ctx contractapi.TransactionContextInterface, BottleID string, PalletID string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)
	//TODO: SETUP TO LIMIT ONLY Botteling

	if msp == "bottling-supply-com" {
		
		if len(BottleID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}
		if len(PalletID) == 0 {
			return fmt.Errorf("PalletID field must be a non-empty string")
		}

		OrderID := ("BOTTLE"+BottleID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", BottleID)
		}

		orderToUpdate := Bottling{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		

		orderToUpdate.PalletID = PalletID
		orderToUpdate.Status = "Palleted"

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) ReadBottle(ctx contractapi.TransactionContextInterface, BottleID string) (*Bottling, error) {

	OrderID := ("BOTTLE"+BottleID) 
	orderJSON, err := ctx.GetStub().GetState(OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from order %s", err.Error())
	}
	if orderJSON == nil {
		return nil, fmt.Errorf("%s does not exist", BottleID)
	}

	order := new(Bottling)
	_ = json.Unmarshal(orderJSON, order)

	return order, nil
}

func (s *SmartContract) BottleLife(ctx contractapi.TransactionContextInterface, BottleID string) (*BottleLifeModel, error) {

	lifecycle := new(BottleLifeModel)

	//Get Bottle Info 
	BottleIDN := ("BOTTLE"+BottleID) 
	BottleJSON, err := ctx.GetStub().GetState(BottleIDN)
	if err != nil {
		return nil, fmt.Errorf("failed to read from bottle %s", err.Error())
	}
	if BottleJSON == nil {
		return nil, fmt.Errorf("%s does not exist", BottleID)
	}

	bottle := new(Bottling)
	_ = json.Unmarshal(BottleJSON, bottle)

	lifecycle.BottleID = bottle.BottleID

	lifecycle.Casks = make([]CaskLifeModel, 0)

	//Get Cask Info 
	for i, iCask := range bottle.CaskID{
		CASKIDN := ("CASK"+iCask) 
		CaskJSON, err := ctx.GetStub().GetState(CASKIDN)
		if err != nil {
			return nil, fmt.Errorf("failed to read from cask %s", err.Error())
		}
		if CaskJSON == nil {
			return nil, fmt.Errorf("%s does not exist", CASKIDN, i)
		}

		cask := new(Maturation)
		_ = json.Unmarshal(CaskJSON, cask)




		//Get Batch Info 
		BatchIDN := ("BATCH"+cask.BatchID) 
		BatchJSON, err := ctx.GetStub().GetState(BatchIDN)
		if err != nil {
			return nil, fmt.Errorf("failed to read from bottle %s", err.Error())
		}
		if BatchJSON == nil {
			return nil, fmt.Errorf("%s does not exist", BatchIDN)
		}

		batch := new(Distillation)
		_ = json.Unmarshal(BatchJSON, batch)

		//Get Malt Info 
		MaltIDN := ("MALT"+batch.MaltOrderID) 
		MaltJSON, err := ctx.GetStub().GetState(MaltIDN)
		if err != nil {
			return nil, fmt.Errorf("failed to read from bottle %s", err.Error())
		}
		if MaltJSON == nil {
			return nil, fmt.Errorf("%s does not exist", MaltIDN)
		}

		malt := new(MaltOrder)
		_ = json.Unmarshal(MaltJSON, malt)

		//Get Barley Info 
		BarleyIDN := ("BARLEY"+malt.BarleyOrderID) 
		BarleyJSON, err := ctx.GetStub().GetState(BarleyIDN)
		if err != nil {
			return nil, fmt.Errorf("failed to read from bottle %s", err.Error())
		}
		if BarleyJSON == nil {
			return nil, fmt.Errorf("%s does not exist", BarleyIDN)
		}

		barley := new(BarleyOrder)
		_ = json.Unmarshal(BarleyJSON, barley)

		thisCask := new(CaskLifeModel)
		thisCask.CaskID = cask.CaskID
		thisCask.QCPass = cask.QCPass
		thisCask.BatchID = cask.BatchID
		thisCask.BatchQCPass = batch.QCPass
		thisCask.MaltID = batch.MaltOrderID
		thisCask.MaltQCPass = malt.QCPass
		thisCask.BarleyOrderID = malt.BarleyOrderID
		thisCask.Producer = barley.Producer
		thisCask.SoilPH = barley.SoilPH
		thisCask.GeoLocation = barley.GeoLocation
		
		lifecycle.Casks = append(lifecycle.Casks, *thisCask)
	}

	return lifecycle, nil
}

func (s *SmartContract) InitPalletOrder(ctx contractapi.TransactionContextInterface) error {
	//TODO: MSP MUST BE Retailer
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}

	transientPalletOrderJSON, ok := transMap["InputJSON"]
	if !ok {
		return fmt.Errorf("Malt Order not found in the transient map")
	}

	var orderInput RetailerOrder

	err = json.Unmarshal(transientPalletOrderJSON, &orderInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	if len(orderInput.RetailerOrderID) == 0 {
		return fmt.Errorf("Order ID field must be a non-empty string")
	}
	if len(orderInput.Size) == 0 {
		return fmt.Errorf("Producer field must be a non-empty string")
	}
	if len(orderInput.Shop) == 0 {
		return fmt.Errorf("Shop field must be a non-empty string")
	}

	// ==== Check if order already exists ====
	OrderID := ("PALLET"+orderInput.RetailerOrderID) 
	orderAsBytes, err := ctx.GetStub().GetState(OrderID)
	if err != nil {
		return fmt.Errorf("Failed to get order: " + err.Error())
	} else if orderAsBytes != nil {
		fmt.Println("This order already exists: " + orderInput.RetailerOrderID)
		return fmt.Errorf("This order already exists: " + orderInput.RetailerOrderID)
	}

	// ==== Create order object, marshal to JSON, and save to state ====
	order := &RetailerOrder{
		ObjectType: "RetailerOrder",
		RetailerOrderID:       orderInput.RetailerOrderID,
		Size:       orderInput.Size,
		Status:      "Ordered",
		Shop:	orderInput.Shop,
	}

	orderJSONasBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// === Save order to state ===
	err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put Order: %s", err.Error())
	}

	return nil
}

func (s *SmartContract) ConfirmRetailerOrder(ctx contractapi.TransactionContextInterface) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}
	fmt.Println("Failed: ", msp)
	if msp == "distillery-supply-com" {
		transMap, err := ctx.GetStub().GetTransient()
		if err != nil {
			return fmt.Errorf("Error getting transient: " + err.Error())
		}

		transientJSON, ok := transMap["InputJSON"]
		if !ok {
			return fmt.Errorf("Order not found in the transient map")
		}

		type OrderTransientInput struct {
			RetailerOrderID  string `json:"RetailerOrderID"`
			PalletID	string	`json:"PalletID"`
			Price int `json:"Price"`
			InvoiceID int `json:"InvoiceID"`
		}

		var OrderInput OrderTransientInput
		err = json.Unmarshal(transientJSON, &OrderInput)
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		OrderID := ("PALLET"+ OrderInput.RetailerOrderID) 

		if len(OrderInput.RetailerOrderID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}
		if len(OrderInput.PalletID) == 0 {
			return fmt.Errorf("Pallet ID field must be a non-empty string")
		}
		if OrderInput.Price == 0 {
			return fmt.Errorf("price field must not be nil")
		}
		if OrderInput.InvoiceID == 0 {
			return fmt.Errorf("InvoiceID field must not be nil")
		}

		orderAsBytes, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("Failed to get order:" + err.Error())
		} else if orderAsBytes == nil {
			return fmt.Errorf("Order does not exist: " + OrderInput.RetailerOrderID)
		}

		orderToUpdate := RetailerOrder{}
		err = json.Unmarshal(orderAsBytes, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		orderToUpdate.Status = "Confirmed"
		orderToUpdate.PalletID = OrderInput.PalletID

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

		privOrder := &RetailerPrivateOrder{
			ObjectType: "RetailerPrivateOrder",
			RetailerOrderID:       OrderInput.RetailerOrderID,
			Price:       OrderInput.Price,
			InvoiceID:	OrderInput.InvoiceID,
		}

		orderPrivJSONasBytes, err := json.Marshal(privOrder)
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		//TODO: - SPLIT FOR THE TWO DIFFERENT RETAILERS WITH AN IF BASED ON OWNER
		// OF ORDER MAYBE ADD MSP TO ORIGIONAL ORDER JSON

		err = ctx.GetStub().PutPrivateData("collectionPrivateRetailer1-Orders", OrderID, orderPrivJSONasBytes)

		if err != nil {
			return fmt.Errorf("failed to put Order: %s", err.Error())
		}

		//ELSE ERROR

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) ShipRetailerOrder(ctx contractapi.TransactionContextInterface, RetailerOrderID string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)
	if msp == "bottling-supply-com" {
		
		if len(RetailerOrderID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}

		OrderID := ("PALLET"+RetailerOrderID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", RetailerOrderID)
		}

		orderToUpdate := RetailerOrder{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		orderToUpdate.Status = "Shipped"

		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) DeliveredRetailerOrder(ctx contractapi.TransactionContextInterface, RetailerOrderID string) error {
	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)
	if (msp == "retailer1-supply-com" || msp == "retailer2-supply-com") {
		
		if len(RetailerOrderID) == 0 {
			return fmt.Errorf("ID field must be a non-empty string")
		}
		
		OrderID := ("PALLET"+RetailerOrderID) 
		orderJSON, err := ctx.GetStub().GetState(OrderID)
		if err != nil {
			return fmt.Errorf("failed to read from order %s", err.Error())
		}
		if orderJSON == nil {
			return fmt.Errorf("%s does not exist", RetailerOrderID)
		}

		orderToUpdate := RetailerOrder{}
		err = json.Unmarshal(orderJSON, &orderToUpdate) 
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
		}
		
		orderToUpdate.Status = "Delivered"


		orderJSONasBytes, _ := json.Marshal(orderToUpdate)
		err = ctx.GetStub().PutState(OrderID, orderJSONasBytes)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Wrong MSP - Access Deinied")
	}
	return nil
}

func (s *SmartContract) ReadRetailerOrder(ctx contractapi.TransactionContextInterface, RetailerOrderID string) (*RetailerOrder, error) {

	OrderID := ("PALLET"+RetailerOrderID) 
	orderJSON, err := ctx.GetStub().GetState(OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from order %s", err.Error())
	}
	if orderJSON == nil {
		return nil, fmt.Errorf("%s does not exist", RetailerOrderID)
	}

	order := new(RetailerOrder)
	_ = json.Unmarshal(orderJSON, order)

	return order, nil
}

func (s *SmartContract) ReadPrivateRetailerOrder(ctx contractapi.TransactionContextInterface, RetailerOrderID string) (*RetailerPrivateOrder, error) {
	OrderID := ("PALLET"+RetailerOrderID) 


	msp, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("Error getting MSPID: " + err.Error())
	}
	fmt.Println("Failed: ", msp)

	if (msp == "distillery-supply-com") {
		orderJSON, err := ctx.GetStub().GetPrivateData("collectionPrivateRetailer1-Orders", OrderID)
		if orderJSON == nil {
			orderJSON, err = ctx.GetStub().GetPrivateData("collectionPrivateRetailer2-Orders", OrderID)
		}
	}

	if (msp == "retailer1-supply-com") {
		orderJSON, err := ctx.GetStub().GetPrivateData("collectionPrivateRetailer1-Orders", OrderID)
	}
	if (msp == "retailer2-supply-com") {
		orderJSON, err := ctx.GetStub().GetPrivateData("collectionPrivateRetailer2-Orders", OrderID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read from order %s", err.Error())
	}
	if orderJSON == nil {
		return nil, fmt.Errorf("%s does not exist", RetailerOrderID)
	}

	Privorder := new(RetailerPrivateOrder)
	_ = json.Unmarshal(orderJSON, Privorder)

	return Privorder, nil
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
