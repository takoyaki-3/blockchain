package link

import (
	"encoding/json"
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