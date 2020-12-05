package main

import (
	"fmt"
	"./pkg/block"
	"./pkg/block/encoder"
	"./pkg/block/decoder"
)

func main(){
	fmt.Println("start")

	bytes := []byte("hello")
	fp := block.Properties{}
	files := []block.File{block.File{fp,bytes}}
	properties := block.Properties{}
	properties["hello"]="world"

	block := block.Block{properties,uint64(1),files}

	fmt.Println(block)

	encoder.Write(block,"./test.block")

	block = decoder.Read("./test.block")
	fmt.Println(block)

}

