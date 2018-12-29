# Reverse proxy w/ basic authentication

![circleci status](https://circleci.com/gh/pottava/basicauth-reverseproxy.svg?style=shield&circle-token=e15e15a99a3ad48806369829ab87e447aed7fbe7)

[![pottava/proxy](http://dockeri.co/image/pottava/proxy)](https://hub.docker.com/r/pottava/proxy/)

Supported tags and respective `Dockerfile` links:  
・latest ([prod/Dockerfile](https://github.com/pottava/basicauth-reverseproxy/blob/master/prod/Dockerfile))  
・1.1 ([prod/Dockerfile](https://github.com/pottava/basicauth-reverseproxy/blob/master/prod/Dockerfile))  
・1 ([prod/Dockerfile](https://github.com/pottava/basicauth-reverseproxy/blob/master/prod/Dockerfile))  

## Description

This is a reverse proxy, which is able to provide basic authentication as well.  
([日本語はこちら](https://github.com/pottava/basicauth-reverseproxy/blob/master/README-ja.md))

## Usage

### 1. Set environment variables

Environment Variables     | Description                                       | Required
------------------------- | ------------------------------------------------- | ---------------------
PROXY_PATTERNS            | Both virtual host & virtual path can be specified.| * (or PROXY_URL)
PROXY_URL                 | The URL to be proxied with this app.              | * (or PROXY_PATTERNS)
BASIC_AUTH_USER           | User for basic authentication.                    | 
BASIC_AUTH_PASS           | Password for basic authentication.                | 
APP_PORT                  | The port number to be assigned for listening.     | 
SSL_CERT_PATH             | TLS: cert.pem file path.                          | 
SSL_KEY_PATH              | TLS: key.pem file path.                           | 
CORS_ALLOW_ORIGIN         | CORS: a URI that may access the resource.         | 
CORS_ALLOW_METHODS        | CORS: Comma-delimited list of the allowed [HTTP request methods](https://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html). | 
CORS_ALLOW_HEADERS        | CORS: Comma-delimited list of the supported request headers. | 
CORS_MAX_AGE              | CORS: Maximum number of seconds the results of a preflight request can be cached. | 
ACCESS_LOG                | Send access logs to /dev/stdout. (default: false) | 
CONTENT_ENCODING          | Compress response data if the request allows. (default: false) |
HEALTHCHECK_PATH          | If it's specified, the path always returns 200 OK |
USE_X_FORWARDED_HOST      | Set original host to X-Forwarded-Host. (default: true) |

### 2. Run the application

`$ docker run -d -p 8080:80 -e PROXY_URL pottava/proxy`

* with basic auth:  

`$ docker run -d -p 8080:80 -e PROXY_URL -e BASIC_AUTH_USER -e BASIC_AUTH_PASS pottava/proxy`

* with TLS:  

`$ docker run -d -p 8080:80 -e PROXY_URL -e SSL_CERT_PATH -e SSL_KEY_PATH pottava/proxy`

* with CORS:

`$ docker run -d -p 8080:80 -e PROXY_URL -e CORS_ALLOW_ORIGIN -e CORS_ALLOW_METHODS -e CORS_ALLOW_HEADERS -e CORS_MAX_AGE pottava/proxy`

* with virtual hosts:  

`$ export PROXY_PATTERNS="/static=http://assets.cdn/,*.example.com=http://app.io/,*=http://sorry.com/"`  
`$ docker run -d -p 8080:80 -e PROXY_PATTERNS pottava/proxy`

* with docker-compose.yml:  

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

* with docker-compose.yml (virtual hosts & SSL/TLS):  

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

## Copyright and license

Code released under the [MIT license](https://github.com/pottava/basicauth-reverseproxy/blob/master/LICENSE).
