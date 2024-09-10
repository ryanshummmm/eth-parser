# Ethereum Parser API
## Base URL: http://localhost:8080
## Endpoints

### Get Current Block

- GET /current-block
- Response: { "currentBlock": 12345 }


### Get Subscribe List

- GET /subscribe-list
- Response: { "subscribedAddresses": ["0x123...", "0x456..."] }


### Subscribe Address

- POST /subscribe
- Body: { "address": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e" }
- Response: { "subscribed": true }


### Unsubscribe Address

- POST /unsubscribe
- Body: { "address": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e" }
- Response: { "unsubscribed": true }


### Get Transactions

- GET /transactions?address=0x742d35Cc6634C0532925a3b844Bc454e4438f44e
- Response: [{ "from": "0x123...", "to": "0x456...", "value": "1000000000000000000" }, ...]



## Notes

1. Addresses are case-insensitive
2. Background task updates current block and processes new transactions for subscribed addresses