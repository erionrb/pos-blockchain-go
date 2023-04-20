package node

import (
	"crypto/sha256"
	"encoding/hex"
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
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}
