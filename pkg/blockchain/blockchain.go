package blockchain

import (
	"os"
	"log"
	"time"
	"./block"
	"./encoder"
	"./decoder"
	"io/ioutil"
	"strings"
	"strconv"
	"fmt"
	"sync"
	"./link"
)

type BlockChain struct{
	Latest int	// 最新ブロックのインデックス
	LatestHex map[int]map[string]string // 最新世代ブロックのハッシュ
}

func LoadChain()BlockChain{

	bc := BlockChain{}

	_, err := os.Stat("hoge.txt")
	if err != nil {
    if err := os.Mkdir("./blocks", 0777); err != nil {
			fmt.Println(err)
		}
	}

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

func NewBlock(filepaths []string)block.Block{
	b := block.Block{}

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

func NewBlockFromRowfile(rowfile []byte,filename string)block.Block{
	f := block.File{}
	f.RowData = rowfile
	f.Properties = block.Properties{"filename":filename}
	return block.Block{block.Properties{"type":"files"},uint64(1),[]block.File{f}}
}
func NewBlockFromString(str string)block.Block{
	f := block.File{}
	f.RowData = []byte(str)
	f.Properties = block.Properties{"type":"string"}
	return block.Block{block.Properties{"type":"string"},uint64(1),[]block.File{f}}
}

var mu sync.Mutex

func AddBlock(bc *BlockChain,b block.Block)string{

	// ブロックチェーンに組み込む為のプロパティ付与
	b.Properties = block.Properties{}
	for k,v := range bc.LatestHex[bc.Latest]{
		b.Properties["previous_hash_"+k] = v
	}
	b.Properties["created"] = time.Now().String()

	mu.Lock()
	bc.Latest+=1
	index := encoder.Write(b,bc.Latest)
	if _,ok:=bc.LatestHex[bc.Latest];!ok{
		bc.LatestHex[bc.Latest]=map[string]string{}
	}
	s := strings.Split(index, ".")
	bc.LatestHex[bc.Latest][s[0]+"."+s[1]] = s[1]
	if _,ok:=bc.LatestHex[bc.Latest-5];ok{
		delete(bc.LatestHex,bc.Latest-5)
	}
	mu.Unlock()
	return index
}

func ReadBlock(index string)block.Block{
	s := strings.Split(index, ".")
	block := decoder.Read("./blocks/"+s[0]+"."+s[1]+".block")
	return block
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

func MakeLinkBlock(bc *BlockChain,publick_chain string){
	// 何世代目から作成するか
	paths := dirwalk("./blocks")

	indexes := [][]string{}
	for _,v:=range paths{
		s := strings.Split(v, ".")
		gen,_ := strconv.Atoi(s[0])
		for len(indexes) <= gen{
			indexes = append(indexes,[]string{})
		}
		indexes[gen] = append(indexes[gen],v)
	}

	start := len(indexes)-1
	linkBlock := link.LinkBlock{}

	for ;start>=0;start--{
		flag := false
		for _,v:=range indexes[start]{
			s := strings.Split(v, ".")
			b := link.Block{}
			b.Id = v
			b.Hex = s[1]
			linkBlock.Blocks = append(linkBlock.Blocks,b)

			fmt.Println(v)
			block := decoder.Read("./blocks/"+v)
			if block.Properties["type"] == publick_chain+"_"+"link_block"{
				flag = true
			}
		}
		if flag{
			break
		}
	}

	linkBlock.Created = time.Now().String()
	rowJSON := string(link.JsonEncode(linkBlock))

	fmt.Println(rowJSON)

	fmt.Println(encoder.Hex([]byte(rowJSON)))
}