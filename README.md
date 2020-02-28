# ca_tech_dojo_server_side

オンライン版　CA Tech Dojo サーバサイド (Go)編を進めていく

使い方
- docker-composeでmysqlを立ち上げる

- curl -H 'Content-Type:application/json' -X POST  -d '{"name":"yurina hirate"}' localhost:8080/user/create　で登録

    jsonでjwtトークンが返ってくる

- curl localhost:8080/user/get -H "x-token: もらったトークン"　でデータベース照会

    jsonでデータベースの情報が返ってくる

- curl -H 'Content-Type:application/json' -H "x-token: もらったトークン" -X PUT　-d '{"name":"nana komatu"}' localhost:8080/user/update　名前を更新