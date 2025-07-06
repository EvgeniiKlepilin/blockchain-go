package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Transaction struct {
	Source      string
	Destination string
	Amount      float64
}

type Block struct {
	Index     int
	Timestamp string
	Hash      string
	PrevHash  string
	Transaction
}

var Blockchain []Block

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		emptyTransaction := Transaction{"", "", 0}
		startingBlock := Block{0, t.String(), "", "", emptyTransaction}
		spew.Dump(startingBlock)
		Blockchain = append(Blockchain, startingBlock)
	}()

	log.Fatal(run())
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.Source + block.Destination + strconv.FormatFloat(block.Amount, 'f', -1, 64) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(prevBlock Block, transaction Transaction) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = prevBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Transaction = transaction
	newBlock.PrevHash = prevBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func validateBlock(currentBlock Block, prevBlock Block) bool {
	if currentBlock.Index != prevBlock.Index+1 {
		return false
	}
	if currentBlock.PrevHash != prevBlock.Hash {
		return false
	}
	if calculateHash(currentBlock) != currentBlock.Hash {
		return false
	}
	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var t Transaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	lastBlock := Blockchain[len(Blockchain)-1]
	newBlock, err := generateBlock(lastBlock, t)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, t)
		return
	}
	if validateBlock(newBlock, lastBlock) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
		spew.Dump(Blockchain)
		respondWithJSON(w, r, http.StatusCreated, newBlock)
	} else {
		respondWithJSON(w, r, http.StatusInternalServerError, newBlock)
	}
}

func run() error {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")

	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on port ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        muxRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
