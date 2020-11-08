package main

import (
	"fmt"
	"bytes"
	"crypto/sha256"
	"time"
	"strconv"
	"reflect"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nuonce		  int64
}

func (b *Block) SetHash() {
	var nuonce_i int64
	nuonce_i = 0
	for {
		timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
		nuonce := []byte(strconv.FormatInt(nuonce_i, 10))
		headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp, nuonce}, []byte{})
		hash := sha256.Sum256(headers)
	
		b.Hash = hash[:]
		if reflect.DeepEqual(b.Hash[0:3],[]byte("000")[0:3]){
			fmt.Println(b.Hash[0:3])
			fmt.Println([]byte("000")[0:3])

			b.Nuonce = nuonce_i
			break
		}
		nuonce_i++
	}
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{},0}
	block.SetHash()
	return block
}

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}

