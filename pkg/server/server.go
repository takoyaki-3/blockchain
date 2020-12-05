package main

import (
	"fmt"
	"io"
	"html/template"
	"net/http"
	"mime/multipart"
	"../blockchain"
	"../blockchain/encoder"
	"bytes"
)

func main () {
	var mux *http.ServeMux;
	mux = http.NewServeMux();
	mux.HandleFunc("/hello", func (writer http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(writer, "hello, world.");
	})
	var hf http.HandlerFunc;
	hf = func (writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "HandlerFunc型を定義->Handleメソッドに渡す");
	}
	mux.Handle("/hf", hf);
	mux . HandleFunc("/upload", upload);
	mux . HandleFunc("/", index);
	
	var server *http.Server;
	// http.Serverのオブジェクトを確保
	// &をつけること構造体ではなくポインタを返却
	server = &http.Server{}; // or new (http.Server);
	server.Addr = ":11180";
	server.Handler = mux;
	server.ListenAndServe();
}

func index (writer http.ResponseWriter , request *http.Request) {
	var t *template.Template;

	// テンプレートをロード
	t, _ = template.ParseFiles("./index.html");
	t.Execute(writer, struct{}{});
}

func upload ( w http.ResponseWriter, r *http.Request) {
	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if  (r.Method != "POST") {
		fmt.Fprintln(w, "Please access by POST.");
		return;
	}
	var file multipart.File;
	var fileHeader *multipart.FileHeader;
	var e error;
	var uploadedFileName string;
	// POSTされたファイルデータを取得する
	file , fileHeader , e = r.FormFile("file");
	if (e != nil) {
		fmt.Fprintln(w, "error occurred in uploading file.");
		return;
	}
	uploadedFileName = fileHeader.Filename;

	bc := blockchain.LoadChain()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
    // return nil, err
	}
	block := blockchain.NewBlockFromRowfile(bc,buf.Bytes(),uploadedFileName)
	encoder.Write(block,"./s.block")

	defer file.Close();
}