package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// Test data
var testTransaction = Transaction{
	Source:      "Alice",
	Destination: "Bob",
	Amount:      10.5,
}

var testBlock = Block{
	Index:     1,
	Timestamp: "2023-01-01 12:00:00",
	Hash:      "testhash",
	PrevHash:  "prevhash",
	Transaction: Transaction{
		Source:      "Alice",
		Destination: "Bob",
		Amount:      10.5,
	},
}

// Helper function to create a test block with calculated hash
func createTestBlock(index int, prevHash string, transaction Transaction) Block {
	t := time.Now()
	block := Block{
		Index:       index,
		Timestamp:   t.String(),
		PrevHash:    prevHash,
		Transaction: transaction,
	}
	block.Hash = calculateHash(block)
	return block
}

// Test calculateHash function
func TestCalculateHash(t *testing.T) {
	tests := []struct {
		name     string
		block    Block
		expected string
	}{
		{
			name: "Basic hash calculation",
			block: Block{
				Index:     1,
				Timestamp: "2023-01-01 12:00:00",
				PrevHash:  "prevhash",
				Transaction: Transaction{
					Source:      "Alice",
					Destination: "Bob",
					Amount:      10.5,
				},
			},
			expected: func() string {
				record := "1" + "2023-01-01 12:00:00" + "Alice" + "Bob" + "10.5" + "prevhash"
				h := sha256.New()
				h.Write([]byte(record))
				return hex.EncodeToString(h.Sum(nil))
			}(),
		},
		{
			name: "Empty transaction hash",
			block: Block{
				Index:     0,
				Timestamp: "2023-01-01 12:00:00",
				PrevHash:  "",
				Transaction: Transaction{
					Source:      "",
					Destination: "",
					Amount:      0,
				},
			},
			expected: func() string {
				record := "0" + "2023-01-01 12:00:00" + "" + "" + "0" + ""
				h := sha256.New()
				h.Write([]byte(record))
				return hex.EncodeToString(h.Sum(nil))
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateHash(tt.block)
			if result != tt.expected {
				t.Errorf("calculateHash() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test generateBlock function
func TestGenerateBlock(t *testing.T) {
	prevBlock := Block{
		Index:     0,
		Timestamp: "2023-01-01 11:00:00",
		Hash:      "prevhash",
		PrevHash:  "",
		Transaction: Transaction{
			Source:      "",
			Destination: "",
			Amount:      0,
		},
	}

	transaction := Transaction{
		Source:      "Alice",
		Destination: "Bob",
		Amount:      25.0,
	}

	newBlock, err := generateBlock(prevBlock, transaction)
	if err != nil {
		t.Errorf("generateBlock() returned error: %v", err)
	}

	// Verify block properties
	if newBlock.Index != prevBlock.Index+1 {
		t.Errorf("New block index = %d, expected %d", newBlock.Index, prevBlock.Index+1)
	}

	if newBlock.PrevHash != prevBlock.Hash {
		t.Errorf("New block PrevHash = %s, expected %s", newBlock.PrevHash, prevBlock.Hash)
	}

	if newBlock.Transaction != transaction {
		t.Errorf("New block transaction = %+v, expected %+v", newBlock.Transaction, transaction)
	}

	// Verify hash is calculated correctly
	expectedHash := calculateHash(newBlock)
	if newBlock.Hash != expectedHash {
		t.Errorf("New block hash = %s, expected %s", newBlock.Hash, expectedHash)
	}

	// Verify timestamp is set
	if newBlock.Timestamp == "" {
		t.Error("New block timestamp is empty")
	}
}

// Test validateBlock function
func TestValidateBlock(t *testing.T) {
	validPrevBlock := Block{
		Index:       0,
		Timestamp:   "2023-01-01 11:00:00",
		Hash:        "validhash",
		PrevHash:    "",
		Transaction: Transaction{},
	}

	tests := []struct {
		name         string
		currentBlock Block
		prevBlock    Block
		expected     bool
	}{
		{
			name: "Valid block",
			currentBlock: func() Block {
				block := Block{
					Index:     1,
					Timestamp: "2023-01-01 12:00:00",
					PrevHash:  "validhash",
					Transaction: Transaction{
						Source:      "Alice",
						Destination: "Bob",
						Amount:      10.0,
					},
				}
				block.Hash = calculateHash(block)
				return block
			}(),
			prevBlock: validPrevBlock,
			expected:  true,
		},
		{
			name: "Invalid index",
			currentBlock: func() Block {
				block := Block{
					Index:     2, // Wrong index
					Timestamp: "2023-01-01 12:00:00",
					PrevHash:  "validhash",
					Transaction: Transaction{
						Source:      "Alice",
						Destination: "Bob",
						Amount:      10.0,
					},
				}
				block.Hash = calculateHash(block)
				return block
			}(),
			prevBlock: validPrevBlock,
			expected:  false,
		},
		{
			name: "Invalid previous hash",
			currentBlock: func() Block {
				block := Block{
					Index:     1,
					Timestamp: "2023-01-01 12:00:00",
					PrevHash:  "wronghash", // Wrong previous hash
					Transaction: Transaction{
						Source:      "Alice",
						Destination: "Bob",
						Amount:      10.0,
					},
				}
				block.Hash = calculateHash(block)
				return block
			}(),
			prevBlock: validPrevBlock,
			expected:  false,
		},
		{
			name: "Invalid hash",
			currentBlock: Block{
				Index:     1,
				Timestamp: "2023-01-01 12:00:00",
				Hash:      "wronghash", // Wrong hash
				PrevHash:  "validhash",
				Transaction: Transaction{
					Source:      "Alice",
					Destination: "Bob",
					Amount:      10.0,
				},
			},
			prevBlock: validPrevBlock,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateBlock(tt.currentBlock, tt.prevBlock)
			if result != tt.expected {
				t.Errorf("validateBlock() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test replaceChain function
func TestReplaceChain(t *testing.T) {
	// Setup original blockchain
	originalBlockchain := []Block{
		{Index: 0, Hash: "genesis"},
		{Index: 1, Hash: "block1"},
	}
	Blockchain = originalBlockchain

	tests := []struct {
		name          string
		newBlocks     []Block
		expectedLen   int
		shouldReplace bool
	}{
		{
			name: "Replace with longer chain",
			newBlocks: []Block{
				{Index: 0, Hash: "genesis"},
				{Index: 1, Hash: "block1"},
				{Index: 2, Hash: "block2"},
				{Index: 3, Hash: "block3"},
			},
			expectedLen:   4,
			shouldReplace: true,
		},
		{
			name: "Don't replace with shorter chain",
			newBlocks: []Block{
				{Index: 0, Hash: "genesis"},
			},
			expectedLen:   4, // Should remain the same as previous test
			shouldReplace: false,
		},
		{
			name: "Don't replace with equal length chain",
			newBlocks: []Block{
				{Index: 0, Hash: "genesis"},
				{Index: 1, Hash: "block1"},
				{Index: 2, Hash: "block2"},
				{Index: 3, Hash: "block3"},
			},
			expectedLen:   4,
			shouldReplace: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replaceChain(tt.newBlocks)
			if len(Blockchain) != tt.expectedLen {
				t.Errorf("Blockchain length = %d, expected %d", len(Blockchain), tt.expectedLen)
			}
		})
	}
}

// Test HTTP handlers
func TestHandleGetBlockchain(t *testing.T) {
	// Setup test blockchain
	Blockchain = []Block{
		{Index: 0, Hash: "genesis", Timestamp: "2023-01-01"},
		{Index: 1, Hash: "block1", Timestamp: "2023-01-02"},
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleGetBlockchain)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response []Block
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 blocks in response, got %d", len(response))
	}
}

func TestHandleWriteBlock(t *testing.T) {
	// Setup test blockchain with genesis block
	Blockchain = []Block{
		{
			Index:     0,
			Hash:      "genesis",
			Timestamp: "2023-01-01",
			PrevHash:  "",
			Transaction: Transaction{
				Source:      "",
				Destination: "",
				Amount:      0,
			},
		},
	}

	transaction := Transaction{
		Source:      "Alice",
		Destination: "Bob",
		Amount:      15.5,
	}

	jsonData, err := json.Marshal(transaction)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleWriteBlock)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response Block
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response: %v", err)
	}

	if response.Index != 1 {
		t.Errorf("Expected block index 1, got %d", response.Index)
	}

	if response.Transaction.Source != "Alice" {
		t.Errorf("Expected source 'Alice', got %s", response.Transaction.Source)
	}

	if len(Blockchain) != 2 {
		t.Errorf("Expected blockchain length 2, got %d", len(Blockchain))
	}
}

func TestHandleWriteBlockInvalidJSON(t *testing.T) {
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleWriteBlock)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// Benchmark tests
func BenchmarkCalculateHash(b *testing.B) {
	block := Block{
		Index:     1,
		Timestamp: "2023-01-01 12:00:00",
		PrevHash:  "prevhash",
		Transaction: Transaction{
			Source:      "Alice",
			Destination: "Bob",
			Amount:      10.5,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculateHash(block)
	}
}

func BenchmarkGenerateBlock(b *testing.B) {
	prevBlock := Block{
		Index:       0,
		Timestamp:   "2023-01-01 11:00:00",
		Hash:        "prevhash",
		PrevHash:    "",
		Transaction: Transaction{},
	}

	transaction := Transaction{
		Source:      "Alice",
		Destination: "Bob",
		Amount:      25.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateBlock(prevBlock, transaction)
	}
}

func BenchmarkValidateBlock(b *testing.B) {
	prevBlock := Block{
		Index:       0,
		Timestamp:   "2023-01-01 11:00:00",
		Hash:        "validhash",
		PrevHash:    "",
		Transaction: Transaction{},
	}

	currentBlock := Block{
		Index:     1,
		Timestamp: "2023-01-01 12:00:00",
		PrevHash:  "validhash",
		Transaction: Transaction{
			Source:      "Alice",
			Destination: "Bob",
			Amount:      10.0,
		},
	}
	currentBlock.Hash = calculateHash(currentBlock)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateBlock(currentBlock, prevBlock)
	}
}

func BenchmarkReplaceChain(b *testing.B) {
	originalChain := make([]Block, 100)
	for i := 0; i < 100; i++ {
		originalChain[i] = Block{Index: i, Hash: "hash" + strconv.Itoa(i)}
	}

	newChain := make([]Block, 150)
	for i := 0; i < 150; i++ {
		newChain[i] = Block{Index: i, Hash: "newhash" + strconv.Itoa(i)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Blockchain = originalChain
		replaceChain(newChain)
	}
}

func BenchmarkHandleGetBlockchain(b *testing.B) {
	// Setup test blockchain
	Blockchain = make([]Block, 100)
	for i := 0; i < 100; i++ {
		Blockchain[i] = Block{
			Index:     i,
			Hash:      "hash" + strconv.Itoa(i),
			Timestamp: "2023-01-01",
		}
	}

	req, _ := http.NewRequest("GET", "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleGetBlockchain)
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkHandleWriteBlock(b *testing.B) {
	transaction := Transaction{
		Source:      "Alice",
		Destination: "Bob",
		Amount:      15.5,
	}

	jsonData, _ := json.Marshal(transaction)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Setup fresh blockchain for each iteration
		Blockchain = []Block{
			{
				Index:       0,
				Hash:        "genesis",
				Timestamp:   "2023-01-01",
				PrevHash:    "",
				Transaction: Transaction{},
			},
		}

		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleWriteBlock)
		handler.ServeHTTP(rr, req)
	}
}
