## Installation
```
git clone https://github.com/ryanshummmm/eth-parser.git
cd eth-parser
go mod tidy
```

## Configuration
Set the Ethereum node URL in internal/common/constants.go:
```
const CloudFlareRpcUrl = "https://cloudflare-eth.com"
```

## Usage

1. Start the server
`go run cmd/main.go`


2. API Endpoints:

- Get Current Block: GET /current-block
- Subscribe: POST /subscribe
- Unsubscribe: POST /unsubscribe
- Get Transactions: GET /transactions?address=0x...
- Get Subscribe List: GET /subscribe-list