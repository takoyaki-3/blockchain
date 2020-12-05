package decoder

import (
	"os"
	"log"
	"encoding/binary"
	"../block"
)

// Bytes2uint converts []byte to uint64
func Bytes2uint(bytes []byte) uint64 {
	padding := make([]byte, 8-len(bytes))
	i := binary.LittleEndian.Uint64(append(padding, bytes...))
	return i
}

func DecodeProperty(rowData []byte)map[string]string{
	props := map[string]string{}
	numProps := Bytes2uint(rowData[0:8])
	p:=uint64(8)
	for i:=uint64(0);i<numProps;i++{
		ksize := Bytes2uint(rowData[p:p+8])
		p+=8
		k := string(rowData[p:p+ksize])
		p+=ksize
		vsize := Bytes2uint(rowData[p:p+8])
		p+=8
		v := string(rowData[p:p+vsize])
		p+=vsize
		props[k]=v
	}
	return props
}

func DecodeFile(rowData []byte)block.File{

	// プロパティサイズ (8bit) 0~7 prosize
	propertiesSize := Bytes2uint(rowData[0:8])
	// fileSize := Bytes2uint(rowData[8:16])

	f := block.File{DecodeProperty(rowData[16:16+propertiesSize]),rowData[16+propertiesSize:]}
	return f
}

func DecodeBlock(buf []byte)block.Block{

	block := block.Block{}

	// プロパティサイズ (8bit) 0~7 prosize
	propertiesSize := Bytes2uint(buf[0:8])

	// ファイル数 (8bit) 8~15 filenum
	numberOfFile := Bytes2uint(buf[8:16])

	pos := uint64(16)

	filesizes := []uint64{}
	// ファイルサイズの配列 (8bit * filenum) 16~8*filenum-1
	for i:=uint64(0);i<numberOfFile;i++{
		u := Bytes2uint(buf[pos:pos+8])
		filesizes = append(filesizes,u)
		pos+=8
	}

	// ファイル終了位置読み込み
	end := Bytes2uint(buf[pos:pos+8])
	pos+=8
	filesizes = append(filesizes,end)

	// プロパティ本体 8*(filenum+1)~8*filenum+prosize-1
	block.Properties = DecodeProperty(buf[pos:pos+propertiesSize])
	pos+=propertiesSize
	block.NumberOfFiles = numberOfFile

	// ファイル本体 8*filenum+prosize~
	for i:=uint64(0);i<numberOfFile;i++{
		size := filesizes[i+1]-filesizes[i]
		block.Files = append(block.Files,DecodeFile(buf[pos:pos+size]))
		pos += size
	}

	return block
}

func Read(path string)block.Block{
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("error occured 'os.Open()'")
		panic(err)
	}

	buf := make([]byte, 8)

	// ブロックサイズの取得
	file.Read(buf)
	blockSize := Bytes2uint(buf)

	// ブロック取得
	blockBuf := make([]byte, blockSize-8)
	file.Read(blockBuf)

	return DecodeBlock(blockBuf)
}