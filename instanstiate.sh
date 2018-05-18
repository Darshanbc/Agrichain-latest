jq --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
	echo
	exit 1
fi

starttime=$(date +%s)

# Print the usage message
function printHelp () {
  echo "Usage: "
  echo "  ./testAPIs.sh -l golang|node"
  echo "    -l <language> - chaincode language (defaults to \"golang\")"
}
# Language defaults to "golang"
LANGUAGE="golang"

# Parse commandline args
while getopts "h?l:" opt; do
  case "$opt" in
    h|\?)
      printHelp
      exit 0
    ;;
    l)  LANGUAGE=$OPTARG
    ;;
  esac
done

##set chaincode path
function setChaincodePath(){
	LANGUAGE=`echo "$LANGUAGE" | tr '[:upper:]' '[:lower:]'`
	case "$LANGUAGE" in
		"golang")
		CC_SRC_PATH="github.com/example_cc/go"
		;;
		"node")
		CC_SRC_PATH="$PWD/artifacts/src/github.com/example_cc/node"
		;;
		*) printf "\n ------ Language $LANGUAGE is not supported yet ------\n"$
		exit 1
	esac
}

setChaincodePath

echo "POST request Enroll on Org1  ..."
echo
ORG1_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Darshan123456&orgName=Org1&role="Farmer"')
echo $ORG1_TOKEN
ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".token" | sed "s/\"//g")



# echo "POST invoke KYC on peers of Org1 and Org2"
# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{
#   		"fcn":"newuser",
#   		"args":["Darshan","BC","25","bangalore","darshan@gmail.com","8904374405"] 	
# }')
# echo "Transacton ID is $TRX_ID"




# echo "POST invoke PlotRegisteration on peers of Org1 and Org2"

# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{

# 	"fcn":"PlotRegisteration",
# 	"args":["{\"survey_no\":\"Eg456-456-5677\",\"soil_type\":\"sasa\",\"co_ordinates\":[{\"lattitude\":45.15445,\"longitude\":32.45455},{\"lattitude\":45.15445,\"longitude\":32.45455},{\"lattitude\":45.15445,\"longitude\":32.45455},{\"lattitude\":45.15445,\"longitude\":32.45455}]}"]
 	
# }')
# echo "Transacton ID is $TRX_ID"


# echo "POST invoke insert crop history on peers of Org1 and Org2"
# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{

# 	"fcn":"CropDetails",
# 	"args":["Plot0","{\"crop_name\":\"rice\",\"crop_cycle\":{\"from_month\":\"june\",\"to_month\":\"september\"},\"fertilzer_used\":[{\"fertlizer_name\":\"urea\",\"fertlizer_id\":\"U35\",\"quantity\":2.5}]}","1"]
 	
# }')

# echo "Transacton ID is $TRX_ID"
# echo
# echo

# echo "POST invoke certRegistration on peers of Org1 and Org2"
# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{

# 	"fcn":"CropDetails",
# 	"args":["Plot0","{\"crop_name\":\"rice\",\"crop_cycle\":[{\"from_month\":\"june\",\"to_month\":\"september\"}],\"fertilzer_used\":[]}","certRegistration"]
 	
# }')

# echo "Transacton ID is $TRX_ID"
# echo
# echo


# echo "POST invoke addFertilizerToStore on peers of Org1 and Org2"
# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{

# 	"fcn":"addFertilizerToStore",
# 	"args":["kycid","[{\"fertlizer_name\":\"urea\",\"fertlizer_id\":\"U35\",\"quantity\":2.5}"]
 	
# }')

# echo "Transacton ID is $TRX_ID"
# echo
# echo


# echo "POST invoke addFertilizerToCrop on peers of Org1 and Org2"
# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{

# 	"fcn":"addFertilizerToCrop",
# 	"args":["Plot0","{\"fertlizer_name\":\"urea\",\"fertlizer_id\":\"U35\",\"quantity\":2.5}"]
 	
# }')

# echo "Transacton ID is $TRX_ID"
# echo
# echo


# echo "GET query chaincode on peer1 of Org1"
# echo
# key="DarshanBC6"
# curl -s -X GET \
#   "http://localhost:4000/channels/mychannel/chaincodes/dtwin?peer=peer0.org1.example.com&fcn=query&args=%5B%22$key%22%5D" \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json"
# echo
# echo

# # TRX_ID=0a025c8701db596bf997bd6f31952bf0586faec2426f6857130879f2759aa5a6
# echo "GET query Transaction by TransactionID"
# echo
# curl -s -X GET http://localhost:4000/channels/mychannel/transactions/$TRX_ID?peer=peer0.org1.example.com \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json"
# echo
# echo
# key="yes"
# echo "GET Approve or deny by TransactionID"
# curl -s -X GET \
# "http://localhost:4000/channels/mychannel/chaincodes/dtwin/$TRX_ID?peer=peer0.org1.example.com&fcn=ApproveOrDenyFertilizer&permission=$key"\
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json"
# echo "Transacton ID is $TRX_ID"
# echo
# echo


# echo "POST invoke  Harvest on peers of Org1 and Org2"
# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{

#   "fcn":"Harvest",
#   "args":["Plot0"]
  
# }')

# echo "Transacton ID is $TRX_ID"
# echo
# echo

# ORG2_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjU5ODg2NjIsInVzZXJuYW1lIjoiZmVydCAxIiwib3JnTmFtZSI6Ik9yZzEiLCJpYXQiOjE1MjU5NTI2NjJ9.MqYKcVikwQQFZDNxe5ylhPlWwOmGu1l5sXNLtxspoDs
# echo
TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjY2NTQ3MjgsInVzZXJuYW1lIjoiZmFybSIsIm9yZ05hbWUiOiJPcmcxIiwiaWF0IjoxNTI2NjE4NzI4fQ.1mgyh9u8V-1gXpmS2_KIQB3Yu6WvdKRPuiS21kTdyio
key=""
# TRX_ID="cb799da8256834a4ad2bbebbef1febd0e47c44d28a47f6e0b2ed1bda9e54e2a4"
echo "GET query chaincode on peer1 of Org1"
echo $key
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/dtwin?peer=peer0.org1.example.com&fcn=getAssetList&args=%5B%22$key%22%5D" \
  -H "authorization: Bearer $TOKEN" \
  -H "content-type: application/json"
echo
echo


# echo "GET query write sets Transaction by TransactionID"
# echo
# curl -s -X GET http://localhost:4000/channels/mychannel/transactions/$TRX_ID/writeSet?peer=peer0.org1.example.com \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json"
# echo
# echo

# echo "GET query Transaction by TransactionID"
# echo
# curl -s -X GET http://localhost:4000/channels/mychannel/transactions/$TRX_ID?peer=peer0.org1.example.com \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json"
# echo
# echo

# echo "POST invoke  add Stock"
# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $TOKEN" \
#   -H "content-type: application/json" \
#   -d '{

#   "fcn":"addStock",
#   "args":["[{\"fertlizer_name\":\"asd\",\"fertlizer_id\":\"1\",\"quantity\":1}]","KYC1001"]
  
# }')



# {"asset":{"cropID":"CROP1001","fertList":[{"fertlizer_name":"dfg","fertlizer_id":"a","quantity":1},{"fertlizer_name":"sef","fertlizer_id":"d","quantity":2}]},"status":"pending"}

# [{"key":"CROP1001","is_delete":false,"value":"{"crop_name":"as","type":"","crop_cycle":{"from_month":"01-05-2018","to_month":"14-05-2018"},"fertilzer_used":null,"fertilzer_req":null,"cert_request_tx_id":"734fe06f946b6c782ba67a2ba05af5c2791c60ad8c886fa1ba9d5230e4d9f5ca","cert":{"certficate_id":"","digital_signature":"","issue_date":"","type":""}}"},{"key":"Criyagen","is_delete":false,"value":"{"KYC_ID":"KYC1001","CertificateRequests":["734fe06f946b6c782ba67a2ba05af5c2791c60ad8c886fa1ba9d5230e4d9f5ca"],"FertilizerRequests":["d22f78b213b8ba0ee29252030d578cb53010d5dd8256ac07db95f0dfecf9b7db","c58a3ee03ca7a7f3651b2e028baf0ed7ea4e641cf0a39a11e3390c59a09d56fa"]}"},{"key":"c58a3ee03ca7a7f3651b2e028baf0ed7ea4e641cf0a39a11e3390c59a09d56fa","is_delete":false,"value":"{"asset":{"cropID":"CROP1001","fertList":[{"fertlizer_name":"dfg","fertlizer_id":"a","quantity":1},{"fertlizer_name":"sef","fertlizer_id":"d","quantity":2}]},"status":"pending"}"}]