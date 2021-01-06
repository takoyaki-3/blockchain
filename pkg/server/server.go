package server

import (
	"fmt"
	"io"
	"html/template"
	"net/http"
	"mime/multipart"
	"../blockchain"
	"../blockchain/encoder"
	"bytes"
	"strings"
	"strconv"
	"../blockchain/link"
)

var bc blockchain.BlockChain

func Init (bcn blockchain.BlockChain) {
	bc = bcn
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
	mux . HandleFunc("/addstring", addstring);
	mux . HandleFunc("/read_string", read_string);
	mux . HandleFunc("/get_file", get_file);
	mux . HandleFunc("/link_blockchain",makeLinkBlock)
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

// ファイルアップロードによる登録
func upload ( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )

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
	// fmt.Println(r.FormFile("file"))
	file , fileHeader , e = r.FormFile("file");
	if (e != nil) {
		fmt.Fprintln(w, "error occurred in uploading file.");
		return;
	}
	uploadedFileName = fileHeader.Filename;

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
    // return nil, err
	}
	block := blockchain.NewBlockFromRowfile(buf.Bytes(),uploadedFileName)
	// encoder.Write(block,"./s.block")
	index := blockchain.AddBlock(&bc,block)
	fmt.Fprintln(w, "{\"index\":\""+index+".0\"}");

	defer file.Close();
}

// 文字列の登録
func addstring ( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )

	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if  (r.Method != "POST" && r.Method != "GET") {
		fmt.Fprintln(w, "Please access by POST or GET.");
		return;
	}

	// クエリパラメータを取得する
	queryparm := r.URL.Query()

	if v,ok:=queryparm["data"];ok{
		// 出力
		block := blockchain.NewBlockFromString(v[0])
		index := blockchain.AddBlock(&bc,block)
		fmt.Fprintln(w, "{\"index\":\""+index+".0\"}");
	} else {
		fmt.Fprintln(w, "data must.");
	}
}

func read_string( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )

	// クエリパラメータを取得する
	queryparm := r.URL.Query()

	if v,ok:=queryparm["index"];ok{
		// 出力
		s := strings.Split(v[0], ".")
		block := blockchain.ReadBlock(v[0])
		i,_ := strconv.Atoi(s[2])
		str := string(block.Files[i].RowData)
		fmt.Fprintln(w, "{\"data\":\""+str+"\"}");
	} else {
		fmt.Fprintln(w, "index must.");
	}
}

func get_file( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )

	// クエリパラメータを取得する
	queryparm := r.URL.Query()

	if v,ok:=queryparm["index"];ok{
		// 出力
		s := strings.Split(v[0], ".")
		block := blockchain.ReadBlock(v[0])
		i,_ := strconv.Atoi(s[2])

		w.Header().Set("Content-Disposition", "attachment; filename="+block.Files[i].Properties["filename"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		writer := bytes.NewBuffer(block.Files[i].RowData)
		io.Copy(w, writer)

		// fmt.Fprintln(w, "{\"data\":\""+str+"\"}");
	} else {
		fmt.Fprintln(w, "index must.");
	}
}

func makeLinkBlock( w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )

	queryparm := r.URL.Query()
	if v,ok:=queryparm["key"];ok{
		if v[0]=="iccd" {
			rowJSON := blockchain.MakeLinkBlock(&bc,"LTC")
			hex := encoder.Hex([]byte(rowJSON))
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			tx := link.Write2PublicBlockchain(hex)
			fmt.Println(tx)

			// POSTされたファイルデータを取得する
			buf := []byte(rowJSON)
			block := blockchain.NewBlockFromRowfile(buf,"link_block")
			block.Properties["type"]="_link_block"
			index := blockchain.AddBlock(&bc,block)

			fmt.Fprintln(w, "{\"transaction\":\""+tx+"\",\"JSON\":"+rowJSON+",\"index\":\""+index+"\"}");
		} else {
			fmt.Fprintln(w, "key is uncollect.");
		}
	} else {
		fmt.Fprintln(w, "Please input key.");
	}
}