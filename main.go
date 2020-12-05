package main

import (
	"fmt"
	"./pkg/server"
	"./pkg/blockchain"
	"./pkg/blockchain/decoder"
)

func main(){
	fmt.Println("start")

	bc := blockchain.LoadChain()
	server.Init(bc)

	// block := blockchain.NewBlock(bc,[]string{"./main.go"})
	// fmt.Println(block)
	// index := blockchain.AddBlock(bc,block)
	// block = decoder.Read("./blocks/"+index+".block")
	// fmt.Println(block)
}

