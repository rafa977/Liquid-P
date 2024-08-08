/*
 * This is the main file of chaincode for the project of Liquid 
 * Author: Rafail Kiloudis
 */
 
 package main

 import (
	 "fmt"
	 "github.com/hyperledger/fabric-chaincode-go/shim"
	 sc "github.com/hyperledger/fabric-protos-go/peer"
 )


  // Define the Smart Contract structure
type SmartContract struct {
}

type Type struct {
	Type 		string `json:"type"`
}

type Message struct {
	Type 		string `json:"type"`
	Apptmsg 	string `json:"apptmsg"`
	Financmsg  	string `json:"financmsg"`
	Date		string `json:"date"`
}

type ReceivingMessage struct {
	TitleOfReq 	string `json:"titleofreq"`
	FromAcc		string `json:"fromacc"`
	Message		string `json:"message"`
}

type Account struct {
	Type 		string `json:"type"`
	ID			string `json:"id"`
	Username	string `json:"username"`
	Firstname 	string `json:"firstname"`
	Lastname  	string `json:"lastname"`
	Email  		string `json:"email"`
	UserType 	string `json:"usertype"`
	Date		string `json:"date"`
	IsDelete	string `json:"isdelete"`
}

type Document struct {
	Type 	 		string `json:"type"`
	ID				string `json:"id"`
	Filename		string `json:"filename"`
	DocType			string `json:"doctype"`
	Signature		string `json:"signature"`
	Hash			string `json:"hash"`
	Owner			string `json:"owner"`	//user id
	AddedBy  		string `json:"addedby"` //username
	Date			string `json:"date"`  
	UserLastChange 	string `json:"userlastchange"` //username
}

type AccessControl struct {
	Type 	 	string `json:"type"`
	FileID		string `json:"fileid"`
	Filename	string `json:"filename"`
	Creator		string `json:"creator"`
	UserID		string `json:"userid"` //user id
	ValidFrom	string `json:"validfrom"`
	ValidTo		string `json:"validto"`
	Access		string `json:"access"`
	Status 		string `json:"status"`
}

type Company struct { 
	Type		string `json:"type"`
	Name 		string `json:"name"`
	Qty			string `json:"qty"`
	Available	string `json:"available"`
}

type Transaction struct {
	Type		string `json:"type"`
	FromAcc 	string `json:"fromacc"`
	ToAcc 		string `json:"toacc"`
	Company 	string `json:"company"`
	Qty 		string `json:"qty"`
	Operation	string `json:"operation"`
	TxId		string `json:"txid"`
	Date		string `json:"date"`
}

type LoanRequest struct {
	Type			string `json:"type"`
	Title			string `json:"title"`
	Description		string `json:"description"`	
	FromAcc 		string `json:"fromacc"`
	ToAcc			string `json:"toacc"`
	AcceptedAcc		string `json:"acceptedacc"`
	Amount	 		string `json:"amount"`
	Currency 		string `json:"currency"`
	TxId			string `json:"txid"`
	Status			string `json:"status"`
	ValidFrom 		string `json:"validfrom"`
	ValidTo 		string `json:"validto"`
	Date			string `json:"date"`
	Note			string `json:"note"`
	UserMsg			string `json:"usermsg"`
	DocumentLt		string `json:"documentlt"`
	UserLastChange 	string `json:"userlastchange"`
}

type Request struct {
	Type		string `json:"type"`
	FromAcc 	string `json:"fromacc"`
	AcceptedAcc	string `json:"acceptedacc"`
	Company 	string `json:"company"`
	Qty 		string `json:"qty"`
	Operation	string `json:"operation"`
	TxId		string `json:"txid"`
	Status		string `json:"status"`
	ValidFrom 	string `json:"validfrom"`
	ValidTo 	string `json:"validto"`
	Date		string `json:"date"`
}

type HistoryArray struct {
	Records		[]History `json:"records"`
}


type History struct {
	Value		string `json:"value"`
	TxId		string `json:"txid"`
	Timestamp	string `json:"timestamp"`
	IsDelete 	string `json:"isdelete"`
}


 func main() {
 
	 // Create a new Smart Contract
	 err := shim.Start(new(SmartContract))
	 if err != nil {
		 fmt.Printf("Error creating new Smart Contract: %s", err)
	 }
 }
 

/*
   * The Init method is called when the Smart Contract is instantiated by the blockchain network
   */
   func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "ts"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "registerAccount" {
		return s.registerAccount(APIstub, args)
	} else if function == "updateAccount" {
		return s.updateAccount(APIstub, args)
	} else if function == "getAccountById" {
	   return s.getAccountById(APIstub, args)
   } else if function == "permDeleteAccount" {
	   return s.permDeleteAccount(APIstub, args)
   } else if function == "markDeleteAccount" {
	   return s.markDeleteAccount(APIstub, args)
   } else if function == "addDocument" {
		return s.addDocument(APIstub, args)
	} else if function == "updateDocument" {
	   return s.updateDocument(APIstub, args)
   } else if function == "getDocumentById" {
	   return s.getDocumentById(APIstub, args)
   } else if function == "deleteDocument" {
	   return s.deleteDocument(APIstub, args)
   } else if function == "getDocumentHistory" {
	   return s.getDocumentHistory(APIstub, args)
   }else if function == "addAccessControl" {
	   return s.addAccessControl(APIstub, args)
   } else if function == "getAllFinancers" {
	   return s.getAllFinancers(APIstub, args)
   } else if function == "getAllAccessControlsByAccId"{
	   return s.getAllAccessControlsByAccId(APIstub, args)
   } else if function == "grantAccessAll"{ 
	   return s.grantAccessAll(APIstub, args)
   } else if function == "getLoanRequestById" { 
	   return s.getLoanRequestById(APIstub, args)
   } else if function == "addLoanRequest" {
	   return s.addLoanRequest(APIstub,args)
   } else if function == "updateLoanRequest" { 
	   return s.updateLoanRequest(APIstub, args)
   } else if function == "getRequestsByAccId" {
	   return s.getRequestsByAccId(APIstub, args)
   } else if function == "actionLoanRequest" {
	   return s.actionLoanRequest(APIstub, args)
   } else if function == "deleteLoanRequest" {
	   return s.deleteLoanRequest(APIstub, args)
   } else if function == "getRequestHistory" {
	   return s.getRequestHistory(APIstub, args)
   } else if function == "checkIdExists" {
	   return s.checkIdExists(APIstub, args)
   } else if function == "getAllDocumentsByUserId"{
	   return s.getAllDocumentsByUserId(APIstub, args)
   } else if function == "getUsersHistory" {
	   return s.getUsersHistory(APIstub, args)
   } else if function == "deleteAccessControl"{
	   return s.deleteAccessControl(APIstub, args)
   } else if function == "getMessagesByReqId" {
	   return s.getMessagesByReqId(APIstub, args)
   } else if function == "sendMessage" {
	   return s.sendMessage(APIstub, args)
   } else if function == "getReqHistoryByIdAuditor" {
	   return s.getReqHistoryByIdAuditor(APIstub, args)
   } else if function == "getAllKeys"{
	   return s.getAllKeys(APIstub, args)
   } else if function == "getAllApplicantIds" {
	   return s.getAllApplicantIds(APIstub, args)
   } else if function == "getAllKeysBasedOnApplicantId"{
	   return s.getAllKeysBasedOnApplicantId(APIstub, args)
   }

	return shim.Error("Invalid Smart Contract function name.")
}