/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main 

import (
        "encoding/json"
        "fmt"
        "github.com/hyperledger/fabric-chaincode-go/shim"
        "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract of this fabric sample
type SmartContract struct {
        contractapi.Contract
}

type ContractWork struct {
	ContractId string `json:"contractid"` 
	Name       string `json:"name"`    
	Brief      string `json:"brief"`
	LastDate   string `json:"lastdate"`
	Status     string `json:"status"`
}

type VendorBid struct {
	VendorId   int	  `json:"vendorid"`
	ContractId string `json:"contractid"`
	BidAmount  int    `json:"bidamt"`
	BidDate    string `json:"biddate"`
}

func (s *SmartContract) CreateContract(ctx contractapi.TransactionContextInterface) error {

	// Get new contract work data from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient: %v", err)
	}
	
	// Bid details are private, therefore they get passed in transient field, instead of func args
	contractJSON, ok := transientMap["contract_details"]
	if !ok {
		//log error to stdout
		return fmt.Errorf("Contract work not found in the transient map input")
	}
	
	type contractTransientBid struct {
		ContractId string `json:"contractid"` 
		Name       string `json:"name"`    
		Brief      string `json:"brief"`
		LastDate   string `json:"lastdate"`
		Status     string `json:"status"`
	}
	
	var contractInput contractTransientBid
	err = json.Unmarshal(contractJSON, &contractInput)
	
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if contractInput.ContractId == "" {
		return fmt.Errorf("ContractId field must not be empty")
	}
	if contractInput.Name == "" {
		return fmt.Errorf("Contract Name field must not be empty")
	}
	if contractInput.Brief == "" {
		return fmt.Errorf("Contract Brief field must not be empty")
	}
	if contractInput.LastDate == "" {
		return fmt.Errorf("Contract LastDate field must not be empty")
	}
	if contractInput.Status == "" {
		return fmt.Errorf("Contract Status field must not be empty")
	}
	
	// Get client org id and verify it matches peer org id.
	// In this scenario, client is only authorized to read/write private data from its own peer.
	clientOrgID, err := getClientOrgID(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to get verified OrgID: %s", err.Error())
	}
	
	// Persist private immutable marble properties to owner's private data collection
	collection := "_implicit_org_" + clientOrgID
	
	
	// Check if asset already exists
	contractAsBytes, err := ctx.GetStub().GetPrivateData(collection, contractInput.ContractId)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	} else if contractAsBytes != nil {
		fmt.Println("Contract already exists: " + contractInput.ContractId)
		return fmt.Errorf("this contract already exists: " + contractInput.ContractId)
	}
	
	// Put agreed value in the org specifc private data collection
	err = ctx.GetStub().PutPrivateData(collection, contractInput.ContractId, []byte(contractJSON))
	if err != nil {
		return fmt.Errorf("failed to put asset bid: %v", err)
	}

        return nil	
}

// QueryContractPrivate returns the Contract details from owner's private data collection
func (s *SmartContract) QueryContractPrivate(ctx contractapi.TransactionContextInterface, contractId string) (string, error) {

	// Get client org id and verify it matches peer org id.
	// In this scenario, client is only authorized to read/write private data from its own peer.
	clientOrgID, err := getClientOrgID(ctx, true)
	if err != nil {
		return "", fmt.Errorf("failed to get verified OrgID: %s", err.Error())
	}

	collection := "_implicit_org_" + clientOrgID

	contractDetails, err := ctx.GetStub().GetPrivateData(collection, contractId)
	if err != nil {
		return "", fmt.Errorf("failed to read bid private properties from client org's collection: %s", err.Error())
	}
	if contractDetails == nil {
		return "", fmt.Errorf("bid private details does not exist in client org's collection: %s", contractId)
	}

	return string(contractDetails), nil
}

func (s *SmartContract) CreateBid(ctx contractapi.TransactionContextInterface) error {

	// Get new bid data from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient: %v", err)
	}
	
	// Bid details are private, therefore they get passed in transient field, instead of func args
	bidJSON, ok := transientMap["bid_details"]
	if !ok {
		//log error to stdout
		return fmt.Errorf("Bid not found in the transient map input")
	}
	
	type vendorTransientBid struct {
	VendorId   int	  `json:"vendorid"`
	ContractId string `json:"contractid"`
	BidAmount  int    `json:"bidamt"`
	BidDate    string `json:"biddate"`
	}
	
	var vendorBidInput vendorTransientBid
	err = json.Unmarshal(bidJSON, &vendorBidInput)
	
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if vendorBidInput.VendorId == 0 {
		return fmt.Errorf("VendorId field must be a non-zero number")
	}
	if vendorBidInput.ContractId == "" {
		return fmt.Errorf("ContractId field must be a non-zero number")
	}
	if vendorBidInput.BidAmount == 0 {
		return fmt.Errorf("BidAmount field must be a non-zero number")
	}
	if vendorBidInput.BidDate == "" {
		return fmt.Errorf("BidDate field must be a non-zero value")
	}
	
	// Get client org id and verify it matches peer org id.
	// In this scenario, client is only authorized to read/write private data from its own peer.
	clientOrgID, err := getClientOrgID(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to get verified OrgID: %s", err.Error())
	}
	
	// Persist private immutable marble properties to owner's private data collection
	collection := "_implicit_org_" + clientOrgID
	
	
	// Check if asset already exists
	vendorBidAsBytes, err := ctx.GetStub().GetPrivateData(collection, vendorBidInput.ContractId)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	} else if vendorBidAsBytes != nil {
		fmt.Println("Bid for this Contract already exists: " + vendorBidInput.ContractId)
		return fmt.Errorf("this bid already exists: " + vendorBidInput.ContractId)
	}
	
	// Put agreed value in the org specifc private data collection
	err = ctx.GetStub().PutPrivateData(collection, vendorBidInput.ContractId, []byte(bidJSON))
	if err != nil {
		return fmt.Errorf("failed to put asset bid: %v", err)
	}
      
        return nil
	
}
	
// QueryBidPrivate returns the Bid details from owner's private data collection
func (s *SmartContract) QueryBidPrivate(ctx contractapi.TransactionContextInterface, contractId string) (string, error) {

	// Get client org id and verify it matches peer org id.
	// In this scenario, client is only authorized to read/write private data from its own peer.
	clientOrgID, err := getClientOrgID(ctx, true)
	if err != nil {
		return "", fmt.Errorf("failed to get verified OrgID: %s", err.Error())
	}

	collection := "_implicit_org_" + clientOrgID

	bidDetails, err := ctx.GetStub().GetPrivateData(collection, contractId)
	if err != nil {
		return "", fmt.Errorf("failed to read bid private properties from client org's collection: %s", err.Error())
	}
	if bidDetails == nil {
		return "", fmt.Errorf("bid private details does not exist in client org's collection: %s", contractId)
	}

	return string(bidDetails), nil
}


// verifyClientOrgMatchesPeerOrg is an internal function used verify client org id and matches peer org id.
func verifyClientOrgMatchesPeerOrg(ctx contractapi.TransactionContextInterface) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	peerMSPID, err := shim.GetMSPID()
	if err != nil {
		return fmt.Errorf("failed getting the peer's MSPID: %v", err)
	}

	if clientMSPID != peerMSPID {
		return fmt.Errorf("client from org %v is not authorized to read or write private data from an org %v peer", clientMSPID, peerMSPID)
	}

	return nil
}

	
// getClientOrgID gets the client org ID.
// The client org ID can optionally be verified against the peer org ID, to ensure that a client from another org doesn't attempt to read or write private data from this peer.
// The only exception in this scenario is for TransferMarble, since the current owner needs to get an endorsement from the buyer's peer.
func getClientOrgID(ctx contractapi.TransactionContextInterface, verifyClientOrgMatchesPeerOrg bool) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed getting client's orgID: %s", err.Error())
	}

	if verifyClientOrgMatchesPeerOrg {
		peerOrgID, err := shim.GetMSPID()
		if err != nil {
			return "", fmt.Errorf("failed getting peer's orgID: %s", err.Error())
		}

		if clientOrgID != peerOrgID {
			return "", fmt.Errorf("client from org %s is not authorized to read or write private data from an org %s peer", clientOrgID, peerOrgID)
		}
	}

	return clientOrgID, nil
}

func main() {

        chaincode, err := contractapi.NewChaincode(new(SmartContract))

        if err != nil {
                fmt.Printf("Error creating bidcontract chaincode: %s", err.Error())
                return
        }

        if err := chaincode.Start(); err != nil {
                fmt.Printf("Error starting bidcontract chaincode: %s", err.Error())
        }
}
