# ca_tech_dojo_server_side

オンライン版　CA Tech Dojo サーバサイド (Go)編を進めていく

使い方
- docker-compose up でmysqlを立ち上げる

- curl -H 'Content-Type:application/json' -X POST  -d '{"name":"yurina hirate"}' localhost:8080/user/create　で登録

    jsonでjwtトークンが返ってくる

- curl localhost:8080/user/get -H "x-token: もらったトークン"　でデータベース照会

    jsonでデータベースの情報が返ってくる

- curl -H 'Content-Type:application/json' -H "x-token: もらったトークン" -X PUT　-d '{"name":"nana komatu"}' localhost:8080/user/update　名前を更新

- curl -X POST "http://localhost:8080/gacha/draw" -H  "accept: application/json" -H  "x-token: もらったトークン" -H  "Content-Type: application/json" -d "{  \"times\": 引きたい回数}"

    ガチャを引く

- curl -X GET "http://localhost:8080/character/list" -H  "accept: application/json" -H  "x-token: もらったトークン"

    ユーザの所持キャラクタが返って来る