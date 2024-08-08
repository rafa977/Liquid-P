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
	 //"time"
	 //"strconv"
 )
 

const (
    layoutISO = "2006-01-02 15:04:05"
)
 
 
// ===========================================================================================
// Register user data (args: user Id, json with all data)
// Company Admin and Super Admin can register a user in the blockchain
// ===========================================================================================
 func (s *SmartContract) registerAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 // check all args are provived
	 if len(args) != 2 {
		 return shim.Error("Incorrect number of arguments. Expecting 2 args")
	 }
	 accId := args[0]
	 jsonData := args[1]
	 
	 fmt.Println(jsonData)

	 errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}
 
	 // check that the uuid does not exist in the  WS ledger
	 checkExistsAsBytes, err := APIstub.GetState(accId)
 
	 // If the key does not exist in the state database, (nil, nil) is returned.
	 if checkExistsAsBytes != nil {
		 return shim.Error("The provided ID exists already in the WS ledger. ")
	 }
 
	 var account Account
	 json.Unmarshal([]byte(jsonData), &account)
	 
	 account.Username = caller
	 account.Type = "user"
	 account.IsDelete = "false"

	 // Save validity dates and operation to WS state
	 accountAsBytes, _ := json.Marshal(account)
	 err =APIstub.PutState(accId, accountAsBytes)
	 if err != nil {
		 return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
	 }
 
 
	 return shim.Success(nil) 
 }
 
// ===========================================================================================
// Register user data (args: user Id, json with all data)
// Company Admin and Super Admin can register a user in the blockchain
// ===========================================================================================
func (s *SmartContract) checkIdExists(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2 args")
	}
	accId := args[0]
	
	//Check that the user is Super Admin of the network
	isSuperAdmin, message := isSuperAdmin(APIstub);
	if(!isSuperAdmin){
		return shim.Error(message)
	}

	// check that the uuid does not exist in the  WS ledger
	checkExistsAsBytes, _ := APIstub.GetState(accId)

	// If the key does not exist in the state database, (nil, nil) is returned.
	if checkExistsAsBytes != nil {
		return shim.Error("The provided ID exists already in the WS ledger. ")
	}


	return shim.Success(nil) 
}


// ===========================================================================================
// Get information of an account (args: user Id)
// Any user can get the information of his account
// ===========================================================================================
func (s *SmartContract) getAccountById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting an ID")
	}

	accId := args[0]

	callerData, err := APIstub.GetState(accId)
	if err != nil {
		return shim.Error("Something went wrong")
	}
	if callerData == nil {
		return shim.Error("Something went wrong")
	}

	var account Account
	json.Unmarshal(callerData, &account)

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	if(caller != account.Username){
		//Check that the user is Super Admin of the network
		isSuperAdmin, _ := isSuperAdmin(APIstub);
		if(!isSuperAdmin){
			return shim.Error("You are not authorized to access these data")
		}
	}

	return shim.Success(callerData) 
}

// ===========================================================================================
// Update user data (args: json with all data)
// Only user calling can update his/her data in the blockchain --> User Id retrieved from APIstub
// ===========================================================================================
func (s *SmartContract) updateAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2 args")
	}

	jsonData := args[0]
	
	fmt.Println(jsonData)

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

	var updatedAccount Account
	json.Unmarshal([]byte(jsonData), &updatedAccount)

	var account Account
	json.Unmarshal([]byte(callerData), &account)

	account.Firstname = updatedAccount.Firstname
	account.Lastname = updatedAccount.Lastname
	account.Email = updatedAccount.Email
	account.ID = callerId
	account.Username = caller

	// Save validity dates and operation to WS state
	accountAsBytes, _ := json.Marshal(account)
	err = APIstub.PutState(callerId, accountAsBytes)
	if err != nil {
		return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
	}


	return shim.Success(nil) 
}

// ===========================================================================================
// Mark account as deleted (args: user Id)
// Only owner & Super Admin can do this action
// ===========================================================================================
func (s *SmartContract) markDeleteAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1 args")
	}

	accId := args[0]
	delete := args[1]

	accData, err := APIstub.GetState(accId)
	if err != nil {
		return shim.Error("Something went wrong")
	}
	if accData == nil {
		return shim.Error("Something went wrong")
	}

	var account Account
	json.Unmarshal(accData, &account)

	errorB, caller := getUserid(APIstub);
	if(errorB){
		return shim.Error(caller)
	}

	if(caller != account.Username){
		//Check that the user is Super Admin of the network
		isSuperAdmin, message := isSuperAdmin(APIstub);
		if(!isSuperAdmin){
			return shim.Error(message)
		}
	}

	if(delete == "true"){
		account.IsDelete = "true"
	}else if(delete == "false"){
		account.IsDelete = "false"
	}

	// Mark as Deleted
	accountAsBytes, _ := json.Marshal(account)
	err =APIstub.PutState(accId, accountAsBytes)
	if err != nil {
		return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
	}


	return shim.Success(nil) 
}

// ===========================================================================================
// Permanent deletion of an account (args: user Id)
// Only owner & Super Admin can do this action
// ===========================================================================================
func (s *SmartContract) permDeleteAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	// check all args are provived
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 args")
	}

	accId := args[0]

	accData, err := APIstub.GetState(accId)
	if err != nil {
		return shim.Error("Something went wrong")
	}
	if accData == nil {
		return shim.Error("Something went wrong")
	}
	
	var account Account
	json.Unmarshal(accData, &account)

	//Check that the user is Super Admin of the network
	isSuperAdmin, message := isSuperAdmin(APIstub);
	if(!isSuperAdmin){
		return shim.Error(message)
	}
	

	// Permanent delete from WS and mark as Deleted in Ledger
	err =APIstub.DelState(accId)
	if err != nil {
		return shim.Error("Error occured when trying to put account infromation in WS ledger: "+err.Error())
	}

	return shim.Success([]byte(account.Username)) 
}

// ===========================================================================================
// Get users history for a specific account (args: user Id)
// Admin
// ===========================================================================================
func (s *SmartContract) getUsersHistory(APIstub shim.ChaincodeStubInterface , args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting an ID")
	}

	accId := args[0]

	//Check that the user is Super Admin of the network
	isSuperAdmin, _ := isSuperAdmin(APIstub);
	if(!isSuperAdmin){
		return shim.Error("You are not authorized to access these data")
	}

	resultsIterator, err := APIstub.GetHistoryForKey(accId)
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


// ===========================================================================================
// Get All Financers that are registered in the blockchain (args: no args)
// ===========================================================================================
func (s *SmartContract) getAllFinancers(APIstub shim.ChaincodeStubInterface , args []string) sc.Response {

	//Check that the user is Super Admin of the network
	isSuperAdmin, _ := isSuperAdmin(APIstub);
	if(!isSuperAdmin){
		return shim.Error("You are not authorized to access these data")
	}

	queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"usertype\":{\"$eq\":\"financer\"}}}")

	resultsIterator, err := getQueryResultForQueryString(APIstub, queryString)
    if err != nil {
    	return shim.Error("Error occured: "+err.Error())
    }

	return shim.Success(resultsIterator)

 }


// ===========================================================================================
// Get All Financers that are registered in the blockchain (args: no args)
// ===========================================================================================
func (s *SmartContract) getAllUsers(APIstub shim.ChaincodeStubInterface , args []string) sc.Response {

	//Check that the user is Super Admin of the network
	isSuperAdmin, _ := isSuperAdmin(APIstub);
	if(!isSuperAdmin){
		return shim.Error("You are not authorized to access these data")
	}
	
	queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"user\"}}}")

	resultsIterator, err := getQueryResultForQueryString(APIstub, queryString)
    if err != nil {
    	return shim.Error("Error occured: "+err.Error())
    }

	return shim.Success(resultsIterator)

 }




// ===========================================================================================
// Get All Financers that are registered in the blockchain (args: no args)
// ===========================================================================================
func (s *SmartContract) getAllDocumentsById(APIstub shim.ChaincodeStubInterface , args []string) sc.Response {

	queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"doc\"}}}")

	resultsIterator, err := getQueryResultForQueryString(APIstub, queryString)
    if err != nil {
    	return shim.Error("Error occured: "+err.Error())
    }

	return shim.Success(resultsIterator)

 }




// Generates the response json with the provided results iterator
func createHistoryResultsJson(resultsIterator shim.HistoryQueryIteratorInterface, APIstub shim.ChaincodeStubInterface) ([]byte, error) {
 
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	var count int
	buffer.WriteString("[")

	isDeleted :=""
	txValue :=""

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
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
	return buffer.Bytes(), nil

}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// ===========================================================================================
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		fmt.Println(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return &buffer, nil
}
