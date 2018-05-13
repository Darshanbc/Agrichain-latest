package main

import (

	"encoding/json"
	//"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"log"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
)

var logger = shim.NewLogger("farmer_agricert")

type Fertilizer struct {
	FertlizerName string
	FertlizerID   string
	Quantity      float32
}

type SmartContract struct {
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {

		logger.Info("-------------------Invoke----------")
		// Retrieve the requested Smart Contract function and arguments
		function, args := stub.GetFunctionAndParameters()
		// Route to the appropriate handler function to interact with the ledger appropriately
		logger.Info("fucntion")
		logger.Info(function)
		logger.Info("args")
		logger.Info(args)
		if function =="addStock"{
			return s.addStock(stub)
		}
		return shim.Error("Invalid Smart contract name")
}

func (s *SmartContract) addStock(stub shim.ChaincodeStubInterface)sc.Response{
	//org2:farmer org3:fert
	org,_:=cid.GetMSPID(stub)
	enrollID,_:=cid.GetID(stub)
	fert:=mapstruct(stub)
	userAsBytes,_:=stub.GetState(enrollID)
	if userAsBytes== nil && org=="Org3MSP"{
		newfert:=[]Fertilizer{}
		newfert=append(newfert,fert)
		stockasbytes,_:=json.Marshal(newfert)
		stub.PutState(enrollID,stockasbytes)
		return shim.Success(nil)
	}
	if userAsBytes!=nil && org=="Org3MSP" {
		fertlistAsbytes,_:=stub.GetState(enrollID)
		newfertList:=[]Fertilizer{}
		json.Unmarshal(fertlistAsbytes,&newfertList)
		flag:=0
		var index int
		var fertl Fertilizer

		for index, fertl.FertlizerID = range newfertList{
			if fertl.FertlizerID== fert.FertlizerID{
				flag=1
				break
			}
		}

		if flag==0{
			newfertList=append(newfertList,fert)
			fertlistAsbytes,_=json.Marshal(newfertList)
			stub.PutState(enrollID,fertlistAsbytes)
			return shim.Success(nil)
		}else{
			newfertList[index].Quantity=newfertList[index].Quantity+fert.Quantity
			fertlistAsbytes,_=json.Marshal(newfertList)
			stub.PutState(enrollID,fertlistAsbytes)
			return shim.Success(nil)
		}

	}
	return shim.Error("Wrong smartcontract invoked")
}

func mapstruct(stub shim.ChaincodeStubInterface) Fertilizer{
	args:=stub.GetArgs()
	fert:=Fertilizer{}
	if err := json.Unmarshal([]byte(args[1]), &fert); err != nil {
		log.Fatal(err)
	}
	return  fert
}