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
done by using **implicit private data collection**. More the vendor organizations onboard, explicit private data collection
will become maintenance hassle and hence preferred **implicit** way.

=> In unusual/contingency cases, last date for submitting bids can be changed by Government but only after all vendors 
agree for the same.

# Chaincode operations:
Government announces new project/contract work **[CreateContract]**

List active projects announced by Government **[ListAllContracts]**

List Submitted bids by a vendor **[QueryBidPrivate]**

Change bid date.

List projects for which bids submitted by a particular Vendor organization **[ListAllBids]**

Submit bid along with quote **[CreateBid]**

update bid with revised quote

# Asset model to be tracked in Blockchain ledger:
project id, project name, description, last date, state(active/closed)

private data collection (between Vendor and government) :
Vendor org id, project id, bid amount, bid date
