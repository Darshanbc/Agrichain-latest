package main

import (
	// "bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	sc "github.com/hyperledger/fabric/protos/peer"
	"log"
	strconv2 "strconv"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
)

var logger = shim.NewLogger("farmer_agricert")

type SmartContract struct {
}

type Farmer struct {
	KYC_ID string
	Plots  []Plot       `json:"plot"`
	Store  []Fertilizer `json:"store"`
}
type FertAgent struct {
	KYC_ID string
	Store  []Fertilizer `json:"store"`
}
type CertAgent struct {
	KYC_ID              string
	CertificateRequests [] string
	FertilizerRequests  []string
}

type KYC struct {
	First_Name string `json:"first name"`
	Last_Name  string `json:"last name"`
	Age        string `json:"age"`
	Address    string `json:"address"`
	Email      string `json:"email"`
	MobileNo   string `json:"mobile_no"`
}

type Plot struct {
	PlotId       string
	Co_ordinates []Latlong `json:"co_ordinates"`
	Survey_no    string    `json:"survey_no"`
	Crop_history []string  `json:"crop_history"`
	SoilType     string    `json:"soil_type"`
	Current_crop string    `json:"current_crop"`
}

type Latlong struct {
	Lattitude float64 `json:"lattitude"`
	Longitude float64 `json:"longitude"`
}

type Fertilizer struct {
	FertlizerName string  `json:"fertlizer_name"`
	FertlizerID   string  `json:"fertlizer_id"`
	Quantity      float32 `json:"quantity"`
}

type Crop struct {
	CropName        string       `json:"crop_name"`
	Type            string       `json:"type"`
	CropCycle       Cycle        `json:"crop_cycle"`
	FertilzerUsed   []Fertilizer `json:"fertilzer_used"`
	FertilzerReq    []Fertilizer `json:"fertilzer_req"`
	CertRequestTxID string       `json:"cert_request_tx_id"`
	Cert            Certificate  `json:"cert"`
}

type Cycle struct {
	FromMonth string `json:"from_month"`
	ToMonth   string `json:"to_month"`
}

//0-----------------------------------Cert Agency-----
//

type Requests struct {
	Asset  interface{} `json:"asset"` //cert or fert
	Status string      `json:"status"`
}

type Certificate struct {
	CertficateID     string `json:"certficate_id"`
	DigitalSignature string `json:"digital_signature"`
	IssueDate        string `json:"issue_date"`
	Type             string `json:"type"`
	//OrganicPercentage float32
}
type IDs struct {
	KYC_ID  string
	Crop_ID string
	PlotID  string
}

//----------------------

//var logger = shim.NewLogger("example_cc0")

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("-----------Instantiating---------------------------")
	id := []IDs{
		IDs{
			KYC_ID:  "KYC1001",
			Crop_ID: "CROP1001",
			PlotID:  "Plot1001",
		},
	}
	idsAsBytes, _ := json.Marshal(id[0])
	stub.PutState("ids", idsAsBytes)
	FertAgents := []string{}
	fertagentAsbytes, _ := json.Marshal(FertAgents)
	stub.PutState("FertAgents", fertagentAsbytes)

	//kyc := KYC{First_Name: "Criyagen", Address: "Bangalore", Email: "agricert@criyagen.com"}
	//idsAsBytes, _ = json.Marshal(kyc)
	//stub.PutState("1000", idsAsBytes)
	//
	//agent := CertAgent{KYC_ID: "1000"}
	//idsAsBytes, _ = json.Marshal(agent)
	//stub.PutState("Criyagen", idsAsBytes)

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
	if function == "newuser" {
		return s.KYCRegistration(stub, args)
	} else if function == "PlotRegisteration" {
		return s.PlotRegisteration(stub)
	} else if function == "query" {
		return s.query(stub, args)
	} else if function == "CropDetails" {
		return s.CreateCrop(stub)
	} else if function == "addFertilizerToCrop" {
		return s.addFertilizerToCrop(stub)
	} else if function == "ApproveOrDenyFertilizer" {
		return s.ApproveOrDenyFertilizer(stub, args)
	} else if function == "GetHistoryOfCrop" {
		return s.GetHistoryOfCrop(stub, args)
	} else if function == "queryPlot" {
		return s.GetPlot(stub)
	} else if function == "addStock" {
		return s.addStock(stub)
	} else if function == "queryStock" {
		return s.queryStock(stub, args)
	} else if function == "listDealers" {
		return s.listDealers(stub)
	} else if function == "queryPendingRequests" {
		return s.queryPendingRequests(stub, args)
	} else if function =="getAssetList"{
		return s.getAssetList(stub)
	}else {
		return shim.Error("Invalid Smart contract function aname")
	}

}

func (s *SmartContract) KYCRegistration(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Info("########### KYCRegistration ###########")

	enrollID := getEnrollID(stub)
	role, _, _ := cid.GetAttributeValue(stub, "Role")

	userAsbytes, _ := stub.GetState(enrollID)

	if userAsbytes != nil {
		return shim.Error("User already exists")
	}
	fmt.Printf("enrollID: %s", enrollID)
	//mspID := id.GetMspid()

	newuser := KYC{First_Name: args[0], Last_Name: args[1], Age: args[2], Address: args[3], Email: args[4], MobileNo: args[5]}
	//generating user ID
	idsAsBytes, _ := stub.GetState("ids")
	ids := IDs{}
	json.Unmarshal(idsAsBytes, &ids)
	logger.Debug("ids are ", ids)
	kycID := ids.KYC_ID
	logger.Info("kycID ", kycID)
	j := string([]rune(ids.KYC_ID)[3:])
	logger.Debug("value of J in KYC_ID is ", j)
	tempKYCno, _ := strconv2.Atoi(j)
	tempKYCno = tempKYCno + 1
	ids.KYC_ID = "KYC" + strconv2.Itoa(tempKYCno)
	idsAsBytes, _ = json.Marshal(ids)
	stub.PutState("ids", idsAsBytes)
	// fmt.Println("enrollID", enrollID)
	fmt.Println("user object ", newuser)
	//------------------------
	newuserAsBytes, _ := json.Marshal(newuser)
	stub.PutState(kycID, newuserAsBytes)
	if role == "Fertlizer Agency" {
		newFertagent := FertAgent{}
		newFertagent.KYC_ID = kycID
		listfertAgent := []string{}
		varfertAgentAsbytes, _ := stub.GetState("FertAgents")
		json.Unmarshal(varfertAgentAsbytes, &listfertAgent)
		listfertAgent = append(listfertAgent, kycID)
		varfertAgentAsbytes, _ = json.Marshal(listfertAgent)
		stub.PutState("FertAgents", varfertAgentAsbytes)
		newFertagentAsbytes, _ := json.Marshal(newFertagent)
		stub.PutState(enrollID, newFertagentAsbytes)
		indexname := "enroll~kyc"
		CompositeID, _ := stub.CreateCompositeKey(indexname, []string{kycID, enrollID})
		logger.Info("composite key", CompositeID)
		value := []byte{0x00}
		stub.PutState(CompositeID, value)

	}
	if role == "Certificate Agency" {
		newCertAgent := CertAgent{}
		newCertAgent.KYC_ID = kycID
		newCertAgentAsbytes, _ := json.Marshal(newCertAgent)
		stub.PutState(enrollID, newCertAgentAsbytes)
	}
	if role == "Farmer" {
		newfarmer := Farmer{}
		newfarmer.KYC_ID = kycID
		newfarmerAsBytes, _ := json.Marshal(newfarmer)
		stub.PutState(enrollID, newfarmerAsBytes)
	}
	return shim.Success(nil)
}

func (s *SmartContract) addStock(stub shim.ChaincodeStubInterface) sc.Response {
	logger.Info("---------------------add stock")
	args := stub.GetArgs()
	enrollID := getEnrollID(stub)
	userAsbytes, _ := stub.GetState(enrollID)
	role, _, _ := cid.GetAttributeValue(stub, "Role")
	logger.Info("Role", role)
	if role == "Fertlizer Agency" {
		fert := Fertilizer{}
		if err := json.Unmarshal([]byte(args[1]), &fert); err != nil {
			log.Fatal(err)
		}

		newFertagent := FertAgent{}
		json.Unmarshal(userAsbytes, &newFertagent)
		var index int
		flag := 0
		for index, _ = range newFertagent.Store {
			if newFertagent.Store[index].FertlizerID == fert.FertlizerID {
				flag = 1
				break
			}
		}
		if flag == 0 {
			logger.Info("New fertilizer ")
			newFertagent.Store = append(newFertagent.Store, fert)
			logger.Info("fertlizer Agent", newFertagent)

			userAsbytes, _ = json.Marshal(newFertagent)
			stub.PutState(enrollID, userAsbytes)
			return shim.Success(nil)
		} else {
			newFertagent.Store[index].Quantity = newFertagent.Store[index].Quantity + fert.Quantity
			logger.Info("old fertilzer")
			logger.Info("fertlizer Agent", newFertagent)
			userAsbytes, _ = json.Marshal(newFertagent)
			stub.PutState(enrollID, userAsbytes)
			return shim.Success(nil)
		}

	} else if role == "Farmer" {
		//[[fertlizer0,fertlizer1,...],kycid,]
		fertArr := []Fertilizer{}
		logger.Info(string(args[1]))

		if err := json.Unmarshal([]byte(args[1]), &fertArr); err != nil {
			log.Fatal(err)
		}
		tempFarmer := Farmer{}
		json.Unmarshal(userAsbytes, &tempFarmer)
		kycId := string(args[2])

		logger.Info("reached line 295")
		indexName := "enroll~kyc"
		it, _ := stub.GetStateByPartialCompositeKey(indexName, []string{kycId})
		defer it.Close()
		locationRange, err := it.Next()
		_, compositeKeyParts, err := stub.SplitCompositeKey(locationRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		logger.Info("reached line 304")
		fertenrollID := compositeKeyParts[1]
		logger.Info("feert enroll id ", fertenrollID)
		tempFertAgent := FertAgent{}
		userAsbytes, _ := stub.GetState(fertenrollID)
		json.Unmarshal(userAsbytes, &tempFertAgent)
		logger.Info("reached line 309")
		var farmerIndex int
		var fertIndex int
		farmflag := 0
		for _, fert := range fertArr {
			logger.Info("fertlizerId of the fertilizer request", fert.FertlizerID)
			for fertIndex, _ = range tempFertAgent.Store {
				if tempFertAgent.Store[fertIndex].FertlizerID == fert.FertlizerID {
					//flag=1
					break
				}
			}
			logger.Info("reached line 320")
			for farmerIndex, _ = range tempFarmer.Store {
				if tempFarmer.Store[farmerIndex].FertlizerID == fert.FertlizerID {
					farmflag = 1
					break
				}
			}
			logger.Info("reached line 327")
			if farmflag == 0 { //farmer  doesnt have a stock
				logger.Info("farmer doesnt have the stock")
				tempFarmer.Store = append(tempFarmer.Store, fert)
				logger.Info("farmer store", tempFarmer.Store)
				logger.Info("fert index", fertIndex)
				logger.Info("fertAgent", tempFertAgent)
				tempFertAgent.Store[fertIndex].Quantity = tempFertAgent.Store[fertIndex].Quantity - fert.Quantity

			} else {
				logger.Info("farmer have the stock")
				tempFertAgent.Store[fertIndex].Quantity = tempFertAgent.Store[fertIndex].Quantity - fert.Quantity
				tempFarmer.Store[farmerIndex].Quantity = tempFarmer.Store[farmerIndex].Quantity + fert.Quantity
			}
			farmflag = 0
		}
		userAsbytes, _ = json.Marshal(tempFarmer)
		stub.PutState(enrollID, userAsbytes)
		userAsbytes, _ = json.Marshal(tempFertAgent)
		stub.PutState(fertenrollID, userAsbytes)
		return shim.Success(nil)

	}

	return shim.Error("Unauthorized Access attempted")
}

func (s *SmartContract) listDealers(stub shim.ChaincodeStubInterface) sc.Response {
	FertAgents, _ := stub.GetState("FertAgents")
	var maparr []map[string]string
	var list []string
	json.Unmarshal(FertAgents, &list)
	kyc := KYC{}
	full_name := ""
	for _, element := range list {
		userAsbytes, _ := stub.GetState(element)
		json.Unmarshal(userAsbytes, &kyc)
		full_name = kyc.First_Name + " " + kyc.Last_Name
		mp := make(map[string]string)
		mp["kyc_id"] = element
		mp["full_name"] = full_name
		maparr = append(maparr, mp)
	}
	mapAsbytes, _ := json.Marshal(maparr)
	return shim.Success(mapAsbytes)

}

func (s *SmartContract) queryStock(stub shim.ChaincodeStubInterface, args []string) sc.Response {
logger.Info("=========================================querystock===================")
	enrollID := getEnrollID(stub)
	userAsbytes, _ := stub.GetState(enrollID)
	//role,_,_:=cid.GetAttributeValue(stub,"Role")
	logger.Info("args; ",args)
	if args[0] != "" {
		kycId := args[0]
		indexName := "enroll~kyc"
		it, _ := stub.GetStateByPartialCompositeKey(indexName, []string{kycId})
		defer it.Close()
		locationRange, err := it.Next()
		_, compositeKeyParts, err := stub.SplitCompositeKey(locationRange.Key)
		logger.Info("composite key", compositeKeyParts)
		if err != nil {
			return shim.Error(err.Error())
		}
		fertenrollID := compositeKeyParts[1]
		tempFertAgent := FertAgent{}
		userAsbytes, _ := stub.GetState(fertenrollID)
		json.Unmarshal(userAsbytes, &tempFertAgent)
		stock := tempFertAgent.Store
		stockAsbytes, _ := json.Marshal(stock)
		return shim.Success(stockAsbytes)
	} else {
		tempFarmer := Farmer{}
		json.Unmarshal(userAsbytes, &tempFarmer)
		stockAsbytes, _ := json.Marshal(tempFarmer.Store)
		return shim.Success(stockAsbytes)
		//plotId := args[0]
		//var farmerIndex int
		////flag:=0
		//for farmerIndex, _:= range tempFarmer.Plots {
		//	if plotId==tempFarmer.Plots[farmerIndex].PlotId{
		//		//flag=1
		//		break
		//	}
		//}

	}
	//return shim.Success([]byte("Un Authorized Access"))
}
func (s *SmartContract) PlotRegisteration(stub shim.ChaincodeStubInterface) sc.Response {
	// {"co_ordinates":[{"lattitude":45.15445,"longitude":32.45455},{"lattitude":45.15445,"longitude":32.45455},{"lattitude":45.15445,"longitude":32.45455},{"lattitude":45.15445,"longitude":32.45455}],
	// "survey_no":"Eg456-456-5677",
	// "soil_type":"sasa"}
	logger.Info("-----------------------PlotRegistration--------------------")
	//retreiving the args
	args := stub.GetArgs()

	//-------------------------------
	// creatinng plot object
	tempPlot := Plot{}

	if err := json.Unmarshal([]byte(args[1]), &tempPlot); err != nil {
		log.Fatal(err)
		return shim.Error(err.Error())
	}

	//-----------------------------------

	enrollID := getEnrollID(stub)
	userAsbytes, _ := stub.GetState(enrollID)
	tempfarmer := Farmer{}
	json.Unmarshal(userAsbytes, &tempfarmer)
	//---------------------------------------

	// PlotId geneartion and assigning to Plot
	length := len(tempfarmer.Plots)
	tempPlot.PlotId = "Plot" + strconv2.Itoa(length)

	tempfarmer.Plots = append(tempfarmer.Plots, tempPlot)
	userAsbytes, _ = json.Marshal(tempfarmer)
	stub.PutState(enrollID, userAsbytes)
	return shim.Success(nil)
}

func getEnrollID(stub shim.ChaincodeStubInterface) string {
	creator, _ := stub.GetCreator()
	id := &mspprotos.SerializedIdentity{}
	_ = proto.Unmarshal(creator, id)
	block, _ := pem.Decode(id.GetIdBytes())
	cert, _ := x509.ParseCertificate(block.Bytes)
	enrollID := cert.Subject.CommonName
	return enrollID
}

//fill crop history

func (s *SmartContract) CreateCrop(stub shim.ChaincodeStubInterface) sc.Response {
	//["plotID","crop details obj ","certRegistration"]
	fmt.Print("-------------------cropDetails------------------")
	args := stub.GetArgs()
	plotID := args[1]
	logger.Info("plotId", string(plotID))
	enrollID := getEnrollID(stub)
	tempfarmer := Farmer{}
	farmerAsbytes, _ := stub.GetState(enrollID)
	json.Unmarshal(farmerAsbytes, &tempfarmer)
	flag := 0
	var plot Plot
	var index int
	for index, plot = range tempfarmer.Plots {
		if plot.PlotId == string(plotID) {
			flag = 1
			break
		}
	}
	if flag == 0 {
		return shim.Error("Invalid plot ID")
	}
	if tempfarmer.Plots[index].Current_crop != "" {
		return shim.Error("Crop already registred or need to be harvested")
	}
	logger.Info("plot is ", plot)
	idsAsBytes, _ := stub.GetState("ids")
	ids := IDs{}
	json.Unmarshal(idsAsBytes, &ids)
	logger.Debug("ids are ", ids)
	cropID := ids.Crop_ID
	logger.Info("crop", cropID)
	j := string([]rune(ids.Crop_ID)[4:])
	logger.Debug("value of J in cropID is ", j)
	tempcropno, _ := strconv2.Atoi(j)
	tempcropno = tempcropno + 1
	ids.Crop_ID = "CROP" + strconv2.Itoa(tempcropno)
	idsAsBytes, _ = json.Marshal(ids)
	stub.PutState("ids", idsAsBytes)

	crop := Crop{}
	if err := json.Unmarshal([]byte(args[2]), &crop); err != nil {
		log.Fatal(err)
		return shim.Error(err.Error())
	}
	logger.Info("crop deatils sent is ", crop)

	if string(args[3]) == "certRegistration" {
		logger.Info("Arg[3] pased is ", string(args[3]))
		crop.CertRequestTxID = stub.GetTxID()
		tempfarmer.Plots[index].Current_crop = cropID
		req := Requests{}
		req.Status = "Pending"
		req.Asset = cropID
		reqAsBytes, _ := json.Marshal(req)
		stub.PutState(stub.GetTxID(), reqAsBytes)
		userAsBytes, _ := stub.GetState("Criyagen")
		agent := CertAgent{}
		json.Unmarshal(userAsBytes, &agent)
		agent.CertificateRequests = append(agent.CertificateRequests, stub.GetTxID())
		userAsBytes, _ = json.Marshal(agent)
		stub.PutState("Criyagen", userAsBytes)

	} else {
		tempfarmer.Plots[index].Crop_history = append(tempfarmer.Plots[index].Crop_history, cropID)

	}

	logger.Info("crop before map ", args[2])
	logger.Info("Crop details ", crop)
	cropAsBytes, _ := json.Marshal(crop)
	stub.PutState(cropID, cropAsBytes)

	logger.Info("after adding crop history", plot)
	farmerAsbytes, _ = json.Marshal(tempfarmer)
	stub.PutState(enrollID, farmerAsbytes)
	return shim.Success(nil)

}

func (s *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Info("args", args[0])
	result, _ := stub.GetState(args[0])

	return shim.Success(result)
}

func (s *SmartContract) addFertilizerToStore(stub shim.ChaincodeStubInterface) sc.Response {
	//[{fertlizer_id:"jhjkh",fertlizer_name:"ghjghg",quantity:4}] for fert agent
	//[kycID_of_dealer,[{fertlizer_id:"jhjkh",fertlizer_name:"ghjghg",quantity:4}]]for farmer
	args := stub.GetArgs()
	fert := Fertilizer{}
	if err := json.Unmarshal([]byte(args[1]), &fert); err != nil {
		log.Fatal(err)
		return shim.Error(err.Error())
	}
	enrollID := getEnrollID(stub)
	idAsBytes, _ := stub.GetState(enrollID)
	farmer := Farmer{}
	json.Unmarshal(idAsBytes, &farmer)
	farmer.Store = append(farmer.Store, fert)
	idAsBytes, _ = json.Marshal(farmer)
	stub.PutState(enrollID, idAsBytes)
	return shim.Success(nil)

}
func (s *SmartContract) addFertilizerToCrop(stub shim.ChaincodeStubInterface) sc.Response {
	//[plotId,[{fertlizer_name:"name",fertlizer_id:"asd",quantity:2.56}]]
	//get plot
	enrollID := getEnrollID(stub)
	idAsBytes, _ := stub.GetState(enrollID)
	farmer := Farmer{}
	json.Unmarshal(idAsBytes, &farmer)
	args := stub.GetArgs()
	plotId := string(args[1])
	var index int
	flag := 0
	var cropid string
	for index, _ = range farmer.Plots {

		if farmer.Plots[index].PlotId == plotId {
			flag = 1
			break;
		}
	}
	//check if plot exists and get current crop id
	if flag == 1 {
		cropid = farmer.Plots[index].Current_crop

	} else {
		return shim.Error("plot doesnot belong to you")
	}

	logger.Info("current crop", cropid)
	cropAsBytes, _ := stub.GetState(cropid)
	crop := Crop{}
	json.Unmarshal(cropAsBytes, &crop)
	logger.Info("Crop object: ", crop)
	certAgentAsBytes, _ := stub.GetState("Criyagen")
	agent := CertAgent{}
	json.Unmarshal(certAgentAsBytes, &agent)
	logger.Info("Criyagen agent obj: ", agent)
	fertArr := []Fertilizer{}
	if err := json.Unmarshal([]byte(args[2]), &fertArr); err != nil {
		log.Fatal(err)
		return shim.Error(err.Error())
	}

	flag = 0
	for _, fert := range fertArr {
		flag, _ = validateFertPresence(fert, farmer)
		if flag == 0 {
			return shim.Error("Fertilizer not present in the Store or Insufficient Quantity")
		}
	}

	reqobj := make(map[string]interface{})
	reqobj["cropID"] = crop
	reqobj["fertList"] = fertArr

	txid := stub.GetTxID()
	req := Requests{Status: "pending"}
	req.Asset = reqobj
	reqAsBytes, _ := json.Marshal(req)
	agent.FertilizerRequests = append(agent.FertilizerRequests, txid)
	logger.Info("criyagen after adding request txid")
	stub.PutState(txid, reqAsBytes)
	cropAsBytes, _ = json.Marshal(crop)
	stub.PutState(cropid, cropAsBytes)
	idAsBytes, _ = json.Marshal(agent)
	stub.PutState("Criyagen", idAsBytes)
	return shim.Success(nil)

}

func (s *SmartContract) getAssetList(stub shim.ChaincodeStubInterface) sc.Response {
	logger.Info("============================getAssetList===================================")
	enrollID := getEnrollID(stub)
	role, _, _ := cid.GetAttributeValue(stub, "Role")
	userAsbytes,_:=stub.GetState(enrollID)
	if role == "Fertlizer Agency" {
		//list fertlizer
		agent:=FertAgent{}
		json.Unmarshal(userAsbytes,&agent)
		 return s.queryStock(stub,[]string{agent.KYC_ID})

	} else if role == "Farmer" {
		//list plot current crop and fertstock
		farmer:=Farmer{}
		json.Unmarshal(userAsbytes,&farmer)
		stockResp:=make(map[string]interface{})
		logger.Info("kycid",farmer.KYC_ID)
		resp:=s.queryStock(stub,[]string{""})
		logger.Info("payload",resp)
		stockArrAsbytes:=resp.Payload
		stockArr:=[]Fertilizer{}
		json.Unmarshal(stockArrAsbytes,&stockArr)
		logger.Info(stockArr)
		stockResp["store"]=stockArr
		plotlist:=[]map[string]string{}
		for index,_:=range farmer.Plots{
			mp:=make(map[string]string)
			mp["plot"]=farmer.Plots[index].PlotId
			cropId:=farmer.Plots[index].Current_crop
			cropAsbytes,_:=stub.GetState(cropId)
			crop:=Crop{}
			json.Unmarshal(cropAsbytes,&crop)
			mp["crop"]=crop.CropName
			plotlist=append(plotlist,mp)
		}
		assetList:=make(map[string]interface{})
		assetList["plotList"]=plotlist
		assetList["stockList"]=stockResp
		assetListAsbytes,_:=json.Marshal(assetList)
		return shim.Success(assetListAsbytes)
	} else {
		return shim.Error("Unauthorized Access")
	}

}

func validateFertPresence(fert Fertilizer, farmer Farmer) (int, int) {
	var index int
	flag := 0
	for index, _ = range farmer.Store {
		if farmer.Store[index].FertlizerID == fert.FertlizerID && farmer.Store[index].Quantity >= fert.Quantity {
			flag = 1
			break;
		}
	}
	return flag, index
}
func (s *SmartContract) ApproveOrDenyFertilizer(stub shim.ChaincodeStubInterface, args [] string) sc.Response {
	//	[ 'CROP1002',
	//	'DarshanBC6',
	//	'dfd44d11d76d4f2bcfce965b31f7a8ddee2c89a9f7e689b39a31ce6fcf688f2d',
	//'yes' ]

	req := Requests{}
	fertArr := []Fertilizer{}
	reqAsbytes, _ := stub.GetState(args[2]) //txid
	json.Unmarshal(reqAsbytes, &req)
	req.Status = args[3] //status update if yes or no
	assetAsBytes, _ := json.Marshal(req.Asset)
	json.Unmarshal(assetAsBytes, &fertArr)
	if args[3] == "yes" {
		farmer := Farmer{}
		idsAsBytes, _ := stub.GetState(args[1]) //enrollid
		json.Unmarshal(idsAsBytes, &farmer)
		var index int
		flag := 0
		//==========================
		for _, fert := range fertArr {
			for index, _ = range fertArr {
				flag, index = validateFertPresence(fert, farmer)
				if flag == 1 {
					farmer.Store[index].Quantity = farmer.Store[index].Quantity - fert.Quantity
				} else {
					return shim.Error("Insufficient Stock")
				}

			}
		}
		//=============

		certAgentAsBytes, _ := stub.GetState("Criyagen")
		agent := CertAgent{}
		json.Unmarshal(certAgentAsBytes, &agent)

		for index, element := range agent.FertilizerRequests {
			if strings.Compare(strings.ToLower(element), args[2]) == 0 { //txid

				agent.FertilizerRequests = append(agent.FertilizerRequests[:index], agent.FertilizerRequests[index+1:]...)
			}
		}
		//upload agent
		certAgentAsBytes, _ = json.Marshal(agent)
		stub.PutState("Criyagen", certAgentAsBytes)
		//upload farmer
		usersAsBytes, _ := json.Marshal(farmer)
		stub.PutState(args[1], usersAsBytes)
	}
	//read Crop
	crop, _ := stub.GetState(args[0])
	//write crop
	stub.PutState(args[0], crop)
	return shim.Success(nil)

}

func (s *SmartContract) GetHistoryOfCrop(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	//["plotId"]
	enrollId := getEnrollID(stub)
	farmer := Farmer{}
	farmerAsBytes, _ := stub.GetState(enrollId)
	json.Unmarshal(farmerAsBytes, &farmer)
	var cropId string
	var index int
	var txlist []string
	//var buffer bytes.Buffer
	for index, _ = range farmer.Plots {
		if farmer.Plots[index].PlotId == args[0] {
			cropId = farmer.Plots[index].Current_crop
		}
	}
	//bArrayMemberAlreadyWritten := false
	resultsIterator, err := stub.GetHistoryForKey(cropId)

	//buffer.WriteString("[")
	if err != nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}
	//var modification *queryresult.KeyModification
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		txlist = append(txlist, response.TxId)
	}

	txIdlistAsBytes, _ := json.Marshal(txlist)
	//farmer.Plots[index].Crop_history=append(farmer.Plots[index].Crop_history,farmer.Plots[index].Current_crop)
	//farmer.Plots[index].Current_crop=""
	//farmerAsBytes,_=json.Marshal(farmer)
	//stub.PutState(enrollId,farmerAsBytes)

	return shim.Success(txIdlistAsBytes)
}

func (s *SmartContract) GetPlot(stub shim.ChaincodeStubInterface) sc.Response {

	enrollID := getEnrollID(stub)
	idAsBytes, _ := stub.GetState(enrollID)
	farmer := Farmer{}
	json.Unmarshal(idAsBytes, &farmer)
	var plots []string
	for index, _ := range farmer.Plots {
		plots = append(plots, farmer.Plots[index].PlotId)
	}
	plotsAsbytes, _ := json.Marshal(plots)
	return shim.Success(plotsAsbytes)
}
func (s *SmartContract) queryPendingRequests(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	enrollID := getEnrollID(stub)
	role, _, _ := cid.GetAttributeValue(stub, "Role")
	//list:=[]string{}
	var transactionIDs []string
	if role == "Certificate Agency" {
		userAsBytes, _ := stub.GetState(enrollID)
		agent := CertAgent{}
		json.Unmarshal(userAsBytes, &agent)
		if args[0] == "fertilizer" {
			transactionIDs = agent.FertilizerRequests
		} else if args[0] == "certificate" {
			transactionIDs = agent.CertificateRequests
		} else {
			return shim.Error("Inappropriate Request")
		}
		//for index,_:=range  transactionIDs{
		//	txid:=agent.FertilizerRequests[index]
		//	req:=Requests{}
		//	reqAsbytes,_:=stub.GetState(txid)
		//	json.Unmarshal(reqAsbytes,&req)
		//	if req.Status=="pending"{
		//		list=append(list,txid)
		//	}
		//}
		listAsbytes, _ := json.Marshal(transactionIDs)
		return shim.Success(listAsbytes)
	} else {
		return shim.Error("Unauthorized Access")
	}
}
func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		logger.Error("Error starting Simple chaincode: %s", err)
	}
}

//----------------------fertilizer agency--------------

//type Fertilizer struct {
//	FertlizerName string
//	FertlizerID   string
//	Quantity      float32
//}
