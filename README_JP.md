# Go Todo

![hero](https://res.cloudinary.com/ichtrojan/image/upload/v1574958373/Screenshot_2019-11-28_at_17.22.25_gyegdr.png)

## はじめに

Goで書かれたシンプルなTodoリストアプリケーション

## 必要条件
* MySQLがインストールされていること
* Goがインストールされていること

## インストール方法

* このリポジトリをクローン

git clone https://github.com/ichtrojan/go-todo.git

* ディレクトリを移動

cd go-todo

* `.env`ファイルを作成

cp .env.example .env

* `.env`ファイルに正しいデータベース認証情報と希望のポート番号を設定してください

## 使用方法

このアプリケーションを実行するには、以下のコマンドを実行します：

go run main.go

アプリケーションには`http://127.0.0.1:4040`でアクセスできるはずです

>**注意**<br>
>`.env`ファイルでポート番号を変更した場合は、設定したポート番号でアクセスしてください

## まとめ

このプロジェクトは、デフォルトの`database/sql`パッケージを使用したCRUD操作と、HTMLテンプレートを適切に提供する方法を学ぶための例です。

このプロジェクトに追加したい内容がある場合は、PRを送信してください。[私](https://github.com/ichtrojan)による積極的なメンテナンスは終了しています。