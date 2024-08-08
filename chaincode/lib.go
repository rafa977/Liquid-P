/*
 * This is the lib file of chaincode for the project of Liquid 
 * Author: Rafail Kiloudis
 */
 
package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	//"encoding/json"
	//"strings"
	"fmt"
)

const (
	layout string = "2006/01/02"
)

func getBlockchainId(APIstub shim.ChaincodeStubInterface, caller string) (bool, string){
	queryString := fmt.Sprintf("{\"selector\":{\"_id\": {\"$gt\":null},\"type\":{\"$eq\":\"user\"},\"username\":{\"$eq\":\"%v\"}},\"sort\":[{\"_id\":\"asc\"}]}", caller)
	fmt.Println(queryString);

	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return false, "Something wrong has happened."
	}

	callerId := ""
	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return false, err.Error()
		}

		callerId =  responseRange.Key
		fmt.Println(callerId);
	}
	return true, callerId
}

func getMspID(APIstub shim.ChaincodeStubInterface) (bool, string) {

	// Get the client ID object
	id, err := cid.New(APIstub)
	if err != nil {
		return false, "Error occured when trying to initiate the CID library"
	}
	
	//Get the MSP ID of the caller and check if belongs to the org we want to give access
	mspid, err := id.GetMSPID()
	if err != nil {
		return false, "Error occured when trying to get the MSPID of the proposed transaction"
	}

	return true, mspid
}

func isSuperAdmin(APIstub shim.ChaincodeStubInterface) (bool, string) {

	// Get the client ID object
	id, err := cid.New(APIstub)
	if err != nil {
		return false, "Error occured when trying to initiate the CID library"
	}
	
	//Get the MSP ID of the caller and check if belongs to the org we want to give access
	mspid, err := id.GetMSPID()
	if err != nil {
		return false, "Error occured when trying to get the MSPID of the proposed transaction"
	} else if mspid != "LiquidMSP" {
		return false, "Only Liquid members are allowed to execute this transaction, you come from this Org MSPID: "+mspid
	}

	// check that the user has the superadmin app role
	err = id.AssertAttributeValue("role", "superadmin")
	if err != nil {
		return false, "You have to be a Super Admin in order to access this information: "+ err.Error()
	}

	return true, ""
}

func getUserid(APIstub shim.ChaincodeStubInterface) (bool, string){
	// Get the client ID object
	id, err := cid.New(APIstub)
	if err != nil {
		return true, "Error occured when trying to initiate the CID library: "
	}

	// check that the user is either a Doctor or a BDC Clerk based his role
	userId,err := id.GetX509Certificate()
	if err != nil {
		return true,"Error occured when trying to get the ID of the proposed transaction."
	}

	fmt.Println("this is the id of the user: " + userId.Subject.CommonName)

	return false, userId.Subject.CommonName
}