/*
 * This is the liquid loan requests smart contract file of chaincode for the project of Liquid 
 * Author: Rafail Kiloudis
 */

 package main

 import (
	 "encoding/json"
	 "strings"
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
	 "time"
	 //"strconv"
 )

// ===========================================================================================
// Add User Document (args: user Id, json with all data)
// ===========================================================================================
func (s *SmartContract) addLoanRequest(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 args")
	}

	jsonData := args[0]
	
	fmt.Println(jsonData)

	var loanRequest LoanRequest
	json.Unmarshal([]byte(jsonData), &loanRequest)

	if(loanRequest.ToAcc == ""){
		return shim.Error("Financer account is required!")
	}

	financerUserId := loanRequest.ToAcc
	// check that does financer exists in the  WS ledger
	financerInfo, err := APIstub.GetState(financerUserId)
	if err != nil {
		return shim.Error("Financer not found")
	}
	if financerInfo == nil {
		return shim.Error("Financer not found")
	}

	var financerAccount Account
	json.Unmarshal(financerInfo, &financerAccount)
	
	if(financerAccount.UserType != "financer"){
		return shim.Error("The account is not a Financer account")
	}

	docsList := loanRequest.DocumentLt
	if(docsList == ""){
		return shim.Error("Please specify a list of documents")
	}
	docs := strings.Split(docsList, "/")

	for index,element := range docs{
	fmt.Println(index)
	fmt.Println(element)  

		fileID := "file." + element
		//check that the document exist
		docInfo, err := APIstub.GetState(fileID)
		if err != nil {
			return shim.Error("Document not found")
		}
		if docInfo == nil {
			return shim.Error("Document not found")
		}

		//Get Document from World State
		var document Document
		json.Unmarshal(docInfo, &document)

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
			return shim.Error("Some of the files are not owned by the user who made the transaction")
		}

		//Store json data tp an access control object
		var access AccessControl
		json.Unmarshal([]byte(jsonData), &access)
		
		accessID := "access." + fileID + "." + financerUserId
		access.UserID = financerUserId
		access.Type = "access"
		access.FileID = element
		access.Filename = document.Filename
		access.ValidFrom = loanRequest.ValidFrom
		access.ValidTo = loanRequest.ValidTo
		access.Access = "true"
		access.Status = "active"
		
		// Save to WS state
		accessAsBytes, _ := json.Marshal(access)
		err =APIstub.PutState(accessID, accessAsBytes)
		if err != nil {
			return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
		}

		
	}

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	requestId := "request." + loanRequest.Title + "." + caller
	txntmsp, err := APIstub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Error occured when trying to get the timestmap."+err.Error())
	}

	txID := APIstub.GetTxID();

	loanRequest.TxId = txID;
	loanRequest.Date = time.Unix(txntmsp.Seconds, int64(txntmsp.Nanos)).String()
	loanRequest.FromAcc = caller;
	loanRequest.Type = "loanreq";
	loanRequest.Status = "pending";
	loanRequest.UserLastChange = caller

	// check that does not exist in the  WS ledger
	checkExistsAsBytes, err := APIstub.GetState(requestId)

	// If the key does not exist in the state database, (nil, nil) is returned.
	if checkExistsAsBytes != nil {
		return shim.Error("The provided ID exists already in the WS ledger. ")
	}

 
	// Save to WS state
	accountAsBytes, _ := json.Marshal(loanRequest)
	err =APIstub.PutState(requestId, accountAsBytes)
	if err != nil {
		return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
	}


	return shim.Success(nil) 
}

// ===========================================================================================
// Update Loan Request (args: user Id, json with all data)
// ===========================================================================================
func (s *SmartContract) updateLoanRequest(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 args")
	}

	titleOfReq := args[0]
	jsonData := args[1]
	
	var loanRequest LoanRequest
	json.Unmarshal([]byte(jsonData), &loanRequest)

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	reqId := "request." + titleOfReq + "." + caller

	value, err := APIstub.GetState(reqId)
	if err != nil {
		return shim.Error("Something went wrong")
	}
	if value == nil {
		return shim.Error("Request does not exists")
	}

	//Get Loan Request from World State
	var loanRequestWS LoanRequest
	json.Unmarshal(value, &loanRequestWS)

	if(loanRequestWS.FromAcc != caller){
		return shim.Error("You are not authorized to access these data")
	}

	if(loanRequestWS.Status != "pending" && loanRequestWS.Status != "actionrequired"){
		return shim.Error("You are not authorized to perfrorm this action")
	}

	txntmsp, err := APIstub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Error occured when trying to get the timestmap."+err.Error())
	}

	txID := APIstub.GetTxID();

	loanRequestWS.TxId = txID;
	loanRequestWS.Date = time.Unix(txntmsp.Seconds, int64(txntmsp.Nanos)).String()
	loanRequestWS.Status = "pending";
	loanRequestWS.Description = loanRequest.Description;
	loanRequestWS.Amount = loanRequest.Amount;
	loanRequestWS.ValidFrom = loanRequest.ValidFrom;
	loanRequestWS.ValidTo = loanRequest.ValidTo;
	loanRequestWS.Currency = loanRequest.Currency;
	loanRequestWS.Note = loanRequest.Note
	loanRequestWS.UserMsg = loanRequest.UserMsg
	// loanRequestWS.DocumentLt = loanRequest.DocumentLt
	loanRequestWS.UserLastChange = caller

	var documents = loanRequest.DocumentLt

	//Old Documents
	oldDocsList := loanRequestWS.DocumentLt
	if(oldDocsList == ""){
		return shim.Error("Please specify a list of documents")
	}
	oldDocs := strings.Split(oldDocsList, "/")
	currentNumberOfDocs := len(oldDocs)


	//New Documents
	docsList := loanRequest.DocumentLt
	if(docsList == ""){
		return shim.Error("Please specify a list of documents")
	}
	docs := strings.Split(docsList, "/")
	newNumberOfDocs := len(docs)

	if(currentNumberOfDocs > newNumberOfDocs){
		for index,element := range oldDocs{
			fmt.Println(index)
			fmt.Println(element)  
		
			if(strings.Contains(docsList, element)){
				continue
			}else{
				fileID := "file." + element					
				accessID := "access." + fileID + "." + loanRequestWS.ToAcc

				err = APIstub.DelState(accessID)
				if err != nil {
					continue
				}
			}
		}
	}

	for index,element := range docs{
	fmt.Println(index)
	fmt.Println(element)  

		if(strings.Contains(documents, element)){
			continue
		}else{

			fileID := "file." + element
			//check that the document exist
			docInfo, err := APIstub.GetState(fileID)
			if err != nil {
				return shim.Error("Document not found")
			}
			if docInfo == nil {
				return shim.Error("Document not found")
			}

			//Get Document from World State
			var document Document
			json.Unmarshal(docInfo, &document)

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
				return shim.Error("Some of the files are not owned by the user who made the transaction")
			}

			//Store json data tp an access control object
			var access AccessControl
			json.Unmarshal([]byte(jsonData), &access)
			
			accessID := "access." + fileID + "." + loanRequestWS.ToAcc
			access.UserID = loanRequestWS.ToAcc
			access.Type = "access"
			access.FileID = element
			access.Filename = document.Filename
			access.ValidFrom = loanRequest.ValidFrom
			access.ValidTo = loanRequest.ValidTo
			access.Access = "true"
			access.Status = "active"
			
			// Save to WS state
			accessAsBytes, _ := json.Marshal(access)
			err =APIstub.PutState(accessID, accessAsBytes)
			if err != nil {
				return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
			}
		}
	}

	// Save to WS state
	loanRequestAsBytes, _ := json.Marshal(loanRequestWS)
	err =APIstub.PutState(reqId, loanRequestAsBytes)
	if err != nil {
		return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
	}

	return shim.Success(nil) 
}

// ===========================================================================================
// Action on a Loan Request (args: user Id, json with all data)
// ===========================================================================================
func (s *SmartContract) actionLoanRequest(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 args")
	}

	titleOfReq := args[0]
	jsonData := args[1]
	var reqId = ""

	var loanRequestUpdated LoanRequest
	json.Unmarshal([]byte(jsonData), &loanRequestUpdated)

	reqOwner := loanRequestUpdated.FromAcc

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

	var request LoanRequest

	if(account.UserType != "financer"){
		return shim.Error("You are not authorized to perfrorm this action.")
	}

	queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"loanreq\"},\"fromacc\":{\"$eq\":\"%v\"} ,\"toacc\":{\"$eq\":\"%v\"}, \"title\":{\"$eq\":\"%v\"}},\"sort\":[{\"_id\":\"asc\"}]}",reqOwner, callerId, titleOfReq)

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
			
		reqValue, _ := APIstub.GetState(key_ws_rec)

		json.Unmarshal(reqValue, &request)

		reqId = "request." + titleOfReq + "." + request.FromAcc

	}
	
	
	if(request.Status != "pending"){
		return shim.Error("You can not to perfrorm this action")
	}

	if(loanRequestUpdated.Status != "accepted" && loanRequestUpdated.Status != "rejected" && loanRequestUpdated.Status != "actionrequired"){
		return shim.Error("You are not authorized to perfrorm this action")
	}

	txntmsp, err := APIstub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Error occured when trying to get the timestmap."+err.Error())
	}

	txID := APIstub.GetTxID();

	request.TxId = txID;
	request.Date = time.Unix(txntmsp.Seconds, int64(txntmsp.Nanos)).String()
	request.Status = loanRequestUpdated.Status
	request.UserLastChange = caller
	request.UserMsg = loanRequestUpdated.UserMsg
 
	// Save to WS state
	requestAsBytes, _ := json.Marshal(request)
	err =APIstub.PutState(reqId, requestAsBytes)
	if err != nil {
		return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
	}

	return shim.Success(nil) 
}

// ===========================================================================================
// Get information of an account (args: user Id)
// Company Admin and the user itself can get the information of the user asked
// ===========================================================================================
func (s *SmartContract) deleteLoanRequest(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting an ID")
	}

	titleOfReq := args[0]

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	reqId := "request." + titleOfReq + "." + caller

	value, err := APIstub.GetState(reqId)
	if err != nil {
		return shim.Error("Something went wrong")
	}
	if value == nil {
		return shim.Error("Request does not exists")
	}

	//Get Document from World State
	var request LoanRequest
	json.Unmarshal(value, &request)


	if(request.FromAcc != caller){
		return shim.Error("You are not authorized to access these data")
	}
		
	// Delete from World State and mark deleted in Ledger
	err = APIstub.DelState(reqId)
	if err != nil {
		return shim.Error("Error occured when trying delete information: "+err.Error())
	}
  
	return shim.Success(nil) 
}

// ===========================================================================================
// Get information of an account (args: user Id)
// Company Admin and the user itself can get the information of the user asked
// ===========================================================================================
func (s *SmartContract) getLoanRequestById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting ID and Owner")
	}

	titleOfReq := args[0]
	ownerOfReq := args[1]
	var reqId = ""

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

	var request LoanRequest

	if(account.UserType == "financer"){
		queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"loanreq\"}, \"fromacc\":{\"$eq\":\"%v\"},\"toacc\":{\"$eq\":\"%v\"}, \"title\":{\"$eq\":\"%v\"}},\"sort\":[{\"_id\":\"asc\"}]}", ownerOfReq, callerId, titleOfReq)

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
				
			reqValue, _ := APIstub.GetState(key_ws_rec)
	
			json.Unmarshal(reqValue, &request)

			reqId = "request." + titleOfReq + "." + request.FromAcc

		}
	}else{


		reqId = "request." + titleOfReq + "." + caller

		value, err := APIstub.GetState(reqId)
		if err != nil {
			return shim.Error("Something went wrong")
		}
		if value == nil {
			return shim.Error("Request does not exists")
		}

		//Get Document from World State
		json.Unmarshal(value, &request)
	}


	if(request.FromAcc != caller && request.Status == "Expired"){
		return shim.Error("You are not authorized to access these data")
	}

	if(request.FromAcc != caller && request.ToAcc != callerId){
		return shim.Error("You are not authorized to access these data")
	}
		
	validFrom, _ := time.Parse(layoutISO, request.ValidFrom)
	validTo, _ := time.Parse(layoutISO, request.ValidTo)

	fmt.Println(validFrom)                  // 1999-12-31 00:00:00 +0000 UTC
	fmt.Println(validTo) 
	
	tm := time.Now()

	g1 := tm.Before(validTo)
	g2 := tm.After(validFrom)

	if((!g1 || !g2) && (request.Status != "Expired") && (request.FromAcc != caller)) {
		request.Status = "Expired"

		// Save to WS state
		requestAsBytes, _ := json.Marshal(request)
		err = APIstub.PutState(reqId, requestAsBytes)
		if err != nil {
			return shim.Error("Error occured when trying to update loan infromation in WS ledger: "+err.Error())
		}

		return shim.Error("You are not authorized to get these data. The request has expired")
	}
	
	requestAsBytesToDeliver, _ := json.Marshal(request)

	return shim.Success(requestAsBytesToDeliver) 
}

// ===========================================================================================
// List all access rights to any document (args: none)
// Only for Financer has access to these data
// ===========================================================================================
func (s *SmartContract) getRequestsByAccId(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
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

	queryString := ""

	if(account.UserType == "financer"){
		queryString = fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"loanreq\"},\"toacc\":{\"$eq\":\"%v\"}},\"sort\":[{\"_id\":\"asc\"}]}", callerId)
	}else{
		queryString = fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"loanreq\"},\"fromacc\":{\"$eq\":\"%v\"}},\"sort\":[{\"_id\":\"asc\"}]}", caller)
	}

	resultsIterator, err := getQueryResultForQueryString(APIstub, queryString)
    if err != nil {
    	return shim.Error("Error occured: "+err.Error())
    }

	return shim.Success(resultsIterator)

}

// ===========================================================================================
// Get document history for a document from a specific account (args: user Id, file Id)
// Owner of the document 
// ===========================================================================================
func (s *SmartContract) getRequestHistory(APIstub shim.ChaincodeStubInterface , args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting an ID")
	}

	titleOfReq := args[0]

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	reqId := "request." + titleOfReq + "." + caller
	
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
