package main

import (
	"fmt"
	"./pkg/blockchain"
	"./pkg/blockchain/encoder"
	"./pkg/blockchain/decoder"
)

func main(){
	fmt.Println("start")

	bc := blockchain.LoadChain()
	block := blockchain.NewBlock(bc,[]string{"./main.go"})

	fmt.Println(block)

	encoder.Write(block,"./test.block")

	block = decoder.Read("./test.block")
	fmt.Println(block)

}

