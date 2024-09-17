# FileSync

## 概要

FileSyncは、HTTPリクエストを通じてファイルのアップロード、ダウンロード、削除、リスト取得、解凍、鍵の変更などの操作を行うシンプルなファイル同期システムです。

## 特徴

* 簡易なHTTP APIでファイル操作が可能
* 署名による認証でセキュリティを確保
* zip、tar.gz形式のファイル解凍に対応

## インストール方法

1. Go言語の開発環境を構築してください。
2. `go get github.com/takoyaki-3/filesync` を実行してリポジトリをクローンします。

## 使い方

### 設定ファイル

`config.json` でサーバの設定を行います。

```json
{
  "port": 11182,
  "hostname": "c3.d.takoyaki3.com"
}
```

* **port**: サーバがListenするポート番号
* **hostname**: サーバのホスト名

### APIエンドポイント

APIエンドポイントは `http://<hostname>:<port>/` です。

* 例: `http://c3.d.takoyaki3.com:11182/`

### 認証

APIリクエストには、署名による認証が必要です。署名は `pkg.Sign()` 関数で生成できます。

```
sign := pkg.Sign()
```

署名は `sign` パラメータとしてAPIリクエストに付与します。

* 例: `http://c3.d.takoyaki3.com:11182/auth?sign=<sign>`

### API一覧

| メソッド | エンドポイント | 説明 | パラメータ |
|---|---|---|---|
| GET | `/auth` | 認証テスト | `sign` |
| POST | `/upload` | ファイルアップロード | `sign`, `path` |
| GET | `/download` | ファイルダウンロード | `sign`, `path` |
| GET | `/remove` | ファイル削除 | `sign`, `path` |
| GET | `/remover` | ディレクトリ削除 | `sign`, `path` |
| GET | `/getlist` | ファイルリスト取得 | `sign`, `path` |
| GET | `/chagekey` | 鍵の変更 | `sign` |
| GET | `/unzip` | zipファイル解凍 | `sign`, `path`, `dist` |
| GET | `/untargz` | tar.gzファイル解凍 | `sign`, `path` |

### コマンド実行例

#### ファイルアップロード

```
go run upload.go
```

#### ファイルダウンロード

```
go run cliant.go
```

#### ファイルリスト取得

```
curl -X GET "http://<hostname>:<port>/getlist?sign=<sign>&path=<path>"
```

#### zipファイル解凍

```
curl -X GET "http://<hostname>:<port>/unzip?sign=<sign>&path=<path>&dist=<dist>"
```


## ファイルツリー

```
├── auth-test.go
├── cliant.go
├── cliant_odpt.go
├── config.json
├── go.mod
├── go.sum
├── key
├── pkg
│   └── common.go
├── server.go
├── unzip.go
└── upload.go
```

### ファイルの説明

* **auth-test.go**: 認証テストを行うプログラム
* **cliant.go**: ファイルダウンロードを行うプログラム
* **cliant_odpt.go**: ODPT用ファイルダウンロードプログラム
* **config.json**: サーバ設定ファイル
* **go.mod**: Goモジュールファイル
* **go.sum**: Goモジュールチェックサムファイル
* **key**: 認証鍵ファイル
* **pkg/common.go**: 共通関数
* **server.go**: APIサーバ
* **unzip.go**: zipファイル解凍プログラム
* **upload.go**: ファイルアップロードプログラム

## 注意点

* 本システムはセキュリティを考慮した設計にはなっていません。機密情報の取り扱いには注意してください。
