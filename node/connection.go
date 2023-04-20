package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type Connection struct {
	RealState *State
}

func newConnection(state *State) *Connection {
	return &Connection{
		RealState: state,
	}
}

func (c *Connection) handleConn(conn net.Conn) {
	defer conn.Close()

	go func() {
		for {
			msg := <-c.RealState.Announcements
			io.WriteString(conn, msg)
		}
	}()

	var address string

	stakeTokens(conn, c.RealState, address)
	setBPM(conn, c.RealState, address)
	broadcast(conn, c.RealState)
}

func stakeTokens(conn net.Conn, state *State, address string) {
	io.WriteString(conn, "Enter token balance:")
	scanBalance := bufio.NewScanner(conn)

	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}
		t := time.Now()
		address = calculateHash(t.String())
		state.Validators[address] = balance
		fmt.Println(state.Validators)
		break
	}
}

func setBPM(conn net.Conn, state *State, address string) {
	io.WriteString(conn, "\nEnter a new BPM:")

	scanBPM := bufio.NewScanner(conn)

	go func() {
		for {
			// take in BPM from stdin and add it to blockchain after conducting necessary validation
			for scanBPM.Scan() {
				bpm, err := strconv.Atoi(scanBPM.Text())
				// if malicious party tries to mutate the chain with a bad input, delete them as a validator and they lose their staked tokens
				if err != nil {
					log.Printf("%v not a number: %v", scanBPM.Text(), err)
					delete(state.Validators, address)
					conn.Close()
				}

				state.Mutex.Lock()
				oldLastIndex := state.Blockchain[len(state.Blockchain)-1]
				state.Mutex.Unlock()

				// create newBlock for consideration to be forged
				newBlock, err := generateBlock(oldLastIndex, bpm, address)
				if err != nil {
					log.Println(err)
					continue
				}
				if isBlockValid(newBlock, oldLastIndex) {
					state.CandidateBlocks <- newBlock
				}
				io.WriteString(conn, "\nEnter a new BPM:")
			}
		}
	}()
}

func broadcast(conn net.Conn, state *State) {
	// simulate receiving broadcast
	for {
		time.Sleep(time.Minute)
		state.Mutex.Lock()
		output, err := json.Marshal(state.Blockchain)
		state.Mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, string(output)+"\n")
	}
}
