/*
 * This is the liquid messaging smart contract file of chaincode for the project of Liquid 
 * Author: Rafail Kiloudis
 */

 package main

 import (
	 "encoding/json"
	 "github.com/hyperledger/fabric-chaincode-go/shim"
	 sc "github.com/hyperledger/fabric-protos-go/peer"
 
	 "fmt"
	 "time"
 )

// ===========================================================================================
// Send Message (args: title of request, message)
// ===========================================================================================
func (s *SmartContract) sendMessage(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2 args")
	}

	jsonData := args[0]
	
	fmt.Println(jsonData)

	var receivedMessage ReceivingMessage
	json.Unmarshal([]byte(jsonData), &receivedMessage)

	titleOfReq := receivedMessage.TitleOfReq
	reqOwner := receivedMessage.FromAcc
	message := receivedMessage.Message

	if(reqOwner == ""){
		return shim.Error("Owner of the request is required.")
	}


	if(titleOfReq == ""){
		return shim.Error("Title request is required.")
	}

	if(message == ""){
		return shim.Error("Message is required.")
	}

	//get caller username
	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	//get caller ID
	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	//get caller data
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
	var reqId = ""

	//get transaction timestamp
	txntmsp, err := APIstub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Error occured when trying to get the timestmap."+err.Error())
	}

	//if a financer --> get loan request
	if(account.UserType == "financer"){
		queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"loanreq\"}, \"fromacc\":{\"$eq\":\"%v\"},\"toacc\":{\"$eq\":\"%v\"}, \"title\":{\"$eq\":\"%v\"}},\"sort\":[{\"_id\":\"asc\"}]}", reqOwner, callerId, titleOfReq)

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

	//if an applicant --> owner --> get loan request
	}else{
		reqId = "request." + titleOfReq + "." + caller

		value, err := APIstub.GetState(reqId)
		if err != nil {
			return shim.Error("Something went wrong")
		}
		if value == nil {
			return shim.Error("Request does not exists")
		}

		json.Unmarshal(value, &request)

	}

	if(request.FromAcc != caller && request.Status == "Expired"){
		return shim.Error("You are not authorized to access these data")
	}

	if(request.FromAcc != caller && request.ToAcc != callerId){
		return shim.Error("You are not authorized to access these data")
	}

	msgID := "messages." + reqId

	msgValue, err := APIstub.GetState(msgID)
	if err != nil {
		return shim.Error("Something went wrong")
	}
	//if no messages has been sent yet add new record
	if msgValue == nil {
		var newMessage Message

		if(account.UserType == "financer"){
			newMessage.Financmsg = message
		}else{
			newMessage.Apptmsg = message
		}
		newMessage.Date = time.Unix(txntmsp.Seconds, int64(txntmsp.Nanos)).String()
		newMessage.Type = "message"

		// Save to WS state
		newMessageAsBytes, _ := json.Marshal(newMessage)
		err =APIstub.PutState(msgID, newMessageAsBytes)
		if err != nil {
			return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
		}
	//update record if messages has been sent already
	}else{

		var retrievedMsgs Message
		//Get messages from World State
		json.Unmarshal(msgValue, &retrievedMsgs)

		//if financer has sent a message and the applicant has not respond, update same record
		if(retrievedMsgs.Apptmsg == "" && retrievedMsgs.Financmsg != "" && account.UserType == "financer"){

			retrievedMsgs.Financmsg = retrievedMsgs.Financmsg + " | " + message
		
		//if applicant has sent a message and the financer has not respond, update same record
		}else if(retrievedMsgs.Apptmsg != "" && retrievedMsgs.Financmsg == "" && account.UserType == "applicant"){

			retrievedMsgs.Apptmsg = retrievedMsgs.Apptmsg + " | " + message

		}else{
			if(account.UserType == "financer"){
				retrievedMsgs.Financmsg = message
			}else{
				retrievedMsgs.Apptmsg = message
			}
		}
		retrievedMsgs.Date = time.Unix(txntmsp.Seconds, int64(txntmsp.Nanos)).String()
		retrievedMsgs.Type = "message"

		// Save to WS state
		retrievedMsgsAsBytes, _ := json.Marshal(retrievedMsgs)
		err =APIstub.PutState(msgID, retrievedMsgsAsBytes)
		if err != nil {
			return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
		}
	}

	return shim.Success(nil) 
}

// ===========================================================================================
// Get message by a Loan Request Id (args: Loan Request Id)
// ===========================================================================================
func (s *SmartContract) getMessagesByReqId(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 args")
	}

	titleOfReq := args[0]
	reqOwner := args[1]

	if(reqOwner == ""){
		return shim.Error("Owner of the request is required.")
	}
	
	var reqId = ""

	//get caller username
	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	//get caller ID
	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	//get caller data
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
		queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"loanreq\"}, \"fromacc\":{\"$eq\":\"%v\"},\"toacc\":{\"$eq\":\"%v\"}, \"title\":{\"$eq\":\"%v\"}},\"sort\":[{\"_id\":\"asc\"}]}", reqOwner, callerId, titleOfReq)

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
	}

	msgID := "messages." + reqId

	resultsIterator, err := APIstub.GetHistoryForKey(msgID)
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

