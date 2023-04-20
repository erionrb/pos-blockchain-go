package main

import (
	"github.com/erionrb/pos-blockchain-go/node"
)

var blockchainState node.State

func main() {
	blockchainState = node.InitNetwork()
}
