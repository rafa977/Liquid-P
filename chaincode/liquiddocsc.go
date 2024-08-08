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
	 "github.com/golang/protobuf/ptypes"
	 "time"
	 //"strconv"
 )
 


// ===========================================================================================
// Add User Document (args: user Id, json with all data)
// ===========================================================================================
func (s *SmartContract) addDocument(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 args")
	}

	fileId := args[0]
	jsonData := args[1]
	
	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}
	
	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	fmt.Println(jsonData)

	var document Document
	json.Unmarshal([]byte(jsonData), &document)

	txntmsp, err := APIstub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Error occured when trying to get the timestmap."+err.Error())
	}

	document.ID = "file." + fileId
	document.Type = "document"
	document.AddedBy = caller
	document.Date = time.Unix(txntmsp.Seconds, int64(txntmsp.Nanos)).String()
	document.UserLastChange = caller

	isSuperAdmin, _ := isSuperAdmin(APIstub);
	if(!isSuperAdmin){
		document.Owner = callerId
	}

	// check that does not exist in the  WS ledger
	checkExistsAsBytes, err := APIstub.GetState(document.ID)

	// If the key does not exist in the state database, (nil, nil) is returned.
	if checkExistsAsBytes != nil {
		return shim.Error("The provided ID exists already in the WS ledger. ")
	}

	// Save to WS state
	accountAsBytes, _ := json.Marshal(document)
	err =APIstub.PutState(document.ID, accountAsBytes)
	if err != nil {
		return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
	}


	return shim.Success(nil) 
}

// ===========================================================================================
// Update Document (args: user Id, json with all data)
// ===========================================================================================
func (s *SmartContract) updateDocument(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 args")
	}

	fileId := args[0]
	jsonData := args[1]
	
	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}
	
	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	var document Document
	json.Unmarshal([]byte(jsonData), &document)

	fmt.Println(jsonData)

	//create the file ID
	fileID := "file." + fileId 
	
	value, err := APIstub.GetState(fileID)
	if err != nil {
		return shim.Error("Something went wrong")
	}
	if value == nil {
		return shim.Error("Document does not exists")
	}

	//Get Stroed Document from World State
	var storedDocument Document
	json.Unmarshal(value, &storedDocument)


	document.Owner = storedDocument.Owner
	document.AddedBy = storedDocument.AddedBy
	document.Date = storedDocument.Date
	document.Type = storedDocument.Type
	
	document.ID = fileID
	document.UserLastChange = caller

	//Check if user is the owner
	if(callerId != document.Owner){
		//Check that the user is Super Admin of the network
		isSuperAdmin, _ := isSuperAdmin(APIstub);
		if(!isSuperAdmin){
			return shim.Error("You are not authorized to perform this action")
		}
	}

	// Save to WS state
	accountAsBytes, _ := json.Marshal(document)
	err = APIstub.PutState(fileID, accountAsBytes)
	if err != nil {
		return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
	}


	return shim.Success(nil) 
}

// ===========================================================================================
// Update Document (args: user Id, json with all data)
// ===========================================================================================
func (s *SmartContract) deleteDocument(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting an ID")
	}

	docId := args[0]

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	fmt.Println(callerId)

	//create the file ID
	fileID := "file." + docId
	
	value, err := APIstub.GetState(fileID)
	if err != nil {
		return shim.Error("Something went wrong")
	}
	if value == nil {
		return shim.Error("Document does not exists")
	}

	//Get Document from World State
	var document Document
	json.Unmarshal(value, &document)

	if(document.Owner != callerId){
		//Check that the user is Super Admin of the network
		isSuperAdmin, _ := isSuperAdmin(APIstub);
		if(!isSuperAdmin){
			return shim.Error("You are not authorized to perform this action")
		}
	}
	
	// Delete from World State and mark deleted in Ledger
	err = APIstub.DelState(fileID)
	if err != nil {
		return shim.Error("Error occured when trying delete information: "+err.Error())
	}

	return shim.Success(nil) 
}

// ===========================================================================================
// Get information of an account (args: user Id)
// Company Admin and the user itself can get the information of the user asked
// ===========================================================================================
func (s *SmartContract) getDocumentById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting an ID")
	}

	docId := args[0]

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	fmt.Println(callerId)

	//create the file ID
	fileID := "file." + docId
	
	value, err := APIstub.GetState(fileID)
	if err != nil {
		return shim.Error("Document does not exists")
	}
	if value == nil {
		return shim.Error("Document does not exists")
	}

	//Get Document from World State
	var document Document
	json.Unmarshal(value, &document)

	if(document.Owner != callerId){
		
		accessID := "access." + fileID + "." + callerId
		fmt.Println(accessID)
		
		accessValue, err := APIstub.GetState(accessID)
		if err != nil {
			return shim.Error("You are not authorized to access these data")
		}
		if accessValue == nil {
			return shim.Error("You are not authorized to access these data")
		}

		//Get access control from World State
		var access AccessControl
		json.Unmarshal(accessValue, &access)

		if(access.Access == "false"){
			return shim.Error("You don not have access to this document")
		}

		validFrom, _ := time.Parse(layoutISO, access.ValidFrom)
		validTo, _ := time.Parse(layoutISO, access.ValidTo)

		fmt.Println(validFrom)                  // 1999-12-31 00:00:00 +0000 UTC
		fmt.Println(validTo) 
		
		tm := time.Now()

		g1 := tm.Before(validTo)
		g2 := tm.After(validFrom)

		if(!g1 || !g2) {
			return shim.Error("The access has been expired")
		}
		
	}
	   
	return shim.Success(value) 
}

// ===========================================================================================
// Get all documents by a user Id (args: user Id)
// A owner can get all documents & A not owner user gets all documents who only has access
// ===========================================================================================
func (s *SmartContract) getAllDocumentsByUserId(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting an ID")
	}

	userId := args[0]

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}

	fmt.Println(callerId)

	if(userId != callerId){
		queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"document\"},\"owner\":{\"$eq\":\"%v\"}}}", userId)

		resultsIterator, err := APIstub.GetQueryResult(queryString)
		if err != nil {
			return shim.Error("Error occured: "+err.Error())
		}

		//var documentsArray []Document
		
		var buffer bytes.Buffer
		buffer.WriteString("[")	
		bArrayMemberAlreadyWritten := false

		var i int
		for i = 0; resultsIterator.HasNext(); i++ {
			responseRange, err := resultsIterator.Next()
			if err != nil {
				return shim.Error(err.Error())
			}

			key_ws_rec := responseRange.Key
				
			docValue, _ := APIstub.GetState(key_ws_rec)
	
			//Store json data
			var document Document
			json.Unmarshal(docValue, &document)
			
			//create the file ID
			fileID := document.ID
	

			accessID := "access." + fileID + "." + callerId
			fmt.Println(accessID)
			
			accessValue, err := APIstub.GetState(accessID)
			if err != nil {
				continue
			}
			if accessValue == nil {
				continue
			}

			//Get access control from World State
			var access AccessControl
			json.Unmarshal(accessValue, &access)

			if(access.Access == "false"){
				continue
			}

			validFrom, _ := time.Parse(layoutISO, access.ValidFrom)
			validTo, _ := time.Parse(layoutISO, access.ValidTo)

			fmt.Println(validFrom)                  // 1999-12-31 00:00:00 +0000 UTC
			fmt.Println(validTo) 
			
			tm := time.Now()

			g1 := tm.Before(validTo)
			g2 := tm.After(validFrom)

			if(!g1 || !g2) {
				continue
			}
			
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("{\"Record\":")
			// Record is a JSON object, so we write as-is
			buffer.WriteString(string(docValue))

			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true

		}

		buffer.WriteString("]")

		var jj = &buffer;


		return shim.Success(jj.Bytes());

	}else{
		queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"document\"},\"owner\":{\"$eq\":\"%v\"}}}", callerId)

		resultsIterator, err := getQueryResultForQueryString(APIstub, queryString)
		if err != nil {
			return shim.Error("Error occured: "+err.Error())
		}

		return shim.Success(resultsIterator)
	}

	return shim.Error("Something went wrong")
}

// ===========================================================================================
// Get document history for a document from a specific account (args: file Id)
// Owner of the document 
// ===========================================================================================
func (s *SmartContract) getDocumentHistory(APIstub shim.ChaincodeStubInterface , args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting an ID")
	}

	docId := args[0]

	fileID := "file." + docId

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	errorID, callerId := getBlockchainId(APIstub, caller);
	if(!errorID){
		return shim.Error(callerId)
	}


	resultsIterator, err := APIstub.GetHistoryForKey(fileID)
	if err != nil {
		return shim.Error("You do not have the permissions to access this data")
	}
	defer resultsIterator.Close()


	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	var count int
	buffer.WriteString("[")

	isDeleted :=""
	txValue :=""
	checked := false

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Println(err.Error())
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")

		}

		timestamp, err := ptypes.Timestamp(queryResponse.Timestamp)
		isDeleted = "No"
		txValue = string(queryResponse.Value)
		if (queryResponse.IsDelete) {
			isDeleted = "Yes"
			txValue ="\"\""
		}

		if(checked == false){
			var historyDocument Document
			json.Unmarshal(queryResponse.Value, &historyDocument)
		
			fmt.Println(historyDocument)
	
			if(historyDocument.Owner != callerId){
				return shim.Error("You are not authorized to access this information")
			}	

			checked = true;
		}
		

		// Record is a JSON object, so we write as-is
		buffer.WriteString("{\"value\":")
		buffer.WriteString(txValue)
		buffer.WriteString(", \"TxId\":\"")
		buffer.WriteString(string(queryResponse.TxId))
		buffer.WriteString("\", \"Timestamp\":\"")
		buffer.WriteString(timestamp.String())
		buffer.WriteString("\", \"IsDelete\":\"")
		buffer.WriteString(isDeleted)
		buffer.WriteString("\"}")
		count++
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Println("createResultsJson():", buffer.String())
	fmt.Println("createResultsJson():", count)

	return shim.Success(buffer.Bytes())

}
