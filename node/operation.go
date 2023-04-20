package node

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

func calculateHash(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

func calculateBlockHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	return calculateHash(record)
}

func generateBlock(oldBlock Block, BPM int, address string) (Block, error) {

	t := time.Now()

	newBlock := Block{
		Index:     oldBlock.Index + 1,
		Timestamp: t.String(),
		BPM:       BPM,
		PrevHash:  oldBlock.Hash,
	}
	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.Validator = address

	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock Block) bool {
	// Verify if the  the block being added to the blockchain is the next block in sequence
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	//Verify if the previous hash in the new block matches the hash of the old block
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	// Verify if the hash of the new block matches the generated hash */
	if calculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func InitNetwork() State {
	conn := newConnection(&State{
		Blockchain:      make([]Block, 0),
		TempBlocks:      make([]Block, 0),
		CandidateBlocks: make(chan Block),
		Announcements:   make(chan string),
		Mutex:           &sync.Mutex{},
		Validators:      make(map[string]int),
	})

	return *conn.RealState
}
