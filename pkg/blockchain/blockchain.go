package blockchain

import (
	"os"
	"log"
	"time"
	"./block"
	"./encoder"
	"io/ioutil"
	"strings"
	"strconv"
	"fmt"
)

type BlockChain struct{
	Latest int	// 最新ブロックのインデックス
	LatestHex map[int]map[string]string // 最新世代ブロックのハッシュ
}

func LoadChain()BlockChain{

	bc := BlockChain{}

	// 既にブロックチェーンが存在するか
	paths := dirwalk("./blocks")

	bc.Latest = -1
	bc.LatestHex = map[int]map[string]string{}

	for _,v:=range paths{
		slice := strings.Split(v, ".")
		i, _ := strconv.Atoi(slice[0])
		if i > bc.Latest{
			bc.Latest = i
		}
	}
	bc.LatestHex[bc.Latest] = map[string]string{}
	for _,v:=range paths{
		s := strings.Split(v, ".")
		i, _ := strconv.Atoi(s[0])
		if i != bc.Latest{
			continue
		}
		bc.LatestHex[bc.Latest][s[0]+"."+s[1]] = s[1]
	}

	fmt.Println(bc)

	return bc
}

func NewBlock(bc BlockChain,filepaths []string)block.Block{
	b := block.Block{}

	// プロパティ
	b.Properties = block.Properties{}
	b.Properties["type"] = "files"
	for k,v := range bc.LatestHex[bc.Latest]{
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
	for k,v := range bc.LatestHex[bc.Latest]{
		b.Properties["previous_hash_"+k] = v
	}
	b.Properties["created"] = time.Now().String()

	f := block.File{}
	f.RowData = rowfile
	f.Properties = block.Properties{"filename":filename}
	b.Files = append(b.Files,f)

	return b
}

func AddBlock(bc BlockChain,block block.Block)string{
	bc.Latest+=1
	index := encoder.Write(block,bc.Latest)
	if _,ok:=bc.LatestHex[bc.Latest];!ok{
		bc.LatestHex[bc.Latest]=map[string]string{}
	}
	s := strings.Split(index, ".")
	bc.LatestHex[bc.Latest][s[0]+"."+s[1]] = s[1]
	return index
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
			panic(err)
	}

	var paths []string
	for _, file := range files {
		paths = append(paths, file.Name())
	}

	return paths
}