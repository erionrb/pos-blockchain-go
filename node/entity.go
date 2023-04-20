package node

import (
	"sync"
)

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Validator string
}

type State struct {
	Blockchain      []Block
	TempBlocks      []Block
	Announcements   chan string
	CandidateBlocks chan Block
	Mutex           *sync.Mutex
	Validators      map[string]int
}
