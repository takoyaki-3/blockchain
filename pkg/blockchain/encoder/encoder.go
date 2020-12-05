package encoder

import (
	"log"
	"os"
	"strconv"
	"fmt"
	"encoding/binary"
	"crypto/sha256"
	"../block"
)

func uint642bytes(u uint64)[]byte{
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, u)
	return b
}

func propertiesEncoeder(properties block.Properties)[]byte{
	b := []byte{}
	b = append(b,uint642bytes(uint64(len(properties)))...)
	for k,v := range properties{
		b = append(b,uint642bytes(uint64(len(k)))...)
		b = append(b,[]byte(k)...)
		b = append(b,uint642bytes(uint64(len(v)))...)
		b = append(b,[]byte(v)...)
	}
	return b
}

func File2Bytes(file block.File)[]byte{
	b := []byte{}
	propRow := propertiesEncoeder(file.Properties)
	b = append(b,uint642bytes(uint64(len(propRow)))...)
	b = append(b,uint642bytes(uint64(len(file.RowData)))...)
	b = append(b,propRow...)
	b = append(b,file.RowData...)
	return b
}

func Block2Bytes(block block.Block)[]byte{
	p := []byte{}
	b := []byte{}

	// ブロックプロパティの追加
	p = append(p,propertiesEncoeder(block.Properties)...)

	// プロパティサイズ
	b = append(b,uint642bytes(uint64(len(p)))...)
	// ファイル数の追加
	b = append(b,uint642bytes(uint64(len(block.Files)))...)

	// ファイルのエンコード
	files := [][]byte{}
	for _,v := range block.Files{
		files = append(files,File2Bytes(v))
	}

	// ファイル開始位置までのバイト数をファイル数分追加
	i := uint64(16+(len(files)+1)*8+len(p))
	for _,v:=range files{
		b = append(b,uint642bytes(i)...)
		i += uint64(len(v))
	}
	b = append(b,uint642bytes(i)...)

	// プロパティ本体
	b = append(b,p...)

	// ファイル本体
	for _,v := range files{
		b = append(b,v...)
	}

	// プロパティサイズ (8bit) 0~7 prosize
	// ファイル数 (8bit) 8~15 filenum
	// ファイルサイズの配列 (8bit * filenum) 16~8*filenum-1
	// ファイル終了位置 (8bit)
	// プロパティ本体 8*(filenum+1)~8*filenum+prosize-1
	// ファイル本体 8*filenum+prosize~

	return append(uint642bytes(uint64(len(b)+8)),b...)
}

func WriteByte(rowData []byte,generation int)string{
	index := strconv.Itoa(generation)+"."+Hex(rowData)
	path := "./blocks/"+index+".block"
	wf, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer wf.Close()

	// データ部分を書き込み
	_, err = wf.Write(rowData)
	if err != nil {
		log.Fatal(err)
	}
	return index
}

func Write(block block.Block,generation int)string{
	return WriteByte(Block2Bytes(block),generation)
}

func Hex(block []byte)string{
	b := block
	hash := sha256.Sum256(b)
	hash_s := fmt.Sprintf("%x",hash)
	return hash_s
}