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
	 "bytes"
	 //"github.com/golang/protobuf/ptypes"
	 //"time"
	 //"strconv"
 )
 

// ===========================================================================================
// Get All Record Keys without User Accounts that are registered in the blockchain (args: no args)
// ===========================================================================================
func (s *SmartContract) getAllKeys(APIstub shim.ChaincodeStubInterface , args []string) sc.Response {

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}
	
	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	errorMSP, mspID := getMspID(APIstub)
	if(!errorMSP){
		return shim.Error(mspID)
	}
	if(mspID != "AuditorMSP"){
		return shim.Error("You have to be member of Auditor Organization to access these data")
	}

	queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$ne\":\"user\"}}}")

	resultsIterator, err :=	APIstub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error("Error occured: "+err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {
	
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("Error occured when trying to get data from WS: "+err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		var typeS Type
		json.Unmarshal(queryResponse.Value, &typeS)

		if(typeS.Type == "access"){
			// Record is a JSON object, so we write as-is
			var access AccessControl
			json.Unmarshal(queryResponse.Value, &access)
		
			buffer.WriteString(string("access."+"file."+access.FileID+"."+access.UserID))

			fmt.Println(string(queryResponse.Value))

		}else if(typeS.Type == "document"){
			// Record is a JSON object, so we write as-is

			var document Document
			json.Unmarshal(queryResponse.Value, &document)
		
			buffer.WriteString(string(document.ID))
			fmt.Println(string(queryResponse.Value))

		}else if(typeS.Type == "loanreq"){
			// Record is a JSON object, so we write as-is
			var loanRequest LoanRequest
			json.Unmarshal(queryResponse.Value, &loanRequest)
		
			buffer.WriteString(string("request."+loanRequest.Title+"."+loanRequest.FromAcc))
			fmt.Println(string(queryResponse.Value))
		}

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

 }

 // ===========================================================================================
// Get All Record Keys without User Accounts that are registered in the blockchain (args: no args)
// ===========================================================================================
func (s *SmartContract) getAllKeysBasedOnApplicantId(APIstub shim.ChaincodeStubInterface , args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting an ID")
	}

	applicantId := args[0]
	fmt.Println(applicantId)

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}
	
	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	errorMSP, mspID := getMspID(APIstub)
	if(!errorMSP){
		return shim.Error(mspID)
	}
	if(mspID != "AuditorMSP"){
		return shim.Error("You have to be member of Auditor Organization to access these data")
	}

	callerData, err := APIstub.GetState(applicantId)
	if err != nil {
		return shim.Error("Something went wrong")
	}
	if callerData == nil {
		return shim.Error("Applicant does not exists")
	}

	var account Account
	json.Unmarshal(callerData, &account)

	if(account.UserType != "applicant"){
		return shim.Error("The user is not applicant. You are not authorized to access these data.")
	}

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	queryString :=  ""

	for i := 1; i <= 3; i++ {
		if(i == 1){
			queryString = fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"creator\": { \"$eq\":\"%v\"} }}", applicantId)
		}else if( i== 2){
			queryString = fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"owner\": { \"$eq\":\"%v\"} }}", applicantId)
		}else{
			queryString = fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"fromacc\": { \"$eq\":\"%v\"} }}", account.Username)
		}
		//queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":\"user\"}}")
	
		resultsIterator, err :=	APIstub.GetQueryResult(queryString)
		if err != nil {
			return shim.Error("Error occured: "+err.Error())
		}
		defer resultsIterator.Close()
	

		for resultsIterator.HasNext() {
		
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return shim.Error("Error occured when trying to get data from WS: "+err.Error())
			}
			// Add a comma before array members, suppress it for the first array member
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}

			var typeS Type
			json.Unmarshal(queryResponse.Value, &typeS)

			if(typeS.Type == "access"){
				// Record is a JSON object, so we write as-is
				var access AccessControl
				json.Unmarshal(queryResponse.Value, &access)
			
				buffer.WriteString(string("access."+"file."+access.FileID+"."+access.UserID))

				fmt.Println(string(queryResponse.Value))

			}else if(typeS.Type == "document"){
				// Record is a JSON object, so we write as-is

				var document Document
				json.Unmarshal(queryResponse.Value, &document)
			
				buffer.WriteString(string(document.ID))
				fmt.Println(string(queryResponse.Value))

			}else if(typeS.Type == "loanreq"){
				// Record is a JSON object, so we write as-is
				var loanRequest LoanRequest
				json.Unmarshal(queryResponse.Value, &loanRequest)
			
				buffer.WriteString(string("request."+loanRequest.Title+"."+loanRequest.FromAcc))
				fmt.Println(string(queryResponse.Value))
			}

			bArrayMemberAlreadyWritten = true
		}
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

 }



 // ===========================================================================================
// Get All Applicant Record Keys that have registered in the blockchain (args: no args)
// ===========================================================================================
func (s *SmartContract) getAllApplicantIds(APIstub shim.ChaincodeStubInterface , args []string) sc.Response {

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}
	
	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	errorMSP, mspID := getMspID(APIstub)
	if(!errorMSP){
		return shim.Error(mspID)
	}
	if(mspID != "AuditorMSP"){
		return shim.Error("You have to be member of Auditor Organization to access these data")
	}

	queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"user\"},\"usertype\":{\"$eq\":\"applicant\"}}}")

	resultsIterator, err :=	APIstub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error("Error occured: "+err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {
	
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("Error occured when trying to get data from WS: "+err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		var accounts Account
		json.Unmarshal(queryResponse.Value, &accounts)

		buffer.WriteString(string(accounts.ID))
		fmt.Println(string(queryResponse.Value))

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

 }


// ===========================================================================================
// Get loan requests history for a request (args: req Id)
// Admin
// ===========================================================================================
func (s *SmartContract) getReqHistoryByIdAuditor(APIstub shim.ChaincodeStubInterface , args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting an ID")
	}

	reqId := args[0]

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}
	
	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	errorMSP, mspID := getMspID(APIstub)
	if(!errorMSP){
		return shim.Error(mspID)
	}
	if(mspID != "AuditorMSP"){
		return shim.Error("You have to be member of Auditor Organization to access these data")
	}

	callerData, err := APIstub.GetState(reqId)
	if err != nil {
		return shim.Error("Something went wrong")
	}
	if callerData == nil {
		return shim.Error("Data do not exists")
	}

	var typeS Type
	json.Unmarshal(callerData, &typeS)

	if(typeS.Type != "access" && typeS.Type != "loanreq" && typeS.Type != "document"){
		return shim.Error("You are not authorized to access these data.")
	}

	resultsIterator, err := APIstub.GetHistoryForKey(reqId)
	if err != nil {
		return shim.Error("Something wrong has happened.")
	}

	defer resultsIterator.Close()

	
	jsonBytes, err := createHistoryResultsJson(resultsIterator,APIstub)
	if err != nil {
		return shim.Error("Error occured when trying to createHistoryResultsJson: "+ err.Error())
	}

	return shim.Success(jsonBytes)

}
