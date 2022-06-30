/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
        "encoding/json"
        "fmt"
        "time"
        "github.com/hyperledger/fabric-chaincode-go/shim"
        "github.com/hyperledger/fabric-contract-api-go/contractapi"
        "github.com/hyperledger/fabric/common/flogging"
)

// SmartContract of this fabric sample
type SmartContract struct {
        contractapi.Contract
}

// Contract work details
type ContractWork struct {
        ContractId string `json:"contractid"`
        Name       string `json:"name"`
        Brief      string `json:"brief"`
        LastDate   string `json:"lastdate"`
        Status     string `json:"status"`
}

// Bid details submitted by vendor
type VendorBid struct {
        VendorId   string `json:"vendorid"`
        ContractId string `json:"contractid"`
        BidAmount  int    `json:"bidamt"`
        BidDate    string `json:"biddate"`
}

const dateStr = "2006-01-02"
var logger = flogging.MustGetLogger("govtcontract")

func (s *SmartContract) CreateContract(ctx contractapi.TransactionContextInterface, contid string, cname string, brief string, ldate string, status string) error {

        if contid == "" {
                return fmt.Errorf("ContractId field must not be empty")
        }
        if cname == "" {
                return fmt.Errorf("Contract Name field must not be empty")
        }
        if brief == "" {
                return fmt.Errorf("Contract Brief field must not be empty")
        }
        if ldate == "" {
                return fmt.Errorf("Contract LastDate field must not be empty")
        }
        if status == "" {
                return fmt.Errorf("Contract Status field must not be empty")
        }

        // Get client org id and verify it matches peer org id.
        // In this scenario, client is only authorized to read/write private data from its own peer.
        clientOrgID, err := getClientOrgID(ctx, false)
        if err != nil {
                return fmt.Errorf("failed to get verified OrgID: errorstring %s, orgid: %s", err.Error(), clientOrgID)
        }

        // Only Org1MSP is allowed to create new contract
        if clientOrgID != "Org1MSP" {
                return fmt.Errorf("%s org is not allowed to create new contract", clientOrgID)
        }

        // Check if contract asset already exists
        contractAsBytes, err := ctx.GetStub().GetState(contid)
        if err != nil {
                return fmt.Errorf("failed to get asset: %v", err)
        } else if contractAsBytes != nil {
                fmt.Println("Contract already exists: " + contid)
                return fmt.Errorf("this contract already exists: " + contid)
        }

        cwork := ContractWork{
                ContractId: contid,
                Name: cname,
                Brief: brief,
                LastDate: ldate,
                Status: status,
        }

        contJSON, err := json.Marshal(&cwork)

        if err != nil {
                return fmt.Errorf("CreateContract: failed to marshal contract work : %v", err)
        }


        // Put contract details in the ledger
        err = ctx.GetStub().PutState(contid, contJSON)
        if err != nil {
                return fmt.Errorf("failed to put asset bid: %v", err)
        }

        return nil
}

// QueryContract returns the Contract details from world state
func (s *SmartContract) QueryContract(ctx contractapi.TransactionContextInterface, contractId string) (string, error) {

        // Get client org id and verify it matches peer org id.
        // In this scenario, client is only authorized to read/write private data from its own peer.
        clientOrgID, err := getClientOrgID(ctx, true)
        if err != nil {
                return "", fmt.Errorf("failed to get verified OrgID: errostr: %s, orgid: %s", err.Error(), clientOrgID)
        }

        contractDetails, err := ctx.GetStub().GetState(contractId)
        if err != nil {
                return "", fmt.Errorf("failed to read bid private properties from client org's collection: %s", err.Error())
        }
        if contractDetails == nil {
                return "", fmt.Errorf("bid private details does not exist in client org's collection: %s", contractId)
        }

        return string(contractDetails), nil
}

// ListAllContracts returns all Contract details from org's world state
func (s *SmartContract) ListAllContracts(ctx contractapi.TransactionContextInterface) ([]ContractWork, error) {

        contractIterator, err := ctx.GetStub().GetStateByRange("", "")
        if err != nil {
                return nil, fmt.Errorf("failed to read contract list from govt org's collection: %s", err.Error())
        }
        if contractIterator == nil {
                return nil, fmt.Errorf("contract details does not exist in govt org's collection")
        }

        defer contractIterator.Close()

        var allcontracts []ContractWork
        for contractIterator.HasNext() {
                entrycont, err := contractIterator.Next()
                if err != nil {
                        return nil, err
                }

                var contvar ContractWork
                err = json.Unmarshal(entrycont.Value, &contvar)
                if err != nil {
                        return nil, err
                }

                allcontracts = append(allcontracts, contvar)

        }
	
        return allcontracts, nil
}

// QueryBidPrivate returns the Bid details from vendor's private data collection
func (s *SmartContract) QueryBidPrivate(ctx contractapi.TransactionContextInterface, vendorId string, contractId string) (string, error) {

        // Get client org id and verify it matches peer org id.
        // In this scenario, client is only authorized to read/write private data from its own peer.
        clientOrgID, err := getClientOrgID(ctx, true)
        if err != nil {
                return "", fmt.Errorf("failed to get verified OrgID: %s", err.Error())
        }

        collection := "_implicit_org_" + clientOrgID

        bidconkey, err := ctx.GetStub().CreateCompositeKey(vendorId, []string{contractId})
	
        bidDetails, err := ctx.GetStub().GetPrivateData(collection, bidconkey)
        if err != nil {
                return "", fmt.Errorf("failed to read bid private properties from client org's collection: %s", err.Error())
        }
        if bidDetails == nil {
                logger.Infof("vendorId : %s", vendorId)
                logger.Infof("contractId : %s", []string{contractId})
                logger.Infof("collection : %s", collection)
                return "", fmt.Errorf("bid private details does not exist in client org's collection: %s", contractId)
        }

        return string(bidDetails), nil
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
        VendorId   string `json:"vendorid"`
        ContractId string `json:"contractid"`
        BidAmount  int    `json:"bidamt"`
        BidDate    string `json:"biddate"`
        }

        var vendorBidInput vendorTransientBid
        err = json.Unmarshal(bidJSON, &vendorBidInput)

        if err != nil {
                return fmt.Errorf("failed to unmarshal JSON: %v", err)
        }

        if vendorBidInput.VendorId == "" {
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
        clientOrgID, err := getClientOrgID(ctx, false)
        if err != nil {
               return fmt.Errorf("failed to get verified OrgID: %s", err.Error())
        }

        // Persist private immutable bid details to private data collection
        collection := "_implicit_org_" + clientOrgID


        // Check if bid already submitted by this vendor
        bidconkey, err := ctx.GetStub().CreateCompositeKey(vendorBidInput.VendorId, []string{vendorBidInput.ContractId})
        if err != nil {
                return fmt.Errorf("failed to create composite key: %v", err)
        }

        vendorBidAsBytes, err := ctx.GetStub().GetPrivateData(collection, bidconkey)
        if err != nil {
                return fmt.Errorf("failed to get asset: %v", err)
        } else if vendorBidAsBytes != nil {
                fmt.Println("Bid for this Contract already exists: " + vendorBidInput.ContractId)
		return fmt.Errorf("this bid already exists: " + vendorBidInput.ContractId)
        }

        // Org1MSP is Govt org.
        govcollection := "_implicit_org_" + "Org1MSP"

        contdata, err := ctx.GetStub().GetState(vendorBidInput.ContractId)
        if err != nil {
                return fmt.Errorf("failed to read contract details from govt org's collection: %v", err)
        }
        if contdata == nil {
                return fmt.Errorf("contract details does not exist in govt org's collection: %s", vendorBidInput.ContractId)
        }
        var contentry ContractWork
        err = json.Unmarshal(contdata, &contentry)
        if err != nil {
                return fmt.Errorf("Could not unmarshal: %v", err)
        }

        if contentry.Status != "open" {
                return fmt.Errorf("Contract is not open to submit bid")
        }

        biddate, err := time.Parse(dateStr, vendorBidInput.BidDate)
        if err != nil {
                return fmt.Errorf("failed to parse biddate: %v", err)
        }
        clastdate, err := time.Parse(dateStr, contentry.LastDate)
        if err != nil {
                return fmt.Errorf("failed to parse lastdate: %v", err)
        }
        if biddate.After(clastdate) {
                return fmt.Errorf("Last date for submitting bid is over")
        }

        // Put bid details in the org specifc private data collection
        err = ctx.GetStub().PutPrivateData(collection, bidconkey, []byte(bidJSON))
        if err != nil {
                return fmt.Errorf("failed to put asset bid: %v", err)
        }

        // Bid details of each vendor is saved in Govt ledger
        err = ctx.GetStub().PutPrivateData(govcollection, bidconkey, []byte(bidJSON))
        if err != nil {
                return fmt.Errorf("failed to put asset bid in Gov collection: %v", err)
        }

        return nil

}

// ListAllBids returns all Contract details from world state
func (s *SmartContract) ListAllBids(ctx contractapi.TransactionContextInterface) ([]VendorBid, error) {

        // Get client org id and verify it matches peer org id.
       // In this scenario, client is only authorized to read/write private data from its own peer.
        clientOrgID, err := getClientOrgID(ctx, true)
        if err != nil {
                return nil, fmt.Errorf("failed to get verified OrgID: %s", err.Error())
        }

        collection := "_implicit_org_" + clientOrgID
        var vendorList []string

        if clientOrgID == "Org1MSP" {
           vendorList = []string{"300", "400"}
        } else if clientOrgID == "Org2MSP" {
           vendorList = []string{"300"}
        } else if clientOrgID == "Org3MSP" {
           vendorList = []string{"400"}
        }

        myMSPID, err := ctx.GetClientIdentity().GetMSPID()
        logger.Infof("myMSPID: %s", myMSPID)
        var allbids []VendorBid

        for _, vendorId := range vendorList {

        BidIterator, err := ctx.GetStub().GetPrivateDataByPartialCompositeKey(collection, vendorId, []string{})
        defer BidIterator.Close()

        if err != nil {
                logger.Infof("ListAllBids error: %s", err.Error())
                return nil, fmt.Errorf("failed to read bid list: %s error: %s", err.Error())
        }
        if BidIterator == nil {
                logger.Infof("ListAllBids : null iterator for %s ", vendorId)
                return nil, fmt.Errorf("bid private details does not exist ")
        }

        logger.Infof("ListAllBids in govtcontract: no error")

        for BidIterator.HasNext() {
                logger.Infof("Iterator has element for %s: vendorId")

                entrybid, err := BidIterator.Next()
                if err != nil {
                        return nil, err
                }

                var bidvar VendorBid
                err = json.Unmarshal([]byte(entrybid.Value), &bidvar)
                if err != nil {
                      return nil, err
                }

                allbids = append(allbids, bidvar)
                logger.Infof("Iterator element: %s", string(entrybid.Value))

        }
        }
        logger.Infof("Iterator traversed")

        return allbids, nil
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
func getClientOrgID(ctx contractapi.TransactionContextInterface, verifyClientOrgMatchesPeerOrg bool) (string, error) {

        clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
        if err != nil {
                return "", fmt.Errorf("failed getting client's orgID: %s", err.Error())
        }
        fmt.Printf("invoking getClientOrgID with : %v", verifyClientOrgMatchesPeerOrg)

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
                fmt.Printf("Error creating govtcontract chaincode: %s", err.Error())
                return
        }

        if err := chaincode.Start(); err != nil {
                fmt.Printf("Error starting govtcontract chaincode: %s", err.Error())
        }
}
