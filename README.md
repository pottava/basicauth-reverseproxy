# Reverse proxy w/ basic authentication

![circleci status](https://circleci.com/gh/pottava/basicauth-reverseproxy.svg?style=shield&circle-token=e15e15a99a3ad48806369829ab87e447aed7fbe7)

## Description

This is a reverse proxy, which is able to provide basic authentication as well.  
([日本語はこちら](https://github.com/pottava/basicauth-reverseproxy/blob/master/README-ja.md))


## Usage

### 1. Set environment variables

Environment Variables     | Description                                       | Required
------------------------- | ------------------------------------------------- | ---------
PROXY_URL                 | The URL to be proxied with this app.              | *
BASIC_AUTH_USER           | User for basic authentication.                    | 
BASIC_AUTH_PASS           | Password for basic authentication.                | 
APP_PORT                  | The port number to be assigned for listening.     | 
SSL_CERT_PATH             | TLS: cert.pem file path.                          | 
SSL_KEY_PATH              | TLS: key.pem file path.                           | 

### 2. Run the application

`docker run -d -p 8080:80 -e PROXY_URL pottava/proxy`

with basic auth:  
`docker run -d -p 8080:80 -e PROXY_URL -e BASIC_AUTH_USER -e BASIC_AUTH_PASS pottava/proxy`


## Copyright and license

Code released under the [MIT license](https://github.com/pottava/basicauth-reverseproxy/blob/master/LICENSE).
