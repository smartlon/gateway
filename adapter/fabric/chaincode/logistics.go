package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type mamchannel struct {
	Root    string     `json:",Root"`
	SideKey string     `json:"SideKey"`
}

// logisticstrans type
type logisticstrans struct {
	//product might be food,fish,phone,other itmes
	//Product id should be unique such as FISH123,Prawns456,ICECREAM789
	ProductID         string       `json:"ProductID"`
	ProductType       string       `json:"ProductType"`
	SellerID          string       `json:"SellerID"`
	SellerLocation    string       `json:"SellerLocation"`
	BuyerID           string       `json:"BuyerID"`
	BuyerLocation     string       `json:"BuyerLocation"`
	LogisticsID       string       `json:"LogisticsID"`
	LogisticsLocation string       `json:"LogisticsLocation"`
	JourneyStartTime  string       `json:",JourneyStartTime"`
	JourneyEndTime    string       `json:",JourneyEndTime"`
	Status            string       `json:"Status"`
	MAMChannel    mamchannel `json:"MAMChannel"`
}

func main() {

	err := shim.Start(new(logisticstrans))
	if err != nil {
		fmt.Println("Error with chaincode")
	} else {
		fmt.Println("Chaincode installed successfully")
	}
}

//Init logisticstrans
func (t *logisticstrans) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Initiate the chaincode")
	return shim.Success(nil)
}

//Invoke logisticstrans
func (t *logisticstrans) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fun, args := stub.GetFunctionAndParameters()
	fmt.Println("Arguements for function  ", fun)
	switch fun {
	case "RequestLogistic":
		return t.RequestLogistic(stub, args)
	case "TransitLogistics":
		return t.TransitLogistics(stub, args)
	case "InTransitLogistics":
		return t.InTransitLogistics(stub, args)
	case "DeliveryLogistics":
		return t.DeliveryLogistics(stub, args)
	case "SignLogistics":
		return t.SignLogistics(stub, args)
	case "QueryLogistics":
		return t.QueryLogistics(stub, args)
	}
	fmt.Println("Function not found!")
	return shim.Error("Recieved unknown function invocation!")
}

//Genlogistics for

func (t *logisticstrans) RequestLogistic(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var ProductID string

	if len(args) < 1 {
		fmt.Println("Invalid number od arguements")
		return shim.Error(err.Error())
	}
	if err != nil {
		return shim.Error("Invalid Request Number")
	}
	ProductID = args[0]
	var logobj = logisticstrans{ProductID: ProductID, ProductType: args[1], BuyerID: args[2], BuyerLocation: args[3], SellerID: args[4], SellerLocation: args[5]}
	logobj.Status = "Requested"

	logobjasBytes, _ := json.Marshal(logobj)
	stub.PutState(args[0], logobjasBytes)

	return shim.Success(nil)
}

//TransitLogistics at the same time measuring the temp details from logistics
func (t *logisticstrans) TransitLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting Minimum 7")
	}

	if err != nil {
		return shim.Error("Invalid ")
	}
	logisticsAsBytes, _ := stub.GetState(args[0])

	var logisticobj logisticstrans
	json.Unmarshal(logisticsAsBytes, &logisticobj)
	logisticobj.ProductID = args[0]
	logisticobj.LogisticsID = args[1]
	logisticobj.LogisticsLocation = args[2]
	logisticobj.JourneyStartTime = args[3]
	logisticobj.MAMChannel.SideKey = args[4]
	if logisticobj.Status != "Requested" {
		fmt.Println("we cannnot transit  the product which was not requested")
		return shim.Error("we cannnot transit  the product which was not requested")
	}

	logisticobj.Status = "Ready-Transit"
	logisticsAsBytes, _ = json.Marshal(logisticobj)
	stub.PutState(args[0], logisticsAsBytes)
	err = stub.SetEvent(`{"From":"Fabric","To":"Iota","Func":"CreateChannel"}`, logisticsAsBytes)
	if err != nil {
		fmt.Println("Could not set event for loan application creation", err)
	}
	return shim.Success(nil)
}

func (t *logisticstrans) InTransitLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting Minimum 7")
	}

	if err != nil {
		return shim.Error("Invalid ")
	}
	logisticsAsBytes, _ := stub.GetState(args[0])

	var logisticobj logisticstrans
	json.Unmarshal(logisticsAsBytes, &logisticobj)
	logisticobj.ProductID = args[0]
	logisticobj.MAMChannel.Root = args[1]

	if logisticobj.Status != "Ready-Transit" {
		fmt.Println("we cannnot transit  the product which was not Ready_Transit")
		return shim.Error("we cannnot transit  the product which was not Ready_Transit")
	}

	logisticobj.Status = "In-Transit"
	logisticsAsBytes, _ = json.Marshal(logisticobj)
	stub.PutState(args[0], logisticsAsBytes)

	return shim.Success(nil)
}

func (t *logisticstrans) DeliveryLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Invalid   no of arg for delivery function ")

	}
	logisticsasbytes1, _ := stub.GetState(args[0])
	var logisticobj1 logisticstrans

	json.Unmarshal(logisticsasbytes1, &logisticobj1)

	if logisticobj1.Status != "In-Transit" {
		fmt.Println("we cannnot delivery the product which is not in In_Transit")
		return shim.Error("we cannnot delivery the product which is not in In_Transit")
	}
	logisticobj1.JourneyEndTime = args[1]
	logisticobj1.Status = "Wait-Sign"
	logisticsasbytes1, _ = json.Marshal(logisticobj1)
	stub.PutState(args[0], logisticsasbytes1)
	err := stub.SetEvent(`{"From":"Fabric","To":"Iota","Func":"DeliveryLogistics"}`, logisticsasbytes1)
	if err != nil {
		fmt.Println("Could not set event for loan application creation", err)
	}
	return shim.Success(nil)
}

func (t *logisticstrans) SignLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Invalid   no of arg for Sign function ")

	}
	logisticsasbytes1, _ := stub.GetState(args[0])
	var logisticobj1 logisticstrans

	json.Unmarshal(logisticsasbytes1, &logisticobj1)

	if logisticobj1.Status != "Wait-Sign" {
		fmt.Println("we cannnot delivery the product which is not in Wait_Sign")
		return shim.Error("we cannnot delivery the product which is not in Wait_Sign")
	}
	fmt.Println("length of the logibj journry in  device", logisticobj1.JourneyEndTime)
	fmt.Println("length of the logibj  journey out timefrrom device", logisticobj1.JourneyStartTime)

	count := 0
	tempStr := strings.Split(args[1],",")
	for i := 0; i < len(tempStr); i++ {
		temp, _ := strconv.Atoi(tempStr[i])
		if temp >= 20 {
			//fmt.Println("Temperature from array is :", logisticobj1.Timefromdevice[i].Temperature)
			//fmt.Println("status of temp from array is :", logisticobj1.Status)
			count++
		} else {
			count = 0
		}
		//fmt.Println("Count is  from for loop:", count)
		//fmt.Println("status of temp  is :", logisticobj1.Status)

		if count >= 3 {
			logisticobj1.Status = "Rejected from Buyer"
			break

		} else {
			logisticobj1.Status = "Accepted  from Buyer"
		}

	}

	logisticsasbytes1, _ = json.Marshal(logisticobj1)
	stub.PutState(args[0], logisticsasbytes1)

	return shim.Success(nil)
}


func (t *logisticstrans) QueryLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Invalid   no of arg for Query function ")
	}
	logisticsasbytes1, _ := stub.GetState(args[0])
	return shim.Success(logisticsasbytes1)
}

