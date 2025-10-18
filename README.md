# Todoist Inbox Cleaner

Todoistのインボックスから、URLを含むタスクを自動で別プロジェクトに移動するCLIツールです。

## 概要

Pocket等の「あとで読む」サービスの代わりにTodoistを使用する際、インボックスに記事のURLを含むタスクを追加すると、他のタスクと混ざってしまいます。このツールは、URLを含むタスクを自動的に専用プロジェクトに移動させることで、インボックスを整理します。

## 主な機能

- インボックス内のタスクを自動スキャン
- タイトルにURLが含まれるタスクを検出
- 指定したプロジェクトに自動移動
- 設定ファイル(config.toml)による柔軟な設定管理
- APIレート制限対策(自動リトライ、待機時間)
- デバッグコマンドによる詳細情報の確認

## 技術スタック

- **Go**: プログラミング言語
- **Cobra**: CLIコマンド構造の構築
- **Viper**: 設定ファイル(TOML形式)管理
- **Todoist API**: タスク管理
- **VS Code Dev Containers**: 開発環境のコンテナ化

## 前提条件

- TodoistアカウントとAPIトークン
- Docker Desktop または Docker Engine (開発環境を使用する場合)
- VS Code と Dev Containers 拡張機能 (開発環境を使用する場合)

## セットアップ

### 1. リポジトリのクローン

```bash
git clone https://github.com/yourusername/todoist-inbox-cleaner.git
cd todoist-inbox-cleaner
```

### 2. Todoist APIトークンの取得

1. [Todoist Settings](https://todoist.com/prefs/integrations) にアクセス
2. 「Developer」タブを開く
3. APIトークンをコピー

### 3. 設定ファイルの編集

[config.toml](config.toml) を編集して、以下の設定を行います:

```toml
[todoist]
apiToken = "YOUR_TODOIST_API_TOKEN"  # 取得したAPIトークンを設定
targetProjectName = "あとで読む"      # 移動先のプロジェクト名を設定
```

### 4. Todoistでプロジェクトを作成

Todoistアプリで、`config.toml` で指定したプロジェクト名のプロジェクトを作成してください。

### 5. 開発環境のセットアップ (オプション)

VS Code Dev Containersを使用する場合:

1. VS Codeでプロジェクトを開く
2. コマンドパレット (`Ctrl+Shift+P` または `Cmd+Shift+P`) を開く
3. "Dev Containers: Reopen in Container" を選択
4. コンテナのビルドと起動を待つ

## プロジェクト構成

```bash
.
├── .devcontainer/          # Dev Containers 設定
├── cmd/                    # Cobraコマンド定義
│   ├── clean.go           # cleanコマンド(メイン機能)
│   ├── debug.go           # debugコマンド(デバッグ用)
│   ├── info.go            # infoコマンド
│   ├── root.go            # ルートコマンド
│   └── version.go         # versionコマンド
├── internal/              # 内部パッケージ
│   ├── todoist/          # Todoist APIクライアント
│   │   └── client.go     # クライアント、レート制限対策、リトライ処理
│   └── utils/            # ユーティリティ
│       └── url.go        # URL検出機能
├── config.toml            # 設定ファイル
├── go.mod                 # Go Modules
├── go.sum
└── main.go                # エントリーポイント
```

### 主要ファイルの説明

- **[cmd/clean.go](cmd/clean.go)**: インボックスからURLを含むタスクを移動するメインロジック
- **[cmd/debug.go](cmd/debug.go)**: デバッグ情報表示(プロジェクト一覧、タスク一覧)
- **[internal/todoist/client.go](internal/todoist/client.go)**: Todoist API操作、レート制限対策、自動リトライ処理
- **[internal/utils/url.go](internal/utils/url.go)**: URL検出用の正規表現パターンマッチング

## ビルド方法

開発コンテナに接続後、またはローカル環境で以下のコマンドを実行します:

```bash
go build -o todoist-inbox-cleaner .
```

成功すると、`todoist-inbox-cleaner` 実行ファイルが生成されます。

## 使い方

### cleanコマンド (メイン機能)

インボックス内のURLを含むタスクを別プロジェクトに移動します:

```bash
./todoist-inbox-cleaner clean
```

実行すると:
1. インボックス内の全タスクをスキャン
2. タイトルにURLが含まれるタスクを検出
3. 検出したタスクを設定ファイルで指定したプロジェクトに移動
4. 移動したタスク数を表示

**実行例:**
```
$ ./todoist-inbox-cleaner clean
Using config file: /workspaces/todoist-inbox-cleaner/config.toml
移動先プロジェクト: あとで読む👀 (ID: 2361532527)
インボックス内のタスク数: 154

  [1/154] 移動中: [記事タイトル](https://example.com)
  [2/154] 移動中: https://example.com/article
  ...

完了: 154 件のタスクを 'あとで読む👀' プロジェクトに移動しました
```

**注意事項:**
- レート制限対策として、各タスク移動後に200ミリ秒の待機時間が入ります
- 大量のタスク(100件以上)を移動する場合、数分かかることがあります
- APIレート制限により失敗した場合、自動的に最大3回までリトライします
- それでも失敗する場合は、少し待ってから再実行してください

### debugコマンド (デバッグ用)

Todoistのプロジェクトとタスクの情報を確認します:

```bash
./todoist-inbox-cleaner debug
```

実行すると:
- 全プロジェクトの一覧(ID、名前、インボックスフラグ)
- インボックス内のタスク一覧(最初の20件)
- インボックス内のタスク総数

このコマンドは、インボックスのタスクが正しく認識されているか確認する際に便利です。

### その他のコマンド

**infoコマンド**: 現在の設定内容を確認

```bash
./todoist-inbox-cleaner info
```

**versionコマンド**: アプリケーションのバージョンを確認

```bash
./todoist-inbox-cleaner version
```

## 環境変数による設定

環境変数でも設定を上書きできます:

```bash
# APIトークンを環境変数で設定
export TODOIST_CLEANER_TODOIST_APITOKEN="your_api_token_here"

# 移動先プロジェクト名を環境変数で設定
export TODOIST_CLEANER_TODOIST_TARGETPROJECTNAME="あとで読む"

./todoist-inbox-cleaner clean
```

## 自動実行の設定

定期的に実行したい場合は、cron等を使用できます:

```bash
# 毎日午前9時に実行
0 9 * * * /path/to/todoist-inbox-cleaner clean
```

## トラブルシューティング

### APIトークンエラー

```
エラー: Todoist APIトークンが設定されていません
```

→ `config.toml` の `todoist.apiToken` を確認してください

### プロジェクトが見つからないエラー

```
エラー: プロジェクト 'あとで読む' が見つかりません
```

→ Todoistアプリで指定した名前のプロジェクトを作成してください

### インボックスのタスクが0件と表示される

```
インボックス内のタスク数: 0
```

→ 以下を確認してください:
1. `debug` コマンドを実行して、実際のインボックスのタスク数を確認
2. インボックスのタスクが完了済み(チェック済み)になっていないか確認
3. タスクが既に他のプロジェクトに移動済みでないか確認

### APIレート制限エラー

```
警告: タスクの移動に失敗しました: failed to move task: bad request: 429 Too Many Requests
```

→ Todoist APIのレート制限に達しました。対処方法:
1. **自動リトライ**: 通常は自動的に3回までリトライします
2. **待機してから再実行**: それでも失敗する場合は、15分程度待ってから再実行してください
3. **レート制限**: Todoistは450リクエスト/15分の制限があります

### デバッグ方法

問題が発生した場合は、以下の順序で確認してください:

1. **設定確認**
   ```bash
   ./todoist-inbox-cleaner info
   ```

2. **デバッグ情報確認**
   ```bash
   ./todoist-inbox-cleaner debug
   ```

   これで以下が確認できます:
   - インボックスプロジェクトのID
   - インボックス内のタスク一覧
   - 各タスクのProjectIDとContent

## 技術的な詳細

### インボックスの識別方法

Todoistでは、インボックスも通常のプロジェクトとして扱われます:
- プロジェクトの `InboxProject` フラグが `true` のものがインボックス
- インボックスに属するタスクは、インボックスプロジェクトのIDを `ProjectID` として持つ

### URL検出ロジック

以下のパターンをURLとして検出します:
- `http://` または `https://` で始まる文字列
- `www.` で始まる文字列
- Markdownリンク形式 `[タイトル](URL)` も検出対象

正規表現: `https?://[^\s]+|www\.[^\s]+`

### APIレート制限対策

**Todoist APIの制限:**
- 450リクエスト/15分(平均30リクエスト/分)

**実装している対策:**
1. **待機時間**: 各タスク移動後に200ミリ秒の待機
   - これにより、1秒あたり最大5リクエストに制限
   - 15分で最大450リクエスト = 制限内に収まる

2. **自動リトライ**: 429エラー時の対応
   - 最大3回までリトライ
   - エクスポネンシャルバックオフ: 2秒 → 4秒 → 8秒

3. **エラーハンドリング**
   - レート制限エラーとその他のエラーを区別
   - レート制限エラーのみリトライ対象

### パフォーマンス

タスク移動の所要時間(目安):
- 10件: 約2秒
- 50件: 約10秒
- 100件: 約20秒
- 150件: 約30秒

※ レート制限対策の待機時間(200ms/タスク)を含む

## ライセンス

このプロジェクトは自由にご利用いただけます。

## 貢献

バグ報告や機能改善の提案は、GitHubのIssuesでお願いします。
