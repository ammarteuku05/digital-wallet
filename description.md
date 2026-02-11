# Digital Wallet 
### Start Application
```bash
cd /Users/ammar/Desktop/digital-wallet
go run main.go
# or
make run
```

---

## List API

### 1. Check Balance
```bash
curl -X GET http://localhost:8080/v1/wallet/balance/user_id_here \
  -H "Content-Type: application/json"
```

### 2. Withdraw 
```bash
curl -X POST http://localhost:8080/v1/wallet/withdraw \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_id_here",
    "amount": 100.0,
    "description": "Test withdrawal"
  }'
```

### 3. View Transaction History
```bash
curl -X GET http://localhost:8080/v1/wallet/user_id_here/transactions?limit=10&offset=2' \
  --header 'Content-Type: application/json'
```