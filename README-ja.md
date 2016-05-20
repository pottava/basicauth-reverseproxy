# Basic 認証も可能なリバースプロキシ


## 概要

指定した URL にリバースプロキシするだけのサービスです。  
オプションでフロントに Basic 認証がかけられます。


## 使い方

### 1. 環境変数をセットします

環境変数            | 説明                                             | 必須
------------------ | ----------------------------------------------- | -------------------
PROXY_PATTERNS     | 仮想ホストや仮想パスと、そのリバースプロキシ先          | * (または PROXY_URL)
PROXY_URL          | リバースプロキシ先の `URL`                         | * (または PROXY_PATTERNS)
BASIC_AUTH_USER    | Basic 認証をかけるなら、その `ユーザー名`            | 
BASIC_AUTH_PASS    | Basic 認証をかけるなら、その `パスワード`            | 
APP_PORT           | このサービスが待機する `ポート番号` （デフォルト 80番） | 
SSL_CERT_PATH      | TLS を有効にしたいなら、その `cert.pem` へのパス     | 
SSL_KEY_PATH       | TLS を有効にしたいなら、その `key.pem` へのパス      | 
ACCESS_LOG         | 標準出力へアクセスログを送る (初期値: false)          | 

### 2. アプリを起動します

`$ docker run -d -p 8080:80 -e PROXY_URL pottava/proxy`

* Basic 認証をつけるなら:  

`$ docker run -d -p 8080:80 -e PROXY_URL -e BASIC_AUTH_USER -e BASIC_AUTH_PASS pottava/proxy`

* TLS を有効にしたいなら:  

`$ docker run -d -p 8080:80 -e PROXY_URL -e SSL_CERT_PATH -e SSL_KEY_PATH pottava/proxy`

* 仮想ホストでのリバースプロキシ:

`$ export PROXY_PATTERNS="/static=http://assets.cdn/,*.example.com=http://app.io/,*=http://sorry.com/"`  
`$ docker run -d -p 8080:80 -e PROXY_PATTERNS pottava/proxy`

* docker-compose.yml として使うなら:  

```
proxy:
  image: pottava/proxy
  ports:
    - 80:80
  links:
    - web
  environment:
    - PROXY_URL=http://web/
    - BASIC_AUTH_USER=admin
    - BASIC_AUTH_PASS=password
    - ACCESS_LOG=true
  container_name: proxy
```

* docker-compose.yml（仮想ホストと SSL/TLS 有効）として使うなら:  

```
proxy:
  image: pottava/proxy
  ports:
    - 443:80
  links:
    - web
  environment:
    - PROXY_PATTERNS="/static=http://assets.cdn/,*.example.com=http://app.io/,*=http://sorry.com/"
    - SSL_CERT_PATH=/etc/certs/cert.pem
    - SSL_KEY_PATH=/etc/certs/key.pem
  volumes:
    - ./certs:/etc/certs
  container_name: proxy
```
