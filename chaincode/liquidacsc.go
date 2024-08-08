/*
 * This is the liquidsc file of chaincode for the project of Liquid 
 * Author: Rafail Kiloudis
 */

 package main

 import (
	 "encoding/json"
	 //"strings"
	 //"crypto/rand"
	 "github.com/hyperledger/fabric-chaincode-go/shim"
	 //"github.com/hyperledger/fabric-contract-api-go/contractapi"
	 // sc "github.com/hyperledger/fabric/protos/peer"
	 sc "github.com/hyperledger/fabric-protos-go/peer"
 
	 //"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	 //"github.com/hyperledger/fabric-chaincode-go/pkg/statebased"
	 //"log"
	 "fmt"
	 //"bytes"
	 //"github.com/golang/protobuf/ptypes"
	 //"time"
	 //"strconv"
 )
 

// ===========================================================================================
// Give access rights to a financer (args: user Id, json with all data)
// ===========================================================================================
func (s *SmartContract) addAccessControl(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 args")
	}

	accId := args[0]
	jsonData := args[1]
	
	fmt.Println(jsonData)

	//Store json data tp an access control object
	var access AccessControl
	json.Unmarshal([]byte(jsonData), &access)


	//create the file ID
	fileID := "file." + access.FileID
	
	accessID := "access." + fileID + "." + accId
	access.UserID = accId
	access.Type = "access"

	value, err := APIstub.GetState(fileID)
	if err != nil {
		return shim.Error("Failed to get data from WS key :"+ fileID)
	}
	if value == nil {
		return shim.Error("Failed to get data from WS key :"+ fileID)
	}

	//Get Document from World State
	var document Document
	json.Unmarshal(value, &document)

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	
	//Check if owner is performing this action
	if(document.Owner != callerId){
		return shim.Error("You are not authorized to perform this action")
	}

	access.Creator = callerId

	// Save to WS state
	accessAsBytes, _ := json.Marshal(access)
	err =APIstub.PutState(accessID, accessAsBytes)
	if err != nil {
		return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
	}


	return shim.Success(nil) 
}

// ===========================================================================================
// Give access rights to a financer (args: user Id, json with all data)
// ===========================================================================================
func (s *SmartContract) grantAccessAll(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 args")
	}

	accId := args[0]
	jsonData := args[1]

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"document\"},\"owner\":{\"$eq\":\"%v\"}}}", callerId)

	resultsIterator, err := APIstub.GetQueryResult(queryString)
    if err != nil {
    	return shim.Error("Error occured: "+err.Error())
    }

	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		key_ws_rec := responseRange.Key
			
		docValue, _ := APIstub.GetState(key_ws_rec)

		//Store json data tp an access control object
		var document Document
		json.Unmarshal(docValue, &document)
		
		//create the file ID
		fileID := document.ID

		//Store json data tp an access control object
		var access AccessControl
		json.Unmarshal([]byte(jsonData), &access)
		
		access.FileID = fileID
		access.Filename = document.Filename
		access.UserID = accId
		accessID := "access." + fileID + "." + accId
		access.Type = "access"
		access.Creator = callerId

		// Save to WS state
		accessAsBytes, _ := json.Marshal(access)
		err =APIstub.PutState(accessID, accessAsBytes)
		if err != nil {
			return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
		}
	}



	return shim.Success(nil) 
}

// ===========================================================================================
// List all access rights to any document (args: none)
// Only for Financer has access to these data
// ===========================================================================================
func (s *SmartContract) getAllAccessControlsByAccId(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	// if len(args) != 2 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 2 args")
	// }
	
	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	callerData, err := APIstub.GetState(callerId)
	if err != nil {
		return shim.Error("Something went wrong")
	}
	if callerData == nil {
		return shim.Error("Something went wrong")
	}

	var account Account
	json.Unmarshal(callerData, &account)
	var queryString = ""

	if(account.UserType == "financer"){
		queryString = fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"access\"},\"userid\":{\"$eq\":\"%v\"}},\"sort\":[{\"_id\":\"asc\"}]}", callerId)
	}else if(account.UserType == "applicant"){
		queryString = fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"access\"},\"creator\":{\"$eq\":\"%v\"}},\"sort\":[{\"_id\":\"asc\"}]}", callerId)
	}


	resultsIterator, err := getQueryResultForQueryString(APIstub, queryString)
    if err != nil {
    	return shim.Error("Error occured: "+err.Error())
    }

	return shim.Success(resultsIterator)

}

// ===========================================================================================
// List all access rights to any document (args: none)
// Only for Financer has access to these data
// ===========================================================================================
func (s *SmartContract) deleteAccessControl(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 args")
	}

	docId := args[0]
	accId := args[1]
	
	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}



	fileID := "file." + docId

	docValue, err := APIstub.GetState(fileID)
	if err != nil {
		return shim.Error("Failed to get data from WS key"+ fileID)
	}
	if docValue == nil {
		return shim.Error("Failed to get data from WS key :"+ fileID)
	}

	//Store json data
	var document Document
	json.Unmarshal(docValue, &document)


	if(document.Owner != callerId){
		return shim.Error("You are not authorized to perform this action")
	}

	accessID := "access." + fileID + "." + accId

	value, err := APIstub.GetState(accessID)
	if err != nil {
		return shim.Error("Failed to get data from WS key"+ accessID)
	}
	if value == nil {
		return shim.Error("Failed to get data from WS key :"+ accessID)
	}


	err = APIstub.DelState(accessID)
	if err != nil {
		return shim.Error("Failed to delete data from WS key :"+ accessID)
	}

	return shim.Success(nil)

}
