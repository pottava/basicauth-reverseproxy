# Basic 認証も可能なリバースプロキシ


## 概要

指定した URL にリバースプロキシするだけのサービスです。  
オプションでフロントに Basic 認証がかけられます。


## 使い方

### 1. 環境変数をセットします

環境変数            | 説明                                             | 必須
------------------ | ----------------------------------------------- | ---------
PROXY_URL          | リバースプロキシ先の `URL`                         | *
BASIC_AUTH_USER    | Basic 認証をかけるなら、その `ユーザー名`            | 
BASIC_AUTH_PASS    | Basic 認証をかけるなら、その `パスワード`            | 
APP_PORT           | このサービスが待機する `ポート番号` （デフォルト 80番） | 

### 2. アプリを起動します

`docker run -d -p 8080:80 -e PROXY_URL pottava/proxy`