# ブロックチェーン

## Usage
```
go run main.go
```

## 使い方

- GUIによるファイルアップロード
``/``

- ファイル書き込み
``/upload``
-- file : アップロードするファイル

- ファイル取得
``/get_file``
-- index : トランザクションID

- 文字列書き込み
``/add_string``
-- data : 書き込む文字列

- 文字列読み取り
``/read_string``
-- index : トランザクションID

## To Do

- 仕様策定（検索をどう落とし込むか、イーサリアムとの同期）
- ブロックの保存、読み込み
- 複数マシンでの同期
- イーサリアムとの同期

### ブロックチェーンに含まれるデータ

#### ブロック本体
- ブロックプロパティ（連想配列）
- 含まれるファイル数
- ファイル配列

#### ファイル
-- サイズ
-- ファイルプロパティ（連想配列）
-- バイナリデータデータサイズ
-- バイナリデータ本体

#### プロパティ連想配列

# 元URL
- https://qiita.com/seita-uc/items/553d315b84b0e7bfc4d0
- https://blockchain-jp.com/tech/2011
- https://github.com/olegpankov/go_blockchain
