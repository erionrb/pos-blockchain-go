package main

import (
	"github.com/erionrb/pos-blockchain-go/node"
)

var Blockchain []node.Block
var tempBlocks []node.Block

var candidateBlocks = make(chan node.Block)
var announcements = make(chan string)

var mutex = &sync.Mutex{}

var validators = make(map[string]int)

