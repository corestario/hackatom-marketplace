## Dgaming Hackathon Marketplace hub

A cosmos zone with built-in simplistic NFT marketplace and IBC connectivity. The main functions available are setting NFT (i.e. transferred from the zone) on sale with a set price in avstract currency and buying NFT on sale. 

## Hub RPC endpoint (used by the webapp):

1. TransferNFTokenToZone(ConnectionID, TokenID) TransferID //Initiates a transfer of NF token through the open connection
2. GetTransferStatus(TransferID) Status //Get initiated transfer status (in flight, success, fail)
3. PutNFTokenOnTheMarket(TokenID, Price) Status 
4. BuyNFToken(TokenID) Status 
5. GetNFTokenData(TokenID) TokenData
6. GetNFTokensOnSaleList() []TokenData // Get list of all tokens on sale
7. MakeDeposit(Amount) Status // Currency faucet


## General scenario

User X makes a token in Zone A. User makes a "transfer and set price" IBC transactions to benefit from atomic fire-n-forget UI for selling NFTs.

User Y buys a token on the Marketplace, user X gets some currency on the marketplace.


## What is done

Most basic types and service functions are implemented/stubbed according to the specification, along with some commands.  

## TODOs

* Rewrite `./x/hh/client/rest/rest.go`, I haven't removed the `nameservice` code there yet.
* Add transaction commands to `./x/hh/client/cli/tx.go`, there is only one command implemented.
* All the Keeper / Handler logic. Search for `// TODO: ` comments throughout the project.  
