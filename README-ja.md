# Basic 認証も可能なリバースプロキシ

[![pottava/proxy](http://dockeri.co/image/pottava/proxy)](https://hub.docker.com/r/pottava/proxy/)

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
CORS_ALLOW_ORIGIN  | CORS を有効にしたいなら、リソースへのアクセスを許可する URI. | 
CORS_ALLOW_METHODS | CORS を有効にしたいなら、許可する [HTTP request methods](https://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html)のカンマ区切りのリスト. | 
CORS_ALLOW_HEADERS | CORS を有効にしたいなら、サポートするヘッダーのカンマ区切りのリスト. | 
CORS_MAX_AGE       | CORS における preflight リクエスト結果のキャッシュ上限時間(秒). (デフォルト 600秒) | 
ACCESS_LOG         | 標準出力へアクセスログを送る (初期値: false)          | 
CONTENT_ENCODING   | リクエストが許可して入ればレスポンスを圧縮します (初期値: false) |
HEALTHCHECK_PATH   | 指定すると Basic 認証設定の有無などに依らず 200 OK を返します |
USE_X_FORWARDED_HOST | Host 名は X-Forwarded-Host に格納してバックエンドに流します. (初期値: true) |

### 2. アプリを起動します

`$ docker run -d -p 8080:80 -e PROXY_URL pottava/proxy`

* Basic 認証をつけるなら:  

`$ docker run -d -p 8080:80 -e PROXY_URL -e BASIC_AUTH_USER -e BASIC_AUTH_PASS pottava/proxy`

* TLS を有効にしたいなら:  

`$ docker run -d -p 8080:80 -e PROXY_URL -e SSL_CERT_PATH -e SSL_KEY_PATH pottava/proxy`

* CORS を有効にしたいなら:

`$ docker run -d -p 8080:80 -e PROXY_URL -e CORS_ALLOW_ORIGIN -e CORS_ALLOW_METHODS -e CORS_ALLOW_HEADERS -e CORS_MAX_AGE pottava/proxy`

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
