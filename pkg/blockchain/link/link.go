package link

import (
	"fmt"
	"encoding/json"
	// "math/big"
	"github.com/blockcypher/gobcy"
)

type LinkBlock struct {
	Created string `json:"created"`
	LinkedChain string `json:"linked_chain"`
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Id string `json:"id"`
	Hex string `json:"hex"`
}

func JsonEncode(linkBlock LinkBlock)string{
	rowJson,_ := json.Marshal(linkBlock)
	return string(rowJson)
}

func Write2PublicBlockchain(str string)string{

	bc := gobcy.API{}
	bc.Token = "f198074f2de843d3b4cd0b496dd1c057"
	bc.Coin = "ltc" //options: "btc","bcy","ltc","doge","eth"
	bc.Chain = "main" //depending on coin: "main","test3","test"

	// var keys1, keys2 gobcy.AddrKeychain
	// keys1,_ := bc.GenAddrKeychain()
	// keys2,_ := bc.GenAddrKeychain()
	// keys1 := gobcy.AddrKeychain{}
	// keys1.Address = "LfsLoRK81a8rDPJzSA2rVJo61GiP7WqbyV"
	// keys1.Private = "a588d171e537041175c1b80d66aa9fdcf901bc73cd2b3168d2cd7e55d5044668"
	// keys1.Public = "03e858d72ecc23562520b92e6826a96c95ea1aaead7d0d46404f56d4979ee9ed4f"
	// keys1.Wif = "T8bkiMADzryVUxcCPwvnJvcn87M9czXT3EuYVetGuvbd7BCVgUWp"
	// keys2 := gobcy.AddrKeychain{}
	// keys2.Address = "LLm89UqPZWDgBSQATLq1qt6bz9UjrQbVFH"
	// keys2.Private = "7205edabe1edca9d114ee40fcbe987c5da6325fa666d3598de5d9586ddd617f9"
	// keys2.Public = "02d3b54c22e9ffc54e6ce2bbbf2d6c35bc74fa18fa9a472979f9fb2adc626d26a9"
	// keys2.Wif = "T6sd6wyWsnzh1ZmbNRHCEvxwqEjyFhPGXigA2QTdkiJsNnXWa7eD"


	keys1 := gobcy.AddrKeychain{}
	keys1.Address = "LWUpzu2yrY9jz7L657QUkQf1zZ8QKiBQfW"
	keys1.Private = "ee53fcf91c600ca2ee6052faab7482f7e09f1141e357d6660ec20fa3b92f1ddd"
	keys1.Public = "03ac2643003a844bc286898b3c8bd92fb4c1dc414fb339e66494af75c9c46373ff"
	keys1.Wif = "TB3FnknoBL6RUx2HKcvHfhcG46mHTXHxfHvehW8ypvSrZt6qUzfR"
	keys2 := gobcy.AddrKeychain{}
	keys2.Address = "LU5A7xJ987JSB7KwyAnoxSPWcgsdWp3Gay"
	keys2.Private = "4cf7d4229427e69a22af1fd36437487e29635e84f34e809417370f194b0ae46f"
	keys2.Public = "03780f2d2fcab3e57a0456085b820fd35e4ec27c6c207480bc31a4c9ceb5885709"
	keys2.Wif = "T5dbMGAHAJpTkyncDYbDjPpYMguYa4QjmtkWzngGAhDFQZmakkkN"

	fmt.Println(keys1)
	fmt.Println(keys2)

	fmt.Println("hello")

	// a:=0

	// for a==0{
	// 	fmt.Scanf("%d", &a)
	// }

	input := gobcy.TXInput{}
	output1 := gobcy.TXOutput{}
	// output2 := gobcy.TXOutput{}

	input.Addresses = []string{keys2.Address}
	// input.ScriptType = "null-data"
	output1.ScriptType = "null-data"
	output1.Addresses = []string{keys2.Address}
	// output1.DataHex = "48656c6c6f2c204920616d2074616b6f79616b6933"
	output1.DataString = str//"I am takoyaki3"

	// output2.Addresses = []string{"LbgqbgeEEnc9eEFE6AjWqy9XFoahzobMRM"}

	// tx := gobcy.TempNewTX(keys2.Address,keys1.Address,*big.NewInt(10))
	tx := gobcy.TX{}
	tx.Inputs = []gobcy.TXInput{input}
	tx.Outputs = []gobcy.TXOutput{output1}

	fmt.Println(tx)

	txskel,err := bc.NewTX(tx,true)
	// txskel.PubKeys = []string{"4c62677162676545456e63396545464536416a5771793958466f61687a6f624d524d"}
	// txskel.Signatures = []string{"367631366b5645555947423379784743337961377750356a3541364a4e6f5332574c427a794357674e673650664b4538527637"}
	fmt.Println("################")
	fmt.Println(txskel)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("################")

	err = txskel.Sign([]string{keys2.Private})	// 必要に応じて、ここの配列の数は調整が必要
	fmt.Println(err)
	fmt.Println(txskel)
	fmt.Println("################")

	txskel, err = bc.SendTX(txskel)

	fmt.Println(txskel.Trans.Hash)
	fmt.Println(err)

	fmt.Println("--------------------------------------------------------------------")

	return txskel.Trans.Hash
}