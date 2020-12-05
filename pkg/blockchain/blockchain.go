package blockchain

import (
	"os"
	"log"
	"time"
	"./block"
)

type BlockChain struct{
	Latest string	// 最新ブロックのインデックス
	LatestHex map[string]string // 最新世代ブロックのハッシュ
}

func LoadChain()BlockChain{

	// 既にブロックチェーンが存在するか

	return BlockChain{}
}

func NewBlock(bc BlockChain,filepaths []string)block.Block{
	b := block.Block{}

	// プロパティ
	b.Properties = block.Properties{}
	b.Properties["type"] = "files"
	for k,v := range bc.LatestHex{
		b.Properties["previous_hash_"+k] = v
	}
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

func NewBlockFromRowfile(bc BlockChain,rowfile []byte,filename string)block.Block{
	b := block.Block{}

	// プロパティ
	b.Properties = block.Properties{}
	b.Properties["type"] = "files"
	for k,v := range bc.LatestHex{
		b.Properties["previous_hash_"+k] = v
	}
	b.Properties["created"] = time.Now().String()

	f := block.File{}
	f.RowData = rowfile
	f.Properties = block.Properties{"filename":filename}
	b.Files = append(b.Files,f)

	return b
}