package main

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
)

type Block struct {
    Index        int         `json:"index"`
    Timestamp    string      `json:"timestamp"`
    Transactions []Vote      `json:"transactions"`
    Hash         string      `json:"hash"`
    PrevHash     string      `json:"prevHash"`
    Nonce        int         `json:"nonce"`
}

type Vote struct {
    VoterID      string    `json:"voterId"`
    CandidateID  string    `json:"candidateId"`
    Timestamp    time.Time `json:"timestamp"`
    Signature    string    `json:"signature"`
}

type Blockchain struct {
    Blocks []Block `json:"blocks"`
}

var blockchain Blockchain
var pendingVotes []Vote

func createGenesisBlock() Block {
    return Block{
        Index:        0,
        Timestamp:    time.Now().String(),
        Transactions: []Vote{},
        Hash:         "0",
        PrevHash:     "0",
        Nonce:        0,
    }
}

func calculateHash(block Block) string {
    record := fmt.Sprintf("%d%s%v%s%d", 
        block.Index, 
        block.Timestamp, 
        block.Transactions, 
        block.PrevHash, 
        block.Nonce)

    h := sha256.New()
    h.Write([]byte(record))
    hashed := h.Sum(nil)
    return hex.EncodeToString(hashed)
}

func proofOfWork(block *Block) {
    for {
        hash := calculateHash(*block)
        if hash[:4] == "0000" {
            block.Hash = hash
            break
        }
        block.Nonce++
    }
}

func addBlock(votes []Vote) Block {
    prevBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
    newBlock := Block{
        Index:        prevBlock.Index + 1,
        Timestamp:    time.Now().String(),
        Transactions: votes,
        PrevHash:     prevBlock.Hash,
        Nonce:        0,
    }

    proofOfWork(&newBlock)
    blockchain.Blocks = append(blockchain.Blocks, newBlock)

    return newBlock
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(blockchain)
}

func addVote(w http.ResponseWriter, r *http.Request) {
    var vote Vote
    json.NewDecoder(r.Body).Decode(&vote)
    vote.Timestamp = time.Now()

    pendingVotes = append(pendingVotes, vote)

    if len(pendingVotes) >= 5 {
        newBlock := addBlock(pendingVotes)
        pendingVotes = []Vote{}

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(newBlock)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(vote)
}

func validateBlockchain(w http.ResponseWriter, r *http.Request) {
    isValid := true

    for i := 1; i < len(blockchain.Blocks); i++ {
        currentBlock := blockchain.Blocks[i]
        prevBlock := blockchain.Blocks[i-1]

        if currentBlock.Hash != calculateHash(currentBlock) {
            isValid = false
            break
        }

        if currentBlock.PrevHash != prevBlock.Hash {
            isValid = false
            break
        }
    }

    response := map[string]bool{"isValid": isValid}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func getStats(w http.ResponseWriter, r *http.Request) {
    totalVotes := 0
    for _, block := range blockchain.Blocks {
        totalVotes += len(block.Transactions)
    }

    stats := map[string]interface{}{
        "totalBlocks": len(blockchain.Blocks),
        "totalVotes":  totalVotes,
        "pendingVotes": len(pendingVotes),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

func main() {
    blockchain.Blocks = append(blockchain.Blocks, createGenesisBlock())

    router := mux.NewRouter()
    router.HandleFunc("/blockchain", getBlockchain).Methods("GET")
    router.HandleFunc("/vote", addVote).Methods("POST")
    router.HandleFunc("/validate", validateBlockchain).Methods("GET")
    router.HandleFunc("/stats", getStats).Methods("GET")

    router.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        })
    })

    fmt.Println("Blockchain node başlatıldı: http://localhost:8001")
    log.Fatal(http.ListenAndServe(":8001", router))
}
