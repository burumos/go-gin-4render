# Go Gin Web サーバー

このリポジトリは、[Render](https://render.com/) で [Go](https://golang.org/) Web アプリケーションをデプロイするためのスタートポイントとして使用できます。

[Gin](https://github.com/gin-gonic/gin) Web フレームワークをベースにした [リアルタイムチャット](https://github.com/gin-gonic/examples/tree/master/realtime-advanced) のサンプルです。

サンプルアプリは https://go-gin.onrender.com で公開されています。

## デプロイ

https://render.com/docs/deploy-go-gin のガイドを参照してください。

## ローカル開発

### 前提条件
- Go 1.21+
- Docker (PostgreSQL 用)

### セットアップ
1. リポジトリをクローン。
2. Docker で PostgreSQL を起動:
   ```bash
   docker compose up -d
   ```
3. `.env` ファイルに環境変数を設定:
   ```
   DATABASE_URL=postgres://user:password@localhost:5432/testdb?sslmode=disable
   ```
4. サーバーをビルドして実行:
   ```bash
   ./build.sh
   PORT=8080 ./app
   ```
5. `http://localhost:8080/login` にアクセスしてログイン。

### ユーザー管理

ユーザーを作成するには、CLI ツールを使用:

1. CLI をビルド:
   ```bash
   ./build-create-user.sh
   ```
2. ユーザーを作成:
   ```bash
   ./create-user username:password
   ```
   例:
   ```bash
   ./create-user admin:secret123
   ```

または、Go で直接実行:
```bash
go run ./cmd/create-user username:password
```

ユーザーは `users` テーブルにハッシュ化されたパスワードで保存されます。
