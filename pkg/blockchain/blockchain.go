package blockchain

import (
	"os"
	"log"
	"time"
	"./block"
)

type BlockChain struct{
	Latest string	// 最新ブロックのインデックス
	LatestHex string // 最新ブロックのハッシュ
}

func LoadChain()BlockChain{
	return BlockChain{}
}

func NewBlock(bc BlockChain,filepaths []string)block.Block{
	b := block.Block{}

	// プロパティ
	b.Properties = map[string]string{}
	b.Properties["type"] = "files"
	b.Properties["previous_hash_0"] = "hex"
	b.Properties["created"] = time.Now().String()

	// ファイルをブロックへ追加
	for _,v:=range filepaths{
		f := block.File{}

		// ファイル読み込み
		file, err := os.Open(v)
		if err != nil {
			log.Fatal("error occured 'os.Open()'")
			panic(err)
		}

		stats, statsErr := file.Stat()
    if statsErr != nil {
			return b
    }

		var size int64 = stats.Size()
    bytes := make([]byte, size)
	
		_,err = file.Read(bytes)
		
		if err != nil{
			return b
		}
	
		f.RowData = bytes
		f.Properties = block.Properties{}
		b.Files = append(b.Files,f)
	}
	return b
}