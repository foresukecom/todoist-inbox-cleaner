# Go CLI Sample with Dev Containers, Cobra, and Viper

このリポジトリは、VS Code Dev Containers を使用して、Go言語で以下の要素を取り入れたコマンドラインインターフェース（CLI）アプリケーションを開発するためのサンプルです。

- **VS Code Dev Containers**: ローカル環境を汚さずに開発環境をコンテナ内に構築
- **Go Modules**: Goの依存関係管理
- **Cobra**: CLIコマンド構造の構築
- **Viper**: 設定ファイル（TOML形式）および環境変数からの設定読み込みと管理

## 主な特徴

- Dockerコンテナ内で開発環境が完結するため、ローカルへのGoや開発ツールのインストールが最小限で済みます。
- Cobraによる構造化されたコマンドラインインターフェースの例を示します。
- Viperによる柔軟な設定管理（TOMLファイル、環境変数）を実演します。
- `info` コマンドで読み込まれた設定内容を確認できます。
- `version` コマンドでアプリケーションのバージョンを確認できます。
- **自動ビルドとリリース**: git タグをプッシュすると、GitHub Actions により各アーキテクチャ向けのバイナリが自動的にビルドされ、リリースが作成されます。

## 前提条件

- Docker Desktop または Docker Engine が開発用端末にインストールされていること。
- VS Code が開発用端末にインストールされていること。
- VS Code 拡張機能 **Dev Containers** がインストールされていること。（旧称: Remote - Containers）

## はじめるには

1.  このリポジトリを開発用端末にクローンします。
```
git clone <このリポジトリのURL> # ★★★ ここをあなたのリポジトリのURLに置き換えてください ★★★
cd <クローンしたディレクトリ>
```
2.  VS Codeでクローンしたプロジェクトのルートフォルダを開きます。
3.  VS Codeが `.devcontainer` フォルダを検出し、「フォルダーをコンテナーで再度開きますか？」のようなプロンプトが表示されたら **「Reopen in Container」** を選択します。
    プロンプトが表示されない場合は、コマンドパレット (`Ctrl+Shift+P` または `Cmd+Shift+P`) を開いて「Dev Containers: Reopen in Container」を実行します。
4.  Dockerイメージのビルドとコンテナの起動が自動的に行われます。これには初回時間がかかる場合があります。
5.  VS Codeがコンテナに接続されたら開発準備完了です。VS Codeの左下隅が緑色になり、「Dev Container: ...」のような表示が出ます。

## プロジェクト構成

```bash
.
├── .devcontainer/      # Dev Containers 設定ファイルと Dockerfile
│   ├── Dockerfile
│   ├── compose.yml
│   └── devcontainer.json
├── README.md           # このファイル
├── cmd/                # Cobra コマンド定義
│   ├── info.go         # info コマンドの実装
│   ├── root.go         # ルートコマンドの定義、Viper 初期化
│   └── version.go      # version コマンドの実装
├── config.toml         # 設定ファイル (TOML形式)
├── go.mod              # Go Modules ファイル
├── go.sum              # Go Modules チェックサムファイル
└── main.go             # アプリケーションエントリーポイント
```

- `.devcontainer/`: VS Code Dev Containers がコンテナを構築・接続するための設定が含まれています。`Dockerfile` で開発環境となるコンテナイメージを定義し、`compose.yml` でコンテナの起動設定（ボリュームマウントなど）を行い、`devcontainer.json` でVS Codeの接続方法を構成します。
- `cmd/`: Cobraコマンドの定義ファイルが含まれます。`root.go` がアプリケーションの起点となるルートコマンドを定義し、`info.go` や `version.go` がそれぞれのサブコマンドを定義します。
- `config.toml`: アプリケーションが使用する設定値を記述するファイルです。Viperによって読み込まれます。
- `go.mod` / `go.sum`: プロジェクトのGoモジュール情報と依存ライブラリのチェックサム情報です。
- `main.go`: アプリケーションの起動時に最初に実行されるファイルです。`cmd.Execute()` を呼び出し、Cobraコマンドを実行します。

## ビルド方法

開発コンテナに接続後、VS Codeのターミナルを開き、プロジェクトのルートディレクトリで以下のコマンドを実行します。

```bash
go build .
```
成功すると、プロジェクトルートに実行ファイル（例: `your_cli_app` またはOSに応じた名前）が生成されます。

## 使い方

ビルドして生成された実行ファイルをターミナルで実行します (`your_cli_app` はビルドされた実行ファイル名に置き換えてください)。

```bash
# アプリケーションの実行 (ルートコマンド)
./your_cli_app

# info コマンドの実行
./your_cli_app info

# version コマンドの実行
./your_cli_app version
```

### `info` コマンド

Viperによって現在読み込まれている設定内容をJSON形式で出力します。`config.toml` の内容や、環境変数で上書きされた値を確認するのに役立ちます。  
`password`などの機密情報が含まれる項目はマスクされます。  

```bash
./your_cli_app info
```

### `version` コマンド

アプリケーションのバージョン情報を出力します。
GitHub Actions でビルドされたバイナリの場合、バージョンはビルド時に使用された git タグ（例: `v1.0.0`）になります。
ローカルで `go build` した場合は、ビルド時にバージョン情報が埋め込まれません。


```bash
./your_cli_app version
```

## 設定について

本サンプルアプリケーションは、Viper を使用して設定を管理します。

### `config.toml`

プロジェクトルートにある `config.toml` ファイルが、アプリケーションの主要な設定ファイルです。現在の内容は以下のようになっています。

```toml
# 作者
author = "xxx"

# デバッグモードかどうか
debug = true

# RDBセクション (実際のサンプルでは使用しません)
[rdB]
type = "mysql"
host = "localhost"
port = 3306
user = "myuser"
password = "mypassword"
dbName = "mydatabase"
sslmode = "disable"
maxConnections = 10
timeout = "5s"
```

### 環境変数による設定上書き

Viperは、設定ファイルだけでなく環境変数からも設定を読み込むように構成されています (`cmd/root.go` の `initViper` 関数を参照)。

`viper.SetEnvPrefix("YOURAPP")` と `viper.AutomaticEnv()` の設定により、`YOURAPP_` で始まる環境変数がスキャンされ、対応する設定キーにマッピングされます（環境変数名の `_` は設定キーの `.` に対応します）。

**環境変数で設定された値は、`config.toml` で設定された値やViperのデフォルト値よりも優先されます。**

例として、`author` の値を環境変数で上書きするには、CLIツールを実行する前に `YOURAPP_AUTHOR` 環境変数を設定します。

```bash
# 環境変数 YOURAPP_AUTHOR を設定して info コマンドを実行
YOURAPP_AUTHOR="新しい作者名" ./your_cli_app info
```

このコマンドを実行すると、`info` コマンドの出力で `author` の値が `"xxx"` から `"新しい作者名"` に変わっていることが確認できます。

同様に、RDBのホストを環境変数で上書きしたい場合は `YOURAPP_RDB_HOST="your_db_host"` のように設定します。

この機能により、機密情報（パスワードなど）を設定ファイルに書かずに環境変数で渡したり、開発環境と本番環境で設定を簡単に切り替えたりすることが可能になります。


## ビルドとリリースプロセス (GitHub Actions)
このリポジトリでは、GitHub Actions を利用して、リリースプロセスを自動化しています。 
1. 開発が完了し、新しいバージョンをリリースする準備ができたら、適切なバージョン番号で git タグを作成します (例: git tag v1.0.0)。 
2. 作成したタグをリモートリポジトリ (GitHub) にプッシュします (例: git push origin v1.0.0)。 
3. タグのプッシュをトリガーとして、GitHub Actions のワークフローが実行されます。 
4. ワークフローは、複数のOSおよびアーキテクチャ（Linux, macOS, Windows / amd64, arm64など）に対応したバイナリをビルドします。 
5. ビルドされたバイナリは、GitHub Releases に自動的にアップロードされ、新しいリリースが作成されます。 
このプロセスにより、version コマンドで表示されるバージョン番号は、プッシュされた git タグの名前 が使用されます。また、ビルド時のコミットハッシュやビルド日時も埋め込まれます。

## ライセンス

このサンプルコードは、必要に応じて適切なオープンソースライセンスを適用してください。
