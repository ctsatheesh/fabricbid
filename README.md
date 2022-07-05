# fabricbid
Bidding system for Government projects using Hyperledger Fabric Blockchain platform

# Problem statement: 
Governments worldwide invites bids from vendors for contract work on a continuous basis. 
In order to safeguard interest of the vendors, we need a system which provides a level playing field
for bids submitted by them. 

# Proposal:
Hyperledger Fabric provides a blockchain platform ideally suited for consortium kind of usecases
where multiple vendor organizations participate in the system along with Government departments.
Private data collections, a feature in Hyperledger Fabric comes handy here to secure bid data
from not sharing with competing vendor organizations.

# Life cycle of transactions:
=> Government departments announces new project/contract for inviting bids and announces last date for submitting the same.

=> Companies review them and submit their bids before last date.

=> Bids can be updated any number of times before the last date, giving them the flexibility based on ground work
done by vendors.

=> Bids submitted by individual vendors will not be shared with other vendors participating in bidding process. This is
done by using ***implicit private data collection***. More the vendor organizations onboard, explicit private data collection
will become maintenance hassle and hence preferred ***implicit*** way.

=> In unusual/contingency cases, last date for submitting bids can be changed by Government but only after all vendors 
agree for the same.

# Chaincode operations:
Government announces new project/contract work **[CreateContract]**

List active projects announced by Government **[ListAllContracts]**

List Submitted bids by a vendor **[QueryBidPrivate]**

List projects for which bids submitted by a particular Vendor organization **[ListAllBids]**
***To Do: Submitted bids should be viewed by Government only after last date***

Submit bid along with quote **[CreateBid]**

***To be Implemented***

update bid with revised quote

Change bid date

Finalize contract and award to successful Vendor

# Asset model to be tracked in Blockchain ledger:
project id, project name, description, last date, state(active/closed)

private data collection (between Vendor and government) :
Vendor org id, project id, bid amount, bid date

# Commands Example / Usage:
```
peer lifecycle chaincode package govtcontract.tar.gz --path . --lang golang --label govtcontract_1.0

peer lifecycle chaincode install govtcontract.tar.gz

peer lifecycle chaincode queryinstalled

export CC_PACKAGE_ID=

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --channelID <channel> --name govtcontract --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --signature-policy "AND('Org1MSP.peer','Org2MSP.peer','Org3MSP.peer')"
 
peer lifecycle chaincode checkcommitreadiness --channelID <channel> --name govtcontract --version 1.0 --sequence 1 --signature-policy "AND('Org1MSP.peer','Org2MSP.peer','Org3MSP.peer')" --output json

# commit

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --channelID <channel> --name govtcontract --version 1.0  --sequence 1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" --peerAddresses localhost:9151 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt"  --signature-policy "AND('Org1MSP.peer','Org2MSP.peer','Org3MSP.peer')"

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C <channel> -n govtcontract --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" --peerAddresses localhost:9151 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt" -c '{"function":"CreateContract","Args":["1033", "Project1", "Execute Project1", "2022-05-09","open"]}'

export BID_DETAILS=$(echo -n "{\"vendorid\":\"400\",\"contractid\":\"1033\",\"bidamt\":150,\"biddate\":\"2022-05-05\",\"salt\":\"1234ab1234\"}" | base64 | tr -d \\n)

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C <channel> -n govtcontract --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt  --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9151 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt -c '{"function":"CreateBid","Args":[]}' --transient "{\"bid_details\":\"$BID_DETAILS\"}"
 
peer chaincode query -C mych -n govtcontract -c '{"function":"QueryContract","Args":["1033"]}'
 
peer chaincode query -C mych -n govtcontract -c '{"function":"ListAllContracts","Args":[""]}'
 
peer chaincode query -C mych -n govtcontract -c '{"function":"QueryBidPrivate","Args":["400", "1033"]}'
 
peer chaincode query -C mych -n govtcontract -c '{"function":"ListAllBids","Args":[""]}'
```
