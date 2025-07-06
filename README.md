# Blockchain Go

[![Go Version](https://img.shields.io/badge/Go-1.24.2-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/EvgeniiKlepilin/blockchain-go)](https://goreportcard.com/report/github.com/EvgeniiKlepilin/blockchain-go)
[![Coverage](https://img.shields.io/badge/coverage-61%25-brightgreen.svg)](https://github.com/EvgeniiKlepilin/blockchain-go)
[![Go](https://github.com/EvgeniiKlepilin/blockchain-go/actions/workflows/go.yml/badge.svg)](https://github.com/EvgeniiKlepilin/blockchain-go/actions/workflows/go.yml)

A simple blockchain implementation in Go that demonstrates the core concepts of blockchain technology. This project showcases a basic ledger system with transaction handling, block validation, and a REST API for blockchain operations.

## Features

- **Simple Blockchain Structure**: Basic blockchain with transaction support
- **SHA-256 Hashing**: Secure block hashing using SHA-256 algorithm
- **Block Validation**: Comprehensive validation of blocks and chain integrity
- **REST API**: HTTP endpoints for blockchain operations
- **Transaction Support**: Handle transactions between different parties
- **Chain Replacement**: Longest chain rule implementation
- **Comprehensive Testing**: Unit tests, benchmarks, and HTTP handler tests

## Data Structures

### Transaction
```go
type Transaction struct {
    Source      string  // Source address
    Destination string  // Destination address
    Amount      float64 // Transaction amount
}
```

### Block
```go
type Block struct {
    Index       int         // Block index in the chain
    Timestamp   string      // Block creation timestamp
    Hash        string      // Block hash
    PrevHash    string      // Previous block hash
    Transaction Transaction // Transaction data
}
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/`      | Get the entire blockchain |
| POST   | `/`      | Add a new block with transaction |

### Example API Usage

**Get Blockchain:**
```bash
curl http://localhost:8080/
```

**Add New Block:**
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "Source": "Alice",
    "Destination": "Bob",
    "Amount": 25.50
  }'
```

## Prerequisites

- Go 1.24.2 or higher
- Git

## Installation

1. **Clone the repository:**
```bash
git clone https://github.com/EvgeniiKlepilin/blockchain-go.git
cd blockchain-go
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Create environment file:**
```bash
echo "ADDR=8080" > .env
```

## Running the Project

### Start the Server
```bash
go run main.go
```

The server will start on port 8080 (or the port specified in your `.env` file).

### Alternative: Build and Run
```bash
# Build the binary
go build -o blockchain-go

# Run the binary
./blockchain-go
```

## Testing

### Run All Tests
```bash
go test -v
```

### Run Tests with Coverage
```bash
go test -v -cover
```

### Run Specific Test
```bash
go test -v -run TestCalculateHash
```

### Run Benchmarks
```bash
go test -bench=.
```

### Run Benchmarks with Memory Statistics
```bash
go test -bench=. -benchmem
```

## Core Functions

- **`calculateHash(block Block) string`**: Calculates SHA-256 hash for a block
- **`generateBlock(prevBlock Block, transaction Transaction) (Block, error)`**: Creates a new block
- **`validateBlock(currentBlock Block, prevBlock Block) bool`**: Validates block integrity
- **`replaceChain(newBlocks []Block)`**: Implements longest chain rule

## Dependencies

- **[gorilla/mux](https://github.com/gorilla/mux)**: HTTP router and URL matcher
- **[joho/godotenv](https://github.com/joho/godotenv)**: Environment variable loader
- **[davecgh/go-spew](https://github.com/davecgh/go-spew)**: Deep pretty printer for debugging

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Potential Future Enhancements

- [ ] Proof of Work consensus mechanism
- [ ] Persistent storage (database integration)
- [ ] Wallet functionality
- [ ] Transaction fees
- [ ] Network synchronization
- [ ] Block explorer web interface
- [ ] Transaction signatures and verification

## Author

**Evgenii Klepilin**

## Acknowledgments

- Thanks to the Go community for excellent tooling and libraries
- Inspired by blockchain technology and cryptocurrency concepts
