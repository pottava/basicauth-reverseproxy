# Reverse proxy w/ basic authentication


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

### 2. Run the application

`docker run -d -p 8080:80 -e PROXY_URL pottava/proxy`


## Copyright and license

Code released under the [MIT license](https://github.com/pottava/basicauth-reverseproxy/blob/master/LICENSE).
