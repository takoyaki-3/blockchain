package main

import (
	"fmt"
	"bytes"
	"crypto/sha256"
	"time"
	"strconv"
//	"reflect"
	"math/rand"
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
	nuonce_i = int64(rand.Intn(100000000))

	keta := 6

	zeros := ""
	for i:=0;i<keta;i++{
		zeros += "0"
	}

	for {
		timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
		nuonce := []byte(strconv.FormatInt(nuonce_i, 10))
		headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp, nuonce}, []byte{})
		hash := sha256.Sum256(headers)
	
		hash_s := fmt.Sprintf("%x",hash)
		
		if hash_s[0:keta] == zeros[0:keta] {
			b.Hash = hash[:]
			b.Nuonce = nuonce_i
			fmt.Println(hash_s)
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
	rand.Seed(time.Now().UnixNano())

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

