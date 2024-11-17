# Ethereum Transaction Parser

## Functionality

Exposes http API with routes:

- GET /current_block - gets last parsed block
- POST /subscribe/{address} - adds ethereum wallet address to monitor when parsing next blocks
- GET /transactions/{address} - gets list of parsed transactions for specified address

## Project structure

````
txparser/
├── api/  # api response structs
├── cmd/  # entrypoint
├── internal/
│   ├── client/  # client to interact with eth JSON-RPC
│   ├── config/  # config parsing
│   ├── domain/  # domain models and errors
│   ├── handlers/  # handlers layer, handler funcs that parse api requests and write responses 
│   ├── logger/  # logger setup
│   ├── repositories/  # data store layer
│   ├── router/  # router that defines api routes and matches them to handler funcs
│   ├── server/  # server that runs http api
│   └── services/  # business logic layer where core parser functionality is implemented
````

Each package is modular, any external dependencies defined as interfaces within each package. Special packages are `cmd`, `logger`, and `domain`. They are meant to either import or be imported in other packages.

Dependencies are initialised in `main.go` and passed into constructors of other dependencies.

Modular structure and interfaces allow to easily generate mocks and add automated tests for each layer when needed.

Logger package provides exportable interface that implements basic methods. Implementation of logger uses `slog`, configured to output json with several levels to stdout.

Error handling uses `errors.Join` to combine errors on different layers with domain errors. Later error can be analysed with `errors.Is` to determine correct status codes for api responses or logs.

## General flow

App builds dependencies in `main.go`. NotifyContext is used to listen to OS interrupt signals and shutdown application gracefully.

HTTP API is started as a server, parser subprocess is set to query chain each 10 seconds and parse every missing block.

- New addresses are saved in inmemory store
- Parser workers save transactions from block where the address is a sender or a receiver

## Running and manual testing

I found [random transaction](https://etherscan.io/tx/0xb3c5c2dbef5174c6a7e83cbc64b986239b5c93afe18533561c08143b11080bca) on Etherscan and set `start_block` in config to be previous of that transaction.

After starting the app, add From and To addresses to subscription list (the parser process will sleep for 10 secs in the beginning). *You don't have to use proposed addresses, pick any recent tx, add addresses and adjust start_block in config to be not too far from now.   

To start the app `make run`

To add addresses

```
curl -X POST "localhost:8080/subscription/0x1f9090aaE28b8a3dCeaDf281B0F12828e676c326
curl -X POST "localhost:8080/subscription/0x388C818CA8B9251b393131C08a736A67ccB19297
```

To get transactions

```
curl "localhost:8080/transactions/0x1f9090aaE28b8a3dCeaDf281B0F12828e676c326"
curl "localhost:8080/transactions/0x388C818CA8B9251b393131C08a736A67ccB19297"
```

To get last parsed block

```
curl "localhost:8080/current_block"
```


## Future improvements

- Implement rate limiter in the client to restrict number of requests to rpc
- Add pagination for get transactions
- Implement retries in the client to deal with network failures when interacting with rpc